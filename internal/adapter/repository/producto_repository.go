package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
)

type ProductoRepository struct {
	pool *pgxpool.Pool
}

func (p ProductoRepository) ObtenerProductoById(ctx context.Context, id *uuid.UUID) (*domain.ProductoDetail, error) {
	fullHostname := ctx.Value("fullHostname").(string)
	fullHostname = fmt.Sprintf("%s%s", fullHostname, "/uploads/productos")
	query := `SELECT p.id,p.nombre_comercial,p.forma_farmaceutica,p.laboratorio,p.precio_venta,p.stock_min,p.stock,p.fotos,p.created_at,p.deleted_at,p.estado,p.categorias,p.principio_activos,p.precio_compra,p.presentacion,p.unidades_presentacion FROM obtener_producto_detalle_by_id($1,$2) p;`
	var item domain.ProductoDetail
	err := p.pool.QueryRow(ctx, query, id.String(), fullHostname).Scan(&item.Id, &item.NombreComercial, &item.FormaFarmaceutica,
		&item.Laboratorio, &item.PrecioVenta, &item.StockMin, &item.Stock, &item.UrlFotos, &item.CreatedAt, &item.DeletedAt, &item.Estado, &item.Categorias,
		&item.PrincipiosActivos, &item.PrecioCompra, &item.Presentacion, &item.UnidadesPresentacion)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Producto no encontrado")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &item, nil
}

