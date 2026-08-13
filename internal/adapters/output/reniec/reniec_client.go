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
	"unicode"

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

	// Log de depuración: el layout del arreglo no está documentado, así que
	// se registra completo para poder mapear las posiciones de cada campo.
	log.Printf("reniec consulta dni=%s operacion=%s resultado=%q", dni, operacion, resultado)

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
		datos.Sexo = detectarSexo(resultado)
		domicilio := extraerDomicilio(resultado)
		datos.Departamento = domicilio.Departamento
		datos.Provincia = domicilio.Provincia
		datos.Distrito = domicilio.Distrito
		datos.Direccion = domicilio.Direccion
		datos.Ubigeo = domicilio.Ubigeo
		datos.EstadoCivil = detectarEstadoCivil(resultado)
		// Posiciones conocidas de obtenerDatosCompletos: [30] nombre del padre,
		// [31] nombre de la madre, [16]/[17]/[18] ubigeo de nacimiento.
		datos.NombrePadre = idx(30)
		datos.NombreMadre = idx(31)
		nacimiento := extraerNacimiento(resultado)
		datos.DepartamentoNacimiento = nacimiento.Departamento
		datos.ProvinciaNacimiento = nacimiento.Provincia
		datos.DistritoNacimiento = nacimiento.Distrito
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

// detectarSexo busca el valor de sexo dentro del arreglo crudo de
// obtenerDatosCompletos. Posición conocida: [20] con el código 1/2; si el
// layout varía, recorre el arreglo buscando un token que parezca sexo
// ("MASCULINO", "FEMENINO", "M", "F" o dígitos 1/2) descartando las primeras
// posiciones reservadas a identidad. Si no encuentra nada concluyente,
// devuelve vacío (el frontend lo deja editable).
func detectarSexo(resultado []string) string {
	switch strings.ToUpper(strings.TrimSpace(tokenEn(resultado, 20))) {
	case "MASCULINO", "M", "1":
		return "MASCULINO"
	case "FEMENINO", "F", "2":
		return "FEMENINO"
	}
	for i := 10; i < len(resultado); i++ {
		token := strings.ToUpper(strings.TrimSpace(resultado[i]))
		switch token {
		case "MASCULINO", "FEMENINO", "M", "F":
			return token
		case "1":
			return "MASCULINO"
		case "2":
			return "FEMENINO"
		}
	}
	return ""
}

// --- Domicilio y estado civil ---

// domicilioReniec agrupa los datos de domicilio inferidos de la respuesta.
type domicilioReniec struct {
	Departamento string
	Provincia    string
	Distrito     string
	Direccion    string
	Ubigeo       string
}

// nacimientoReniec agrupa el ubigeo de nacimiento inferido de la respuesta.
type nacimientoReniec struct {
	Departamento string
	Provincia    string
	Distrito     string
}

// departamentosPeru son los nombres canónicos de los departamentos del Perú,
// usados para validar los valores posicionales de la respuesta RENIEC.
var departamentosPeru = []string{
	"AMAZONAS", "ANCASH", "APURIMAC", "AREQUIPA", "AYACUCHO", "CAJAMARCA",
	"CALLAO", "CUSCO", "HUANCAVELICA", "HUANUCO", "ICA", "JUNIN",
	"LA LIBERTAD", "LAMBAYEQUE", "LIMA", "LORETO", "MADRE DE DIOS",
	"MOQUEGUA", "PASCO", "PIURA", "PUNO", "SAN MARTIN", "TACNA",
	"TUMBES", "UCAYALI",
}

// estadosCivilReniec son los estados civiles que puede devolver RENIEC.
var estadosCivilReniec = []string{
	"SOLTERO", "CASADO", "VIUDO", "DIVORCIADO", "CONVIVIENTE", "SEPARADO",
}

