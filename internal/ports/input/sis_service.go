package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// SisService es el puerto de entrada para consultar el servicio SIS.
type SisService interface {
	// ConsultarAfiliado trae el paciente afiliado por su número de documento.
	ConsultarAfiliado(ctx context.Context, params shared.SISAfiliadoParams) (domain.SisAfiliado, error)
}
