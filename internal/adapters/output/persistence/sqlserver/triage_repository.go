package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type triageRepository struct {
	db *sql.DB
}

// NewTriageRepository construye el adaptador que implementa el puerto de
// salida output.TriageRepository contra el SP webTab_PacienteTriajeAgregar.
func NewTriageRepository(db *sql.DB) output.TriageRepository {
	return &triageRepository{db: db}
}

// Create invoca el procedimiento almacenado webTab_PacienteTriajeAgregar
// con todos los parámetros nombrados (Sp_...), enviados como llamada RPC
// por el driver mssql, por lo que no se concatena SQL. Los parámetros sin
// valor se envían como NULL. El parámetro @Resultado se declara con
// sql.Out para capturar el valor de salida que el SP devuelve.
//
// Los nombres de parámetro conservan exactamente los declarados en el SP,
// incluyendo las variantes con mayúsculas/tipos que la base refleja
// (PrimerNmobre, idEstadollego, telefono, IdfuenteFinanciamiento, etc.),
// ya que el driver los resolve contra el nombre formal del parámetro.
func (r *triageRepository) Create(ctx context.Context, triage *domain.Triage) (string, error) {
	const procedure = `webTab_PacienteTriajeAgregar`

	var resultado string

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdTriaje", triage.IDTriaje),
		sql.Named("IdDocIdentidad", triage.DocIdentityID),
		sql.Named("NroDocumento", triage.DocumentNumber),
		sql.Named("ApellidoPaterno", triage.PaternalSurname),
		sql.Named("ApellidoMaterno", triage.MaternalSurname),
		sql.Named("PrimerNmobre", triage.FirstName),
		sql.Named("SegundoNombre", triage.SecondName),
		sql.Named("TercerNombre", triage.ThirdName),
		sql.Named("IdSexo", triage.SexTypeID),
		sql.Named("FechaNacimiento", triage.DateOfBirth),
		sql.Named("telefono", triage.Phone),
		sql.Named("IdDepartamentoDomicilio", triage.HomeDepartmentID),
		sql.Named("IdProvinciaDomicilio", triage.HomeProvinceID),
		sql.Named("idDistritoDomicilio", triage.HomeDistrictID),
		sql.Named("idComunidadDomicilio", triage.HomeCommunityID),
		sql.Named("Direccion", triage.HomeAddress),
		sql.Named("EsAccidenteTransito", triage.IsTrafficAccident),
		sql.Named("IdfuenteFinanciamiento", triage.FundingSourceID),
		sql.Named("Email", triage.Email),
		sql.Named("IdEstadoCivil", triage.MaritalStatusID),
		sql.Named("FrecCardiaca", triage.HeartRate),
		sql.Named("Temperatura", triage.Temperature),
		sql.Named("PresionArterial", triage.BloodPressure),
		sql.Named("Saturacion", triage.OxygenSaturation),
		sql.Named("FrecRespiratoria", triage.RespiratoryRate),
		sql.Named("FI02", triage.FiO2),
		sql.Named("Peso", triage.Weight),
		sql.Named("Talla", triage.Height),
		sql.Named("IMC", triage.BMI),
		sql.Named("TiempoEvolucionCantidad", triage.EvolutionTimeQuantity),
		sql.Named("TiempoEvolucionCantidadUnidad", triage.EvolutionTimeQuantityUnit),
		sql.Named("EscalaDolor", triage.PainScale),
		sql.Named("EscalaGlasgow", triage.GlasgowScale),
		sql.Named("IdTipoPrioridad", triage.PriorityTypeID),
		sql.Named("IdServicio", triage.ServiceID),
		sql.Named("motivo", triage.Motivo),
		sql.Named("Gestante", triage.IsPregnant),
		sql.Named("IdEstadollego", triage.ArrivalStateID),
		sql.Named("Foto", triage.Photo),
		sql.Named("Idempleado", triage.EmployeeID),
		sql.Named("Resultado", sql.Out{Dest: &resultado}),
	)
	if err != nil {
		return "", fmt.Errorf("calling webTab_PacienteTriajeAgregar: %w", err)
	}

	return resultado, nil
}

// List invoca el procedimiento almacenado ListarTriaje_Emergencia con los
// filtros de rango de fechas, texto de búsqueda, servicio de derivación y
// estado. Los nombres de parámetro coinciden exactamente con los del SP.
// Las columnas que devuelve el SP son desconocidas en build-time, por eso
// se mapean a map[string]any con rowsToMaps (como el proyecto FastAPI).
func (r *triageRepository) List(ctx context.Context, params shared.TriageListParams) ([]map[string]any, error) {
	const procedure = `ListarTriaje_Emergencia`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("fini", params.FechaInicio),
		sql.Named("ffin", params.FechaFin),
		sql.Named("filtro", params.Filtro),
		sql.Named("derivado_a_servicio", params.DerivadoAServicio),
		sql.Named("IdEstado", params.IdEstado),
	)
	if err != nil {
		return nil, fmt.Errorf("calling ListarTriaje_Emergencia: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading triaje list: %w", err)
	}

	return maps, nil
}

