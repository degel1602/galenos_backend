package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type sisUseCase struct {
	client output.SisClient
}

// NewSisUseCase construye el caso de uso de consulta al SIS.
func NewSisUseCase(client output.SisClient) input.SisService {
	return &sisUseCase{client: client}
}

func (uc *sisUseCase) ConsultarAfiliado(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error) {
	if params.DocumentNumber == "" {
		return domain.SisAfiliado{}, domain.ErrInvalidDocumentNumber
	}
	if params.TipoDocumento != 1 && params.TipoDocumento != 3 {
		return domain.SisAfiliado{}, domain.ErrInvalidDocumentType
	}
	if params.Opcion <= 0 {
		params.Opcion = 1
	}

	result, err := uc.client.ConsultarAfiliado(ctx, params)
	if err != nil {
		return domain.SisAfiliado{}, fmt.Errorf("consulting sis: %w", err)
	}

	return result, nil
}
