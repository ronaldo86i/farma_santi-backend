package domain

import "time"

type BackupInfo struct {
	Name string    `json:"name"`
	Size int64     `json:"size"` // Tamaño en bytes
	Date time.Time `json:"date"` // Fecha del archivo
}
