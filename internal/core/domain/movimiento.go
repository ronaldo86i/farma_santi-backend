package domain

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type MovimientoInfo struct {
	Id      int64         `json:"id"`
	Codigo  pgtype.Text   `json:"codigo"`
	Tipo    string        `json:"tipo"`
	Estado  string        `json:"estado"`
	Fecha   time.Time     `json:"fecha"`
	Usuario UsuarioSimple `json:"usuario"`
	Total   float64       `json:"total"`
}
