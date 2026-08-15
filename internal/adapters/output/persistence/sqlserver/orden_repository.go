package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type OrdenRepository struct {
	db *sql.DB
}

func NewOrdenRepository(db *sql.DB) *OrdenRepository {
	return &OrdenRepository{db: db}
}

func (r *OrdenRepository) ListarPorCuenta(ctx context.Context, idCuentaAtencion int) ([]domain.OrdenMedica, error) {
	// SP real: webOrdenesListarIdCuentaAtencion exige 2 parámetros (@idCuentaAtencion, @RecetaAdicional).
	// @RecetaAdicional = -100 devuelve todas (con y sin receta adicional).
	query := "EXEC webOrdenesListarIdCuentaAtencion @IdCuentaAtencion = @p1, @RecetaAdicional = @p2"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idCuentaAtencion), sql.Named("p2", -100))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordenes []domain.OrdenMedica
	for rows.Next() {
		var o domain.OrdenMedica

		// Columnas reales de webOrdenesListarIdCuentaAtencion:
		//   IdProducto, Codigo, Nombre, DescripcionPuntoCarga, CantidadPedida, Precio,
		//   Total, idReceta, idItem, idUNIDDosisReceta, idFrecuencia, Duracion,
		//   TipoProducto, Notificacion
		var idProducto, cantidadPedida, idReceta, idItem, idUnidDosis, idFrecuencia sql.NullInt64
		var codigo, nombre, descripcionPuntoCarga, duracion, tipoProducto, notificacion sql.NullString
		var precio, total sql.NullFloat64

		if err := rows.Scan(
			&idProducto, &codigo, &nombre, &descripcionPuntoCarga, &cantidadPedida, &precio,
			&total, &idReceta, &idItem, &idUnidDosis, &idFrecuencia, &duracion,
			&tipoProducto, &notificacion,
		); err != nil {
			return nil, fmt.Errorf("error escaneando orden: %w", err)
		}

		if idReceta.Valid {
			o.IdOrden = int(idReceta.Int64)
		}
		if tipoProducto.Valid {
			o.Estado = tipoProducto.String
		}
		if notificacion.Valid {
			o.Observacion = notificacion.String
		}
		ordenes = append(ordenes, o)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ordenes, nil
}

func (r *OrdenRepository) CrearOrden(ctx context.Context, orden domain.OrdenMedica, detalles []domain.DetalleOrden) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Crear Cabecera
	queryCabecera := "EXEC webFactOrdenServicioAgregar @IdRegAtencion = @p1, @IdMedico = @p2, @Observacion = @p3"
	var idOrdenGenerada int
	err = tx.QueryRowContext(ctx, queryCabecera, sql.Named("p1", orden.IdRegAtencion), sql.Named("p2", orden.IdMedico), sql.Named("p3", orden.Observacion)).Scan(&idOrdenGenerada)
	if err != nil {
		return err
	}

	// 2. Agregar Detalles
	queryDetalle := "EXEC webFactOrdenServicioDetalleDespachoFinanciamientosAgregar @IdOrden = @p1, @IdServicio = @p2, @Cantidad = @p3"
	for _, det := range detalles {
		_, err = tx.ExecContext(ctx, queryDetalle, sql.Named("p1", idOrdenGenerada), sql.Named("p2", det.IdServicio), sql.Named("p3", det.Cantidad))
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
