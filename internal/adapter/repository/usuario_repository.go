package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/adapter/database"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"log"
	"net/http"
)

type UsuarioRepository struct {
	db *database.DB
}

func (u UsuarioRepository) RestablecerPassword(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error) {
	passwordGenerated, err := password.Generate(8, 3, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}
	hashPassword, err := util.Hash.HashearPassword(passwordGenerated)
	if err != nil {
		return nil, datatype.NewStatusServiceUnavailableError()
	}
	query := `UPDATE usuario SET password = $1,updated_at=CURRENT_TIMESTAMP WHERE id = $2 `
	tx, err := u.db.Pool.Begin(ctx)
	if err != nil {
		return nil, datatype.NewStatusServiceUnavailableError()
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx) // rollback silencioso
		}
	}()
	_, err = tx.Exec(ctx, query, hashPassword, *usuarioId)
	if err != nil {
		return nil, datatype.NewInternalServerError()
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, datatype.NewInternalServerError()
	}
	usuario, err := u.ObtenerUsuarioDetalle(ctx, usuarioId)
	if err != nil {
		return nil, err
	}
	usuario.Password = passwordGenerated
	return usuario, nil
}

func (u UsuarioRepository) HabilitarUsuarioById(ctx context.Context, usuarioId *int) error {

	query := `UPDATE usuario u SET deleted_at=NULL, estado='Activo' WHERE u.id = $1`
	tx, err := u.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx) // rollback silencioso
		}
	}()

	_, err = tx.Exec(ctx, query, *usuarioId)
	if err != nil {
		return datatype.NewInternalServerError()
	}
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}
	return nil
}

func (u UsuarioRepository) DeshabilitarUsuarioById(ctx context.Context, usuarioId *int) error {
	query := `UPDATE usuario u SET deleted_at=CURRENT_TIMESTAMP, estado='Inactivo' WHERE u.id = $1`
	tx, err := u.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableError()
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()
	_, err = tx.Exec(ctx, query, *usuarioId)
	if err != nil {
		return datatype.NewInternalServerError()
	}
	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerError()
	}
	return nil
}

func (u UsuarioRepository) ObtenerUsuarioDetalleByUsername(ctx context.Context, username *string) (*domain.UsuarioDetail, error) {
	queryUsuarioDetalle := `SELECT oud.id, oud.username,oud.estado,oud.created_at,oud.updated_at,oud.deleted_at, oud.persona, oud.roles FROM obtener_usuario_detalle_by_username($1) oud;`
	var usuarioDetalle domain.UsuarioDetail

	err := u.db.Pool.QueryRow(ctx, queryUsuarioDetalle, *username).
		Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Estado, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt, &usuarioDetalle.Persona, &usuarioDetalle.Roles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "No se encontró un usuario con el id proporcionado.",
			}
		}
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al obtener los detalles del usuario desde la base de datos: " + err.Error(),
		}
	}

	return &usuarioDetalle, nil
}

func (u UsuarioRepository) ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error) {
	query := `
	SELECT f.id, f.username,f.estado,f.created_at, f.updated_at, f.deleted_at, f.persona
	FROM obtener_lista_usuarios() f 
	ORDER BY f.id;
	`

	var usuarios []domain.UsuarioInfo

	rows, err := u.db.Pool.Query(ctx, query)
	if err != nil {
		log.Println("Error al listar usuarios", err)
		return nil, datatype.NewInternalServerError()
	}
	defer rows.Close()

	for rows.Next() {
		var usuarioDetalle domain.UsuarioInfo
		err := rows.Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Estado, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt, &usuarioDetalle.Persona)
		if err != nil {
			log.Print("Error al obtener usuario", err.Error())
			return nil, datatype.NewInternalServerError()
		}
		usuarios = append(usuarios, usuarioDetalle)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerError()
	}
	if len(usuarios) == 0 {
		return &[]domain.UsuarioInfo{}, nil
	}
	return &usuarios, nil
}

