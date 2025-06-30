package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type CategoriaRepository struct {
	pool *pgxpool.Pool
}

func (c CategoriaRepository) ListarCategoriasDisponibles(ctx context.Context) (*[]domain.Categoria, error) {

	query := `SELECT c.id,c.nombre,c.estado,c.created_at,c.deleted_at FROM categoria c WHERE c.estado = 'Activo' ORDER BY c.nombre`
	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()
	var categorias = make([]domain.Categoria, 0)
	for rows.Next() {
		var categoria domain.Categoria
		err := rows.Scan(&categoria.Id, &categoria.Nombre, &categoria.Estado, &categoria.CreatedAt, &categoria.DeletedAt)
		if err != nil {
			log.Println(err.Error())
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		categorias = append(categorias, categoria)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &categorias, nil
}

func (c CategoriaRepository) HabilitarCategoria(ctx context.Context, categoriaId *int) error {
	query := `UPDATE categoria c SET deleted_at = NULL,estado='Activo' WHERE id=$1;`
	//Inicializar transacción
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	//Ejecutar transacción
	_, err = tx.Exec(ctx, query, *categoriaId)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	return nil
}

func (c CategoriaRepository) DeshabilitarCategoria(ctx context.Context, categoriaId *int) error {
	query := `UPDATE categoria c SET deleted_at = CURRENT_TIMESTAMP,estado='Inactivo' WHERE id=$1;`
	//Inicializar transacción
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	//Ejecutar transacción
	_, err = tx.Exec(ctx, query, *categoriaId)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	return nil
}

func (c CategoriaRepository) ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error) {
	var categoria domain.Categoria
	query := `SELECT c.id,c.nombre,c.created_at,c.deleted_at FROM categoria c WHERE c.id=$1`
	err := c.pool.QueryRow(ctx, query, *categoriaId).Scan(&categoria.Id, &categoria.Nombre, &categoria.CreatedAt, &categoria.DeletedAt)
	if err != nil {
		// Si no hay registros
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Categoría no encontrada")
		}
		// Error en la consulta a la Base de datos
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	return &categoria, nil
}

func (c CategoriaRepository) ListarCategorias(ctx context.Context) (*[]domain.Categoria, error) {
	var categorias []domain.Categoria
	query := `SELECT c.id,c.nombre,c.estado,c.created_at,c.deleted_at FROM categoria c ORDER BY c.nombre`
	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()

	for rows.Next() {
		var categoria domain.Categoria
		err := rows.Scan(&categoria.Id, &categoria.Nombre, &categoria.Estado, &categoria.CreatedAt, &categoria.DeletedAt)
		if err != nil {
			log.Println(err.Error())
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		categorias = append(categorias, categoria)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	if len(categorias) == 0 {
		return &[]domain.Categoria{}, nil
	}
	return &categorias, nil
}

func (c CategoriaRepository) ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error {
	query := `UPDATE categoria SET nombre = $1 WHERE id = $2`

	// Iniciar transacción
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	_, err = tx.Exec(ctx, query, categoriaRequest.Nombre, *categoriaId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return datatype.NewConflictError("Ya existe una categoría con el mismo nombre")
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmar transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

// RegistrarCategoria registra o restaura una categoría, dependiendo si ya existe
func (c CategoriaRepository) RegistrarCategoria(ctx context.Context, categoriaRequest *domain.CategoriaRequest) error {
	// Primero, verificamos si la categoría ya existe en la base de datos
	queryCheck := `SELECT 1 FROM categoria c WHERE c.nombre = $1`

	// Ejecutar la consulta
	res, err := c.pool.Exec(ctx, queryCheck, categoriaRequest.Nombre)
	if err != nil {
		// Si hay error en la consulta, retornamos un error de servicio no disponible
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	// Verificar si la consulta encontró alguna fila
	rowsAffected := res.RowsAffected()
	if rowsAffected > 0 {
		return datatype.NewConflictError("Ya existe la categoría")
	}

	// Si no existe la categoría, la creamos
	queryInsert := `
	INSERT INTO categoria (nombre, deleted_at)
	VALUES ($1, NULL);
	`

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	// Insertamos la nueva categoría
	_, err = tx.Exec(ctx, queryInsert, categoriaRequest.Nombre)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	// Confirmamos la transacción
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func NewCategoriaRepository(pool *pgxpool.Pool) *CategoriaRepository {
	return &CategoriaRepository{pool: pool}
}

var _ port.CategoriaRepository = (*CategoriaRepository)(nil)