// extraerDomicilio intenta recuperar el domicilio desde el arreglo crudo de
// obtenerDatosCompletos. Posiciones conocidas: [26] departamento, [27]
// provincia, [28] distrito y [36] dirección del domicilio. Como el layout
// puede variar entre versiones del servicio, primero se intenta con etiquetas
// ("DEPARTAMENTO: LIMA" o "DEPARTAMENTO" + valor), luego con las posiciones
// conocidas y al final con heurísticas validadas (ubigeo de 6 dígitos,
// nombres de departamento conocidos, etc.). Solo se devuelven valores
// concluyentes; lo que no se confirma queda vacío para que el frontend lo
// deje editable.
func extraerDomicilio(resultado []string) domicilioReniec {
	var dom domicilioReniec

	// 1) Tokens etiquetados.
	dom.Departamento = valorEtiquetado(resultado, "DEPARTAMENTO")
	dom.Provincia = valorEtiquetado(resultado, "PROVINCIA")
	dom.Distrito = valorEtiquetado(resultado, "DISTRITO")
	dom.Direccion = valorEtiquetado(resultado, "DIRECCION")
	dom.Ubigeo = valorEtiquetado(resultado, "UBIGEO")

	// 2) Posiciones conocidas de obtenerDatosCompletos.
	if dom.Departamento == "" {
		dom.Departamento = tokenEn(resultado, 26)
	}
	if dom.Provincia == "" {
		dom.Provincia = tokenEn(resultado, 27)
	}
	if dom.Distrito == "" {
		dom.Distrito = tokenEn(resultado, 28)
	}
	if dom.Direccion == "" {
		dom.Direccion = tokenEn(resultado, 36)
	}

	// 3) Heurísticas posicionales para lo que aún falta.
	departamentoIdx := -1
	for i := 20; i < len(resultado); i++ {
		token := strings.TrimSpace(resultado[i])
		if token == "" || esSinDatos(token) {
			continue
		}
		switch {
		case dom.Ubigeo == "" && esUbigeo(token):
			dom.Ubigeo = token
		case dom.Departamento == "" && esDepartamento(token):
			dom.Departamento = token
			departamentoIdx = i
		case dom.Direccion == "" && pareceDireccion(token):
			dom.Direccion = token
		}
	}

	if dom.Departamento != "" && (dom.Provincia == "" || dom.Distrito == "") {
		if departamentoIdx < 0 {
			departamentoIdx = buscarIndice(resultado, dom.Departamento)
		}
		// Tras el departamento, los tokens no vacíos siguientes suelen ser
		// provincia y distrito.
		subsiguientes := 0
		for j := departamentoIdx + 1; j < len(resultado) && subsiguientes < 2; j++ {
			candidato := strings.TrimSpace(resultado[j])
			if candidato == "" || esSinDatos(candidato) || esUbigeo(candidato) || esEstadoCivil(candidato) || pareceDireccion(candidato) {
				continue
			}
			if subsiguientes == 0 {
				dom.Provincia = candidato
			} else {
				dom.Distrito = candidato
			}
			subsiguientes++
		}
	}

	return dom
}

// extraerNacimiento recupera el ubigeo de nacimiento desde el arreglo crudo de
// obtenerDatosCompletos. Posiciones conocidas: [16] departamento, [17]
// provincia y [18] distrito. Primero intenta con etiquetas y luego con las
// posiciones; solo devuelve valores concluyentes (lo demás queda vacío).
func extraerNacimiento(resultado []string) nacimientoReniec {
	var nac nacimientoReniec

	// 1) Tokens etiquetados.
	nac.Departamento = valorEtiquetado(resultado, "DEPARTAMENTO NACIMIENTO")
	nac.Provincia = valorEtiquetado(resultado, "PROVINCIA NACIMIENTO")
	nac.Distrito = valorEtiquetado(resultado, "DISTRITO NACIMIENTO")

	// 2) Posiciones conocidas de obtenerDatosCompletos.
	if nac.Departamento == "" {
		nac.Departamento = tokenEn(resultado, 16)
	}
	if nac.Provincia == "" {
		nac.Provincia = tokenEn(resultado, 17)
	}
	if nac.Distrito == "" {
		nac.Distrito = tokenEn(resultado, 18)
	}

	return nac
}

// detectarEstadoCivil busca el estado civil de la persona en el arreglo crudo.
// Al igual que el sexo, la posición no está documentada, así que se valida
// contra los estados civiles conocidos.
func detectarEstadoCivil(resultado []string) string {
	for i := 20; i < len(resultado); i++ {
		token := strings.TrimSpace(resultado[i])
		if token == "" {
			continue
		}
		if esEstadoCivil(token) {
			return token
		}
	}
	return ""
}

