package input

import (
	"context"
	"github.com/galenos-pro/appointments-api/internal/domain"
)

type OrdenService interface {
	ListarPorCuenta(ctx context.Context, idCuentaAtencion int) ([]domain.OrdenMedica, error)
	CrearOrden(ctx context.Context, orden domain.OrdenMedica, detalles []domain.DetalleOrden) error
}
