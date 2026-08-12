package output

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type EvolucionRepository interface {
	Save(ctx context.Context, evolucion *domain.Evolucion) error
	FindByPacienteID(ctx context.Context, pacienteID int64) ([]domain.Evolucion, error)
}
