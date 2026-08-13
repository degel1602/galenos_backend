package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type InterconsultaHandler struct {
	service input.InterconsultaService
}

func NewInterconsultaHandler(service input.InterconsultaService) *InterconsultaHandler {
	return &InterconsultaHandler{service: service}
}

// @Summary Get interconsulta by ID
// @Description Returns the interconsulta details for a given ID
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la Interconsulta"
// @Router /api/v1/interconsultas/{id} [get]
func (h *InterconsultaHandler) HandleObtenerPorId(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	ic, err := h.service.ObtenerPorId(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_GET_ERR", "Error obteniendo interconsulta")
		return
	}

	respondSuccess(c, http.StatusOK, ic)
}

// @Summary List interconsultas by service
// @Description Returns the interconsultas for a given service type
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tipoServicio path string true "Tipo de Servicio"
// @Router /api/v1/interconsultas/servicio/{tipoServicio} [get]
func (h *InterconsultaHandler) HandleListarPorServicio(c *gin.Context) {
	tipoServicio := c.Param("tipoServicio")
	lista, err := h.service.ListarPorServicio(c.Request.Context(), tipoServicio)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_LIST_ERR", "Error listando interconsultas")
		return
	}
	respondSuccess(c, http.StatusOK, lista)
}

type CreateInterconsultaRequest struct {
	IdAtencionOrigen int    `json:"idAtencionOrigen" binding:"required"`
	IdEspecialidad   int    `json:"idEspecialidad" binding:"required"`
	IdMedicoDestino  int    `json:"idMedicoDestino" binding:"required"`
	Motivo           string `json:"motivo" binding:"required"`
}

// @Summary Create an interconsulta
// @Description Creates a new interconsulta request
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateInterconsultaRequest true "Interconsulta data"
// @Router /api/v1/interconsultas [post]
func (h *InterconsultaHandler) HandleCrear(c *gin.Context) {
	var req CreateInterconsultaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	ic := domain.Interconsulta{
		IdAtencionOrigen: req.IdAtencionOrigen,
		IdEspecialidad:   req.IdEspecialidad,
		IdMedicoDestino:  req.IdMedicoDestino,
		Motivo:           req.Motivo,
	}

	err := h.service.Crear(c.Request.Context(), ic)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_CREATE_ERR", "Error creando interconsulta")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Interconsulta creada correctamente"})
}

type UpdateEstadoRequest struct {
	Estado string `json:"estado" binding:"required"`
}

// @Summary Update interconsulta state
// @Description Updates the state of an existing interconsulta
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la Interconsulta"
// @Param request body UpdateEstadoRequest true "Estado data"
// @Router /api/v1/interconsultas/{id}/estado [put]
func (h *InterconsultaHandler) HandleActualizarEstado(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	var req UpdateEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	err = h.service.ActualizarEstado(c.Request.Context(), id, req.Estado)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_UPDATE_ERR", "Error actualizando estado")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Estado actualizado correctamente"})
}

type SignInterconsultaRequest struct {
	DataB64 string `json:"dataB64" binding:"required"`
}

// @Summary Sign an interconsulta
// @Description Saves the signature for an interconsulta
// @Tags Interconsultas
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID de la Interconsulta"
// @Param request body SignInterconsultaRequest true "Firma data"
// @Router /api/v1/interconsultas/{id}/firma [post]
func (h *InterconsultaHandler) HandleGuardarFirma(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID inválido")
		return
	}

	var req SignInterconsultaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	idEmpleado := 1 // Dummy ID si no hay token real en pruebas, en pro debe extraerse del context de JWT
	if val, exists := c.Get("idEmpleado"); exists {
		if ide, ok := val.(float64); ok {
			idEmpleado = int(ide)
		} else if idStr, ok := val.(string); ok {
			ide, _ := strconv.Atoi(idStr)
			idEmpleado = ide
		} else if idInt, ok := val.(int); ok {
			idEmpleado = idInt
		}
	}

	firma := domain.FirmaInterconsulta{
		IdInterconsulta: id,
		DataB64:         req.DataB64,
		IdEmpleadoFirma: idEmpleado,
	}

	err = h.service.GuardarFirma(c.Request.Context(), firma)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTER_SIGN_ERR", "Error guardando firma")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Firma guardada correctamente"})
}
