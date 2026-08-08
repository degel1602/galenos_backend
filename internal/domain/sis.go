package domain

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
