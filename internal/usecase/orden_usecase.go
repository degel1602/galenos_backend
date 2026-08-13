package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type ordenService struct {
	repo output.OrdenRepository
}

func NewOrdenService(repo output.OrdenRepository) *ordenService {
	return &ordenService{repo: repo}
}

func (s *ordenService) ListarPorCuenta(ctx context.Context, idCuentaAtencion int) ([]domain.OrdenMedica, error) {
	return s.repo.ListarPorCuenta(ctx, idCuentaAtencion)
}

func (s *ordenService) CrearOrden(ctx context.Context, orden domain.OrdenMedica, detalles []domain.DetalleOrden) error {
	return s.repo.CrearOrden(ctx, orden, detalles)
}
