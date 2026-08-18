package httpadapter

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/galenos-pro/appointments-api/internal/ports/input"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

// TriageHandler expone el puerto de entrada input.TriageService.
type TriageHandler struct {
	service input.TriageService
}

// NewTriageHandler inyecta el caso de uso de triaje en el adaptador HTTP.
func NewTriageHandler(service input.TriageService) *TriageHandler {
	return &TriageHandler{service: service}
}

// Create maneja POST /api/v1/triaje.
//
// @Summary Registra un triaje de paciente
// @Description Persiste el triaje del paciente invocando el SP webTab_PacienteTriajeAgregar. Devuelve el @Resultado informado por el procedimiento.
// @Tags Triaje
// @Accept json
// @Produce json
// @Param request body createTriajeRequest true "Datos del triaje a registrar"
// @Success 200 {object} apiResponse{data=object} "Resultado informado por el SP"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al registrar el triaje"
// @Router /triaje [post]
func (h *TriageHandler) Create(c *gin.Context) {
	var req createTriajeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	result, err := h.service.CreateTriage(c.Request.Context(), req.toDomain())
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_REGISTER_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

// List maneja GET /api/v1/triaje.
//
// @Summary Lista los triajes registrados
// @Description Devuelve los triajes del período filtrado (SP ListarTriaje_Emergencia). El filtro de texto se aplica en el SP; si no se envía, se pasa vacío.
// @Tags Triaje
// @Produce json
// @Param fini query string true "Fecha inicial (YYYY-MM-DD)"
// @Param ffin query string true "Fecha final (YYYY-MM-DD)"
// @Param filtro query string false "Texto de búsqueda"
// @Param derivadoAServicio query int false "Servicio de derivación (-100 para todos)"
// @Param idEstado query int false "Estado del triaje"
// @Success 200 {object} apiResponse{data=[]object}
// @Failure 400 {object} apiResponse{error=apiError} "Parâmetros inválidos"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar los triajes"
// @Router /triaje [get]
func (h *TriageHandler) List(c *gin.Context) {
	params := shared.TriageListParams{
		FechaInicio:       c.Query("fini"),
		FechaFin:          c.Query("ffin"),
		Filtro:            c.Query("filtro"),
		DerivadoAServicio: -100,
		IdEstado:          -100,
	}

	if params.FechaInicio == "" || params.FechaFin == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fini y ffin son obligatorios (YYYY-MM-DD)")
		return
	}

	if raw := c.Query("derivadoAServicio"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "derivadoAServicio debe ser un entero")
			return
		}
		params.DerivadoAServicio = int(v)
	}

	if raw := c.Query("idEstado"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEstado debe ser un entero")
			return
		}
		params.IdEstado = int(v)
	}

	items, err := h.service.ListTriage(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// ListPendingAdmission maneja GET /api/v1/triaje/pendientes-admision.
//
// @Summary Lista pacientes con triaje sin admisión
// @Description Devuelve los pacientes con triaje que aún no han sido admisionados para una fecha (SP webGestionAtencion_E_H_BusquedaFiltrar). Los filtros opcionales se envían como 0 para no aplicarlos.
// @Tags Triaje
// @Produce json
// @Param fecha query string true "Fecha (YYYY-MM-DD)"
// @Param filtro query string false "Texto de búsqueda (documento o nombre)"
// @Param nroCta query int false "Número de cuenta de atención (0 = todos)"
// @Param idDepartamento query int false "Id departamento (0 = todos)"
// @Param idEspecialidad query int false "Id especialidad (0 = todas)"
// @Param idServicio query int false "Id servicio (0 = todos)"
// @Param idTipoServicio query int false "Id tipo de servicio"
// @Success 200 {object} apiResponse{data=[]object}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar los triajes sin admisión"
// @Router /triaje/pendientes-admision [get]
func (h *TriageHandler) ListPendingAdmission(c *gin.Context) {
	params := shared.TriageAdmisionParams{
		Fecha:  c.Query("fecha"),
		Filtro: c.Query("filtro"),
	}

	if params.Fecha == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fecha es obligatoria (YYYY-MM-DD)")
		return
	}

	intFields := []struct {
		raw   string
		dest  *int
		label string
	}{
		{c.Query("nroCta"), &params.NroCta, "nroCta"},
		{c.Query("idDepartamento"), &params.IdDepartamento, "idDepartamento"},
		{c.Query("idEspecialidad"), &params.IdEspecialidad, "idEspecialidad"},
		{c.Query("idServicio"), &params.IdServicio, "idServicio"},
		{c.Query("idTipoServicio"), &params.IdTipoServicio, "idTipoServicio"},
	}
	for _, f := range intFields {
		if f.raw == "" {
			continue
		}
		v, err := strconv.ParseInt(f.raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", f.label+" debe ser un entero")
			return
		}
		*f.dest = int(v)
	}

	items, err := h.service.ListPendingAdmission(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_PENDING_ADMISSION_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// CreateAdmission maneja POST /api/v1/triaje/admision.
//
// @Summary Admisiona un paciente con triaje
// @Description Crea la atención (admisión) de un paciente desde su triaje invocando el SP WebCrearAtencionDesdeTriaje. Devuelve el @Resultado informado por el procedimiento.
// @Tags Triaje
// @Accept json
// @Produce json
// @Param request body createAdmissionFromTriageRequest true "Datos de la admisión a partir del triaje"
// @Success 200 {object} apiResponse{data=object} "Resultado informado por el SP"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al admisionar al paciente"
// @Router /triaje/admision [post]
func (h *TriageHandler) CreateAdmission(c *gin.Context) {
	var req createAdmissionFromTriageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	domainObj := req.toDomain()
	if idEmpleado := c.GetInt("idEmpleado"); idEmpleado != 0 {
		empID := int64(idEmpleado)
		domainObj.IDEmpleado = &empID
	}

	result, err := h.service.CreateAdmission(c.Request.Context(), domainObj)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ADMISSION_FROM_TRIAGE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

// GetReporte maneja GET /api/v1/triaje/reporte.
//
// @Summary Genera el reporte de triaje
// @Description Devuelve el reporte de triaje (SP WebSelectReporteTriaje). Si el filtro no se aplica, enviar -100.
// @Tags Triaje
// @Produce json
// @Param id query int false "Id del triaje (-100 para todos)"
// @Param idPaciente query int false "Id del paciente (-100 para todos)"
// @Success 200 {object} apiResponse{data=[]object}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 500 {object} apiResponse{error=apiError} "Error al generar el reporte"
// @Router /triaje/reporte [get]
func (h *TriageHandler) GetReporte(c *gin.Context) {
	params := shared.TriageReporteParams{
		IDTriaje:   -100,
		IDPaciente: -100,
	}

	if raw := c.Query("id"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id debe ser un entero")
			return
		}
		params.IDTriaje = int(v)
	}

	if raw := c.Query("idPaciente"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idPaciente debe ser un entero")
			return
		}
		params.IDPaciente = int(v)
	}

	items, err := h.service.GetReporte(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_REPORT_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// GetFichaAdmision maneja GET /api/v1/triaje/ficha-admision.
//
// @Summary Obtiene la ficha de admisión de un paciente
// @Description Devuelve los datos del paciente y datos adicionales para generar la ficha de admisión (SP webFichaEmergencia).
// @Tags Triaje
// @Produce json
// @Param idCuentaAtencion query int true "Id de la cuenta de atención"
// @Success 200 {object} apiResponse{data=object}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 500 {object} apiResponse{error=apiError} "Error al obtener la ficha"
// @Router /triaje/ficha-admision [get]
func (h *TriageHandler) GetFichaAdmision(c *gin.Context) {
	raw := c.Query("idCuentaAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idCuentaAtencion debe ser un entero positivo")
		return
	}

	item, err := h.service.GetFichaAdmision(c.Request.Context(), shared.FichaAdmisionParams{IdCuentaAtencion: id})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "ADMISSION_RECORD_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, item)
}

// ListMedicosPorEspecialidad maneja GET /api/v1/triaje/medicos/:idEspecialidad.
//
// @Summary Lista los médicos de una especialidad
// @Description Devuelve los médicos de la especialidad indicada (SP usp_go_MedicosFiltrarPorIdEspecialidad).
// @Tags Triaje
// @Produce json
// @Param idEspecialidad path int true "Id de la especialidad"
// @Success 200 {object} apiResponse{data=[]domain.MedicoFila}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar los médicos"
// @Router /triaje/medicos/{idEspecialidad} [get]
func (h *TriageHandler) ListMedicosPorEspecialidad(c *gin.Context) {
	raw := c.Param("idEspecialidad")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEspecialidad es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id < 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idEspecialidad debe ser un entero positivo o 0")
		return
	}

	items, err := h.service.ListarMedicosPorEspecialidad(c.Request.Context(), int(id))
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_MEDICOS_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// ListTriajeConsulta maneja GET /api/v1/triaje/consulta.
//
// @Summary Lista la bandeja de triaje de consulta externa
// @Description Devuelve las atenciones de consulta externa con indicador de si tienen triaje (SP AtencionesTriajeFiltro).
// @Tags Triaje
// @Produce json
// @Param fini query string true "Fecha inicial (YYYY-MM-DD)"
// @Param ffin query string true "Fecha final (YYYY-MM-DD)"
// @Param filtro query string false "Texto de búsqueda (documento o nombre)"
// @Param idServicio query int false "Id del consultorio/servicio (0 = todos)"
// @Success 200 {object} apiResponse{data=[]object}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 500 {object} apiResponse{error=apiError} "Error al listar la bandeja"
// @Router /triaje/consulta [get]
func (h *TriageHandler) ListTriajeConsulta(c *gin.Context) {
	params := shared.TriajeConsultaParams{
		FechaInicio: c.Query("fini"),
		FechaFin:    c.Query("ffin"),
		Filtro:      c.Query("filtro"),
	}

	if params.FechaInicio == "" || params.FechaFin == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fini y ffin son obligatorios (YYYY-MM-DD)")
		return
	}

	// Validar el formato de las fechas antes de armar el fragmento WHERE.
	if _, err := time.Parse("2006-01-02", params.FechaInicio); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "fini debe ser YYYY-MM-DD")
		return
	}
	if _, err := time.Parse("2006-01-02", params.FechaFin); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "ffin debe ser YYYY-MM-DD")
		return
	}

	if raw := c.Query("idServicio"); raw != "" {
		v, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || v < 0 {
			respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idServicio debe ser un entero no negativo")
			return
		}
		params.IdServicio = int(v)
	}

	items, err := h.service.ListTriajeConsulta(c.Request.Context(), params)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_LIST_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, items)
}

