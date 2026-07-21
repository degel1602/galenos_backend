// Package output define los puertos de salida (driven ports): el contrato
// que el dominio/casos de uso exigen a cualquier adaptador de persistencia
// (SQL Server, Postgres, memoria...) sin conocer los detalles de la BD.
package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// AppointmentRepository es el puerto de salida para persistencia de citas.
type AppointmentRepository interface {
	// Create persiste una nueva cita. La verificación de disponibilidad del
	// médico (conflicto de horario) se resuelve de forma atómica dentro del
	// adaptador, ya que requiere una transacción contra la base de datos.
	Create(ctx context.Context, appointment *domain.Appointment) error

	// GetByID recupera una cita por su identificador. Debe retornar
	// domain.ErrAppointmentNotFound si no existe.
	GetByID(ctx context.Context, id domain.AppointmentID) (*domain.Appointment, error)
}
