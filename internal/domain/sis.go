package domain

import "time"

// SisAfiliado es la respuesta normalizada de una consulta de afiliación al
// servicio web SOAP del SIS (ConsultarAfiliadoSIS). Mapea los campos de
// ResultQueryAsegurado que devuelve la operación.
type SisAfiliado struct {
	IdError            string `json:"idError"`
	Resultado          string `json:"resultado"`
	TipoDocumento      string `json:"tipoDocumento"`
	NroDocumento       string `json:"nroDocumento"`
	ApePaterno         string `json:"apePaterno"`
	ApeMaterno         string `json:"apeMaterno"`
	Nombres            string `json:"nombres"`
	FecAfiliacion      string `json:"fecAfiliacion"`
	EESS               string `json:"eess"`
	DescEESS           string `json:"descEESS"`
	EESSUbigeo         string `json:"eessUbigeo"`
	DescEESSUbigeo     string `json:"descEessUbigeo"`
	Regimen            string `json:"regimen"`
	TipoSeguro         string `json:"tipoSeguro"`
	DescTipoSeguro     string `json:"descTipoSeguro"`
	Contrato           string `json:"contrato"`
	FecCaducidad       string `json:"fecCaducidad"`
	Estado             string `json:"estado"`
	Tabla              string `json:"tabla"`
	IdNumReg           string `json:"idNumReg"`
	Genero             string `json:"genero"`
	FecNacimiento      string `json:"fecNacimiento"`
	IdUbigeo           string `json:"idUbigeo"`
	Direccion          string `json:"direccion"`
	Disa               string `json:"disa"`
	TipoFormato        string `json:"tipoFormato"`
	NroContrato        string `json:"nroContrato"`
	Correlativo        string `json:"correlativo"`
	IdPlan             string `json:"idPlan"`
	IdGrupoPoblacional string `json:"idGrupoPoblacional"`
	MsgConfidencial    string `json:"msgConfidencial"`
}

// SisAfiliacion representa el registro de afiliación que se persiste
// invocando el SP webSisFiliacionesGestionar. Campos opcionales a NULL.
type SisAfiliacion struct {
	IDSiasis                *int64     `json:"idSiasis"`
	Codigo                  *string    `json:"codigo"`
	AfiliacionDisa          *string    `json:"afiliacionDisa"`
	TipoFormato             *string    `json:"afiliacionTipoFormato"`
	NroFormato              *string    `json:"afiliacionNroFormato"`
	NroIntegrante           *string    `json:"afiliacionNroIntegrante"`
	DocumentoTipo           *string    `json:"documentoTipo"`
	CodigoEstablAdscripcion *string    `json:"codigoEstablAdscripcion"`
	AfiliacionFecha         *time.Time `json:"afiliacionFecha"`
	Paterno                 *string    `json:"paterno"`
	Materno                 *string    `json:"materno"`
	PNombre                 *string    `json:"pNombre"`
	ONombres                *string    `json:"oNombres"`
	Genero                  *string    `json:"genero"`
	FNacimiento             *time.Time `json:"fNacimiento"`
	IdDistritoDomicilio     *string    `json:"idDistritoDomicilio"`
	Estado                  *string    `json:"estado"`
	Fbaja                   *string    `json:"fBaja"`
	DocumentoNumero         *string    `json:"documentoNumero"`
	MotivoBaja              *string    `json:"motivoBaja"`
	FbajaOK                 *time.Time `json:"fBajaOk"`
	DescEESS                *string    `json:"descEESS"`
	DescEESSUbigeo          *string    `json:"descEessUbigeo"`
	Regimen                 *string    `json:"regimen"`
	TipoSeguro              *string    `json:"tipoSeguro"`
	DescTipoSeguro          *string    `json:"descTipoSeguro"`
	Contrato                *string    `json:"contrato"`
	IdPlan                  *string    `json:"idPlan"`
	IdGrupoPoblacional      *string    `json:"idGrupoPoblacional"`
	MsgConfidencial         *string    `json:"msgConfidencial"`
	IdUsuarioAuditoria      *int64     `json:"idUsuarioAuditoria"`
}
