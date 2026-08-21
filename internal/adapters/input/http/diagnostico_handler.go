package httpadapter

import (
	"log"
	"net/http"
	"strconv"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/gin-gonic/gin"
)

type DiagnosticoHandler struct {
	useCase input.DiagnosticoUseCase
}

func NewDiagnosticoHandler(useCase input.DiagnosticoUseCase) *DiagnosticoHandler {
	return &DiagnosticoHandler{useCase: useCase}
}

// SearchDiagnosticos maneja GET /api/v1/diagnosticos/search.
// @Summary Buscar diagnósticos
// @Description Busca diagnósticos en base a un texto (filtro), idAtencion e idPaciente usando el SP usp_go_SelectDiagnosticos
// @Tags Diagnosticos
// @Accept json
// @Produce json
// @Param filtro query string false "Texto a buscar"
// @Param idAtencion query int false "ID de Atención"
// @Param idPaciente query int false "ID de Paciente"
// @Success 200 {object} Response{data=[]domain.DiagnosticoBusqueda}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /diagnosticos/search [get]
func (h *DiagnosticoHandler) SearchDiagnosticos(c *gin.Context) {
	filtro := c.Query("filtro")
	idAtencionStr := c.Query("idAtencion")
	idPacienteStr := c.Query("idPaciente")

	idAtencion, _ := strconv.Atoi(idAtencionStr)
	idPaciente, _ := strconv.Atoi(idPacienteStr)

	log.Printf("Buscando diagnosticos con: filtro=%q, idAtencion=%d, idPaciente=%d", filtro, idAtencion, idPaciente)

	results, err := h.useCase.SearchDiagnosticos(c.Request.Context(), filtro, idAtencion, idPaciente)
	if err != nil {
		log.Printf("Internal error searching diagnosticos: %v", err)
		respondError(c, http.StatusInternalServerError, "ERR_SEARCH_DIAG", "Error buscando diagnósticos: "+err.Error())
		return
	}

	if results == nil {
		results = make([]domain.DiagnosticoBusqueda, 0) // Para devolver [] en vez de null en JSON
	}

	respondSuccess(c, http.StatusOK, results)
}
