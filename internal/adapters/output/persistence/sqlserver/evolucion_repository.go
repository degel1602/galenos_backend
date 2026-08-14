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

func (r *sqlServerEvolucionRepository) ListPatients(ctx context.Context, fini, ffin string) ([]domain.PatientListItem, error) {
	query := `
		SELECT TOP 50
			a.IdPaciente AS idRegAtencion,
			p.IdPaciente,
			ISNULL(CAST(p.NroHistoriaClinica AS VARCHAR(20)), 'N/A') AS historia,
			ISNULL(p.ApellidoPaterno, '') + ' ' + ISNULL(p.ApellidoMaterno, '') + ', ' + ISNULL(p.PrimerNombre, '') AS nombre,
			'N/A' AS edad,
			ISNULL(CAST(p.IdTipoSexo AS VARCHAR(10)), '') AS sexo,
			'Consultorio' AS ubicacion,
			'Atendido' AS estado
		FROM Atenciones a
		INNER JOIN Pacientes p ON a.IdPaciente = p.IdPaciente
		WHERE p.IdEstado = 1
	`
	args := []any{}
	if fini != "" && ffin != "" {
		query += ` AND CONVERT(date, a.FechaIngreso) BETWEEN ? AND ?`
		args = append(args, fini, ffin)
	}
	query += ` ORDER BY a.FechaIngreso DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error querying patients for tray: %w", err)
	}
	defer rows.Close()

	var patients []domain.PatientListItem
	for rows.Next() {
		var p domain.PatientListItem
		if err := rows.Scan(
			&p.IdRegAtencion,
			&p.IdPaciente,
			&p.Historia,
			&p.Nombre,
			&p.Edad,
			&p.Sexo,
			&p.Ubicacion,
			&p.Estado,
		); err != nil {
			return nil, fmt.Errorf("error scanning patient: %w", err)
		}
		patients = append(patients, p)
	}
	return patients, nil
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
