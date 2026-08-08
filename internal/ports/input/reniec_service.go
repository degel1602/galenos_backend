package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// ReniecService es el puerto de entrada para consultar el servicio RENIEC.
type ReniecService interface {
	Consultar(ctx context.Context, dni string, operacion string) (domain.ReniecResult, error)
}
