package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type patientUseCase struct {
	repo output.PatientRepository
}

// NewPatientUseCase construye el caso de uso de consulta de pacientes.
func NewPatientUseCase(repo output.PatientRepository) input.PatientService {
	return &patientUseCase{repo: repo}
}

func (uc *patientUseCase) List(ctx context.Context, page shared.PageRequest) (shared.PageResponse[domain.Patient], error) {
	patients, total, err := uc.repo.List(ctx, page)
	if err != nil {
		return shared.PageResponse[domain.Patient]{}, fmt.Errorf("listing patients: %w", err)
	}
	return shared.NewPageResponse(patients, page, total), nil
}
