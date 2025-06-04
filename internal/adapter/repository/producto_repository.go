package repository

import (
	"context"
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type ProductoRepository struct {
	db *database.DB
}

func (p ProductoRepository) ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error) {
	query := `SELECT um.id,um.nombre,um.abreviatura FROM unidad_medida um ORDER BY um.nombre`
	var list []domain.UnidadMedida

	rows, err := p.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.UnidadMedida
		err := rows.Scan(&item.Id, &item.Nombre, &item.Abreviatura)
		if err != nil {
			return nil, datatype.NewInternalServerError()
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerError()
	}

	if len(list) == 0 {
		return &[]domain.UnidadMedida{}, nil
	}
	return &list, nil
}

func (p ProductoRepository) ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error) {
	query := `SELECT ff.id,ff.nombre FROM forma_farmaceutica ff ORDER BY ff.nombre`
	var list []domain.FormaFarmaceutica

	rows, err := p.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var item domain.FormaFarmaceutica
		err := rows.Scan(&item.Id, &item.Nombre)
		if err != nil {
			return nil, datatype.NewInternalServerError()
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerError()
	}

	if len(list) == 0 {
		return &[]domain.FormaFarmaceutica{}, nil
	}
	return &list, nil
}

func (p ProductoRepository) ListarProductos(ctx context.Context) (*[]domain.ProductInfo, error) {
	fullHostname := ctx.Value("fullHostname").(string)
	fullHostname = fmt.Sprintf("%s%s", fullHostname, "/uploads/productos")

	query := `SELECT p.id, p.nombre_comercial, p.nombre_generico, p.concentracion, p.forma_farmaceutica, p.laboratorio, p.precio_venta, p.stock, p.stock_min, p.url_foto,p.estado FROM listar_productos_info($1) p`
	rows, err := p.db.Pool.Query(ctx, query, fullHostname)
	if err != nil {
		log.Println(err)
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	var list []domain.ProductInfo
	for rows.Next() {
		var item domain.ProductInfo
		err := rows.Scan(&item.Id, &item.NombreComercial, &item.NombreGenerico, &item.Concentracion, &item.FormaFarmaceutica, &item.Laboratorio, &item.PrecioVenta, &item.Stock, &item.StockMin, &item.UrlFoto, &item.Estado)
		if err != nil {
			return nil, datatype.NewInternalServerError()
		}
		list = append(list, item)
	}
	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerError()
	}

	if len(list) == 0 {
		return &[]domain.ProductInfo{}, nil
	}
	return &list, nil
}

func (p ProductoRepository) RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error {
	// Inicio de transacción
	tx, err := p.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}

	query := `INSERT INTO producto(nombre_comercial,nombre_generico,concentracion,forma_farmaceutica_id,precio_compra,precio_venta,estado,stock,stock_min,laboratorio_id,unidad_medida_id) VALUES ($1,$2,$3,$4,0.0,$5,'Activo',0,$6,$7,$8) RETURNING id`

	var id uuid.UUID
	err = tx.QueryRow(ctx, query, request.NombreComercial, request.NombreGenerico, request.Concentracion, request.FormaFarmaceuticaId, request.PrecioVenta, request.StockMin, request.LaboratorioId, request.UnidadMedidaId).Scan(&id)
	if err != nil {
		_ = tx.Rollback(ctx)
		return datatype.NewStatusServiceUnavailableError()
	}

	route := fmt.Sprintf("./public/uploads/productos/%s", id.String())
	defer func() {
		if err != nil {
			_ = util.File.DeleteAllFiles(route)
			_ = tx.Rollback(ctx)
		}
	}()
	err = util.File.MakeDir(route)
	if err != nil {
		log.Println(err)
		return datatype.NewInternalServerError()
	}

	var fotos []string
	for i, fileHeader := range *filesHeader {
		file, err := fileHeader.Open()
		if err != nil {
			return datatype.NewStatusServiceUnavailableError()
		}

		ext := filepath.Ext(fileHeader.Filename)
		nameFile := fmt.Sprintf("%d%s", i+1, ext)

		err = util.File.SaveFile(route, nameFile, file)
		if err != nil {
			_ = util.File.DeleteAllFiles(route)
			return &datatype.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: "Error al guardar fotos",
			}
		}
		_ = file.Close()
		fotos = append(fotos, nameFile)
	}

	query = `INSERT INTO producto_categoria(producto_id, categoria_id) SELECT $1, unnest($2::int[])`
	_, err = tx.Exec(ctx, query, id, pq.Array(request.Categorias))
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}

	query = `UPDATE producto SET fotos = $1 WHERE id = $2`
	_, err = tx.Exec(ctx, query, pq.Array(fotos), id)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}

	// Confirmar la transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}
	return nil
}

func NewProductoRepository(db *database.DB) *ProductoRepository {
	return &ProductoRepository{db: db}
}

var _ port.ProductoRepository = (*ProductoRepository)(nil)
