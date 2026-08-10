package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// ReniecClient es el puerto de salida para consultar el servicio web SOAP
// de RENIEC. La implementación concreta invoca el endpoint externo.
type ReniecClient interface {
	// Consultar consulta los datos de una persona por su número de
	// documento contra el servicio RENIEC. operacion es "basico" o
	// "completo".
	Consultar(ctx context.Context, dni string, operacion string) (domain.ReniecResult, error)
}
