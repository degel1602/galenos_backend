package input

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type EvolucionUseCase interface {
	CreateEvolucion(ctx context.Context, evolucion *domain.Evolucion) error
	GetEvolucionesByPaciente(ctx context.Context, pacienteID int64) ([]domain.Evolucion, error)
}