// valorEtiquetado lee el valor de un token etiquetado. Soporta dos formatos:
// "DEPARTAMENTO: LIMA" (etiqueta y valor juntos) y "DEPARTAMENTO" "LIMA"
// (etiqueta y valor en tokens consecutivos).
func valorEtiquetado(resultado []string, etiqueta string) string {
	for i, token := range resultado {
		t := strings.ToUpper(strings.TrimSpace(token))
		if t != etiqueta && !strings.HasPrefix(t, etiqueta+":") {
			continue
		}
		resto := strings.Trim(strings.TrimPrefix(t, etiqueta), ": ")
		if resto != "" {
			return resto
		}
		if i+1 < len(resultado) {
			return strings.TrimSpace(resultado[i+1])
		}
	}
	return ""
}

// buscarIndice devuelve la posición del primer token que normalizado coincide
// con el valor buscado (desde la posición 20, donde suelen ir los datos de
// domicilio), o -1 si no lo encuentra.
func buscarIndice(resultado []string, valor string) int {
	n := normalizarToken(valor)
	for i := 20; i < len(resultado); i++ {
		if normalizarToken(resultado[i]) == n {
			return i
		}
	}
	return -1
}

// tokenEn devuelve el token en la posición i, vacío si está fuera de rango,
// es "SIN DATOS" o está en blanco (el servicio usa "SIN DATOS" como nulo).
func tokenEn(resultado []string, i int) string {
	if i < 0 || i >= len(resultado) {
		return ""
	}
	v := strings.TrimSpace(resultado[i])
	if v == "" || esSinDatos(v) {
		return ""
	}
	return v
}

// esSinDatos verifica si el token es el nulo que usa RENIEC.
func esSinDatos(token string) bool {
	t := strings.ToUpper(strings.TrimSpace(token))
	return t == "SIN DATOS" || t == "S/D" || t == "SIN DATO"
}

// esDepartamento verifica si el token es uno de los departamentos del Perú.
func esDepartamento(token string) bool {
	n := normalizarToken(token)
	for _, d := range departamentosPeru {
		if n == d {
			return true
		}
	}
	return false
}

// esEstadoCivil verifica si el token es un estado civil reconocido.
func esEstadoCivil(token string) bool {
	n := normalizarToken(token)
	for _, e := range estadosCivilReniec {
		if strings.HasPrefix(n, e) {
			return true
		}
	}
	return false
}

// esUbigeo verifica si el token es un código ubigeo (6 dígitos).
func esUbigeo(token string) bool {
	return len(token) == 6 && isDigit(token)
}

// pareceDireccion verifica si el token tiene forma de dirección: contiene
// dígitos y una vía conocida (AV., JR., MZ., URB., etc.) o texto con número.
func pareceDireccion(token string) bool {
	token = strings.TrimSpace(token)
	if esFecha(token) {
		return false
	}
	n := normalizarToken(token)
	if !strings.ContainsAny(n, "0123456789") {
		return false
	}
	for _, p := range []string{
		"AV", "JR", "MZ", "LT", "URB", "PSJ", "CALLE", "PASAJE",
		"PROLONGACION", "NRO", "BLOCK", "GRUPO", "DPTO", "KM",
	} {
		if strings.Contains(n, p) {
			return true
		}
	}
	return len(n) >= 8
}

// esFecha verifica si el token tiene el formato DD/MM/YYYY.
func esFecha(token string) bool {
	partes := strings.Split(token, "/")
	return len(partes) == 3 &&
		len(partes[2]) == 4 && isDigit(partes[0]) && isDigit(partes[1]) && isDigit(partes[2])
}

// normalizarToken normaliza un token para compararlo: mayúsculas, sin
// acentos, sin signos de puntuación y con espacios simples.
func normalizarToken(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	s = strings.NewReplacer(
		"Á", "A", "É", "E", "Í", "I", "Ó", "O", "Ú", "U", "Ü", "U", "Ñ", "N",
	).Replace(s)
	var b strings.Builder
	ultimoEspacio := false
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			ultimoEspacio = false
		} else if !ultimoEspacio {
			b.WriteByte(' ')
			ultimoEspacio = true
		}
	}
	return strings.TrimSpace(b.String())
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
