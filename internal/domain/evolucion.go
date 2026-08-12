package domain

import "time"

// DiagnosticoEvolucion representa un diagnóstico dentro del SOAP de la evolución médica.
type DiagnosticoEvolucion struct {
	IDEvolucion *int64
	CIE10       string
	Descripcion string
	Tipo        string
	Condicion   string
	Estado      string
}

// Evolucion representa la nota clínica SOAP completa.
// Los campos opcionales son punteros para distinguir "no enviado" (NULL) de valor cero.
type Evolucion struct {
	IDEvolucion *int64
	IDPaciente  *int64
	IDMedico    *int64
	IDCita      *int64
	Fecha       *time.Time

	// Información General y Motivo
	MotivoAtencion *string
	DetalleMotivo  *string

	// Subjetivo (S)
	SubjetivoDetalle *string

	// Objetivo (O) - Signos Vitales
	PresionArterial  *string
	FrecCardiaca     *int
	FrecRespiratoria *int
	Temperatura      *float64
	SaturacionO2     *int
	Peso             *float64
	Talla            *float64
	IMC              *float64
	Glucemia         *int

	// Evaluación (A)
	EstadoClinico *string
	Pronostico    *string

	// Diagnósticos
	Diagnosticos []DiagnosticoEvolucion

	// Plan (P)
	PlanDetalle *string
}
