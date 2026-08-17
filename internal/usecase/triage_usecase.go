package usecase

import (
	"context"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type triageUseCase struct {
	repo output.TriageRepository
}

// NewTriageUseCase construye el caso de uso de registro de triaje.
func NewTriageUseCase(repo output.TriageRepository) input.TriageService {
	return &triageUseCase{repo: repo}
}

// CreateTriage delega la persistencia en el repositorio y devuelve el
// resultado informado por el procedimiento almacenado.
func (uc *triageUseCase) CreateTriage(ctx context.Context, triage *domain.Triage) (string, error) {
	result, err := uc.repo.Create(ctx, triage)
	if err != nil {
		return "", fmt.Errorf("registering triage: %w", err)
	}
	return result, nil
}

// ListTriage delega el listado (SP ListarTriaje_Emergencia) en el
// repositorio y devuelve los registros crudos.
func (uc *triageUseCase) ListTriage(ctx context.Context, params shared.TriageListParams) ([]map[string]any, error) {
	items, err := uc.repo.List(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing triages: %w", err)
	}
	return items, nil
}

// ListPendingAdmission delega en el repositorio (SP
// webGestionAtencion_E_H_BusquedaFiltrar) y devuelve los registros crudos.
func (uc *triageUseCase) ListPendingAdmission(ctx context.Context, params shared.TriageAdmisionParams) ([]map[string]any, error) {
	items, err := uc.repo.ListPendingAdmission(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("listing triages pending admission: %w", err)
	}
	return items, nil
}

// CreateAdmission delega la admisión (SP WebCrearAtencionDesdeTriaje) en
// el repositorio y devuelve el resultado informado por el SP.
func (uc *triageUseCase) CreateAdmission(ctx context.Context, admision *domain.AdmisionDesdeTriaje) (string, error) {
	result, err := uc.repo.CreateAdmission(ctx, admision)
	if err != nil {
		return "", fmt.Errorf("creating admission from triage: %w", err)
	}
	return result, nil
}

// GetReporte delega en el repositorio (SP WebSelectReporteTriaje) y
// devuelve los registros crudos del reporte.
func (uc *triageUseCase) GetReporte(ctx context.Context, params shared.TriageReporteParams) ([]map[string]any, error) {
	items, err := uc.repo.GetReporte(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("getting triage report: %w", err)
	}
	return items, nil
}

// GetFichaAdmision delega en el repositorio (SP webFichaEmergencia) y
// devuelve los datos crudos de la ficha.
func (uc *triageUseCase) GetFichaAdmision(ctx context.Context, params shared.FichaAdmisionParams) (*map[string]any, error) {
	item, err := uc.repo.GetFichaAdmision(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("getting admission record: %w", err)
	}
	return item, nil
}

// ListarMedicosPorEspecialidad delega en el repositorio (SP
// usp_go_MedicosFiltrarPorIdEspecialidad) y devuelve los médicos de la
// especialidad.
func (uc *triageUseCase) ListarMedicosPorEspecialidad(ctx context.Context, idEspecialidad int) ([]domain.MedicoFila, error) {
	items, err := uc.repo.ListarMedicosPorEspecialidad(ctx, idEspecialidad)
	if err != nil {
		return nil, fmt.Errorf("listing doctors by specialty: %w", err)
	}
	return items, nil
}
