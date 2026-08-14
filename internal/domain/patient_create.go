package domain

import "time"

// PatientCreate transporta los datos de un paciente nuevo hacia el SP
// WebPacienteAgregar_E_H. Los campos opcionales son punteros; un puntero nil
// se envía como NULL (o su default) al procedimiento.
type PatientCreate struct {
	PaternalSurname   *string
	MaternalSurname   *string
	FirstName         *string
	SecondName        *string
	DateOfBirth       *time.Time
	DocIdentityID     *int64
	DocumentNumber    *string
	Phone             *string
	HomeAddress       *string
	SexTypeID         *int64
	MaritalStatusID   *int64
	BirthDistrictID   *int64
	HomeDistrictID    *int64
	InsuranceTypeID   *int64
	OriginID          *int64
	EducationDegreeID *int64
	OccupationTypeID  *int64
	HistoryNumber     *string
	StateID           *int64
}
