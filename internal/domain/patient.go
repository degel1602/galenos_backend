package domain

import "time"

// Patient es un modelo de lectura del agregado "Paciente", reflejando
// exactamente las columnas expuestas por la consulta a la tabla pacientes
// (NroDocumento, ApellidoPaterno, ApellidoMaterno, PrimerNombre,
// SegundoNombre, TercerNombre) así como por el SP usp_go_listarpacientes
// (IdPaciente, NroHistoriaClinica, FechaNacimiento, TipoDocumento).
type Patient struct {
	PatientID         int64
	DocumentNumber    string
	DocumentType      *string
	DocIdentityID     *int64
	HistoryNumber     string
	PaternalSurname   string
	MaternalSurname   string
	FirstName         string
	SecondName        string
	ThirdName         string
	DateOfBirth       *time.Time
	HomeDistrictID    *int64
	HomeCenterID      *int64
	SexTypeID         *int64
	MaritalStatusID   *int64
	EducationDegreeID *int64
	HomeAddress       *string
	Phone             *string
}
