package httpadapter

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
)

// AppointmentHandler expone el puerto de entrada input.AppointmentService
// como endpoints REST.
type AppointmentHandler struct {
	service input.AppointmentService
}

// NewAppointmentHandler inyecta el caso de uso en el adaptador HTTP.
func NewAppointmentHandler(service input.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{service: service}
}

// Create maneja POST /api/v1/appointments.
//
// @Summary Agenda una nueva cita médica
// @Description Valida los datos del paciente, médico y horario, y agenda la cita si el médico está disponible.
// @Tags Appointments
// @Accept json
// @Produce json
// @Param request body createAppointmentRequest true "Datos de la cita a agendar"
// @Success 201 {object} apiResponse{data=domain.Appointment}
// @Failure 400 {object} apiResponse{error=apiError} "Validación inválida"
// @Failure 409 {object} apiResponse{error=apiError} "El médico no está disponible en ese horario"
// @Router /appointments [post]
func (h *AppointmentHandler) Create(c *gin.Context) {
	var req createAppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	appointment, err := h.service.Schedule(c.Request.Context(), input.ScheduleAppointmentCommand{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		StartsAt:  req.StartsAt,
		EndsAt:    req.EndsAt,
		Reason:    req.Reason,
	})
	if err != nil {
		status, code := mapDomainError(err)
		respondError(c, status, code, err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, appointment)
}

// GetByID maneja GET /api/v1/appointments/:id.
//
// @Summary Obtiene una cita médica por id
// @Description Devuelve la cita identificada por su id.
// @Tags Appointments
// @Produce json
// @Param id path string true "ID de la cita"
// @Success 200 {object} apiResponse{data=domain.Appointment}
// @Failure 404 {object} apiResponse{error=apiError} "Cita no encontrada"
// @Router /appointments/{id} [get]
func (h *AppointmentHandler) GetByID(c *gin.Context) {
	appointment, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		status, code := mapDomainError(err)
		respondError(c, status, code, err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, appointment)
}

// mapDomainError traduce errores de dominio a códigos HTTP. Los errores no
// reconocidos se tratan como 500 para no filtrar detalles internos.
func mapDomainError(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrInvalidPatient),
		errors.Is(err, domain.ErrInvalidDoctor),
		errors.Is(err, domain.ErrInvalidTimeSlot),
		errors.Is(err, domain.ErrInvalidTimeSlotDuration),
		errors.Is(err, domain.ErrAppointmentInPast),
		errors.Is(err, domain.ErrInvalidAppointmentID):
		return http.StatusBadRequest, "INVALID_APPOINTMENT"
	case errors.Is(err, domain.ErrDoctorNotAvailable):
		return http.StatusConflict, "DOCTOR_NOT_AVAILABLE"
	case errors.Is(err, domain.ErrAppointmentNotFound):
		return http.StatusNotFound, "APPOINTMENT_NOT_FOUND"
	default:
		return http.StatusInternalServerError, "INTERNAL_ERROR"
	}
}
