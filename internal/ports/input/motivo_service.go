package input

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type MotivoService interface {
	ListarMotivos(ctx context.Context, idRegAtencion int) ([]domain.MotivoAtencion, error)
	GuardarMotivo(ctx context.Context, idRegAtencion int, motivo string) error
}
