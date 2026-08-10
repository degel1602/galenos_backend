package domain

import "time"

// Triage representa el registro de triaje que se persiste invocando el
// procedimiento almacenado webTab_PacienteTriajeAgregar. Los campos
// opcionales son punteros para distinguir "no enviado" (NULL) de un valor
// cero; el SP decide cómo interpretarlos.
type Triage struct {
	IDTriaje                  *int64     `json:"idTriaje"`
	DocIdentityID             *int64     `json:"idDocIdentidad"`
	DocumentNumber            *string    `json:"nroDocumento"`
	PaternalSurname           *string    `json:"apellidoPaterno"`
	MaternalSurname           *string    `json:"apellidoMaterno"`
	FirstName                 *string    `json:"primerNombre"`
	SecondName                *string    `json:"segundoNombre"`
	ThirdName                 *string    `json:"tercerNombre"`
	SexTypeID                 *int64     `json:"idSexo"`
	DateOfBirth               *time.Time `json:"fechaNacimiento"`
	Phone                     *string    `json:"telefono"`
	HomeDepartmentID          *int64     `json:"idDepartamentoDomicilio"`
	HomeProvinceID            *int64     `json:"idProvinciaDomicilio"`
	HomeDistrictID            *int64     `json:"idDistritoDomicilio"`
	HomeCommunityID           *int64     `json:"idComunidadDomicilio"`
	HomeAddress               *string    `json:"direccion"`
	IsTrafficAccident         *int64     `json:"esAccidenteTransito"`
	FundingSourceID           *int64     `json:"idFuenteFinanciamiento"`
	Email                     *string    `json:"email"`
	MaritalStatusID           *int64     `json:"idEstadoCivil"`
	HeartRate                 *int64     `json:"frecCardiaca"`
	Temperature               *float64   `json:"temperatura"`
	BloodPressure             *string    `json:"presionArterial"`
	OxygenSaturation          *int64     `json:"saturacion"`
	RespiratoryRate           *int64     `json:"frecRespiratoria"`
	FiO2                      *int64     `json:"fiO2"`
	Weight                    *float64   `json:"peso"`
	Height                    *float64   `json:"talla"`
	BMI                       *float64   `json:"imc"`
	EvolutionTimeQuantity     *int64     `json:"tiempoEvolucionCantidad"`
	EvolutionTimeQuantityUnit *string    `json:"tiempoEvolucionCantidadUnidad"`
	PainScale                 *int64     `json:"escalaDolor"`
	GlasgowScale              *int64     `json:"escalaGlasgow"`
	PriorityTypeID            *int64     `json:"idTipoPrioridad"`
	ServiceID                 *int64     `json:"idServicio"`
	Motivo                    *string    `json:"motivo"`
	IsPregnant                *int64     `json:"idGestante"`
	ArrivalStateID            *int64     `json:"idEstadollego"`
	Photo                     *string    `json:"foto"`
	EmployeeID                *int64     `json:"idEmpleado"`
}
