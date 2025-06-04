package domain

import (
	"time"
)

type Rol struct {
	Id        int32      `json:"id"`
	Nombre    string     `json:"nombre"`
	Estado    string     `json:"estado"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type RolRequest struct {
	// Nombre del rol
	Nombre string
}

type RolInfo struct {
	Id        int32      `json:"id"`
	Nombre    string     `json:"nombre"`
	Estado    string     `json:"estado"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type RolDetail struct {
	Rol
}
