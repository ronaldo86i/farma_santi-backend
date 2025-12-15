package domain

import "time"

type Laboratorio struct {
	Id            int32      `json:"id"`
	Nombre        string     `json:"nombre"`
	Estado        string     `json:"estado"`
	Direccion     *string    `json:"direccion"`
	Representante *string    `json:"representante"`
	Telefono      *int       `json:"telefono,omitzero"`
	Email         *string    `json:"email,omitempty"`
	Celular       *int       `json:"celular,omitzero"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
}

type LaboratorioRequest struct {
	Nombre        string  `json:"nombre"`
	Direccion     *string `json:"direccion"`
	Representante *string `json:"representante"`
	Telefono      *int    `json:"telefono,omitzero"`
	Email         *string `json:"email,omitempty"`
	Celular       *int    `json:"celular,omitzero"`
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
