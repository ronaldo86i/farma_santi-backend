package domain

type ProductoStat struct {
	Id              string   `json:"id"`
	NombreComercial string   `json:"nombreComercial"`
	Fotos           []string `json:"fotos"`
	Cantidad        int      `json:"cantidad"`
}

type VentaDiaria struct {
	Fecha string  `json:"fecha"`
	Total float64 `json:"total"`
}

type DashboardStats struct {
	TotalVentas     float64       `json:"totalVentas"`
	CantidadVentas  int           `json:"cantidadVentas"`
	TotalCompras    float64       `json:"totalCompras"`
	CantidadCompras int           `json:"cantidadCompras"`
	VentasDiarias   []VentaDiaria `json:"ventasDiarias"` // Para gráficas de tendencia
}
