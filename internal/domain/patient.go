package domain

// Patient es un modelo de lectura del agregado "Paciente", reflejando
// exactamente las columnas expuestas por la consulta a la tabla pacientes
// (NroDocumento, ApellidoPaterno, ApellidoMaterno, PrimerNombre,
// SegundoNombre, TercerNombre).
type Patient struct {
	DocumentNumber  string `json:"documentNumber"`
	PaternalSurname string `json:"paternalSurname"`
	MaternalSurname string `json:"maternalSurname"`
	FirstName       string `json:"firstName"`
	SecondName      string `json:"secondName"`
	ThirdName       string `json:"thirdName"`
}
