package domain

import "errors"

// Errores de dominio para la firma digital con el Firmador Web de Firma Perú.
var (
	ErrFirmaPeruDocumentRequired   = errors.New("firma peru document is required")
	ErrInvalidFirmaPeruFormat      = errors.New("invalid firma peru signature format")
	ErrInvalidFirmaPeruLevel       = errors.New("invalid firma peru signature level")
	ErrFirmaPeruDocumentNotFound   = errors.New("firma peru document not found")
	ErrFirmaPeruSignedNotAvailable = errors.New("firma peru signed document not available yet")
	ErrFirmaPeruInvalidParamToken  = errors.New("firma peru param_token is invalid")
)

// Formatos de firma soportados por el Firmador Web.
const (
	FirmaPeruFormatPAdES = "PAdES" // Documentos PDF
	FirmaPeruFormatXAdES = "XAdES" // Documentos XML
	FirmaPeruFormatCAdES = "CAdES" // Firma separada
)

// Niveles de firma soportados por el servicio.
const (
	FirmaPeruLevelB   = "B"   // Firma básica
	FirmaPeruLevelT   = "T"   // Firma con sello de tiempo (requiere TSA)
	FirmaPeruLevelLTA = "LTA" // Firma con datos de validación a largo plazo
)

// FirmaPeruParamWrapper es el objeto JSON (base64) que el frontend pasa a
// startSignature(port, param). El Firmador local, tras decodificarlo, hace
// POST a ParamURL con ParamToken y espera el JSON base64 con los parámetros
// de firma (FirmaPeruSignParams).
type FirmaPeruParamWrapper struct {
	ParamURL          string `json:"param_url"`
	ParamToken        string `json:"param_token"`
	DocumentExtension string `json:"document_extension"`
}

// FirmaPeruSolicitud es el estado que el backend conserva por proceso de
// firma: el token de único uso (param_token), la URL pública con la que el PC
// del firmante alcanza la API y los parámetros de firma definidos por el
// frontend.
type FirmaPeruSolicitud struct {
	ParamToken    string
	PublicBaseURL string
	Params        FirmaPeruSignParams
}

// FirmaPeruSignParams agrupa los parámetros de firma que el Firmador espera
// (objeto JSON codificado en base64 que devuelve el endpoint param_url). El
// endpoint POST de param_url usa el nombre de campo bachtOperation (así, con
// typo, en la API oficial). DocumentToSign, UploadDocumentSigned y Token los
// rellena el caso de uso: DocumentToSign es la URL que sirve el documento por
// firmar, UploadDocumentSigned la URL donde el Firmador devuelve el documento
// firmado y Token el JWT del servicio generate-token.
type FirmaPeruSignParams struct {
	SignatureFormat        string `json:"signatureFormat"`
	SignatureLevel         string `json:"signatureLevel"`
	SignaturePackaging     string `json:"signaturePackaging"`
	DocumentToSign         string `json:"documentToSign"`
	CertificateFilter      string `json:"certificateFilter"`
	WebTsa                 string `json:"webTsa"`
	UserTsa                string `json:"userTsa"`
	PasswordTsa            string `json:"passwordTsa"`
	Theme                  string `json:"theme"`
	VisiblePosition        bool   `json:"visiblePosition"`
	ContactInfo            string `json:"contactInfo"`
	SignatureReason        string `json:"signatureReason"`
	BatchOperation         bool   `json:"bachtOperation"`
	OneByOne               bool   `json:"oneByOne"`
	SignatureStyle         int    `json:"signatureStyle"`
	ImageToStamp           string `json:"imageToStamp"`
	StampTextSize          int    `json:"stampTextSize"`
	StampWordWrap          int    `json:"stampWordWrap"`
	Role                   string `json:"role"`
	StampPage              int    `json:"stampPage"`
	PositionX              int    `json:"positionx"`
	PositionY              int    `json:"positiony"`
	UploadDocumentSigned   string `json:"uploadDocumentSigned"`
	CertificationSignature bool   `json:"certificationSignature"`
	Token                  string `json:"token"`
}

// FirmaPeruDocument es un documento en bytes con su nombre original.
type FirmaPeruDocument struct {
	Name    string
	Content []byte
}

// FirmaPeruInitRequest es la solicitud para iniciar un proceso de firma:
// el documento por firmar, los parámetros de configuración y, opcionalmente,
// la imagen PNG del estampado (p. ej. el logo del hospital) que el Firmador
// dibujará como representación visible de la firma. PublicBaseURL se usa para
// construir param_url/documentToSign/uploadDocumentSigned (debe ser accesible
// desde el equipo del firmante).
type FirmaPeruInitRequest struct {
	Document      []byte
	DocumentName  string
	PublicBaseURL string
	StampImage    []byte
	Params        FirmaPeruSignParams
}

// FirmaPeruInitResult es el resultado de iniciar la firma: la cadena base64
// del wrapper (param_url/param_token/document_extension) que el frontend pasa
// a startSignature y el UUID que identifica el proceso.
type FirmaPeruInitResult struct {
	ParamBase64      string
	DocumentNameUUID string
}

// FirmaPeruLoteInitRequest es la solicitud para iniciar una firma por lote
// (varios documentos firmados con un solo PIN del DNIe). El Firmador exige que
// documentToSign devuelva un archivo 7z con todos los documentos y que
// uploadDocumentSigned reciba el 7z con los firmados.
type FirmaPeruLoteInitRequest struct {
	Documents     []FirmaPeruDocument
	PublicBaseURL string
	StampImage    []byte
	Params        FirmaPeruSignParams
}