// ListPendingAdmission invoca el procedimiento almacenado
// webGestionAtencion_E_H_BusquedaFiltrar con la fecha y los filtros
// opcionales (cuenta de atención, departamento, especialidad, servicio y
// tipo de servicio). Devuelve los pacientes con triaje que aún no han
// sido admisionados. Los parámetros de filtro se envían como 0 cuando no
// se desean aplicar. Las columnas se mapean a map[string]any porque el SP
// las resuelve en runtime.
func (r *triageRepository) ListPendingAdmission(ctx context.Context, params shared.TriageAdmisionParams) ([]map[string]any, error) {
	const procedure = `webGestionAtencion_E_H_BusquedaFiltrar`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("Fecha", params.Fecha),
		sql.Named("filtro", params.Filtro),
		sql.Named("NroCta", params.NroCta),
		sql.Named("IdDepartamento", params.IdDepartamento),
		sql.Named("IdEspecialidad", params.IdEspecialidad),
		sql.Named("IdServicio", params.IdServicio),
		sql.Named("IdtipoServicio", params.IdTipoServicio),
	)
	if err != nil {
		return nil, fmt.Errorf("calling webGestionAtencion_E_H_BusquedaFiltrar: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading triaje pending admission list: %w", err)
	}

	return maps, nil
}

// CreateAdmission invoca el procedimiento almacenado
// WebCrearAtencionDesdeTriaje con los datos del paciente que se admisiona
// desde su triaje. El parámetro @Resultado se declara con sql.Out para
// capturar el valor de salida que el SP devuelve.
func (r *triageRepository) CreateAdmission(ctx context.Context, admision *domain.AdmisionDesdeTriaje) (string, error) {
	const procedure = `WebCrearAtencionDesdeTriaje`

	var resultado string

	_, err := r.db.ExecContext(ctx, procedure,
		sql.Named("IdTriaje", admision.IDTriaje),
		sql.Named("IdPacienteTriaje", admision.IDPacienteTriaje),
		sql.Named("IdEmpleado", admision.IDEmpleado),
		sql.Named("NombreAcompañante", admision.NombreAcompanante),
		sql.Named("TelefonoAcompañante", admision.TelefonoAcompanante),
		sql.Named("DireccionPaciente", admision.DireccionPaciente),
		sql.Named("Observacion", admision.Observacion),
		sql.Named("Resultado", sql.Out{Dest: &resultado}),
	)
	if err != nil {
		return "", fmt.Errorf("calling WebCrearAtencionDesdeTriaje: %w", err)
	}

	return resultado, nil
}

// GetReporte invoca el procedimiento almacenado WebSelectReporteTriaje con
// los filtros por id de triaje e id de paciente. Se envía -100 en ambos
// filtros cuando no se quieren aplicar. Las columnas se mapean a
// map[string]any porque el SP las resuelve en runtime y el frontend puede
// consumir el reporte tal cual.
func (r *triageRepository) GetReporte(ctx context.Context, params shared.TriageReporteParams) ([]map[string]any, error) {
	const procedure = `WebSelectReporteTriaje`

	rows, err := r.db.QueryContext(ctx, procedure,
		sql.Named("Id", params.IDTriaje),
		sql.Named("IdPaciente", params.IDPaciente),
	)
	if err != nil {
		return nil, fmt.Errorf("calling WebSelectReporteTriaje: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading triage report: %w", err)
	}

	return maps, nil
}

// GetFichaAdmision invoca el procedimiento almacenado webFichaEmergencia
// con el id de la cuenta de atención. El SP retorna una única fila (TOP 1)
// con los datos del paciente y los adicionales para la ficha de admisión.
// Devuelve nil si no hay registros.
func (r *triageRepository) GetFichaAdmision(ctx context.Context, params shared.FichaAdmisionParams) (*map[string]any, error) {
	const procedure = `webFichaEmergencia`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("idcuentaatencion", params.IdCuentaAtencion))
	if err != nil {
		return nil, fmt.Errorf("calling webFichaEmergencia: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("reading admission record: %w", err)
	}
	if len(maps) == 0 {
		return nil, nil
	}

	m := maps[0]
	return &m, nil
}
