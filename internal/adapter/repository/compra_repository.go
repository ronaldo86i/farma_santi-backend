package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CompraRepository struct {
	pool *pgxpool.Pool
}

func (c CompraRepository) ObtenerCompraById(ctx context.Context, id *int) (*domain.CompraDetail, error) {
	query := `SELECT v.id,v.codigo,v.estado,v.total,v.comentario,v.laboratorio,v.usuario,v.fecha,deleted_at,v.detalles FROM view_compra_con_detalles v WHERE id = $1`
	var compra domain.CompraDetail
	err := c.pool.QueryRow(ctx, query, *id).Scan(&compra.Id, &compra.Codigo, &compra.Estado, &compra.Total, &compra.Comentario, &compra.Laboratorio, &compra.Usuario, &compra.Fecha, &compra.DeletedAt, &compra.Detalles)
	if err != nil {
		log.Println("Error al obtener compra:", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Registro de compra no encontrada")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &compra, nil
}

func (c CompraRepository) RegistrarOrdenCompra(ctx context.Context, request *domain.CompraRequest) (*uint, error) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error al iniciar transacción: %v", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Generar código de venta de forma más eficiente
	var nextNum int64
	err = tx.QueryRow(ctx, `
        SELECT COALESCE(
            (SELECT MAX(CAST(SUBSTRING(codigo FROM 6) AS INTEGER)) + 1 FROM compra WHERE codigo ~ '^COMP-[0-9]+$'),
            1
        )
    `).Scan(&nextNum)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	codigo := fmt.Sprintf("COMP-%09d", nextNum)

	var total float64
	for _, detalle := range request.Detalles {
		total += detalle.PrecioCompra * float64(detalle.Cantidad)
	}

	query := `INSERT INTO compra(usuario_id,laboratorio_id,comentario,total,codigo) VALUES($1,$2,$3,$4,$5) RETURNING id`

	var compraId uint
	err = tx.QueryRow(ctx, query, request.UsuarioId, request.LaboratorioId, request.Comentario, total, codigo).Scan(&compraId)
	if err != nil {
		log.Println("Ha ocurrido un error al insertar la compra:", err.Error())
		return nil, datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	for _, detalle := range request.Detalles {
		query = `INSERT INTO detalle_compra(cantidad, precio_compra,precio_venta, compra_id, lote_producto_id) VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.Exec(ctx, query, detalle.Cantidad, detalle.PrecioCompra, detalle.PrecioVenta, compraId, detalle.LoteProductoId)
		if err != nil {
			log.Println("Ha ocurrido un error al insertar detalles de la compra:", err.Error())
			return nil, datatype.NewStatusServiceUnavailableErrorGeneric()
		}
	}

	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return &compraId, nil
}

func (c CompraRepository) ModificarOrdenCompra(ctx context.Context, id *int, request *domain.CompraRequest) error {
	// Iniciar transacción
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error al iniciar transacción: %v", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	// Verificar si existe la compra
	var existe bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM compra WHERE id = $1 AND estado='Pendiente')`
	err = tx.QueryRow(ctx, checkQuery, *id).Scan(&existe)
	if err != nil {
		log.Println("Error al verificar existencia de compra:", err)
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	if !existe {
		return datatype.NewNotFoundError("La compra no existe o fue eliminada")
	}

	var total float64
	for _, detalle := range request.Detalles {
		total += detalle.PrecioCompra * float64(detalle.Cantidad)
	}
	// Consulta
	query := `UPDATE compra c SET comentario=$1, total=$2,laboratorio_id=$3,updated_at=CURRENT_TIMESTAMP WHERE id=$4`
	_, err = tx.Exec(ctx, query, request.Comentario, total, request.LaboratorioId, *id)
	if err != nil {
		log.Println("Ha ocurrido un error al modificar la compra:", err.Error())
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	query = `DELETE FROM detalle_compra WHERE compra_id = $1`
	_, err = tx.Exec(ctx, query, *id)
	if err != nil {
		log.Println("Ha ocurrido un error al modificar detalles de compra:", err.Error())
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	for _, detalle := range request.Detalles {
		query = `INSERT INTO detalle_compra(cantidad, precio_compra,precio_venta, compra_id, lote_producto_id) VALUES ($1, $2, $3, $4, $5)`
		_, err := tx.Exec(ctx, query, detalle.Cantidad, detalle.PrecioCompra, detalle.PrecioVenta, *id, detalle.LoteProductoId)
		if err != nil {
			log.Println("Ha ocurrido un error al insertar detalles de la compra:", err.Error())
			return datatype.NewStatusServiceUnavailableErrorGeneric()
		}
	}
	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (c CompraRepository) AnularOrdenCompra(ctx context.Context, id *int) error {
	// Iniciar transacción
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Printf("Error al iniciar transacción: %v", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Verificar si existe la compra y está en estado 'Pendiente'
	var existe bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM compra WHERE id = $1 AND estado = 'Pendiente')`
	err = tx.QueryRow(ctx, checkQuery, *id).Scan(&existe)
	if err != nil {
		log.Println("Error al verificar existencia de compra:", err)
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	if !existe {
		return datatype.NewNotFoundError("La compra no fue encontrada o no está en estado pendiente")
	}

	// Actualizar el estado a 'Anulado'
	updateQuery := `UPDATE compra SET estado = 'Anulado',deleted_at=CURRENT_TIMESTAMP WHERE id = $1`
	_, err = tx.Exec(ctx, updateQuery, *id)
	if err != nil {
		log.Println("Error al anular compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (c CompraRepository) RegistrarCompra(ctx context.Context, id *int) error {
	var compra domain.CompraDAO

	query := `SELECT c.id, c.estado, c.total, c.comentario, c.laboratorio_id, c.usuario_id, c.detalles 
	          FROM view_compras_detalle c 
	          WHERE id = $1 LIMIT 1`

	err := c.pool.QueryRow(ctx, query, *id).Scan(&compra.Id, &compra.Estado, &compra.Total, &compra.Comentario, &compra.LaboratorioId, &compra.UsuarioId, &compra.Detalles)
	if err != nil {
		log.Println("Error al consultar la compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	estado := strings.ToLower(compra.Estado)
	switch estado {
	case "completado":
		log.Println("La compra ya está completada, no se puede volver a registrar.")
		return datatype.NewConflictError("La compra ya fue registrada y completada")
	case "anulado":
		log.Println("La compra ya está anulada, no se puede registrar.")
		return datatype.NewConflictError("La compra ya fue anulada")
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		log.Println("Error al iniciar transacción:", err)
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Lock de compra
	lockCompraQuery := `SELECT id FROM compra WHERE id = $1 FOR UPDATE`
	_, err = tx.Exec(ctx, lockCompraQuery, *id)
	if err != nil {
		log.Println("Error al bloquear compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	for _, detalle := range compra.Detalles {
		// Lock de lote_producto
		lockLoteQuery := `SELECT id FROM lote_producto WHERE id = $1 FOR UPDATE`
		_, err := tx.Exec(ctx, lockLoteQuery, detalle.LoteProductoId)
		if err != nil {
			log.Printf("Error al bloquear lote_producto %d: %v", detalle.LoteProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Lock de producto
		lockProductoQuery := `SELECT stock FROM producto WHERE id = $1 FOR UPDATE`
		var stockActual uint
		err = tx.QueryRow(ctx, lockProductoQuery, detalle.ProductoId).Scan(&stockActual)
		if err != nil {
			log.Printf("Error al obtener stock del producto %d: %v", detalle.ProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Actualizar stock en lote_producto
		updateLoteQuery := `UPDATE lote_producto SET stock = stock + $1 WHERE id = $2`
		_, err = tx.Exec(ctx, updateLoteQuery, detalle.Cantidad, detalle.LoteProductoId)
		if err != nil {
			log.Printf("Error al actualizar stock del lote %d: %v", detalle.LoteProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}

		// Actualizar producto con stock y precios nuevos
		updateProductoQuery := `UPDATE producto 
		                        SET stock = stock + $1, 
		                            precio_compra = $2, 
		                            precio_venta = $3 
		                        WHERE id = $4`
		_, err = tx.Exec(ctx, updateProductoQuery, detalle.Cantidad, detalle.PrecioCompra, detalle.PrecioVenta, detalle.ProductoId)
		if err != nil {
			log.Printf("Error al actualizar producto %d: %v", detalle.ProductoId, err)
			return datatype.NewInternalServerErrorGeneric()
		}
	}

	// Actualizar estado de la compra
	updateEstadoQuery := `UPDATE compra SET estado = 'Completado', fecha = CURRENT_TIMESTAMP WHERE id = $1`
	_, err = tx.Exec(ctx, updateEstadoQuery, *id)
	if err != nil {
		log.Println("Error al actualizar estado de la compra:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	err = tx.Commit(ctx)
	if err != nil {
		log.Println("Error al confirmar transacción:", err)
		return datatype.NewInternalServerErrorGeneric()
	}

	committed = true
	return nil
}

func (c CompraRepository) ObtenerListaCompras(ctx context.Context, filtros map[string]string) (*[]domain.CompraInfo, error) {
	query := `SELECT c.id,c.codigo,c.comentario,c.estado,c.total,c.laboratorio,c.usuario,c.fecha FROM view_compras c`
	var filters []string
	var args []interface{}
	i := 1

	// Sí hay estado en filtros
	if estadoStr := filtros["estado"]; estadoStr != "" {
		filters = append(filters, fmt.Sprintf("c.estado = $%d", i))
		args = append(args, estadoStr)
		i++
	}
	// Filtrar por fechaInicio
	if fechaInicioStr := filtros["fechaInicio"]; fechaInicioStr != "" {
		fechaInicio, err := time.Parse("2006-01-02", fechaInicioStr)
		if err != nil {
			fechaInicio, err = time.Parse(time.RFC3339, fechaInicioStr)
			if err != nil {
				log.Println("Error al convertir fechaInicio:", err)
				return nil, datatype.NewBadRequestError("El valor de fechaInicio no es válido, formatos esperados: YYYY-MM-DD o RFC3339")
			}
		}
		filters = append(filters, fmt.Sprintf("c.fecha >= $%d", i))
		args = append(args, fechaInicio)
		i++
	}

	// Filtrar por fechaFin
	if fechaFinStr := filtros["fechaFin"]; fechaFinStr != "" {
		fechaFin, err := time.Parse("2006-01-02", fechaFinStr)
		if err != nil {
			fechaFin, err = time.Parse(time.RFC3339, fechaFinStr)
			if err != nil {
				log.Println("Error al convertir fechaFin:", err)
				return nil, datatype.NewBadRequestError("El valor de fechaFin no es válido, formatos esperados: YYYY-MM-DD o RFC3339")
			}
		}
		filters = append(filters, fmt.Sprintf("c.fecha <= $%d", i))
		args = append(args, fechaFin)
		i++
	}
	// Aplicar LIMIT (y OFFSET si existe)
	if limitStr := filtros["limit"]; limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return nil, datatype.NewBadRequestError("El valor de limit debe ser un número entero positivo")
		}
		query += fmt.Sprintf(" LIMIT $%d", i)
		args = append(args, limit)
		i++

		if offsetStr := filtros["offset"]; offsetStr != "" {
			offset, err := strconv.Atoi(offsetStr)
			if err != nil || offset < 0 {
				return nil, datatype.NewBadRequestError("El valor de offset debe ser un número entero no negativo")
			}
			query += fmt.Sprintf(" OFFSET $%d", i)
			args = append(args, offset)
			i++
		}
	}

	// Si hay filtros, agregarlos al query
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list = make([]domain.CompraInfo, 0)
	for rows.Next() {
		var item domain.CompraInfo
		err = rows.Scan(&item.Id, &item.Codigo, &item.Comentario, &item.Estado, &item.Total, &item.Laboratorio, &item.Usuario, &item.Fecha)
		if err != nil {
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &list, nil
}

func NewCompraRepository(pool *pgxpool.Pool) *CompraRepository {
	return &CompraRepository{pool: pool}
}

var _ port.CompraRepository = (*CompraRepository)(nil)
