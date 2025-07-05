package repository

import (
	"context"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type VentaRepository struct {
	pool *pgxpool.Pool
}

func (v VentaRepository) FacturarVentaById(ctx context.Context, id *int) error {
	//TODO implement me
	return nil
}

func (v VentaRepository) ObtenerListaVentas(ctx context.Context) (*[]domain.VentaInfo, error) {
	query := `SELECT v.id,v.codigo,v.estado,v.fecha,v.usuario,v.cliente,v.deleted_at FROM view_venta_info v`
	rows, err := v.pool.Query(ctx, query)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()
	var list = make([]domain.VentaInfo, 0)
	for rows.Next() {
		var item domain.VentaInfo
		err = rows.Scan(&item.Id, &item.Codigo, &item.Estado, &item.Fecha, &item.Usuario, &item.Cliente, &item.DeletedAt)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, item)
	}
	return &list, nil
}

func (v VentaRepository) RegistraVenta(ctx context.Context, request *domain.VentaRequest) (*int64, error) {
	// Validar que la venta tenga detalles
	if len(request.Detalles) == 0 {
		return nil, datatype.NewBadRequestError("La venta debe tener al menos un detalle")
	}

	tx, err := v.pool.Begin(ctx)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	// Generar código de venta de forma más eficiente
	var nextNum int64
	err = tx.QueryRow(ctx, `
        SELECT COALESCE(
            (SELECT MAX(CAST(SUBSTRING(codigo FROM 6) AS INTEGER)) + 1 FROM venta WHERE codigo ~ '^VENT-[0-9]+$'),
            1
        )
    `).Scan(&nextNum)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	codigo := fmt.Sprintf("VENT-%09d", nextNum)

	// Crear la venta
	var ventaId int64
	err = tx.QueryRow(ctx, `
        INSERT INTO venta (cliente_id, usuario_id, total, codigo)
        VALUES ($1, $2, 0, $3)
        RETURNING id
    `, request.ClienteId, request.UsuarioId, codigo).Scan(&ventaId)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	totalVenta := 0.0

	// Procesar cada detalle de venta
	for _, item := range request.Detalles {
		if item.Cantidad <= 0 {
			return nil, datatype.NewBadRequestError("La cantidad debe ser mayor a cero")
		}

		// Obtener lotes disponibles ordenados por FEFO (First Expired, First Out) con bloqueo
		rows, err := tx.Query(ctx, `
            SELECT lp.id, lp.stock, p.precio_venta
            FROM lote_producto lp
            JOIN producto p ON p.id = lp.producto_id
            WHERE lp.producto_id = $1 
              AND lp.stock > 0 
              AND lp.estado = 'Activo'
            ORDER BY lp.fecha_vencimiento ASC, lp.id ASC
            FOR UPDATE OF lp
        `, item.ProductoId)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}

		var lotes []domain.VentaLoteProductoDAO
		for rows.Next() {
			var lote domain.VentaLoteProductoDAO
			if err := rows.Scan(&lote.Id, &lote.Stock, &lote.PrecioVenta); err != nil {
				rows.Close()
				return nil, datatype.NewInternalServerErrorGeneric()
			}
			lotes = append(lotes, lote)
		}
		rows.Close()

		if len(lotes) == 0 {
			return nil, datatype.NewNotFoundErrorWithData("Producto sin stock disponible",
				domain.ProductoId{Id: item.ProductoId})
		}

		// Verificar stock total disponible
		stockTotal := uint(0)
		for _, lote := range lotes {
			stockTotal += lote.Stock
		}

		if stockTotal < item.Cantidad {
			return nil, datatype.NewNotFoundErrorWithData(
				fmt.Sprintf("Stock insuficiente. Disponible: %d, Solicitado: %d", stockTotal, item.Cantidad),
				domain.ProductoId{Id: item.ProductoId})
		}

		// Asignar stock desde los lotes usando FEFO
		cantidadRestante := item.Cantidad
		subtotalItem := 0.0

		for _, lote := range lotes {
			if cantidadRestante <= 0 {
				break
			}

			cantidadUsar := cantidadRestante
			if lote.Stock < cantidadUsar {
				cantidadUsar = lote.Stock
			}

			// Crear detalle de venta
			_, err = tx.Exec(ctx, `
                INSERT INTO detalle_venta (venta_id, lote_id, cantidad, precio)
                VALUES ($1, $2, $3, $4)
            `, ventaId, lote.Id, cantidadUsar, lote.PrecioVenta)
			if err != nil {
				return nil, datatype.NewInternalServerErrorGeneric()
			}

			// Actualizar stock del lote con verificación
			result, err := tx.Exec(ctx, `
                UPDATE lote_producto 
                SET stock = stock - $1
                WHERE id = $2 AND stock >= $1
            `, cantidadUsar, lote.Id)
			if err != nil {
				return nil, datatype.NewInternalServerErrorGeneric()
			}
			if result.RowsAffected() == 0 {
				return nil, datatype.NewConflictError("Stock insuficiente en el lote")
			}

			// Actualizar stock del producto principal con verificación
			result, err = tx.Exec(ctx, `
                UPDATE producto 
                SET stock = stock - $1
                WHERE id = $2 AND stock >= $1
            `, cantidadUsar, item.ProductoId)
			if err != nil {
				return nil, datatype.NewInternalServerErrorGeneric()
			}
			if result.RowsAffected() == 0 {
				return nil, datatype.NewConflictError("Stock insuficiente en el producto")
			}

			subtotalItem += float64(cantidadUsar) * lote.PrecioVenta
			cantidadRestante -= cantidadUsar
		}

		totalVenta += subtotalItem
	}

	// Actualizar total de la venta
	_, err = tx.Exec(ctx, `
        UPDATE venta 
        SET total = $1,
            fecha = NOW()
        WHERE id = $2
    `, totalVenta, ventaId)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &ventaId, nil
}

