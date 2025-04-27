package domain

import (
	"time"
)

type Rol struct {
	Id        uint8      `json:"id"`
	Nombre    string     `json:"nombre"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type RolRequest struct {
	Nombre string `json:"nombre"`
}

type RolRequestUpdate struct {
	Nombre    string     `json:"nombre"`
	DeletedAt *time.Time `json:"deletedAt"`
}
type RolInfo struct {
	Id     uint8  `json:"id"`
	Nombre string `json:"nombre"`
}