// CreateTriajeConsulta maneja POST /api/v1/triaje/consulta.
//
// @Summary Registra o actualiza el triaje de consulta externa
// @Description Persiste el triaje de la atención invocando el SP AtencionesTriajeAgregar (inserta si no existe o actualiza el vigente). Devuelve el @Resultado informado por el procedimiento.
// @Tags Triaje
// @Accept json
// @Produce json
// @Param request body createTriajeConsultaRequest true "Datos del triaje de consulta externa"
// @Success 200 {object} apiResponse{data=object} "Resultado informado por el SP"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al registrar el triaje"
// @Router /triaje/consulta [post]
func (h *TriageHandler) CreateTriajeConsulta(c *gin.Context) {
	var req createTriajeConsultaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}

	domainObj := req.toDomain()
	if idEmpleado := c.GetInt("idEmpleado"); idEmpleado != 0 {
		domainObj.IdEmpleado = int64(idEmpleado)
	}

	result, err := h.service.CreateTriajeConsulta(c.Request.Context(), domainObj)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_REGISTER_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]string{"resultado": result})
}

// GetTriajeConsultaPorAtencion maneja GET /api/v1/triaje/consulta/atencion/:idAtencion.
//
// @Summary Obtiene el triaje de consulta externa de una atención
// @Description Devuelve el triaje vigente (UltimoTriaje = 1) de la atención indicada, o null si no existe.
// @Tags Triaje
// @Produce json
// @Param idAtencion path int true "Id de la atención"
// @Success 200 {object} apiResponse{data=object}
// @Failure 400 {object} apiResponse{error=apiError} "Parámetros inválidos"
// @Failure 500 {object} apiResponse{error=apiError} "Error al obtener el triaje"
// @Router /triaje/consulta/atencion/{idAtencion} [get]
func (h *TriageHandler) GetTriajeConsultaPorAtencion(c *gin.Context) {
	raw := c.Param("idAtencion")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "idAtencion debe ser un entero positivo")
		return
	}

	item, err := h.service.GetTriajeConsultaPorAtencion(c.Request.Context(), id)
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_GET_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, item)
}

