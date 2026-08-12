package usecase

import (
	"context"
	"errors"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type evolucionUseCase struct {
	repo output.EvolucionRepository
}

func NewEvolucionUseCase(repo output.EvolucionRepository) input.EvolucionUseCase {
	return &evolucionUseCase{
		repo: repo,
	}
}

func (uc *evolucionUseCase) CreateEvolucion(ctx context.Context, evolucion *domain.Evolucion) error {
	if evolucion == nil {
		return errors.New("evolucion is required")
	}
	if evolucion.IDPaciente == nil || *evolucion.IDPaciente <= 0 {
		return errors.New("IDPaciente is required and must be positive")
	}
	if evolucion.IDMedico == nil || *evolucion.IDMedico <= 0 {
		return errors.New("IDMedico is required and must be positive")
	}

	return uc.repo.Save(ctx, evolucion)
}

func (uc *evolucionUseCase) GetEvolucionesByPaciente(ctx context.Context, pacienteID int64) ([]domain.Evolucion, error) {
	if pacienteID <= 0 {
		return nil, errors.New("invalid paciente ID")
	}

	return uc.repo.FindByPacienteID(ctx, pacienteID)
}
