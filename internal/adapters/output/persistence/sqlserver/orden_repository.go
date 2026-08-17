package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

type OrdenRepository struct {
	db *sql.DB
}

func NewOrdenRepository(db *sql.DB) *OrdenRepository {
	return &OrdenRepository{db: db}
}

func (r *OrdenRepository) ListarPorCuenta(ctx context.Context, idRegAtencion int) ([]domain.OrdenMedica, error) {
	// El frontend envía idRegAtencion = Atenciones.IdAtencion; el SP exige IdCuentaAtencion.
	idCuenta, err := r.resolverCuentaAtencion(ctx, idRegAtencion)
	if err != nil {
		return nil, err
	}

	query := "EXEC webOrdenesListarIdCuentaAtencion @IdCuentaAtencion = @p1, @RecetaAdicional = @p2"

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idCuenta), sql.Named("p2", -100))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordenes []domain.OrdenMedica
	for rows.Next() {
		var o domain.OrdenMedica

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

		det := domain.DetalleOrden{
			IdProducto:     int(idProducto.Int64),
			NombreProducto: nombre.String,
			Codigo:         codigo.String,
			Cantidad:       int(cantidadPedida.Int64),
			Precio:         precio.Float64,
			Total:          total.Float64,
		}

		found := false
		for i := range ordenes {
			if ordenes[i].IdOrden == int(idReceta.Int64) {
				ordenes[i].Detalles = append(ordenes[i].Detalles, det)
				found = true
				break
			}
		}
		if !found {
			o.IdOrden = int(idReceta.Int64)
			if tipoProducto.Valid {
				o.Estado = tipoProducto.String
			}
			if notificacion.Valid {
				o.Observacion = notificacion.String
			}
			o.Detalles = []domain.DetalleOrden{det}
			ordenes = append(ordenes, o)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Completar fecha y médico desde la cabecera real de la receta
	if len(ordenes) > 0 {
		r.completarCabeceras(ctx, idCuenta, ordenes)
	}

	return ordenes, nil
}

// completarCabeceras rellena FechaOrden y Medico de cada receta desde RecetaCabecera.
func (r *OrdenRepository) completarCabeceras(ctx context.Context, idCuentaAtencion int, ordenes []domain.OrdenMedica) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT rc.idReceta, CONVERT(varchar(19), rc.FechaReceta, 120),
		       RTRIM(ISNULL(em.ApellidoPaterno,'')) + ' ' + RTRIM(ISNULL(em.ApellidoMaterno,'')) + ' ' + RTRIM(ISNULL(em.Nombres,''))
		FROM RecetaCabecera rc
		LEFT JOIN Medicos m ON m.IdMedico = rc.idMedicoReceta
		LEFT JOIN Empleados em ON em.IdEmpleado = m.IdEmpleado
		WHERE rc.idCuentaAtencion = @p1`,
		sql.Named("p1", idCuentaAtencion),
	)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var idReceta int
		var fecha, medico string
		if err := rows.Scan(&idReceta, &fecha, &medico); err != nil {
			continue
		}
		for i := range ordenes {
			if ordenes[i].IdOrden == idReceta {
				ordenes[i].FechaOrden = fecha
				ordenes[i].Medico = medico
			}
		}
	}
}

func (r *OrdenRepository) CrearOrden(ctx context.Context, orden domain.OrdenMedica, detalles []domain.DetalleOrden, idEmpleado int) error {
	if len(detalles) == 0 {
		return fmt.Errorf("la orden debe tener al menos un detalle")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Resolver la atención real (el frontend envía IdAtencion como idRegAtencion)
	var idCuentaAtencion, idPaciente, idServicioIngreso, idFormaPago, idFuenteFinanciamiento int
	err = tx.QueryRowContext(ctx, `
		SELECT IdCuentaAtencion, IdPaciente, IdServicioIngreso, IdFormaPago, ISNULL(idFuenteFinanciamiento,0)
		FROM Atenciones WHERE IdAtencion = @p1`,
		sql.Named("p1", orden.IdRegAtencion),
	).Scan(&idCuentaAtencion, &idPaciente, &idServicioIngreso, &idFormaPago, &idFuenteFinanciamiento)
	if err != nil {
		return fmt.Errorf("resolviendo atención %d: %w", orden.IdRegAtencion, err)
	}

	// 2. Resolver el médico real desde el empleado autenticado (JWT)
	var idMedico int
	err = tx.QueryRowContext(ctx, "SELECT TOP 1 IdMedico FROM Medicos WHERE IdEmpleado = @p1",
		sql.Named("p1", idEmpleado),
	).Scan(&idMedico)
	if err != nil {
		return fmt.Errorf("el empleado %d no tiene médico asociado: %w", idEmpleado, err)
	}

	// 3. Crear cabecera de receta (SP real)
	// Params: @Respuesta OUTPUT, @IdPuntoCarga, @idCuentaAtencion, @idServicioReceta,
	//         @idMedicoReceta, @IdProducto, @Idpaciente, @IdUsuarioAuditoria
	// Respuesta: "OK;<IdReceta>" o mensaje de error
	const idPuntoCargaFarmacia = 5 // Punto de carga Farmacia (el más usado en recetas reales)
	var respuesta string
	err = tx.QueryRowContext(ctx, `
		EXEC webRecetaCabeceraAgregar
			@Respuesta = @p1 OUTPUT,
			@IdPuntoCarga = @p2,
			@idCuentaAtencion = @p3,
			@idServicioReceta = @p4,
			@idMedicoReceta = @p5,
			@IdProducto = @p6,
			@Idpaciente = @p7,
			@IdUsuarioAuditoria = @p8`,
		sql.Named("p1", sql.Out{Dest: &respuesta}),
		sql.Named("p2", idPuntoCargaFarmacia),
		sql.Named("p3", idCuentaAtencion),
		sql.Named("p4", idServicioIngreso),
		sql.Named("p5", idMedico),
		sql.Named("p6", detalles[0].IdProducto),
		sql.Named("p7", idPaciente),
		sql.Named("p8", idEmpleado),
	).Scan(&respuesta)
	if err != nil {
		return fmt.Errorf("creando cabecera de receta: %w", err)
	}

	idReceta, ok := parsearRespuestaReceta(respuesta)
	if !ok {
		return fmt.Errorf("el sistema rechazó la receta: %s", respuesta)
	}

	// 4. Agregar cada detalle (SP real)
	// Params: @Mensaje OUTPUT, @idReceta, @IdProducto, @Cantidad, @Precio,
	//         @SaldoEnRegistroReceta, @idDosisRecetada, @observaciones, @IdViaAdministracion,
	//         @CodigoDiagnostico, @Justificacion, @idUNIDDosisReceta, @idFrecuencia,
	//         @DescripcionadicionalReceta, @Duracion, @idCuentaAtencionProxCita
	for _, det := range detalles {
		precio, err := r.precioProducto(ctx, tx, det.IdProducto)
		if err != nil {
			return err
		}

		var mensaje string
		err = tx.QueryRowContext(ctx, `
			EXEC webRecetaDetalleAgregar
				@Mensaje = @p1 OUTPUT,
				@idReceta = @p2,
				@IdProducto = @p3,
				@Cantidad = @p4,
				@Precio = @p5,
				@SaldoEnRegistroReceta = @p6,
				@idDosisRecetada = @p7,
				@observaciones = @p8,
				@IdViaAdministracion = @p9,
				@CodigoDiagnostico = @p10,
				@Justificacion = @p11,
				@idUNIDDosisReceta = @p12,
				@idFrecuencia = @p13,
				@DescripcionadicionalReceta = @p14,
				@Duracion = @p15,
				@idCuentaAtencionProxCita = @p16`,
			sql.Named("p1", sql.Out{Dest: &mensaje}),
			sql.Named("p2", idReceta),
			sql.Named("p3", det.IdProducto),
			sql.Named("p4", det.Cantidad),
			sql.Named("p5", precio),
			sql.Named("p6", nil),
			sql.Named("p7", nil),
			sql.Named("p8", det.Indicaciones),
			sql.Named("p9", nil),
			sql.Named("p10", nil),
			sql.Named("p11", nil),
			sql.Named("p12", nil),
			sql.Named("p13", nil),
			sql.Named("p14", nil),
			sql.Named("p15", nil),
			sql.Named("p16", nil),
		).Scan(&mensaje)
		if err != nil {
			return fmt.Errorf("agregando detalle %d a la receta %d: %w", det.IdProducto, idReceta, err)
		}
		if !strings.HasPrefix(strings.TrimSpace(mensaje), "OK") {
			return fmt.Errorf("el sistema rechazó el producto %d: %s", det.IdProducto, mensaje)
		}
	}

	return tx.Commit()
}

func (r *OrdenRepository) BuscarProductos(ctx context.Context, filtro string, limite int) ([]domain.ProductoCatalogo, error) {
	if limite <= 0 || limite > 100 {
		limite = 50
	}

	query := `
		SELECT TOP (@p3)
			f.IdProducto, f.Codigo, f.Nombre,
			ISNULL(f.Concentracion, ''), ISNULL(f.Presentacion, ''),
			ISNULL(f.FormaFarmaceutica, ''),
			ISNULL((SELECT TOP 1 h.PrecioUnitario FROM FactCatalogoBienesInsumosHosp h
			         WHERE h.IdProducto = f.IdProducto AND h.Activo = 1
			         ORDER BY CASE WHEN h.IdTipoFinanciamiento = @p2 THEN 0 ELSE 1 END, h.IdTipoFinanciamiento), 0)
		FROM FactCatalogoBienesInsumos f
		WHERE f.Estado = 1
		  AND (f.Nombre LIKE @p1 OR ISNULL(f.NombreComercial, '') LIKE @p1 OR f.Codigo LIKE @p1)
		ORDER BY f.Nombre`

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", "%"+filtro+"%"),
		sql.Named("p2", filtro),
		sql.Named("p3", limite),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []domain.ProductoCatalogo
	for rows.Next() {
		var p domain.ProductoCatalogo
		if err := rows.Scan(&p.IdProducto, &p.Codigo, &p.Nombre, &p.Concentracion, &p.Presentacion, &p.FormaFarmaceutica, &p.PrecioVenta); err != nil {
			return nil, fmt.Errorf("error escaneando producto: %w", err)
		}
		productos = append(productos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productos, nil
}

// resolverCuentaAtencion obtiene IdCuentaAtencion desde el IdAtencion enviado
// por el frontend (idRegAtencion). Si el valor ya es una cuenta válida, lo usa directo.
func (r *OrdenRepository) resolverCuentaAtencion(ctx context.Context, idRegAtencion int) (int, error) {
	var idCuenta int
	err := r.db.QueryRowContext(ctx,
		"SELECT IdCuentaAtencion FROM Atenciones WHERE IdAtencion = @p1",
		sql.Named("p1", idRegAtencion),
	).Scan(&idCuenta)
	if err == nil {
		return idCuenta, nil
	}
	if err != sql.ErrNoRows {
		return 0, err
	}
	return idRegAtencion, nil
}

// precioProducto obtiene el precio unitario vigente del catálogo hospitalario
// (por producto y tipo de financiamiento de la atención, con prioridad al tipo de la atención).
func (r *OrdenRepository) precioProducto(ctx context.Context, tx *sql.Tx, idProducto int) (float64, error) {
	var precio float64
	err := tx.QueryRowContext(ctx, `
		SELECT TOP 1 PrecioUnitario FROM FactCatalogoBienesInsumosHosp
		WHERE IdProducto = @p1 AND Activo = 1 ORDER BY IdTipoFinanciamiento`,
		sql.Named("p1", idProducto),
	).Scan(&precio)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("obteniendo precio del producto %d: %w", idProducto, err)
	}
	return precio, nil
}

// parsearRespuestaReceta extrae el IdReceta de una respuesta "OK;<IdReceta>".
func parsearRespuestaReceta(respuesta string) (int, bool) {
	respuesta = strings.TrimSpace(respuesta)
	if !strings.HasPrefix(respuesta, "OK;") {
		return 0, false
	}
	var idReceta int
	if _, err := fmt.Sscanf(strings.TrimPrefix(respuesta, "OK;"), "%d", &idReceta); err != nil {
		return 0, false
	}
	return idReceta, true
}
