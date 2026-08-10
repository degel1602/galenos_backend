// Package reniec implementa el adaptador de salida que consulta el servicio
// web SOAP de RENIEC (wsvmin.minsa.gob.pe), replicando la lógica del
// proyecto FastAPI.
package reniec

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// Config agrupa las credenciales y endpoint del servicio RENIEC.
type Config struct {
	App     string
	Usuario string
	Clave   string
	URL     string
	Timeout time.Duration
}

const (
	baseNS    = "http://schemas.xmlsoap.org/soap/envelope/"
	tempURINS = "http://tempuri.org/"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36"
)

var operaciones = map[string]string{
	"basico":   "obtenerDatosBasicos",
	"completo": "obtenerDatosCompletos",
}

type client struct {
	cfg  Config
	http *http.Client
}

// New crea el cliente SOAP de RENIEC.
func New(cfg Config) *client {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	return &client{cfg: cfg, http: &http.Client{Timeout: cfg.Timeout}}
}

func (c *client) Consultar(ctx context.Context, dni string, operacion string) (domain.ReniecResult, error) {
	metodo, ok := operaciones[operacion]
	if !ok {
		return domain.ReniecResult{}, domain.ErrInvalidReniecOperation
	}

	body := c.construirSOAP(metodo, dni)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.cfg.URL, bytes.NewBuffer(body))
	if err != nil {
		return domain.ReniecResult{}, fmt.Errorf("building reniec request: %w", err)
	}
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", tempURINS+metodo)
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "text/xml, application/xml, */*")

	resp, err := c.http.Do(req)
	if err != nil {
		return domain.ReniecResult{}, fmt.Errorf("calling reniec service: %w", err)
	}
	defer resp.Body.Close()

	contenido, readErr := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if readErr != nil {
		return domain.ReniecResult{}, fmt.Errorf("reading reniec response: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		detalle := truncate(contenido, 800)
		return domain.ReniecResult{}, fmt.Errorf("reniec service responded %d: %s", resp.StatusCode, detalle)
	}

	resultado, parseErr := parsearStrings(contenido)
	if parseErr != nil {
		return domain.ReniecResult{}, fmt.Errorf("parsing reniec response: %w", parseErr)
	}

	if esError(resultado) {
		codigo, mensaje := interpretarError(resultado)
		log.Printf("reniec error %s: %s", codigo, mensaje)
		return domain.ReniecResult{}, fmt.Errorf("reniec error: %s", mensaje)
	}

	return domain.ReniecResult{
		DNI:       dni,
		Operacion: operacion,
		Resultado: resultado,
		Datos:     interpretarDatos(resultado, operacion),
	}, nil
}

// interpretarDatos extrae los campos de la persona desde el arreglo crudo.
// El layout posicional depende de la operación SOAP; el frontend usa "completo".
func interpretarDatos(resultado []string, operacion string) domain.ReniecDatos {
	datos := domain.ReniecDatos{}
	idx := func(i int) string {
		if i >= 0 && i < len(resultado) {
			return strings.TrimSpace(resultado[i])
		}
		return ""
	}

	switch operacion {
	case "basico":
		// obtenerDatosBasicos: [2]=paterno, [3]=materno, [5]=nombres.
		datos.ApellidoPaterno = idx(2)
		datos.ApellidoMaterno = idx(3)
		datos.Nombres = idx(5)
		datos.FechaNacimiento = idx(20)
	default:
		// obtenerDatosCompletos: [4]=paterno, [5]=materno, [7]=nombres,
		// [29]=fecha de nacimiento (DD/MM/YYYY).
		datos.ApellidoPaterno = idx(4)
		datos.ApellidoMaterno = idx(5)
		datos.Nombres = idx(7)
		datos.FechaNacimiento = convertirFecha(idx(29))
	}

	partes := separarNombres(datos.Nombres)
	if len(partes) > 0 {
		datos.PrimerNombre = partes[0]
	}
	if len(partes) > 1 {
		datos.SegundoNombre = partes[1]
	}
	if len(partes) > 2 {
		datos.TercerNombre = partes[2]
	}

	return datos
}

// separarNombres divide los prenombres (e.g. "CARLOS MELECIO") en primer y
// segundo nombre; cualquier parte extra va al tercer nombre.
func separarNombres(nombres string) []string {
	partes := strings.Fields(nombres)
	if len(partes) <= 2 {
		return partes
	}
	tercero := strings.Join(partes[2:], " ")
	return []string{partes[0], partes[1], tercero}
}

// convertirFecha normaliza la fecha RENIEC (DD/MM/YYYY) a formato ISO YYYY-MM-DD
// para que el input tipo date del frontend la consuma de forma directa.
func convertirFecha(fecha string) string {
	fecha = strings.TrimSpace(fecha)
	if fecha == "" {
		return ""
	}
	partes := strings.Split(fecha, "/")
	if len(partes) != 3 {
		return fecha
	}
	dd, mm, aa := partes[0], partes[1], partes[2]
	if mm == "" || dd == "" || len(aa) != 4 || !isDigit(dd) || !isDigit(mm) {
		return fecha
	}
	return aa + "-" + mm + "-" + dd
}

// construirSOAP arma el sobre SOAP con las credenciales en el header,
// igual que el proyecto FastAPI.
func (c *client) construirSOAP(metodo string, nrodoc string) []byte {
	body := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns:xsd="http://www.w3.org/2001/XMLSchema"
 xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Header>
    <Credencialmq xmlns="http://tempuri.org/">
      <app>%s</app>
      <usuario>%s</usuario>
      <clave>%s</clave>
    </Credencialmq>
  </soap:Header>
  <soap:Body>
    <%s xmlns="http://tempuri.org/">
      <nrodoc>%s</nrodoc>
    </%s>
  </soap:Body>
</soap:Envelope>`,
		escapeXML(c.cfg.App),
		escapeXML(c.cfg.Usuario),
		escapeXML(c.cfg.Clave),
		metodo,
		escapeXML(nrodoc),
		metodo,
	)
	return []byte(body)
}

// parsearStrings recorre el XML de la respuesta y extrae el texto de todos
// los elementos <string>, que es el formato que RENIEC devuelve tanto para
// datos como para mensajes de error.
func parsearStrings(contenido []byte) ([]string, error) {
	decoder := xml.NewDecoder(bytes.NewReader(contenido))
	valores := make([]string, 0)

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "string" {
				var valor string
				if err := decoder.DecodeElement(&valor, &t); err != nil {
					return nil, err
				}
				valores = append(valores, valor)
			}
		}
	}

	return valores, nil
}

func esError(resultado []string) bool {
	return len(resultado) >= 2 &&
		resultado[0] != "" &&
		isDigit(resultado[0]) &&
		intVal(resultado[0]) >= 9000
}

func interpretarError(resultado []string) (string, string) {
	mensajes := make([]string, 0, len(resultado)-1)
	for _, m := range resultado[1:] {
		m = strings.TrimSpace(m)
		if m != "" {
			mensajes = append(mensajes, m)
		}
	}
	mensaje := strings.Join(mensajes, " ")
	if mensaje == "" {
		mensaje = "Error desconocido"
	}
	return resultado[0], mensaje
}

func isDigit(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func intVal(s string) int64 {
	v, _ := strconv.ParseInt(s, 10, 64)
	return v
}

func escapeXML(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&apos;",
	)
	return r.Replace(s)
}

func truncate(b []byte, max int) string {
	s := string(b)
	if len(s) > max {
		s = s[:max]
	}
	return s
}