// UpdateEstadoTriajeConsulta maneja PUT /api/v1/triaje/consulta/:id/estado.
//
// @Summary Actualiza el estado de un triaje de consulta externa
// @Description Actualiza el estado del triaje invocando el SP AtencionesTriajeEstado.
// @Tags Triaje
// @Accept json
// @Produce json
// @Param id path int true "Id del triaje"
// @Param request body triajeConsultaEstadoRequest true "Nuevo estado"
// @Success 200 {object} apiResponse{data=object} "Operación exitosa"
// @Failure 400 {object} apiResponse{error=apiError} "Cuerpo inválido"
// @Failure 500 {object} apiResponse{error=apiError} "Error al actualizar el estado"
// @Router /triaje/consulta/{id}/estado [put]
func (h *TriageHandler) UpdateEstadoTriajeConsulta(c *gin.Context) {
	raw := c.Param("id")
	if raw == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id es obligatorio")
		return
	}
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "id debe ser un entero positivo")
		return
	}

	var req triajeConsultaEstadoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
		return
	}
	if req.Estado == "" {
		respondError(c, http.StatusBadRequest, "VALIDATION_ERROR", "estado es obligatorio")
		return
	}

	err = h.service.UpdateEstadoTriajeConsulta(c.Request.Context(), shared.TriajeConsultaEstadoParams{
		IdTriaje: id,
		Estado:   req.Estado,
	})
	if err != nil {
		respondError(c, http.StatusInternalServerError, "TRIAGE_CONSULTA_STATE_FAILED", err.Error())
		return
	}

	respondSuccess(c, http.StatusOK, map[string]bool{"ok": true})
}
