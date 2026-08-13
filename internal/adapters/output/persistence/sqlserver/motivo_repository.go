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
	query := "EXEC webEvolucionMotivoGuardar @IdRegAtencion = ?, @Motivo = ?"
	_, err := r.db.ExecContext(ctx, query, idRegAtencion, motivo)
	return err
}
