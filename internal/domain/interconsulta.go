package domain

type Interconsulta struct {
	IdInterconsulta  int    `json:"idInterconsulta"`
	IdAtencionOrigen int    `json:"idAtencionOrigen"`
	IdEspecialidad   int    `json:"idEspecialidad"`
	IdMedicoDestino  int    `json:"idMedicoDestino"`
	Motivo           string `json:"motivo"`
	FechaSolicitud   string `json:"fechaSolicitud"`
	Estado           string `json:"estado"`
}

type FirmaInterconsulta struct {
	IdInterconsulta int    `json:"idInterconsulta"`
	DataB64         string `json:"dataB64"`
	IdEmpleadoFirma int    `json:"idEmpleadoFirma"`
}