func (u UsuarioRepository) ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error {
	persona := usuarioRequest.Persona
	var pgErr *pgconn.PgError

	tx, err := u.db.Pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerError()
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	usuarioDetalle, err := u.ObtenerUsuarioDetalle(ctx, usuarioId)
	if err != nil {
		return err
	}

	updateUsuarioQuery := `UPDATE usuario SET username = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err = tx.Exec(ctx, updateUsuarioQuery, usuarioRequest.Username, *usuarioId)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			msg := fmt.Sprintf("El %s ya está registrado", "usuario")
			return &datatype.ErrorResponse{
				Code:    http.StatusConflict,
				Message: msg,
			}
		}
		return &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al actualizar usuario: " + err.Error(),
		}
	}

	updatePersonaQuery := `UPDATE persona SET nombres = $1, apellido_paterno = $2, apellido_materno = $3, ci = $4, complemento = $5, genero = $6 WHERE id = $7`
	_, err = tx.Exec(ctx, updatePersonaQuery,
		persona.Nombres,
		persona.ApellidoPaterno,
		persona.ApellidoMaterno,
		persona.Ci,
		persona.Complemento,
		persona.Genero,
		usuarioDetalle.Persona.Id,
	)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			msg := fmt.Sprintf("El %s ya está registrado", "ci")
			return &datatype.ErrorResponse{
				Code:    http.StatusConflict,
				Message: msg,
			}
		}
		return &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al actualizar persona: " + err.Error(),
		}
	}

	// Eliminar roles actuales
	deleteRolesQuery := `DELETE FROM usuario_rol WHERE usuario_id = $1`
	_, err = tx.Exec(ctx, deleteRolesQuery, *usuarioId)
	if err != nil {
		return &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al eliminar roles anteriores: " + err.Error(),
		}
	}

	// Insertar nuevos roles
	insertRolesQuery := `
	INSERT INTO usuario_rol(usuario_id, rol_id)
	SELECT $1, unnest($2::int[])`
	_, err = tx.Exec(ctx, insertRolesQuery, *usuarioId, pq.Array(usuarioRequest.Roles))
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return &datatype.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Algunos roles no existen",
			}
		}
		return &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al insertar los roles: " + err.Error(),
		}
	}

	// Confirmar la transacción
	if err = tx.Commit(ctx); err != nil {
		return &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al confirmar la transacción: " + err.Error(),
		}
	}

	return nil
}

func (u UsuarioRepository) ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error) {
	queryUsuarioDetalle := `SELECT oud.id, oud.username,oud.estado,oud.created_at,oud.updated_at,oud.deleted_at, oud.persona, oud.roles FROM obtener_usuario_detalle_by_id($1) oud;`
	var usuarioDetalle domain.UsuarioDetail

	err := u.db.Pool.QueryRow(ctx, queryUsuarioDetalle, *usuarioId).
		Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Estado, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt, &usuarioDetalle.Persona, &usuarioDetalle.Roles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "No se encontró un usuario con el id proporcionado.",
			}
		}
		return nil, datatype.NewInternalServerError()
	}

	return &usuarioDetalle, nil
}

func (u UsuarioRepository) RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetail, error) {
	passwordGenerated, err := password.Generate(8, 3, 0, false, false)
	if err != nil {
		log.Fatal(err)
	}
	hashPassword, err := util.Hash.HashearPassword(passwordGenerated)
	if err != nil {
		return nil, err
	}
	// Comienza la transacción
	tx, err := u.db.Pool.Begin(ctx)
	if err != nil {
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al iniciar la transacción: " + err.Error(),
		}
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx) // rollback silencioso
		}
	}()

	// Insertar la persona en la tabla `persona`
	queryPersona := `INSERT INTO persona(ci, complemento, nombres, apellido_paterno, apellido_materno, genero) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var personaID int
	err = tx.QueryRow(ctx, queryPersona, usuarioRequest.Persona.Ci, usuarioRequest.Persona.Complemento, usuarioRequest.Persona.Nombres, usuarioRequest.Persona.ApellidoPaterno, usuarioRequest.Persona.ApellidoMaterno, usuarioRequest.Persona.Genero).Scan(&personaID)
	if err != nil {
		var pgErr *pgconn.PgError
		// Verifica si el error es una violación de restricción única (código de error 23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusConflict, // Código HTTP 409: Conflicto
				Message: "La persona ya está registrada",
			}
		}

		return nil, &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al insertar la persona: " + err.Error(),
		}
	}

	// Insertar el usuario en la tabla `usuario`, relacionado con la persona
	queryUsuario := `INSERT INTO usuario(username, password,persona_id) VALUES ($1, $2, $3) RETURNING id, username`

	var usuarioId uint
	var usuarioEmail string
	err = tx.QueryRow(ctx, queryUsuario, usuarioRequest.Username, hashPassword, personaID).Scan(&usuarioId, &usuarioEmail)
	if err != nil {
		var pgErr *pgconn.PgError
		// Verifica si el error es una violación de restricción única (código de error 23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusConflict, // Código HTTP 409: Conflicto
				Message: "El nombre de usuario ya está registrado",
			}
		}
		// Maneja otros tipos de errores
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al insertar el usuario: " + err.Error(),
		}
	}

	// Insertar el usuario en la tabla `usuario`, relacionado con la persona
	queryRol := `INSERT INTO usuario_rol(usuario_id, rol_id) VALUES ($1, unnest($2::int[]));`
	_, err = tx.Exec(ctx, queryRol, usuarioId, pq.Array(usuarioRequest.Roles))
	if err != nil {
		// Manejo de error de roles
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Algunos roles no existen",
			}
		}
		// Otro tipo de error imprevisto
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al insertar los roles: " + err.Error(),
		}
	}

	queryUsuarioDetalle := `SELECT oud.id, oud.username, oud.persona, oud.roles,oud.created_at,oud.updated_at,oud.deleted_at FROM obtener_usuario_detalle_by_id($1) oud;`
	var usuarioDetalle domain.UsuarioDetail
	err = tx.QueryRow(ctx, queryUsuarioDetalle, usuarioId).Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Persona, &usuarioDetalle.Roles, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt)
	if err != nil {
		// Tipo de error imprevisto
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al convertir a json: " + err.Error(),
		}
	}
	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Error al confirmar la transacción: " + err.Error(),
		}
	}
	usuarioDetalle.Password = passwordGenerated
	// Devolver los detalles del usuario creado
	return &usuarioDetalle, nil
}

func (u UsuarioRepository) ObtenerUsuario(ctx context.Context, username *string) (*domain.Usuario, error) {
	query := `SELECT u.id, u.username, u.password, u.deleted_at FROM usuario u WHERE u.username = $1 LIMIT 1`

	var usuario domain.Usuario
	err := u.db.Pool.QueryRow(ctx, query, username).Scan(&usuario.Id, &usuario.Username, &usuario.Password, &usuario.DeletedAt)
	if err != nil {
		// Si no hay registros
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &datatype.ErrorResponse{
				Code:    http.StatusUnauthorized,
				Message: "Usuario o contraseña incorrectos",
			}
		}
		// Error en la consulta a la Base de datos
		return nil, datatype.NewInternalServerError()
	}

	// Si el usuario está eliminado
	if usuario.DeletedAt != nil {
		// Usuario inactivo
		return nil, &datatype.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "Usuario o contraseña incorrectos",
		}
	}

	return &usuario, nil
}

func NewUsuarioRepository(db *database.DB) *UsuarioRepository {
	return &UsuarioRepository{db: db}
}

var _ port.UsuarioRepository = (*UsuarioRepository)(nil)
