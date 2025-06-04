package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"net/http"
)

type CategoriaRepository struct {
	db *database.DB
}

func (c CategoriaRepository) ListarCategoriasDisponibles(ctx context.Context) (*[]domain.Categoria, error) {
	var categorias []domain.Categoria
	query := `SELECT c.id,c.nombre,c.estado,c.created_at,c.deleted_at FROM categoria c WHERE c.estado = 'Activo' ORDER BY c.nombre`
	rows, err := c.db.Pool.Query(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	for rows.Next() {
		var categoria domain.Categoria
		err := rows.Scan(&categoria.Id, &categoria.Nombre, &categoria.Estado, &categoria.CreatedAt, &categoria.DeletedAt)
		if err != nil {
			log.Println(err.Error())
			return nil, datatype.NewInternalServerError()
		}
		categorias = append(categorias, categoria)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerError()
	}
	if len(categorias) == 0 {
		return &[]domain.Categoria{}, nil
	}
	return &categorias, nil
}

func (c CategoriaRepository) HabilitarCategoria(ctx context.Context, categoriaId *int) error {
	query := `UPDATE categoria c SET deleted_at = NULL,estado='Activo' WHERE id=$1;`
	//Inicializar transacción
	tx, err := c.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}
	//Ejecutar transacción
	_, err = tx.Exec(ctx, query, *categoriaId)
	if err != nil {
		return datatype.NewInternalServerError()
	}

	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}
	return nil
}

func (c CategoriaRepository) DeshabilitarCategoria(ctx context.Context, categoriaId *int) error {
	query := `UPDATE categoria c SET deleted_at = CURRENT_TIMESTAMP,estado='Inactivo' WHERE id=$1;`
	//Inicializar transacción
	tx, err := c.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}
	//Ejecutar transacción
	_, err = tx.Exec(ctx, query, *categoriaId)
	if err != nil {
		return datatype.NewInternalServerError()
	}

	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}
	return nil
}

func (c CategoriaRepository) ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error) {
	var categoria domain.Categoria
	query := `SELECT c.id,c.nombre,c.created_at,c.deleted_at FROM categoria c WHERE c.id=$1`
	err := c.db.Pool.QueryRow(ctx, query, *categoriaId).Scan(&categoria.Id, &categoria.Nombre, &categoria.CreatedAt, &categoria.DeletedAt)
	if err != nil {
		// Si no hay registros
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Categoría no encontrada",
			}
		}
		// Error en la consulta a la Base de datos
		return nil, datatype.NewInternalServerError()
	}
	return &categoria, nil
}

func (c CategoriaRepository) ListarCategorias(ctx context.Context) (*[]domain.Categoria, error) {
	var categorias []domain.Categoria
	query := `SELECT c.id,c.nombre,c.estado,c.created_at,c.deleted_at FROM categoria c ORDER BY c.nombre`
	rows, err := c.db.Pool.Query(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	for rows.Next() {
		var categoria domain.Categoria
		err := rows.Scan(&categoria.Id, &categoria.Nombre, &categoria.Estado, &categoria.CreatedAt, &categoria.DeletedAt)
		if err != nil {
			log.Println(err.Error())
			return nil, datatype.NewInternalServerError()
		}
		categorias = append(categorias, categoria)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerError()
	}
	if len(categorias) == 0 {
		return &[]domain.Categoria{}, nil
	}
	return &categorias, nil
}

func (c CategoriaRepository) ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error {
	query := `UPDATE categoria SET nombre = $1 WHERE id = $2`

	// Iniciar transacción
	tx, err := c.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}
	defer func() {
		// Rollback si no se hizo commit
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	_, err = tx.Exec(ctx, query, categoriaRequest.Nombre, *categoriaId)
	if err != nil {
		log.Println("ERROR:", err.Error())

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return &datatype.ErrorResponse{
				Code:    http.StatusConflict,
				Message: "Ya existe una categoría con ese nombre",
			}
		}
		return datatype.NewInternalServerError()
	}

	// Confirmar transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

// RegistrarCategoria registra o restaura una categoría, dependiendo si ya existe
func (c CategoriaRepository) RegistrarCategoria(ctx context.Context, categoriaRequest *domain.CategoriaRequest) error {
	// Primero, verificamos si la categoría ya existe en la base de datos
	queryCheck := `SELECT 1 FROM categoria c WHERE c.nombre = $1`

	// Ejecutar la consulta
	res, err := c.db.Pool.Exec(ctx, queryCheck, categoriaRequest.Nombre)
	if err != nil {
		// Si hay error en la consulta, retornamos un error de servicio no disponible
		return datatype.NewStatusServiceUnavailableError()
	}

	// Verificar si la consulta encontró alguna fila
	rowsAffected := res.RowsAffected()
	if rowsAffected > 0 {
		return &datatype.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Ya existe la categoría",
		}
	}

	// Si no existe la categoría, la creamos
	queryInsert := `
	INSERT INTO categoria (nombre, deleted_at)
	VALUES ($1, NULL);
	`

	tx, err := c.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	// Insertamos la nueva categoría
	_, err = tx.Exec(ctx, queryInsert, categoriaRequest.Nombre)
	if err != nil {
		return datatype.NewInternalServerError()
	}

	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}

	return nil
}

func NewCategoriaRepository(db *database.DB) *CategoriaRepository {
	return &CategoriaRepository{db: db}
}

var _ port.CategoriaRepository = (*CategoriaRepository)(nil)
