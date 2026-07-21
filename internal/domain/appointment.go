// Package domain contiene las entidades y reglas de negocio del dominio
// hospitalario. No debe importar nada de HTTP, SQL ni ningún framework.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// AppointmentID identifica de forma única una cita médica.
type AppointmentID string

// PatientID identifica a un paciente dentro del sistema.
type PatientID string

// DoctorID identifica a un médico dentro del sistema.
type DoctorID string

// AppointmentStatus representa el ciclo de vida de una cita médica.
type AppointmentStatus string

const (
	StatusScheduled AppointmentStatus = "SCHEDULED"
	StatusConfirmed AppointmentStatus = "CONFIRMED"
	StatusCancelled AppointmentStatus = "CANCELLED"
	StatusCompleted AppointmentStatus = "COMPLETED"
)

const (
	minAppointmentDuration = 10 * time.Minute
	maxAppointmentDuration = 4 * time.Hour
)

// TimeSlot es un Value Object: inmutable y validado en su construcción.
type TimeSlot struct {
	StartsAt time.Time `json:"startsAt"`
	EndsAt   time.Time `json:"endsAt"`
}

// NewTimeSlot valida que el rango horario sea coherente antes de crearlo.
func NewTimeSlot(startsAt, endsAt time.Time) (TimeSlot, error) {
	if !endsAt.After(startsAt) {
		return TimeSlot{}, ErrInvalidTimeSlot
	}
	duration := endsAt.Sub(startsAt)
	if duration < minAppointmentDuration || duration > maxAppointmentDuration {
		return TimeSlot{}, ErrInvalidTimeSlotDuration
	}
	return TimeSlot{StartsAt: startsAt, EndsAt: endsAt}, nil
}

// Appointment es la entidad raíz del agregado "Cita Médica".
type Appointment struct {
	ID        AppointmentID     `json:"id"`
	PatientID PatientID         `json:"patientId"`
	DoctorID  DoctorID          `json:"doctorId"`
	Slot      TimeSlot          `json:"slot"`
	Status    AppointmentStatus `json:"status"`
	Reason    string            `json:"reason"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}

// NewAppointment construye una cita nueva aplicando las reglas de negocio
// esenciales del dominio. `now` se inyecta explícitamente para mantener la
// entidad pura y testeable (sin depender de time.Now() internamente).
func NewAppointment(patientID PatientID, doctorID DoctorID, slot TimeSlot, reason string, now time.Time) (*Appointment, error) {
	if patientID == "" {
		return nil, ErrInvalidPatient
	}
	if doctorID == "" {
		return nil, ErrInvalidDoctor
	}
	if slot.StartsAt.Before(now) {
		return nil, ErrAppointmentInPast
	}

	return &Appointment{
		ID:        AppointmentID(uuid.NewString()),
		PatientID: patientID,
		DoctorID:  doctorID,
		Slot:      slot,
		Status:    StatusScheduled,
		Reason:    reason,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// Cancel aplica la regla de negocio: una cita completada no puede cancelarse.
func (a *Appointment) Cancel(now time.Time) error {
	switch a.Status {
	case StatusCompleted:
		return ErrCannotCancelCompleted
	case StatusCancelled:
		return ErrAlreadyCancelled
	}
	a.Status = StatusCancelled
	a.UpdatedAt = now
	return nil
}

// Reschedule aplica un nuevo horario, respetando las mismas reglas que la creación.
func (a *Appointment) Reschedule(newSlot TimeSlot, now time.Time) error {
	if a.Status == StatusCancelled || a.Status == StatusCompleted {
		return ErrCannotRescheduleClosedAppointment
	}
	if newSlot.StartsAt.Before(now) {
		return ErrAppointmentInPast
	}
	a.Slot = newSlot
	a.UpdatedAt = now
	return nil
}
