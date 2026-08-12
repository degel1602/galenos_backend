package domain

// PatientListItem representa a un paciente en la bandeja
type PatientListItem struct {
	IdRegAtencion int    `json:"idRegAtencion"`
	IdPaciente    int    `json:"idPaciente"`
	Historia      string `json:"historia"`
	Nombre        string `json:"nombre"`
	Edad          string `json:"edad"`
	Sexo          string `json:"sexo"`
	Ubicacion     string `json:"ubicacion"`
	Estado        string `json:"estado"`
}

// EvolucionFirma representa una evolución guardada en formato JSON / Base64
type EvolucionFirma struct {
	IdRegAtencion      int    `json:"idRegAtencion"`
	IdFirma            int    `json:"idFirma"`
	NombreDocumento    string `json:"nombreDocumento"`
	DataB64            string `json:"dataB64"`
	IdEmpleadoRegistra int    `json:"idEmpleadoRegistra"`
	FechaRegistro      string `json:"fechaRegistro"`
}
