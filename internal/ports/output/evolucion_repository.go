package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type EvolucionRepository interface {
	ListPatients(ctx context.Context, fini, ffin string) ([]domain.PatientListItem, error)
	ListEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error)
	SaveEvolucion(ctx context.Context, evolucion domain.EvolucionFirma) error
}