func (v VentaRepository) ObtenerVentaById(ctx context.Context, id *int) (*domain.VentaDetail, error) {
	query := `
	SELECT
		v.id,
		v.codigo,
		v.fecha,
		v.estado,
		v.deleted_at,
		v.usuario,
		v.cliente,
		(
			SELECT jsonb_agg(d)
			FROM view_detalle_venta_producto_detail d
			WHERE d.venta_id = v.id
		) AS detalles
	FROM view_venta_info v
	WHERE v.id = $1
	LIMIT 1;
`

	var venta domain.VentaDetail
	err := v.pool.QueryRow(ctx, query, *id).Scan(&venta.Id, &venta.Codigo, &venta.Fecha, &venta.Estado, &venta.DeletedAt, &venta.Usuario, &venta.Cliente, &venta.Detalles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Venta no encontrada")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &venta, nil
}

func (v VentaRepository) AnularVentaById(ctx context.Context, id *int) error {
	tx, err := v.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("Error rolling back transaction: %v", err)
		}
	}()

	// Verificar que la venta existe y obtener su estado
	var existe bool
	var estadoActual string
	query := `SELECT EXISTS(SELECT 1 FROM venta WHERE id = $1), (SELECT estado FROM venta WHERE id = $1) as estado`
	if err := tx.QueryRow(ctx, query, *id).Scan(&existe, &estadoActual); err != nil {
		log.Println("Error al verificar existencia de venta:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	if !existe {
		return datatype.NewNotFoundError("Venta no encontrada")
	}

	if estadoActual == "Anulado" {
		return datatype.NewBadRequestError("La venta ya está anulada")
	}

	// Verificar que existen detalles de venta
	var tieneDetalles bool
	query = `SELECT EXISTS(SELECT 1 FROM detalle_venta dv WHERE dv.venta_id = $1)`
	if err := tx.QueryRow(ctx, query, *id).Scan(&tieneDetalles); err != nil {
		log.Println("Error al verificar existencia de detalles de venta:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	if !tieneDetalles {
		return datatype.NewBadRequestError("No se encontraron detalles de venta")
	}

	// Obtener los lotes y productos involucrados en la venta CON BLOQUEO
	query = `
        SELECT dv.lote_id, dv.cantidad, lp.producto_id
        FROM detalle_venta dv 
        INNER JOIN lote_producto lp ON dv.lote_id = lp.id
        WHERE dv.venta_id = $1
        ORDER BY lp.producto_id, dv.lote_id
        FOR UPDATE OF lp
    `
	rows, err := tx.Query(ctx, query, *id)
	if err != nil {
		log.Println("Error al obtener detalle de venta:", err)
		return datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	productosMap := make(map[string][]domain.VentaLoteProducto)
	var todosLosLotes []domain.VentaLoteProducto

	for rows.Next() {
		var item domain.VentaLoteProducto
		if err := rows.Scan(&item.Id, &item.Cantidad, &item.ProductoId); err != nil {
			log.Println("Error escaneando detalle de venta:", err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Agrupar por producto para procesamiento posterior
		productosMap[item.ProductoId] = append(productosMap[item.ProductoId], item)
		todosLosLotes = append(todosLosLotes, item)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error leyendo filas:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	if len(todosLosLotes) == 0 {
		log.Println("No hay lotes asociados a la venta")
		return datatype.NewBadRequestError("No se encontraron detalles de venta")
	}

	// Procesar cada producto de manera segura
	for productoId, lotes := range productosMap {
		// Bloquear el producto antes de actualizar
		query = `SELECT id FROM producto WHERE id = $1 FOR UPDATE`
		if _, err = tx.Exec(ctx, query, productoId); err != nil {
			log.Printf("Error bloqueando producto %s: %v", productoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Restaurar stock en cada lote del producto
		for _, lote := range lotes {
			// Bloquear el lote antes de actualizar
			query = `SELECT id FROM lote_producto WHERE id = $1 FOR UPDATE`
			if _, err = tx.Exec(ctx, query, lote.Id); err != nil {
				log.Printf("Error bloqueando lote %d: %v", lote.Id, err)
				return datatype.NewInternalServerErrorGeneric()
			}

			query = `UPDATE lote_producto 
                     SET stock = stock + $1
                     WHERE id = $2`
			result, err := tx.Exec(ctx, query, lote.Cantidad, lote.Id)
			if err != nil {
				log.Printf("Error actualizando stock del lote %d: %v", lote.Id, err)
				return datatype.NewInternalServerErrorGeneric()
			}

			if result.RowsAffected() == 0 {
				log.Printf("No se pudo actualizar el lote %d", lote.Id)
				return datatype.NewInternalServerErrorGeneric()
			}
		}

		// Recalcular y actualizar el stock total del producto
		var nuevoStockTotal int64
		query = `SELECT COALESCE(SUM(stock), 0) 
                 FROM lote_producto 
                 WHERE producto_id = $1 AND estado = 'Activo'`
		if err := tx.QueryRow(ctx, query, productoId).Scan(&nuevoStockTotal); err != nil {
			log.Printf("Error al calcular el stock total del producto %s: %v", productoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Actualizar stock del producto con verificación
		query = `UPDATE producto 
                 SET stock = $1
                 WHERE id = $2`
		result, err := tx.Exec(ctx, query, nuevoStockTotal, productoId)
		if err != nil {
			log.Printf("Error al actualizar el stock del producto %s: %v", productoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		if result.RowsAffected() == 0 {
			log.Printf("No se pudo actualizar el producto %s", productoId)
			return datatype.NewInternalServerErrorGeneric()
		}

		log.Printf("Producto %s: stock actualizado a %d", productoId, nuevoStockTotal)
	}

	// Marcar el estado de la venta como 'Anulado'
	query = `UPDATE venta 
             SET estado = 'Anulado', 
                 deleted_at = CURRENT_TIMESTAMP
             WHERE id = $1 AND estado != 'Anulado'`
	result, err := tx.Exec(ctx, query, *id)
	if err != nil {
		log.Println("Error marcando venta como anulada:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	if result.RowsAffected() == 0 {
		return datatype.NewConflictError("No se pudo anular la venta, posiblemente ya está anulada")
	}

	if err := tx.Commit(ctx); err != nil {
		log.Println("Error al hacer commit:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	log.Printf("Venta %d anulada exitosamente", *id)
	return nil
}

func NewVentaRepository(pool *pgxpool.Pool) *VentaRepository {
	return &VentaRepository{pool: pool}
}

var _ port.VentaRepository = (*VentaRepository)(nil)
