package domain

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ClienteInfo struct {
	Id          uint        `json:"id"`
	NitCi       *uint       `json:"nitCi"`
	Complemento pgtype.Text `json:"complemento"`
	Tipo        string      `json:"tipo"`
	Estado      string      `json:"estado"`
	RazonSocial string      `json:"razonSocial"`
	CreatedAt   time.Time   `json:"createdAt"`
}

type ClienteDetail struct {
	Id          uint       `json:"id"`
	NitCi       *uint      `json:"nitCi"`
	Complemento *string    `json:"complemento"`
	Tipo        string     `json:"tipo"`
	RazonSocial string     `json:"razonSocial"`
	Email       string     `json:"email"`
	Telefono    *uint      `json:"telefono"`
	Estado      string     `json:"estado"`
	CreatedAt   time.Time  `json:"createdAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type ClienteRequest struct {
	NitCi       *uint   `json:"nitCi"`
	Complemento *string `json:"complemento"`
	Tipo        string  `json:"tipo"`
	RazonSocial string  `json:"razonSocial"`
	Email       string  `json:"email"`
	Telefono    *uint   `json:"telefono"`
}

type ClienteSimple struct {
	Id          uint        `json:"id"`
	NitCi       *uint       `json:"nitCi"`
	Tipo        string      `json:"tipo"`
	Complemento pgtype.Text `json:"complemento"`
	RazonSocial string      `json:"razonSocial"`
	Email       string      `json:"email,omitempty"`
}

type ClienteId struct {
	Id int `json:"id"`
}
