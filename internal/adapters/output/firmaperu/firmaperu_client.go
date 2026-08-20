// Package firmaperu implementa los adaptadores de salida del Firmador Web de
// la Plataforma Firma Perú (PCM/SGTD). El flujo es interactivo: este paquete
// obtiene un JWT del servicio generate-token y construye el JSON base64 de
// argumentos que el frontend pasa a firmaperu.min.js (startSignature) para
// que el Firmador local del usuario firme el documento con su DNIe.
package firmaperu

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// Config agrupa la URL del servicio generate-token, las credenciales del
// Firmador Web y el timeout de las peticiones HTTP.
type Config struct {
	TokenURL     string
	ClientID     string
	ClientSecret string
	Timeout      time.Duration
}

const (
	defaultTokenURL = "https://apps.firmaperu.gob.pe/admin/api/security/generate-token"
	maxResponse     = 1 << 20 // 1 MB (respuestas JSON de la plataforma)
	tokenMargin     = 60 * time.Second
)

// client implementa output.FirmaPeruClient con caché del JWT en memoria.
type client struct {
	cfg  Config
	http *http.Client

	mu    sync.Mutex
	token string
	exp   time.Time
}

// New crea el cliente del Firmador Web de Firma Perú.
func New(cfg Config) *client {
	if cfg.TokenURL == "" {
		cfg.TokenURL = defaultTokenURL
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 60 * time.Second
	}
	return &client{cfg: cfg, http: &http.Client{Timeout: cfg.Timeout}}
}

// GetToken devuelve un JWT válido del servicio generate-token, reutilizando
// el token en caché mientras no expire.
func (c *client) GetToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.exp.Add(-tokenMargin)) {
		return c.token, nil
	}

	token, exp, err := c.solicitarToken(ctx)
	if err != nil {
		return "", err
	}

	c.token = token
	c.exp = exp
	return token, nil
}

// solicitarToken llama al servicio generate-token con client_id y
// client_secret (form-urlencoded) y devuelve el token con su expiración.
func (c *client) solicitarToken(ctx context.Context) (string, time.Time, error) {
	if c.cfg.ClientID == "" || c.cfg.ClientSecret == "" {
		return "", time.Time{}, fmt.Errorf("firma peru client_id/client_secret are required")
	}

	form := url.Values{
		"client_id":     {c.cfg.ClientID},
		"client_secret": {c.cfg.ClientSecret},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("building token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("calling firma peru generate-token: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, maxResponse))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("reading firma peru token response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", time.Time{}, fmt.Errorf("firma peru generate-token responded %d: %s", resp.StatusCode, truncate(body, 400))
	}

	token := strings.TrimSpace(string(body))
	return token, expiracionJWT(token), nil
}

// BuildParamBase64 serializa los parámetros de firma a JSON y lo codifica en
// base64, tal como lo espera firmaperu.min.js.
func (c *client) BuildParamBase64(params domain.FirmaPeruSignParams) (string, error) {
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("marshalling firma peru params: %w", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// BuildWrapperBase64 serializa el wrapper (param_url/param_token/
// document_extension) a JSON base64, tal como lo espera startSignature(port,
// param) de firmaperu.min.js.
func (c *client) BuildWrapperBase64(wrapper domain.FirmaPeruParamWrapper) (string, error) {
	jsonBytes, err := json.Marshal(wrapper)
	if err != nil {
		return "", fmt.Errorf("marshalling firma peru wrapper: %w", err)
	}
	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// expiracionJWT extrae el claim "exp" (segundos) de un JWT HS256; devuelve
// time.Time{} si no lo puede leer.
func expiracionJWT(token string) time.Time {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return time.Time{}
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return time.Time{}
	}

	var claims struct {
		Exp int64 `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Exp == 0 {
		return time.Time{}
	}

	return time.Unix(claims.Exp, 0)
}

func truncate(b []byte, max int) string {
	s := string(b)
	if len(s) > max {
		s = s[:max]
	}
	return s
}
