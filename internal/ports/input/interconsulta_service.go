package input

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type InterconsultaService interface {
	ObtenerPorId(ctx context.Context, id int) (*domain.Interconsulta, error)
	ListarPorServicio(ctx context.Context, tipoServicio string) ([]domain.Interconsulta, error)
	Crear(ctx context.Context, interconsulta domain.Interconsulta) error
	ActualizarEstado(ctx context.Context, id int, estado string) error
	GuardarFirma(ctx context.Context, firma domain.FirmaInterconsulta) error
}
