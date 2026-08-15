package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sqlServerEvolucionRepository struct {
	db *sql.DB
}

func NewSqlServerEvolucionRepository(db *sql.DB) output.EvolucionRepository {
	return &sqlServerEvolucionRepository{db: db}
}

func (r *sqlServerEvolucionRepository) ListPatients(ctx context.Context, fini, ffin string, idUsuario int) ([]domain.PatientListItem, error) {
	// Se usa el SP usp_go_ListarPacientesSegunTipoServicio, que requiere
	// @IdTipoServicio (2 = Emergencia), @Fecha (datetime), @Filtro y @IdUsuario.
	// El handler ya validó el formato de fini antes de llegar aquí.
	fecha, err := time.Parse("2006-01-02", fini)
	if err != nil {
		return nil, fmt.Errorf("fini inválido: %w", err)
	}
	query := `EXEC [dbo].[usp_go_ListarPacientesSegunTipoServicio] @IdTipoServicio = @p1, @Fecha = @p2, @Filtro = @p3, @IdUsuario = @p4`
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", 2),
		sql.Named("p2", fecha),
		sql.Named("p3", ""),
		sql.Named("p4", idUsuario),
	)
	if err != nil {
		return nil, fmt.Errorf("error querying patients for tray: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("error reading patient tray maps: %w", err)
	}

	var patients []domain.PatientListItem
	for _, m := range maps {
		patients = append(patients, mapToPatientListItem(m))
	}
	return patients, nil
}

func mapToPatientListItem(m map[string]any) domain.PatientListItem {
	var p domain.PatientListItem
	
	// El SP usp_go_ListarPacientesSegunTipoServicio retorna: IdEpisodio,
	// IdAtencion, IdCuentaAtencion, Paciente, NroHistoriaClinica, Servicio,
	// Sexo y cama.
	p.IdRegAtencion = getIntFallback(m, "idRegAtencion", "IdAtencion")
	p.IdPaciente = getIntFallback(m, "IdPaciente", "IdEpisodio")
	p.Historia = getStringFallback(m, "historia", "NroHistoriaClinica", "N/A")
	p.Nombre = getNombrePaciente(m)
	p.Edad = getStringFallback(m, "edad", "", "N/A")
	p.Sexo = getSexoFallback(m)
	p.Ubicacion = getStringFallback(m, "ubicacion", "Servicio", "Emergencia")
	p.Cama = getStringFallback(m, "cama", "", "NS")
	p.Estado = getStringFallback(m, "estado", "", "Pendiente")
	
	return p
}

func getIntFallback(m map[string]any, key1, key2 string) int {
	if val, ok := m[key1]; ok && val != nil {
		return int(val.(int64))
	}
	if key2 != "" {
		if val, ok := m[key2]; ok && val != nil {
			return int(val.(int64))
		}
	}
	return 0
}

func getStringFallback(m map[string]any, key1, key2, fallback string) string {
	if ptr := rowString(m, key1); ptr != nil && *ptr != "" {
		return *ptr
	}
	if key2 != "" {
		if ptr := rowString(m, key2); ptr != nil && *ptr != "" {
			return *ptr
		}
	}
	return fallback
}

func getNombrePaciente(m map[string]any) string {
	if ptr := rowString(m, "nombre", "Paciente"); ptr != nil && *ptr != "" {
		return *ptr
	}
	pat := getStringFallback(m, "ApellidoPaterno", "", "")
	mat := getStringFallback(m, "ApellidoMaterno", "", "")
	nom := getStringFallback(m, "PrimerNombre", "", "")
	return fmt.Sprintf("%s %s, %s", pat, mat, nom)
}

func getSexoFallback(m map[string]any) string {
	if ptr := rowString(m, "sexo"); ptr != nil && *ptr != "" {
		return *ptr
	}
	if ptr := rowInt64(m, "IdTipoSexo"); ptr != nil {
		return fmt.Sprintf("%d", *ptr)
	}
	return "0"
}


func (r *sqlServerEvolucionRepository) ListEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error) {
	query := `EXEC [dbo].[webEvolucionesFirmaListarIdRegAtencion] @IdRegAtencion = ?, @NombreDocumento = 'EvolucionMedica'`
	rows, err := r.db.QueryContext(ctx, query, idRegAtencion)
	if err != nil {
		return nil, fmt.Errorf("error querying evolutions: %w", err)
	}
	defer rows.Close()

	var evolutions []domain.EvolucionFirma
	for rows.Next() {
		var e domain.EvolucionFirma
		var rBase, nArchivo, doc, data, fRegistro sql.NullString
		var idEmpReg, idEmpMod, idEmpAnula, estado sql.NullInt64
		var fMod, fAnul sql.NullTime

		if err := rows.Scan(
			&e.IdRegAtencion,
			&e.IdFirma,
			&rBase,
			&nArchivo,
			&doc,
			&data,
			&idEmpReg,
			&idEmpMod,
			&idEmpAnula,
			&fRegistro,
			&fMod,
			&fAnul,
			&estado,
		); err != nil {
			return nil, fmt.Errorf("error scanning evolution: %w", err)
		}

		if doc.Valid {
			e.NombreDocumento = doc.String
		}
		if data.Valid {
			e.DataB64 = data.String
		}
		if idEmpReg.Valid {
			e.IdEmpleadoRegistra = int(idEmpReg.Int64)
		}
		if fRegistro.Valid {
			e.FechaRegistro = fRegistro.String
		}

		evolutions = append(evolutions, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando evoluciones: %w", err)
	}

	return evolutions, nil
}

func (r *sqlServerEvolucionRepository) SaveEvolucion(ctx context.Context, evolution domain.EvolucionFirma) error {
	query := `
		EXEC [dbo].[Web_sp_GuardarEvolucionFirma] 
			@IdRegAtencion = ?, 
			@RutaBase = ?, 
			@NombreArchivo = ?, 
			@NombreDocumento = ?, 
			@DataB64 = ?, 
			@IdEmpleadoRegistra = ?, 
			@IdEmpleadoModifica = ?, 
			@Estado = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		evolution.IdRegAtencion,
		"evols/",
		fmt.Sprintf("evol_%d_%d.json", evolution.IdRegAtencion, time.Now().Unix()),
		evolution.NombreDocumento,
		evolution.DataB64,
		evolution.IdEmpleadoRegistra,
		evolution.IdEmpleadoRegistra,
		1,
	)
	if err != nil {
		return fmt.Errorf("error saving evolution: %w", err)
	}
	return nil
}
