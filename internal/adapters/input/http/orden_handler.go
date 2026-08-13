package httpadapter

import (
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type OrdenHandler struct {
	service input.OrdenService
}

func NewOrdenHandler(service input.OrdenService) *OrdenHandler {
	return &OrdenHandler{service: service}
}

// @Summary Get ordenes for a patient account
// @Description Returns the medical orders for a given account ID
// @Tags Ordenes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param idCuentaAtencion path int true "ID de la Cuenta de Atencion"
// @Router /api/v1/ordenes/cuenta/{idCuentaAtencion} [get]
func (h *OrdenHandler) HandleListOrdenes(c *gin.Context) {
	idStr := c.Param("idCuentaAtencion")
	idCuentaAtencion, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_ID", "ID de cuenta de atención inválido")
		return
	}

	ordenes, err := h.service.ListarPorCuenta(c.Request.Context(), idCuentaAtencion)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ORDEN_GET_ERR", "Error obteniendo órdenes médicas")
		return
	}

	respondSuccess(c, http.StatusOK, ordenes)
}

type CreateOrdenRequest struct {
	IdRegAtencion int                   `json:"idRegAtencion" binding:"required"`
	IdMedico      int                   `json:"idMedico" binding:"required"`
	Observacion   string                `json:"observacion"`
	Detalles      []domain.DetalleOrden `json:"detalles" binding:"required"`
}

// @Summary Create an orden medica
// @Description Creates a new medical order in a patient's evolution
// @Tags Ordenes
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateOrdenRequest true "Orden data"
// @Router /api/v1/ordenes [post]
func (h *OrdenHandler) HandleCreateOrden(c *gin.Context) {
	var req CreateOrdenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "INVALID_BODY", "Cuerpo de petición inválido")
		return
	}

	orden := domain.OrdenMedica{
		IdRegAtencion: req.IdRegAtencion,
		IdMedico:      req.IdMedico,
		Observacion:   req.Observacion,
	}

	err := h.service.CrearOrden(c.Request.Context(), orden, req.Detalles)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ORDEN_SAVE_ERR", "Error guardando la orden médica")
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"message": "Orden médica guardada correctamente"})
}
