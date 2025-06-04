package domain

import "time"

type Categoria struct {
	Id        int32      `json:"id"`
	Nombre    string     `json:"nombre"`
	Estado    string     `json:"estado"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type CategoriaRequest struct {
	Nombre string `json:"nombre"`
}
