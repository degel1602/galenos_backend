package httpadapter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// SisHandler expone el puerto de entrada input.SisService.
type SisHandler struct {
	service input.SisService
}

// NewSisHandler inyecta el caso de uso de SIS en el adaptador HTTP.
func NewSisHandler(service input.SisService) *SisHandler {
	return &SisHandler{service: service}
}

// ConsultarAfiliado maneja GET /api/v1/sis/afiliado/:nrodoc.
//
// @Summary Consulta el paciente afiliado en el SIS
// @Description Invoca el servicio SOAP externo del SIS (GetSession + ConsultarAfiliadoFuaE) para traer los datos del afiliado.
// @Tags SIS
// @Produce json
// @Param nrodoc path string true "Número de documento del afiliado"
// @Param strTipoDocumento query int false "Tipo de documento: 1 = DNI (default), 3 = Carnet de Extranjería"
// @Param intOpcion query int false "intOpcion del WS (tipo de consulta contratado; default 1)"
// @Param strDisa query string false "strDisa del WS"
// @Param strTipoFormato query string false "strTipoFormato del WS"
// @Param strNroContrato query string false "strNroContrato del WS"
// @Param strCorrelativo query string false "strCorrelativo del WS"
// @Success 200 {object} apiResponse{data=domain.SisAfiliado}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 502 {object} apiResponse{error=apiError} "Error al consultar el SIS"
// @Router /sis/afiliado/{nrodoc} [get]
func (h *SisHandler) ConsultarAfiliado(c *gin.Context) {
	params := shared.SISAfiliadoParams{
		DocumentNumber: c.Param("nrodoc"),
		TipoDocumento:  1,
		Opcion:         1,
		Disa:           c.Query("strDisa"),
		TipoFormato:    c.Query("strTipoFormato"),
		NroContrato:    c.Query("strNroContrato"),
		Correlativo:    c.Query("strCorrelativo"),
	}
	if raw := c.Query("strTipoDocumento"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || (parsed != 1 && parsed != 3) {
			respondError(c, http.StatusBadRequest, "INVALID_SIS_REQUEST", "strTipoDocumento debe ser 1 (DNI) o 3 (Carnet de Extranjería)")
			return
		}
		params.TipoDocumento = int(parsed)
	} else if raw := c.Query("tipoDocumento"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || (parsed != 1 && parsed != 3) {
			respondError(c, http.StatusBadRequest, "INVALID_SIS_REQUEST", "tipoDocumento debe ser 1 (DNI) o 3 (Carnet de Extranjería)")
			return
		}
		params.TipoDocumento = int(parsed)
	}
	if raw := c.Query("intOpcion"); raw != "" {
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || parsed <= 0 {
			respondError(c, http.StatusBadRequest, "INVALID_SIS_REQUEST", "intOpcion debe ser un entero positivo")
			return
		}
		params.Opcion = int(parsed)
	}

	result, err := h.service.ConsultarAfiliado(c.Request.Context(), params)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidDocumentNumber),
			errors.Is(err, domain.ErrInvalidDocumentType):
			respondError(c, http.StatusBadRequest, "INVALID_SIS_REQUEST", err.Error())
		default:
			respondError(c, http.StatusBadGateway, "SIS_UNAVAILABLE", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, result)
}
