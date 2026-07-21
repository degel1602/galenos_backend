package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PatientService es el puerto de entrada para consulta de pacientes.
type PatientService interface {
	List(ctx context.Context, page shared.PageRequest) (shared.PageResponse[domain.Patient], error)

	// GetByDocumentNumber busca un paciente por su número de documento.
	GetByDocumentNumber(ctx context.Context, documentNumber string) (domain.Patient, error)
}
