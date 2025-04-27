package domain

import "time"

type Categoria struct {
	Id        int        `json:"id"`
	Nombre    string     `json:"nombre"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type CategoriaRequest struct {
	Nombre    string     `json:"nombre"`
	DeletedAt *time.Time `json:"deletedAt"`
}
