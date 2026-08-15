package sqlserver

import (
	"context"
	"database/sql"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type MotivoRepository struct {
	db *sql.DB
}

func NewMotivoRepository(db *sql.DB) *MotivoRepository {
	return &MotivoRepository{db: db}
}

func (r *MotivoRepository) ListarPorAtencion(ctx context.Context, idRegAtencion int) ([]domain.MotivoAtencion, error) {
	query := "EXEC webEvolucionMotivoListar @IdRegAtencion = ?"

	rows, err := r.db.QueryContext(ctx, query, idRegAtencion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var motivos []domain.MotivoAtencion
	for rows.Next() {
		var m domain.MotivoAtencion
		var fechaRegistro sql.NullTime
		if err := rows.Scan(&m.IdMotivo, &m.IdRegAtencion, &m.Motivo, &fechaRegistro); err != nil {
			return nil, err
		}
		if fechaRegistro.Valid {
			m.FechaRegistro = fechaRegistro.Time.Format(time.RFC3339)
		}
		motivos = append(motivos, m)
	}
	return motivos, rows.Err()
}

func (r *MotivoRepository) Guardar(ctx context.Context, idRegAtencion int, motivo string) error {
	// SP real: usp_Tab_Atencion_Registro_Guardar (guardado único del registro de atención).
	// Actualiza solo los campos enviados (COALESCE); los demás quedan intactos.
	// Devuelve una fila con Mensaje = "<IdResultado>;<texto>".
	query := `
		EXEC usp_Tab_Atencion_Registro_Guardar
			@IdRegAtencion = ?,
			@Motivo = ?,
			@IdUsuario = NULL,
			@Mensaje = ? OUTPUT
	`
	var mensaje string
	err := r.db.QueryRowContext(ctx, query,
		idRegAtencion,
		motivo,
		sql.Out{Dest: &mensaje},
	).Scan(&mensaje)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}
