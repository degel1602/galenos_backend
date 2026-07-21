// Package usecase implementa los puertos de entrada orquestando el dominio
// y los puertos de salida. Es la única capa que conoce ambos lados.
package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type appointmentUseCase struct {
	repo  output.AppointmentRepository
	clock func() time.Time
}

// NewAppointmentUseCase construye el caso de uso de citas médicas.
func NewAppointmentUseCase(repo output.AppointmentRepository) input.AppointmentService {
	return &appointmentUseCase{repo: repo, clock: time.Now}
}

func (uc *appointmentUseCase) Schedule(ctx context.Context, cmd input.ScheduleAppointmentCommand) (*domain.Appointment, error) {
	slot, err := domain.NewTimeSlot(cmd.StartsAt, cmd.EndsAt)
	if err != nil {
		return nil, err
	}

	appointment, err := domain.NewAppointment(
		domain.PatientID(cmd.PatientID),
		domain.DoctorID(cmd.DoctorID),
		slot,
		cmd.Reason,
		uc.clock(),
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Create(ctx, appointment); err != nil {
		return nil, fmt.Errorf("scheduling appointment: %w", err)
	}

	return appointment, nil
}

func (uc *appointmentUseCase) GetByID(ctx context.Context, id string) (*domain.Appointment, error) {
	if id == "" {
		return nil, domain.ErrInvalidAppointmentID
	}

	appointment, err := uc.repo.GetByID(ctx, domain.AppointmentID(id))
	if err != nil {
		return nil, fmt.Errorf("fetching appointment: %w", err)
	}

	return appointment, nil
}
