package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type reniecUseCase struct {
	client output.ReniecClient
}

// NewReniecUseCase construye el caso de uso de consulta a RENIEC.
func NewReniecUseCase(client output.ReniecClient) input.ReniecService {
	return &reniecUseCase{client: client}
}

func (uc *reniecUseCase) Consultar(ctx context.Context, dni string, operacion string) (domain.ReniecResult, error) {
	if dni == "" {
		return domain.ReniecResult{}, domain.ErrInvalidDocumentNumber
	}
	if operacion == "" {
		operacion = "completo"
	}

	result, err := uc.client.Consultar(ctx, dni, operacion)
	if err != nil {
		return domain.ReniecResult{}, fmt.Errorf("consulting reniec: %w", err)
	}

	return result, nil
}
