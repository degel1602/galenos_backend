package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sqlServerEvolucionRepository struct {
	db *sql.DB
}

func NewSqlServerEvolucionRepository(db *sql.DB) output.EvolucionRepository {
	return &sqlServerEvolucionRepository{
		db: db,
	}
}

func (r *sqlServerEvolucionRepository) Save(ctx context.Context, evolucion *domain.Evolucion) error {
	// Iniciamos transacción transaccional aislando concurrencia.
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// Asumimos un Stored Procedure principal: webTab_EvolucionAgregar
	// Aquí aplicaríamos WITH (UPDLOCK, HOLDLOCK) si hiciéramos un SELECT previo de validación,
	// pero en este caso el SP inserta directamente.
	queryEvolucion := `
		EXEC webTab_EvolucionAgregar 
			@IDPaciente = @p1, 
			@IDMedico = @p2, 
			@IDCita = @p3, 
			@MotivoAtencion = @p4, 
			@DetalleMotivo = @p5, 
			@SubjetivoDetalle = @p6, 
			@PresionArterial = @p7, 
			@FrecCardiaca = @p8, 
			@FrecRespiratoria = @p9, 
			@Temperatura = @p10, 
			@SaturacionO2 = @p11, 
			@Peso = @p12, 
			@Talla = @p13, 
			@IMC = @p14, 
			@Glucemia = @p15, 
			@EstadoClinico = @p16, 
			@Pronostico = @p17, 
			@PlanDetalle = @p18;
	`

	// El SP debe retornar el ID generado de la evolución
	var newIDEvolucion int64
	err = tx.QueryRowContext(ctx, queryEvolucion,
		evolucion.IDPaciente,
		evolucion.IDMedico,
		evolucion.IDCita,
		evolucion.MotivoAtencion,
		evolucion.DetalleMotivo,
		evolucion.SubjetivoDetalle,
		evolucion.PresionArterial,
		evolucion.FrecCardiaca,
		evolucion.FrecRespiratoria,
		evolucion.Temperatura,
		evolucion.SaturacionO2,
		evolucion.Peso,
		evolucion.Talla,
		evolucion.IMC,
		evolucion.Glucemia,
		evolucion.EstadoClinico,
		evolucion.Pronostico,
		evolucion.PlanDetalle,
	).Scan(&newIDEvolucion)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("error executing webTab_EvolucionAgregar: %w", err)
	}
	evolucion.IDEvolucion = &newIDEvolucion

	// Guardar diagnósticos iterando e invocando a otro SP hipotético: webTab_EvolucionDiagnosticoAgregar
	queryDiagnostico := `
		EXEC webTab_EvolucionDiagnosticoAgregar 
			@IDEvolucion = @p1, 
			@CIE10 = @p2, 
			@Descripcion = @p3, 
			@Tipo = @p4, 
			@Condicion = @p5, 
			@Estado = @p6;
	`

	for _, dx := range evolucion.Diagnosticos {
		_, err = tx.ExecContext(ctx, queryDiagnostico,
			newIDEvolucion,
			dx.CIE10,
			dx.Descripcion,
			dx.Tipo,
			dx.Condicion,
			dx.Estado,
		)
		if err != nil {
			return fmt.Errorf("error saving diagnostico CIE10 %s: %w", dx.CIE10, err)
		}
	}

	return tx.Commit()
}

func (r *sqlServerEvolucionRepository) FindByPacienteID(ctx context.Context, pacienteID int64) ([]domain.Evolucion, error) {
	// Asumimos un Stored Procedure: webTab_EvolucionListarPorPaciente
	query := `EXEC webTab_EvolucionListarPorPaciente @IDPaciente = @p1;`

	rows, err := r.db.QueryContext(ctx, query, pacienteID)
	if err != nil {
		return nil, fmt.Errorf("error querying evoluciones: %w", err)
	}
	defer rows.Close()

	var evoluciones []domain.Evolucion

	for rows.Next() {
		var ev domain.Evolucion
		var fecha time.Time

		err := rows.Scan(
			&ev.IDEvolucion,
			&ev.IDPaciente,
			&ev.IDMedico,
			&ev.IDCita,
			&fecha,
			&ev.MotivoAtencion,
			&ev.DetalleMotivo,
			&ev.SubjetivoDetalle,
			&ev.PresionArterial,
			&ev.FrecCardiaca,
			&ev.FrecRespiratoria,
			&ev.Temperatura,
			&ev.SaturacionO2,
			&ev.Peso,
			&ev.Talla,
			&ev.IMC,
			&ev.Glucemia,
			&ev.EstadoClinico,
			&ev.Pronostico,
			&ev.PlanDetalle,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		ev.Fecha = &fecha
		evoluciones = append(evoluciones, ev)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return evoluciones, nil
}
