package domain

import (
	"time"
)

// Usuario representa la información básica de un usuario en el sistema.
// Incluye su ID, nombre de usuario, contraseña, la fecha de eliminación (si está eliminada).
type Usuario struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	PersonaId uint       `json:"-"`
	Persona   Persona    `json:"-"`
	Roles     []Rol      `json:"-"`
}

// LoginRequest se usa para las peticiones de autenticación de los usuarios.
// Contiene el nombre de usuario y la contraseña necesarios para el login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UsuarioInfo se usa para mostrar la información detallada de un usuario.
// Incluye su id, nombre de usuario, contraseña (opcional), y su información personal.
type UsuarioInfo struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Persona   Persona    `json:"persona"`
}

// UsuarioDetalle se usa para mostrar la información detallada de un usuario.
// Incluye su id, nombre de usuario, contraseña (opcional), su información personal y los roles asignados.
type UsuarioDetalle struct {
	Id        uint       `json:"id"`
	Username  string     `json:"username"`
	Password  string     `json:"password,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
	Persona   Persona    `json:"persona"`
	Roles     []RolInfo  `json:"roles"`
}

// UsuarioRequest se usa para las peticiones de creación o modificación de un usuario.
// Contiene el nombre de usuario, información personal y los roles que se asignarán.
type UsuarioRequest struct {
	Id        uint           `json:"id,omitzero"`
	Username  string         `json:"username"`
	Persona   PersonaRequest `json:"persona"`
	Roles     []uint8        `json:"roles"`
	DeletedAt *time.Time     `json:"deletedAt"`
}
