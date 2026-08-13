package output

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type MotivoRepository interface {
	ListarPorAtencion(ctx context.Context, idRegAtencion int) ([]domain.MotivoAtencion, error)
	Guardar(ctx context.Context, idRegAtencion int, motivo string) error
}
