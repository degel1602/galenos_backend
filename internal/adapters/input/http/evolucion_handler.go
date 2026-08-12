package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type EvolucionHandler struct {
	useCase input.EvolucionUseCase
}

func NewEvolucionHandler(useCase input.EvolucionUseCase) *EvolucionHandler {
	return &EvolucionHandler{useCase: useCase}
}

func (h *EvolucionHandler) HandleCreateEvolucion(c *gin.Context) {
	var req EvolucionCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "invalid_request", "JSON inválido o mal formado")
		return
	}

	evolucion := &domain.Evolucion{
		IDPaciente:       &req.IDPaciente,
		IDMedico:         &req.IDMedico,
		IDCita:           req.IDCita,
		MotivoAtencion:   req.MotivoAtencion,
		DetalleMotivo:    req.DetalleMotivo,
		SubjetivoDetalle: req.SubjetivoDetalle,
		PresionArterial:  req.PresionArterial,
		FrecCardiaca:     req.FrecCardiaca,
		FrecRespiratoria: req.FrecRespiratoria,
		Temperatura:      req.Temperatura,
		SaturacionO2:     req.SaturacionO2,
		Peso:             req.Peso,
		Talla:            req.Talla,
		IMC:              req.IMC,
		Glucemia:         req.Glucemia,
		EstadoClinico:    req.EstadoClinico,
		Pronostico:       req.Pronostico,
		PlanDetalle:      req.PlanDetalle,
	}

	for _, dxReq := range req.Diagnosticos {
		evolucion.Diagnosticos = append(evolucion.Diagnosticos, domain.DiagnosticoEvolucion{
			CIE10:       dxReq.CIE10,
			Descripcion: dxReq.Descripcion,
			Tipo:        dxReq.Tipo,
			Condicion:   dxReq.Condicion,
			Estado:      dxReq.Estado,
		})
	}

	if err := h.useCase.CreateEvolucion(c.Request.Context(), evolucion); err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	respondSuccess(c, http.StatusCreated, gin.H{
		"id_evolucion": evolucion.IDEvolucion,
	})
}

func (h *EvolucionHandler) HandleListEvoluciones(c *gin.Context) {
	idParam := c.Param("pacienteId")
	pacienteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid_param", "ID de paciente inválido")
		return
	}

	evoluciones, err := h.useCase.GetEvolucionesByPaciente(c.Request.Context(), pacienteID)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, evoluciones)
}
