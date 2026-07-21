// Package input define los puertos de entrada (driving ports): el contrato
// que cualquier adaptador de entrada (HTTP, CLI, gRPC, jobs...) usa para
// invocar los casos de uso del dominio, sin conocer su implementación.
package input

import (
	"context"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// ScheduleAppointmentCommand transporta la intención de agendar una cita
// desde el adaptador de entrada hacia el caso de uso.
type ScheduleAppointmentCommand struct {
	PatientID string
	DoctorID  string
	StartsAt  time.Time
	EndsAt    time.Time
	Reason    string
}

// AppointmentService es el puerto de entrada para el módulo de citas médicas.
type AppointmentService interface {
	Schedule(ctx context.Context, cmd ScheduleAppointmentCommand) (*domain.Appointment, error)
	GetByID(ctx context.Context, id string) (*domain.Appointment, error)
}
