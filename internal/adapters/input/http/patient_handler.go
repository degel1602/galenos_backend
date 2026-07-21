package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// PatientHandler expone el puerto de entrada input.PatientService.
type PatientHandler struct {
	service input.PatientService
}

// NewPatientHandler inyecta el caso de uso en el adaptador HTTP.
func NewPatientHandler(service input.PatientService) *PatientHandler {
	return &PatientHandler{service: service}
}

// List maneja GET /api/v1/pacientes?page=1&pageSize=20.
func (h *PatientHandler) List(c *gin.Context) {
	page := shared.NewPageRequest(
		queryInt(c, "page", 1),
		queryInt(c, "pageSize", 20),
	)

	result, err := h.service.List(c.Request.Context(), page)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, result)
}

// GetByDocumentNumber maneja GET /api/v1/pacientes/:numDocumento.
func (h *PatientHandler) GetByDocumentNumber(c *gin.Context) {
	documentNumber := c.Param("numDocumento")

	patient, err := h.service.GetByDocumentNumber(c.Request.Context(), documentNumber)
	if err != nil {
		switch err {
		case domain.ErrInvalidDocumentNumber:
			respondError(c, http.StatusBadRequest, "INVALID_DOCUMENT_NUMBER", err.Error())
		case domain.ErrPatientNotFound:
			respondError(c, http.StatusNotFound, "PATIENT_NOT_FOUND", err.Error())
		default:
			respondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		}
		return
	}

	respondSuccess(c, http.StatusOK, patient)
}

func queryInt(c *gin.Context, key string, fallback int) int {
	raw := c.Query(key)
	if raw == "" {
		return fallback
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
