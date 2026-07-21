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

	// GetByDocumentNumber busca un paciente por número de documento
	// invocando el procedimiento almacenado sp_listarPaciente. Retorna
	// domain.ErrPatientNotFound si no existe ningún registro.
	GetByDocumentNumber(ctx context.Context, documentNumber string) (*domain.Patient, error)
}
