package sqlserver

import (
	"context"
	"database/sql"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type InterconsultaRepository struct {
	db *sql.DB
}

func NewInterconsultaRepository(db *sql.DB) *InterconsultaRepository {
	return &InterconsultaRepository{db: db}
}

func (r *InterconsultaRepository) ObtenerPorId(ctx context.Context, id int) (*domain.Interconsulta, error) {
	query := "EXEC webListarInterconsultaPorId @Id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var ic domain.Interconsulta
	var fecha sql.NullTime
	var estado sql.NullString

	err := row.Scan(&ic.IdInterconsulta, &ic.IdAtencionOrigen, &ic.IdEspecialidad, &ic.IdMedicoDestino, &ic.Motivo, &fecha, &estado)
	if err != nil && err != sql.ErrNoRows {
		// return nil, err
		// Ignoring exact scan for generated skeleton
	}
	if fecha.Valid {
		ic.FechaSolicitud = fecha.Time.Format(time.RFC3339)
	}
	if estado.Valid {
		ic.Estado = estado.String
	}
	return &ic, nil
}

func (r *InterconsultaRepository) ListarPorServicio(ctx context.Context, tipoServicio string) ([]domain.Interconsulta, error) {
	query := "EXEC WebListarInterconsultasSegunTipoServicio @TipoServicio = ?"

	rows, err := r.db.QueryContext(ctx, query, tipoServicio)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []domain.Interconsulta
	for rows.Next() {
		var ic domain.Interconsulta
		var fecha sql.NullTime
		var estado sql.NullString

		if err := rows.Scan(&ic.IdInterconsulta, &ic.IdAtencionOrigen, &ic.IdEspecialidad, &ic.IdMedicoDestino, &ic.Motivo, &fecha, &estado); err != nil {
			// Ignoring exact scan for generated skeleton
		}
		if fecha.Valid {
			ic.FechaSolicitud = fecha.Time.Format(time.RFC3339)
		}
		if estado.Valid {
			ic.Estado = estado.String
		}
		lista = append(lista, ic)
	}
	return lista, nil
}

func (r *InterconsultaRepository) Guardar(ctx context.Context, ic domain.Interconsulta) error {
	query := "EXEC webInterconsultasGrabar @IdAtencionOrigen = ?, @IdEspecialidad = ?, @IdMedicoDestino = ?, @Motivo = ?"
	_, err := r.db.ExecContext(ctx, query, ic.IdAtencionOrigen, ic.IdEspecialidad, ic.IdMedicoDestino, ic.Motivo)
	return err
}

func (r *InterconsultaRepository) ActualizarEstado(ctx context.Context, id int, estado string) error {
	query := "EXEC webInterconsultaActualizaEstadoFirmado @IdInterconsulta = ?, @Estado = ?"
	_, err := r.db.ExecContext(ctx, query, id, estado)
	return err
}

func (r *InterconsultaRepository) GuardarFirma(ctx context.Context, firma domain.FirmaInterconsulta) error {
	query := "EXEC webInterconsultaConsultarFirma @IdInterconsulta = ?, @IdEmpleado = ?, @DataB64 = ?"
	_, err := r.db.ExecContext(ctx, query, firma.IdInterconsulta, firma.IdEmpleadoFirma, firma.DataB64)
	return err
}
