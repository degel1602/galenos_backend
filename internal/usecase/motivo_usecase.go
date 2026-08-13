package usecase

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type motivoService struct {
	repo output.MotivoRepository
}

func NewMotivoService(repo output.MotivoRepository) *motivoService {
	return &motivoService{repo: repo}
}

func (s *motivoService) ListarMotivos(ctx context.Context, idRegAtencion int) ([]domain.MotivoAtencion, error) {
	return s.repo.ListarPorAtencion(ctx, idRegAtencion)
}

func (s *motivoService) GuardarMotivo(ctx context.Context, idRegAtencion int, motivo string) error {
	return s.repo.Guardar(ctx, idRegAtencion, motivo)
}
