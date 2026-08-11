package domain

import "time"

// SisAfiliado es la respuesta normalizada de una consulta de afiliación al
// servicio web SOAP del SIS (ConsultarAfiliadoSIS). Mapea los campos de
// ResultQueryAsegurado que devuelve la operación.
type SisAfiliado struct {
	IdError            string
	Resultado          string
	TipoDocumento      string
	NroDocumento       string
	ApePaterno         string
	ApeMaterno         string
	Nombres            string
	FecAfiliacion      string
	EESS               string
	DescEESS           string
	EESSUbigeo         string
	DescEESSUbigeo     string
	Regimen            string
	TipoSeguro         string
	DescTipoSeguro     string
	Contrato           string
	FecCaducidad       string
	Estado             string
	Tabla              string
	IdNumReg           string
	Genero             string
	FecNacimiento      string
	IdUbigeo           string
	Direccion          string
	Disa               string
	TipoFormato        string
	NroContrato        string
	Correlativo        string
	IdPlan             string
	IdGrupoPoblacional string
	MsgConfidencial    string
}

// SisAfiliacion representa el registro de afiliación que se persiste
// invocando el SP webSisFiliacionesGestionar. Campos opcionales a NULL.
type SisAfiliacion struct {
	IDSiasis                *int64
	Codigo                  *string
	AfiliacionDisa          *string
	TipoFormato             *string
	NroFormato              *string
	NroIntegrante           *string
	DocumentoTipo           *string
	CodigoEstablAdscripcion *string
	AfiliacionFecha         *time.Time
	Paterno                 *string
	Materno                 *string
	PNombre                 *string
	ONombres                *string
	Genero                  *string
	FNacimiento             *time.Time
	IdDistritoDomicilio     *string
	Estado                  *string
	Fbaja                   *string
	DocumentoNumero         *string
	MotivoBaja              *string
	FbajaOK                 *time.Time
	DescEESS                *string
	DescEESSUbigeo          *string
	Regimen                 *string
	TipoSeguro              *string
	DescTipoSeguro          *string
	Contrato                *string
	IdPlan                  *string
	IdGrupoPoblacional      *string
	MsgConfidencial         *string
	IdUsuarioAuditoria      *int64
}
