package domain

import "time"

type MovimientoKardex struct {
	IdFila           int       `json:"idFila"`
	FechaMovimiento  time.Time `json:"fechaMovimiento"`
	TipoMovimiento   string    `json:"tipoMovimiento"`
	Documento        string    `json:"documento"`
	CodigoLote       string    `json:"codigoLote"`
	FechaVencimiento time.Time `json:"fechaVencimiento"`
	Usuario          string    `json:"usuario"`
	CantidadEntrada  int       `json:"cantidadEntrada"`
	CantidadSalida   int       `json:"cantidadSalida"`
	CostoUnitario    float64   `json:"costoUnitario"`
	TotalMoneda      float64   `json:"totalMoneda"`
}
