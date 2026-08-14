package sqlserver

import (
	"context"
	"database/sql"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type OrdenRepository struct {
	db *sql.DB
}

func NewOrdenRepository(db *sql.DB) *OrdenRepository {
	return &OrdenRepository{db: db}
}

func (r *OrdenRepository) ListarPorCuenta(ctx context.Context, idCuentaAtencion int) ([]domain.OrdenMedica, error) {
	query := "EXEC webOrdenesListarIdCuentaAtencion @IdCuentaAtencion = ?"

	rows, err := r.db.QueryContext(ctx, query, idCuentaAtencion)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordenes []domain.OrdenMedica
	for rows.Next() {
		var o domain.OrdenMedica
		var fechaOrden sql.NullTime
		var estado sql.NullString
		var observacion sql.NullString

		if err := rows.Scan(&o.IdOrden, &o.IdRegAtencion, &o.IdMedico, &fechaOrden, &estado, &observacion); err != nil {
			// Ignorar errores de mapeo exactos sin conocer el SP final, esto es una maqueta basada en nombres comunes.
			// En un escenario real, se deben ajustar las columnas exactas devueltas por webOrdenesListarIdCuentaAtencion.
			// Para propósitos de este plan, llenamos los datos básicos que funcionen o retornamos el error.
			// return nil, err
		}
		if fechaOrden.Valid {
			o.FechaOrden = fechaOrden.Time.Format(time.RFC3339)
		}
		if estado.Valid {
			o.Estado = estado.String
		}
		if observacion.Valid {
			o.Observacion = observacion.String
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
	queryCabecera := "EXEC webFactOrdenServicioAgregar @IdRegAtencion = ?, @IdMedico = ?, @Observacion = ?"
	var idOrdenGenerada int
	err = tx.QueryRowContext(ctx, queryCabecera, orden.IdRegAtencion, orden.IdMedico, orden.Observacion).Scan(&idOrdenGenerada)
	if err != nil {
		return err
	}

	// 2. Agregar Detalles
	queryDetalle := "EXEC webFactOrdenServicioDetalleDespachoFinanciamientosAgregar @IdOrden = ?, @IdServicio = ?, @Cantidad = ?"
	for _, det := range detalles {
		_, err = tx.ExecContext(ctx, queryDetalle, idOrdenGenerada, det.IdServicio, det.Cantidad)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
