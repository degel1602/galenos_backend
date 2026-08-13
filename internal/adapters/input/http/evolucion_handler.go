package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type EvolucionHandler struct {
	service input.EvolucionService
}

func NewEvolucionHandler(service input.EvolucionService) *EvolucionHandler {
	return &EvolucionHandler{service: service}
}

// @Summary List patients for the medical evolution tray
// @Description Returns a list of patients currently active
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/v1/evoluciones/pacientes [get]
func (h *EvolucionHandler) HandleListPatients(c *gin.Context) {
	patients, err := h.service.GetPatientTray(c.Request.Context())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "EVOL_PATIENTS_ERR", "Error obteniendo pacientes")
		return
	}
	respondSuccess(c, http.StatusOK, patients)
}

// @Summary Get evolutions for a patient
// @Description Returns the saved medical evolutions for a given registration ID
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pacienteId path int true "ID of the Registration / Encounter"
// @Router /api/v1/evoluciones/paciente/{pacienteId} [get]
func (h *EvolucionHandler) HandleListEvoluciones(c *gin.Context) {
	idRegAtencionStr := c.Param("pacienteId")
	idRegAtencion, err := strconv.Atoi(idRegAtencionStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de atención inválido")
		return
	}

	evolutions, err := h.service.GetEvoluciones(c.Request.Context(), idRegAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "EVOL_GET_ERR", "Error obteniendo evoluciones")
		return
	}

	respondSuccess(c, http.StatusOK, evolutions)
}

type SaveEvolucionRequest struct {
	IdRegAtencion int    `json:"idRegAtencion" binding:"required"`
	DataB64       string `json:"dataB64" binding:"required"`
}

// @Summary Save an evolution
// @Description Saves a new medical evolution
// @Tags Evoluciones
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body SaveEvolucionRequest true "Evolution data"
// @Router /api/v1/evoluciones [post]
func (h *EvolucionHandler) HandleCreateEvolucion(c *gin.Context) {
	var req SaveEvolucionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	idEmpleado := 1
	if val, exists := c.Get("idEmpleado"); exists {
		if id, ok := val.(float64); ok {
			idEmpleado = int(id)
		} else if idStr, ok := val.(string); ok {
			id, _ := strconv.Atoi(idStr)
			idEmpleado = id
		} else if idInt, ok := val.(int); ok {
			idEmpleado = idInt
		}
	}

	err := h.service.SaveEvolucion(c.Request.Context(), req.IdRegAtencion, idEmpleado, req.DataB64)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "EVOL_SAVE_ERR", "Error guardando la evolución")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Evolución guardada correctamente"})
}
