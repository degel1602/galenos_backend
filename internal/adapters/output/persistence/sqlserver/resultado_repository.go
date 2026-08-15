package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type ResultadoRepository struct {
	db *sql.DB
}

func NewResultadoRepository(db *sql.DB) *ResultadoRepository {
	return &ResultadoRepository{db: db}
}

func (r *ResultadoRepository) ListarLaboratorioPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	// SP real: webHistorialExamenLaboratorio @IdPaciente
	// Columnas devueltas:
	//   IdProducto, IdPuntoCarga, IdMovimiento, Codigo, Nombre, cantidad,
	//   IdOrden, IdLabEstado, FechaSolicitud, FechaResultado, Resultado
	query := "EXEC webHistorialExamenLaboratorio @IdPaciente = @p1"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []domain.Resultado
	for rows.Next() {
		var res domain.Resultado
		res.TipoResultado = "Laboratorio"
		res.IdPaciente = idPaciente

		var idProducto, idPuntoCarga, idMovimiento, cantidad, idLabEstado, idOrden sql.NullInt64
		var codigo, nombre, fechaSolicitud, fechaResultado, resultado sql.NullString

		if err := rows.Scan(
			&idProducto, &idPuntoCarga, &idMovimiento, &codigo, &nombre, &cantidad,
			&idOrden, &idLabEstado, &fechaSolicitud, &fechaResultado, &resultado,
		); err != nil {
			return nil, fmt.Errorf("error escaneando resultado laboratorio: %w", err)
		}

		if idMovimiento.Valid {
			res.IdResultado = int(idMovimiento.Int64)
		}
		if nombre.Valid {
			res.NombreExamen = nombre.String
		}
		if fechaResultado.Valid {
			res.FechaExamen = fechaResultado.String
		}
		if codigo.Valid {
			res.Detalle = codigo.String
		}
		if resultado.Valid && resultado.String != "" {
			res.Estado = resultado.String
		}
		resultados = append(resultados, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultados, nil
}

func (r *ResultadoRepository) ListarImagenesPorPaciente(ctx context.Context, idPaciente int) ([]domain.Resultado, error) {
	// SP real: webHistorialExamenImageneologia @IdPaciente
	// Columnas devueltas:
	//   IdProducto, Codigo, FechaResultado, FechaRegistro, IdPuntoCarga, IdMovimiento,
	//   Nombre, Cantidad, IdImagEstado, Resultado, IdOrden, DiaAtencion,
	//   MesAtencion, NombreDiaAtencion, AnioAtencion
	query := "EXEC webHistorialExamenImageneologia @IdPaciente = @p1"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultados []domain.Resultado
	for rows.Next() {
		var res domain.Resultado
		res.TipoResultado = "Imagen"
		res.IdPaciente = idPaciente

		var idProducto, idPuntoCarga, idMovimiento, cantidad, idImagEstado, idOrden, dia, anio sql.NullInt64
		var codigo, fechaResultado, nombre, resultado, mes, nombreDia sql.NullString
		var fechaRegistro sql.NullTime

		if err := rows.Scan(
			&idProducto, &codigo, &fechaResultado, &fechaRegistro, &idPuntoCarga, &idMovimiento,
			&nombre, &cantidad, &idImagEstado, &resultado, &idOrden, &dia,
			&mes, &nombreDia, &anio,
		); err != nil {
			return nil, fmt.Errorf("error escaneando resultado imagen: %w", err)
		}

		if idMovimiento.Valid {
			res.IdResultado = int(idMovimiento.Int64)
		}
		if nombre.Valid {
			res.NombreExamen = nombre.String
		}
		if fechaResultado.Valid {
			res.FechaExamen = fechaResultado.String
		}
		if codigo.Valid {
			res.Detalle = codigo.String
		}
		if resultado.Valid && resultado.String != "" {
			res.Estado = resultado.String
		}
		resultados = append(resultados, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return resultados, nil
}
