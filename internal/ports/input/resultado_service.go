package input

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ResultadoService interface {
	ListarResultadosLaboratorio(ctx context.Context, idPaciente int) ([]domain.Resultado, error)
	ListarResultadosImagenes(ctx context.Context, idPaciente int) ([]domain.Resultado, error)
}
