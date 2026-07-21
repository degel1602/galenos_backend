package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PatientRepository es el puerto de salida para lectura de pacientes.
type PatientRepository interface {
	// List retorna la página solicitada de pacientes y el total de
	// registros existentes (para calcular TotalPages en el llamador).
	List(ctx context.Context, page shared.PageRequest) ([]domain.Patient, int, error)
}
