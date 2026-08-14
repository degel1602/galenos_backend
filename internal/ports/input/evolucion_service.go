package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type EvolucionService interface {
	GetPatientTray(ctx context.Context, fini, ffin string) ([]domain.PatientListItem, error)
	GetEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error)
	SaveEvolucion(ctx context.Context, idRegAtencion int, idEmpleadoRegistra int, dataB64 string) error
}
