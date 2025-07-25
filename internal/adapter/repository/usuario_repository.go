package repository

import (
	"context"
	"database/sql"
	"errors"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	"github.com/sethvargo/go-password/password"
	"log"
)

type UsuarioRepository struct {
	pool *pgxpool.Pool
}

func (u UsuarioRepository) RestablecerPassword(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error) {
	passwordGenerated, err := password.Generate(8, 3, 0, false, false)
	if err != nil {
		log.Println("Error al generar contraseña: " + err.Error())
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	hashPassword, err := util.Hash.HashearPassword(passwordGenerated)
	if err != nil {
		return nil, datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	query := `UPDATE usuario SET password = $1,updated_at=CURRENT_TIMESTAMP WHERE id = $2 `
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()
	ct, err := tx.Exec(ctx, query, hashPassword, *usuarioId)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	if ct.RowsAffected() == 0 {
		log.Println("Usuario no encontrado con id:", *usuarioId)
		return nil, datatype.NewNotFoundError("Usuario no encontrado")
	}
	if err := tx.Commit(ctx); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	usuario, err := u.ObtenerUsuarioDetalle(ctx, usuarioId)
	if err != nil {
		return nil, err
	}
	committed = true

	usuario.Password = passwordGenerated
	return usuario, nil
}

func (u UsuarioRepository) HabilitarUsuarioById(ctx context.Context, usuarioId *int) error {

	query := `UPDATE usuario u SET deleted_at=NULL, estado='Activo' WHERE u.id = $1`
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	ct, err := tx.Exec(ctx, query, *usuarioId)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	if ct.RowsAffected() == 0 {
		log.Println("Usuario no encontrado con id:", *usuarioId)
		return datatype.NewNotFoundError("Usuario no encontrado")
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (u UsuarioRepository) DeshabilitarUsuarioById(ctx context.Context, usuarioId *int) error {
	query := `UPDATE usuario u SET deleted_at=CURRENT_TIMESTAMP, estado='Inactivo' WHERE u.id = $1`
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	ct, err := tx.Exec(ctx, query, *usuarioId)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	if ct.RowsAffected() == 0 {
		log.Println("Usuario no encontrado con id:", *usuarioId)
		return datatype.NewNotFoundError("Usuario no encontrado")
	}

	if err := tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (u UsuarioRepository) ObtenerUsuarioDetalleByUsername(ctx context.Context, username *string) (*domain.UsuarioDetail, error) {
	queryUsuarioDetalle := `SELECT oud.id, oud.username,oud.estado,oud.created_at,oud.updated_at,oud.deleted_at, oud.persona, oud.roles FROM obtener_usuario_detalle_by_username($1) oud;`
	var usuarioDetalle domain.UsuarioDetail

	err := u.pool.QueryRow(ctx, queryUsuarioDetalle, *username).
		Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Estado, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt, &usuarioDetalle.Persona, &usuarioDetalle.Roles)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, datatype.NewNotFoundError("Usuario no encontrado")
		}
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &usuarioDetalle, nil
}

func (u UsuarioRepository) ListarUsuarios(ctx context.Context) (*[]domain.UsuarioInfo, error) {
	query := `
	SELECT f.id, f.username,f.estado,f.created_at, f.updated_at, f.deleted_at, f.persona FROM view_lista_usuarios f ORDER BY f.id;`

	rows, err := u.pool.Query(ctx, query)
	if err != nil {
		log.Println("Error al listar usuarios", err)
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	defer rows.Close()
	var usuarios = make([]domain.UsuarioInfo, 0)
	for rows.Next() {
		var usuarioDetalle domain.UsuarioInfo
		err := rows.Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Estado, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt, &usuarioDetalle.Persona)
		if err != nil {
			log.Print("Error al obtener usuario", err.Error())
			return nil, datatype.NewInternalServerErrorGeneric()
		}
		usuarios = append(usuarios, usuarioDetalle)
	}

	// Verifica si hubo algún error durante la iteración
	if err := rows.Err(); err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &usuarios, nil
}

func (u UsuarioRepository) ModificarUsuario(ctx context.Context, usuarioId *int, usuarioRequest *domain.UsuarioRequest) error {
	persona := usuarioRequest.Persona
	var pgErr *pgconn.PgError

	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	usuarioDetalle, err := u.ObtenerUsuarioDetalle(ctx, usuarioId)
	if err != nil {
		return err
	}

	query := `UPDATE usuario SET username = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	ct, err := tx.Exec(ctx, query, usuarioRequest.Username, *usuarioId)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return datatype.NewConflictError("El usuario ya está registrado")
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	if ct.RowsAffected() == 0 {
		log.Println("Usuario no encontrado con id:", *usuarioId)
		return datatype.NewNotFoundError("Usuario no encontrado")
	}

	query = `UPDATE persona SET nombres = $1, apellido_paterno = $2, apellido_materno = $3, ci = $4, complemento = $5, genero = $6 WHERE id = $7`
	_, err = tx.Exec(ctx, query,
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
			return datatype.NewConflictError("El CI ya se encuentra registrado")
		}
		return datatype.NewInternalServerErrorGeneric()
	}

	// Eliminar roles actuales
	query = `DELETE FROM usuario_rol WHERE usuario_id = $1`
	_, err = tx.Exec(ctx, query, *usuarioId)
	if err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}

	// Insertar nuevos roles
	query = `INSERT INTO usuario_rol(usuario_id, rol_id) SELECT $1, unnest($2::int[])`
	_, err = tx.Exec(ctx, query, *usuarioId, pq.Array(usuarioRequest.Roles))
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return datatype.NewBadRequestError("Algunos roles no existen")
		}
		return datatype.NewInternalServerError("Error al registrar roles")
	}

	// Confirmar la transacción
	if err = tx.Commit(ctx); err != nil {
		return datatype.NewInternalServerErrorGeneric()
	}
	committed = true
	return nil
}

func (u UsuarioRepository) ObtenerUsuarioDetalle(ctx context.Context, usuarioId *int) (*domain.UsuarioDetail, error) {
	queryUsuarioDetalle := `SELECT oud.id, oud.username,oud.estado,oud.created_at,oud.updated_at,oud.deleted_at, oud.persona, oud.roles FROM obtener_usuario_detalle_by_id($1) oud;`
	var usuarioDetalle domain.UsuarioDetail

	err := u.pool.QueryRow(ctx, queryUsuarioDetalle, *usuarioId).
		Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Estado, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt, &usuarioDetalle.Persona, &usuarioDetalle.Roles)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, datatype.NewNotFoundError("No se encontró un usuario con el id proporcionado.")
		}

		return nil, datatype.NewInternalServerErrorGeneric()
	}

	return &usuarioDetalle, nil
}

func (u UsuarioRepository) RegistrarUsuario(ctx context.Context, usuarioRequest *domain.UsuarioRequest) (*domain.UsuarioDetail, error) {
	passwordGenerated, err := password.Generate(8, 3, 0, false, false)
	if err != nil {
		return nil, datatype.NewStatusServiceUnavailableErrorGeneric()
	}
	hashPassword, err := util.Hash.HashearPassword(passwordGenerated)
	if err != nil {
		return nil, err
	}
	// Comienza la transacción
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return nil, datatype.NewStatusServiceUnavailableErrorGeneric()
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Insertar la persona en la tabla `persona`
	query := `INSERT INTO persona(ci, complemento, nombres, apellido_paterno, apellido_materno, genero) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	var personaID int
	err = tx.QueryRow(ctx, query, usuarioRequest.Persona.Ci, usuarioRequest.Persona.Complemento, usuarioRequest.Persona.Nombres, usuarioRequest.Persona.ApellidoPaterno, usuarioRequest.Persona.ApellidoMaterno, usuarioRequest.Persona.Genero).Scan(&personaID)
	if err != nil {
		var pgErr *pgconn.PgError
		// Verifica si el error es una violación de restricción única (código de error 23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, datatype.NewConflictError("La persona ya está registrada")
		}

		return nil, datatype.NewInternalServerError("Error al registrar persona")
	}

	// Insertar el usuario en la tabla `usuario`, relacionado con la persona
	query = `INSERT INTO usuario(username, password,persona_id) VALUES ($1, $2, $3) RETURNING id, username`

	var usuarioId uint
	var usuarioEmail string
	err = tx.QueryRow(ctx, query, usuarioRequest.Username, hashPassword, personaID).Scan(&usuarioId, &usuarioEmail)
	if err != nil {
		var pgErr *pgconn.PgError
		// Verifica si el error es una violación de restricción única (código de error 23505)
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, datatype.NewConflictError("El nombre de usuario ya está registrado")
		}
		// Maneja otros tipos de errores
		return nil, datatype.NewInternalServerError("Error al registrar usuario")
	}

	// Insertar el usuario en la tabla `usuario`, relacionado con la persona
	query = `INSERT INTO usuario_rol(usuario_id, rol_id) VALUES ($1, unnest($2::int[]));`
	_, err = tx.Exec(ctx, query, usuarioId, pq.Array(usuarioRequest.Roles))
	if err != nil {
		// Manejo de error de roles
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return nil, datatype.NewBadRequestError("Algunos roles no existen")
		}
		// Otro tipo de error imprevisto
		return nil, datatype.NewInternalServerError("Error al registrar roles")
	}

	query = `SELECT oud.id, oud.username, oud.persona, oud.roles,oud.created_at,oud.updated_at,oud.deleted_at FROM obtener_usuario_detalle_by_id($1) oud;`
	var usuarioDetalle domain.UsuarioDetail
	err = tx.QueryRow(ctx, query, usuarioId).Scan(&usuarioDetalle.Id, &usuarioDetalle.Username, &usuarioDetalle.Persona, &usuarioDetalle.Roles, &usuarioDetalle.CreatedAt, &usuarioDetalle.UpdatedAt, &usuarioDetalle.DeletedAt)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}
	// Confirmar transacción
	err = tx.Commit(ctx)
	if err != nil {
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	committed = true

	usuarioDetalle.Password = passwordGenerated
	return &usuarioDetalle, nil
}

func (u UsuarioRepository) ObtenerUsuario(ctx context.Context, username *string) (*domain.Usuario, error) {
	query := `SELECT u.id, u.username, u.password, u.deleted_at FROM usuario u WHERE u.username = $1 LIMIT 1`

	var usuario domain.Usuario
	err := u.pool.QueryRow(ctx, query, *username).Scan(&usuario.Id, &usuario.Username, &usuario.Password, &usuario.DeletedAt)
	if err != nil {
		// Si no hay registros
		if errors.Is(err, sql.ErrNoRows) {
			return nil, datatype.NewStatusUnauthorizedError("Usuario o contraseña incorrecta")
		}
		// Error en la consulta a la Base de datos
		return nil, datatype.NewInternalServerErrorGeneric()
	}

	// Si el usuario está eliminado
	if usuario.DeletedAt != nil {
		// Usuario inactivo
		return nil, datatype.NewStatusUnauthorizedError("Usuario o contraseña incorrecta")
	}

	return &usuario, nil
}

func NewUsuarioRepository(pool *pgxpool.Pool) *UsuarioRepository {
	return &UsuarioRepository{pool: pool}
}

var _ port.UsuarioRepository = (*UsuarioRepository)(nil)
