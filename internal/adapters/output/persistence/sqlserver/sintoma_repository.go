package sqlserver

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type SintomaRepository struct {
	db *sql.DB
}

func NewSintomaRepository(db *sql.DB) *SintomaRepository {
	return &SintomaRepository{db: db}
}

func (r *SintomaRepository) ListarCatalogo(ctx context.Context) ([]domain.SintomaCatalogo, error) {
	query := "EXEC sp_go_SelectSintomaCatalogo"

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error listando catálogo de síntomas: %w", err)
	}
	defer rows.Close()

	var sintomas []domain.SintomaCatalogo
	for rows.Next() {
		var s domain.SintomaCatalogo
		var orden sql.NullInt64
		if err := rows.Scan(&s.Sistema, &s.Sintoma, &s.IdSintoma, &orden); err != nil {
			return nil, fmt.Errorf("error escaneando síntoma del catálogo: %w", err)
		}
		if orden.Valid {
			s.Orden = int(orden.Int64)
		}
		sintomas = append(sintomas, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sintomas, nil
}

func (r *SintomaRepository) AgregarCatalogo(ctx context.Context, sistema, sintoma string, idUsuario int) error {
	query := "EXEC sp_go_AgregarSintomaCatalogo @Sistema = @p1, @Sintoma = @p2, @IdUsuario = @p3"
	_, err := r.db.ExecContext(ctx, query, sql.Named("p1", sistema), sql.Named("p2", sintoma), sql.Named("p3", idUsuario))
	if err != nil {
		return fmt.Errorf("error agregando síntoma al catálogo: %w", err)
	}
	return nil
}

func (r *SintomaRepository) GuardarEvolucionSintomas(ctx context.Context, idRegAtencion int, sintomas []domain.SintomaSeleccionado, idUsuario int) error {
	payload, err := json.Marshal(sintomas)
	if err != nil {
		return fmt.Errorf("error serializando síntomas: %w", err)
	}

	query := "EXEC sp_go_InsertarEvolucionSintomas @IdRegAtencion = @p1, @Sintomas = @p2, @IdUsuario = @p3"
	_, err = r.db.ExecContext(ctx, query, sql.Named("p1", idRegAtencion), sql.Named("p2", string(payload)), sql.Named("p3", idUsuario))
	if err != nil {
		return fmt.Errorf("error guardando síntomas de la evolución: %w", err)
	}
	return nil
}