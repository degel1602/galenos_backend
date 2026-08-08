package httpadapter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// ReniecHandler expone el puerto de entrada input.ReniecService.
type ReniecHandler struct {
	service input.ReniecService
}

// NewReniecHandler inyecta el caso de uso de RENIEC en el adaptador HTTP.
func NewReniecHandler(service input.ReniecService) *ReniecHandler {
	return &ReniecHandler{service: service}
}

// Consultar maneja GET /api/v1/reniec/:nrodoc?operacion=basico|completo.
//
// @Summary Consulta datos de una persona en RENIEC
// @Description Invoca el servicio SOAP externo de RENIEC. La operación por defecto es "completo".
// @Tags Reniec
// @Produce json
// @Param nrodoc path string true "Número de documento (DNI)"
// @Param operacion query string false "Operación: basico | completo (default completo)"
// @Success 200 {object} apiResponse{data=domain.ReniecResult}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 502 {object} apiResponse{error=apiError} "Error al consultar RENIEC"
// @Router /reniec/{nrodoc} [get]
func (h *ReniecHandler) Consultar(c *gin.Context) {
	dni := c.Param("nrodoc")
	operacion := c.Query("operacion")
	if operacion == "" {
		operacion = "completo"
	}

	result, err := h.service.Consultar(c.Request.Context(), dni, operacion)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDocumentNumber),
			errors.Is(err, domain.ErrInvalidReniecOperation):
			respondError(c, http.StatusBadRequest, "INVALID_RENIEC_REQUEST", err.Error())
		default:
			respondError(c, http.StatusBadGateway, "RENIEC_UNAVAILABLE", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, result)
}