func (p ProductoRepository) HabilitarProducto(ctx context.Context, id *uuid.UUID) error {
	// Inicio de transacción

	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	query := `UPDATE producto SET estado='Activo',deleted_at=NULL WHERE id=$1`
	ct, err := tx.Exec(ctx, query, id.String())
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	if ct.RowsAffected() == 0 {
		return datatype.NewNotFoundError("Producto no encontrado")
	}
	// Confirmar la transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (p ProductoRepository) DeshabilitarProducto(ctx context.Context, id *uuid.UUID) error {
	// Inicio de transacción
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	query := `UPDATE producto SET estado='Inactivo',deleted_at=CURRENT_TIMESTAMP WHERE id=$1`
	ct, err := tx.Exec(ctx, query, id.String())
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	if ct.RowsAffected() == 0 {
		return datatype.NewNotFoundError("Producto no encontrado")
	}

	// Confirmar la transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (p ProductoRepository) RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error {
	// Inicio de transacción
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	query := `INSERT INTO producto(nombre_comercial,forma_farmaceutica_id,precio_compra,precio_venta,estado,stock,stock_min,laboratorio_id,presentacion_id,unidades_presentacion) 
				VALUES ($1,$2,0.0,$3,'Activo',0,$4,$5,$6,$7) RETURNING id`

	var id uuid.UUID
	err = tx.QueryRow(ctx, query, request.NombreComercial, request.FormaFarmaceuticaId, request.PrecioVenta, request.StockMin, request.LaboratorioId, request.PresentacionId, request.UnidadesPresentacion).Scan(&id)
	if err != nil {
		_ = tx.Rollback(ctx)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return datatype.NewConflictError("Ya existe ese producto")
		}
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	route := fmt.Sprintf("./public/uploads/productos/%s", id.String())
	defer func() {
		if err != nil {
			_ = util.File.DeleteAllFiles(route)
		}
	}()
	err = util.File.MakeDir(route)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	var fotos []string
	for i, fileHeader := range *filesHeader {
		file, err := fileHeader.Open()
		if err != nil {
			return datatype.NewStatusServiceUnavailableErrorGeneric()
		}
		ext := filepath.Ext(fileHeader.Filename)
		nameFile := fmt.Sprintf("%d%s", i+1, ext)

		err = util.File.SaveFile(route, nameFile, file)
		if err != nil {
			_ = util.File.DeleteAllFiles(route)
			return datatype.NewInternalServerError("Error al guardar fotos")
		}
		_ = file.Close()
		fotos = append(fotos, nameFile)
	}

	query = `INSERT INTO producto_categoria(producto_id, categoria_id) SELECT $1, unnest($2::int[])`
	_, err = tx.Exec(ctx, query, id, pq.Array(request.Categorias))
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	for _, pa := range request.PrincipiosActivos {
		query = `INSERT INTO producto_principio_activo(producto_id, principio_activo_id, concentracion, unidad_medida_id) 
	          VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(ctx, query, id, pa.PrincipioActivoId, pa.Concentracion, pa.UnidadMedidaId)
		if err != nil {
			return datatype.NewStatusServiceUnavailableErrorGeneric()
		}
	}
	query = `UPDATE producto SET fotos = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, pq.Array(fotos), id)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	// Confirmar la transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (p ProductoRepository) ModificarProducto(ctx context.Context, id *uuid.UUID, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) (err error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Respaldar archivos existentes
	route := fmt.Sprintf("./public/uploads/productos/%s", id.String())
	backupFiles, err := util.File.BackupFiles(route)
	if err != nil {
		log.Println(err)
		return datatype.NewInternalServerErrorGeneric()
	}

	// Ejecutar SQL update
	query := `UPDATE producto SET nombre_comercial=$1,forma_farmaceutica_id=$2,stock_min=$3,laboratorio_id=$4,presentacion_id=$5,precio_venta=$6,unidades_presentacion=$7 WHERE id=$8`
	ct, err := tx.Exec(ctx, query, request.NombreComercial, request.FormaFarmaceuticaId, request.StockMin, request.LaboratorioId, request.PresentacionId, request.PrecioVenta, request.UnidadesPresentacion, id.String())
	if err != nil {
		log.Println(err)
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return datatype.NewConflictError("Ya existe ese producto")
		}
		return datatype.NewInternalServerErrorGeneric()
	}
	if ct.RowsAffected() == 0 {
		return datatype.NewNotFoundError("No existe el producto")
	}
	query = `DELETE FROM producto_categoria WHERE producto_id = $1`
	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	query = `INSERT INTO producto_categoria(producto_id, categoria_id) SELECT $1, unnest($2::int[])`
	_, err = tx.Exec(ctx, query, id, pq.Array(request.Categorias))
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	// Guardar nuevos archivos
	var nuevosArchivos []string
	for i, fileHeader := range *filesHeader {
		log.Println(fileHeader.Filename)
		file, err := fileHeader.Open()
		if err != nil {
			return datatype.NewStatusServiceUnavailableErrorGeneric()
		}

		ext := filepath.Ext(fileHeader.Filename)
		nameFile := fmt.Sprintf("%d%s", i+1, ext)

		saveErr := util.File.SaveFile(route, nameFile, file)
		_ = file.Close()
		if saveErr != nil {
			util.File.DeleteFiles(route, nuevosArchivos)
			_ = util.File.RestoreFiles(backupFiles, route)
			return datatype.NewInternalServerError("Error al guardar fotos")
		}
		nuevosArchivos = append(nuevosArchivos, nameFile)
	}
	query = `UPDATE producto SET fotos = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, pq.Array(nuevosArchivos), id)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	// Eliminar principios activos existentes
	query = `DELETE FROM producto_principio_activo WHERE producto_id = $1`
	_, err = tx.Exec(ctx, query, id)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	// Insertar nuevos principios activos
	for _, pa := range request.PrincipiosActivos {
		query = `INSERT INTO producto_principio_activo(producto_id, principio_activo_id, concentracion, unidad_medida_id)
	         VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(ctx, query, id, pa.PrincipioActivoId, pa.Concentracion, pa.UnidadMedidaId)
		if err != nil {
			return datatype.NewStatusServiceUnavailableErrorGeneric()
		}
	}
	// Confirmar transacción
	if err = tx.Commit(ctx); err != nil {
		util.File.DeleteFiles(route, nuevosArchivos)
		_ = util.File.RestoreFiles(backupFiles, route)
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (p ProductoRepository) ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error) {
	query := `SELECT um.id,um.nombre,um.abreviatura FROM unidad_medida um ORDER BY um.nombre`
	var list = make([]domain.UnidadMedida, 0)

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.UnidadMedida
		err := rows.Scan(&item.Id, &item.Nombre, &item.Abreviatura)
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

func (p ProductoRepository) ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error) {
	query := `SELECT ff.id,ff.nombre FROM forma_farmaceutica ff ORDER BY ff.nombre`
	var list = make([]domain.FormaFarmaceutica, 0)

	rows, err := p.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.FormaFarmaceutica
		err := rows.Scan(&item.Id, &item.Nombre)
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

func (p ProductoRepository) ObtenerListaProductos(ctx context.Context, filtros map[string]string) (*[]domain.ProductoInfo, error) {
	fullHostname := ctx.Value("fullHostname").(string)
	fullHostname = fmt.Sprintf("%s%s", fullHostname, "/uploads/productos")

	var filters []string
	var args []interface{}
	args = append(args, fullHostname) // $1 siempre será la URL base
	i := 2

	if categoriasStr := filtros["categorias"]; categoriasStr != "" {
		ids := strings.Split(categoriasStr, ",")
		var categoriaIDs []int
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, fmt.Errorf("categoría inválida: %s", idStr)
			}
			categoriaIDs = append(categoriaIDs, id)
		}

		if len(categoriaIDs) > 0 {
			filters = append(filters, fmt.Sprintf("pc.categoria_id = ANY($%d)", i))
			args = append(args, categoriaIDs)
			i++
		}
	}
	// Filtro: categoriaId (un solo id)
	if categoriaIDStr := filtros["categoriaId"]; categoriaIDStr != "" {
		categoriaID, err := strconv.Atoi(strings.TrimSpace(categoriaIDStr))
		if err != nil {
			return nil, fmt.Errorf("categoriaId inválido: %s", categoriaIDStr)
		}
		filters = append(filters, fmt.Sprintf("pc.categoria_id = $%d", i))
		args = append(args, categoriaID)
		i++
	}

	// Filtro: laboratorios (puede ser "1,2,3")
	if laboratorioIDsStr := filtros["laboratorios"]; laboratorioIDsStr != "" {
		ids := strings.Split(laboratorioIDsStr, ",")
		var labIDs []int
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, fmt.Errorf("laboratorio inválido: %s", idStr)
			}
			labIDs = append(labIDs, id)
		}

		if len(labIDs) > 0 {
			filters = append(filters, fmt.Sprintf("p.laboratorio_id = ANY($%d)", i))
			args = append(args, labIDs)
			i++
		}
	}

	// Filtro: formasFarmaceuticas (puede ser "1,2,3")
	if formaIDsStr := filtros["formasFarmaceuticas"]; formaIDsStr != "" {
		ids := strings.Split(formaIDsStr, ",")
		var formaIDs []int
		for _, idStr := range ids {
			idStr = strings.TrimSpace(idStr)
			if idStr == "" {
				continue
			}
			id, err := strconv.Atoi(idStr)
			if err != nil {
				return nil, fmt.Errorf("forma farmacéutica inválida: %s", idStr)
			}
			formaIDs = append(formaIDs, id)
		}

		if len(formaIDs) > 0 {
			filters = append(filters, fmt.Sprintf("p.forma_farmaceutica_id = ANY($%d)", i))
			args = append(args, formaIDs)
			i++
		}
	}
	// Filtro: estado
	if estadoStr := filtros["estado"]; estadoStr != "" {
		filters = append(filters, fmt.Sprintf("p.estado = $%d", i))
		args = append(args, estadoStr)
		i++
	}

	// Filtro: laboratorio_id
	if laboratorioID := filtros["laboratorioId"]; laboratorioID != "" {
		filters = append(filters, fmt.Sprintf("p.laboratorio_id = $%d", i))
		args = append(args, laboratorioID)
		i++
	}

	// Filtro: forma_farmaceutica_id
	if formaID := filtros["formaFarmaceuticaId"]; formaID != "" {
		filters = append(filters, fmt.Sprintf("p.forma_farmaceutica_id = $%d", i))
		args = append(args, formaID)
		i++
	}

	// Filtro: nombre_comercial (LIKE para búsqueda parcial)
	if nombre := filtros["search"]; nombre != "" {
		nombre = strings.TrimSpace(nombre)
		filters = append(filters, fmt.Sprintf("LOWER(p.nombre_comercial) LIKE LOWER($%d)", i))
		args = append(args, "%"+nombre+"%")
		i++
	}

	query := `
		SELECT DISTINCT ON (p.id)
			p.id,
			p.nombre_comercial,
			p.forma_farmaceutica,
			p.laboratorio,
			p.precio_venta,
			p.stock,
			p.stock_min,
			p.url_foto,
			p.estado,
			p.deleted_at,
			p.precio_compra,
			p.presentacion,
			p.unidades_presentacion
		FROM listar_productos_info($1) p
		LEFT JOIN producto_categoria pc ON pc.producto_id = p.id
	`

	// Agregar WHERE si hay filtros
	if len(filters) > 0 {
		query += " WHERE " + strings.Join(filters, " AND ")
	}
	// Ordenar por nombre_comercial
	query += " ORDER BY p.id, p.nombre_comercial"

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		log.Println(err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	var list = make([]domain.ProductoInfo, 0)
	for rows.Next() {
		var item domain.ProductoInfo
		err := rows.Scan(
			&item.Id,
			&item.NombreComercial,
			&item.FormaFarmaceutica,
			&item.Laboratorio,
			&item.PrecioVenta,
			&item.Stock,
			&item.StockMin,
			&item.UrlFoto,
			&item.Estado,
			&item.DeletedAt,
			&item.PrecioCompra,
			&item.Presentacion,
			&item.UnidadesPresentacion,
		)
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

func NewProductoRepository(pool *pgxpool.Pool) *ProductoRepository {
	return &ProductoRepository{pool: pool}
}

var _ port.ProductoRepository = (*ProductoRepository)(nil)
