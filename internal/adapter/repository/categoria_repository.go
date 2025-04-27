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
	"time"
)

type CategoriaRepository struct {
	db *database.DB
}

func (c CategoriaRepository) ObtenerCategoriaById(ctx context.Context, categoriaId *int) (*domain.Categoria, error) {
	var categoria domain.Categoria
	query := `SELECT c.id,c.nombre,c.created_at,c.deleted_at FROM negocio.categoria c WHERE c.id=$1`
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
	query := `SELECT c.id,c.nombre,c.created_at,c.deleted_at FROM negocio.categoria c ORDER BY c.nombre`
	rows, err := c.db.Pool.Query(ctx, query)
	if err != nil {
		log.Println(err.Error())
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	for rows.Next() {
		var categoria domain.Categoria
		err := rows.Scan(&categoria.Id, &categoria.Nombre, &categoria.CreatedAt, &categoria.DeletedAt)
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

func (c CategoriaRepository) ModificarEstadoCategoria(ctx context.Context, categoriaId *int) error {
	query := `
	UPDATE negocio.categoria c 
	SET deleted_at = CASE 
	WHEN deleted_at IS NOT NULL THEN NULL
	ELSE CURRENT_TIMESTAMP
	END
	WHERE id=$1;
	`
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

func (c CategoriaRepository) ModificarCategoria(ctx context.Context, categoriaId *int, categoriaRequest *domain.CategoriaRequest) error {
	query := `
		UPDATE negocio.categoria 
		SET nombre = $1, deleted_at = $2 
		WHERE id = $3
	`

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

	_, err = tx.Exec(ctx, query, categoriaRequest.Nombre, categoriaRequest.DeletedAt, *categoriaId)
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
	queryCheck := `SELECT deleted_at FROM negocio.categoria WHERE nombre = $1`
	var deletedAt *time.Time
	err := c.db.Pool.QueryRow(ctx, queryCheck, categoriaRequest.Nombre).Scan(&deletedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// Error al ejecutar la consulta
		return datatype.NewInternalServerError()
	}

	// Si la categoría existe y no está eliminada
	if err == nil && deletedAt == nil {
		return &datatype.ErrorResponse{
			Code:    http.StatusConflict,
			Message: "Ya existe la categoría",
		}
	}

	// Si la categoría existe, pero está eliminada, la restauramos
	if err == nil && deletedAt != nil {
		queryRestaurar := `
		UPDATE negocio.categoria
		SET deleted_at = NULL
		WHERE nombre = $1;
		`
		tx, err := c.db.Pool.Begin(ctx)
		if err != nil {
			return datatype.NewInternalServerError()
		}
		defer func(tx pgx.Tx, ctx context.Context) {
			_ = tx.Rollback(ctx)

		}(tx, ctx)

		// Restauramos la categoría
		_, err = tx.Exec(ctx, queryRestaurar, categoriaRequest.Nombre)
		if err != nil {
			return datatype.NewInternalServerError()
		}

		// Confirmamos la transacción
		if err := tx.Commit(ctx); err != nil {
			return datatype.NewInternalServerError()
		}

		return nil
	}

	// Si no existe la categoría, la creamos
	queryInsert := `
	INSERT INTO negocio.categoria (nombre, deleted_at)
	VALUES ($1, NULL);
	`

	tx, err := c.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerError()
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
