package domain

import "time"

// Patient es un modelo de lectura del agregado "Paciente", reflejando
// exactamente las columnas expuestas por la consulta a la tabla pacientes
// (NroDocumento, ApellidoPaterno, ApellidoMaterno, PrimerNombre,
// SegundoNombre, TercerNombre) así como por el SP usp_go_listarpacientes
// (IdPaciente, NroHistoriaClinica, FechaNacimiento).
type Patient struct {
	PatientID         int64      `json:"patientId"`
	DocumentNumber    string     `json:"documentNumber"`
	HistoryNumber     string     `json:"historyNumber"`
	PaternalSurname   string     `json:"paternalSurname"`
	MaternalSurname   string     `json:"maternalSurname"`
	FirstName         string     `json:"firstName"`
	SecondName        string     `json:"secondName"`
	ThirdName         string     `json:"thirdName"`
	DateOfBirth       *time.Time `json:"dateOfBirth,omitempty"`
	HomeDistrictID    *int64     `json:"homeDistrictId,omitempty"`
	HomeCenterID      *int64     `json:"homeCenterId,omitempty"`
	SexTypeID         *int64     `json:"sexTypeId,omitempty"`
	MaritalStatusID   *int64     `json:"maritalStatusId,omitempty"`
	EducationDegreeID *int64     `json:"educationDegreeId,omitempty"`
	HomeAddress       *string    `json:"homeAddress,omitempty"`
	Phone             *string    `json:"phone,omitempty"`
}
