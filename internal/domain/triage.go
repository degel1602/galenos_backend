package domain

import "time"

// Triage representa el registro de triaje que se persiste invocando el
// procedimiento almacenado webTab_PacienteTriajeAgregar. Los campos
// opcionales son punteros para distinguir "no enviado" (NULL) de un valor
// cero; el SP decide cómo interpretarlos.
type Triage struct {
	IDTriaje                  *int64
	DocIdentityID             *int64
	DocumentNumber            *string
	PaternalSurname           *string
	MaternalSurname           *string
	FirstName                 *string
	SecondName                *string
	ThirdName                 *string
	SexTypeID                 *int64
	DateOfBirth               *time.Time
	Phone                     *string
	HomeDepartmentID          *int64
	HomeProvinceID            *int64
	HomeDistrictID            *int64
	HomeCommunityID           *int64
	HomeAddress               *string
	IsTrafficAccident         *int64
	FundingSourceID           *int64
	Email                     *string
	MaritalStatusID           *int64
	HeartRate                 *int64
	Temperature               *float64
	BloodPressure             *string
	OxygenSaturation          *int64
	RespiratoryRate           *int64
	FiO2                      *int64
	Weight                    *float64
	Height                    *float64
	BMI                       *float64
	EvolutionTimeQuantity     *int64
	EvolutionTimeQuantityUnit *string
	PainScale                 *int64
	GlasgowScale              *int64
	PriorityTypeID            *int64
	ServiceID                 *int64
	Motivo                    *string
	IsPregnant                *int64
	ArrivalStateID            *int64
	Photo                     *string
	EmployeeID                *int64
}
