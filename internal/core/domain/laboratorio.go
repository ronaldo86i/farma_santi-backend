package domain

import "time"

type Laboratorio struct {
	Id        int32      `json:"id"`
	Nombre    string     `json:"nombre"`
	Estado    string     `json:"estado"`
	Direccion *string    `json:"direccion"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type LaboratorioRequest struct {
	Nombre    string  `json:"nombre"`
	Direccion *string `json:"direccion"`
}

type LaboratorioInfo struct {
	Id        int32      `json:"id"`
	Nombre    string     `json:"nombre"`
	Estado    string     `json:"estado"`
	Direccion *string    `json:"direccion"`
	CreatedAt time.Time  `json:"createdAt"`
	DeletedAt *time.Time `json:"deletedAt"`
}

type LaboratorioDetail struct {
	Laboratorio
}

type LaboratorioSimple struct {
	Id     int32  `json:"id"`
	Nombre string `json:"nombre"`
}
