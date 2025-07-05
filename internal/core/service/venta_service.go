package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
)

type VentaService struct {
	ventaRepository port.VentaRepository
}

func (v VentaService) ObtenerListaVentas(ctx context.Context) (*[]domain.VentaInfo, error) {
	return v.ventaRepository.ObtenerListaVentas(ctx)
}

func (v VentaService) RegistraVenta(ctx context.Context, request *domain.VentaRequest) (*int64, error) {
	val := ctx.Value(util.ContextUserIdKey)
	userIdFloat, ok := val.(int)
	if !ok {
		return nil, datatype.NewBadRequestError("ID de usuario inválido o no encontrado en el contexto")
	}

	request.UsuarioId = uint(userIdFloat)
	return v.ventaRepository.RegistraVenta(ctx, request)
}

func (v VentaService) ObtenerVentaById(ctx context.Context, id *int) (*domain.VentaDetail, error) {
	return v.ventaRepository.ObtenerVentaById(ctx, id)
}

func (v VentaService) AnularVentaById(ctx context.Context, id *int) error {
	return v.ventaRepository.AnularVentaById(ctx, id)
}

func (v VentaService) FacturarVentaById(ctx context.Context, id *int) error {
	return v.ventaRepository.FacturarVentaById(ctx, id)
}

func NewVentaService(ventaRepository port.VentaRepository) *VentaService {
	return &VentaService{ventaRepository: ventaRepository}
}

var _ port.VentaService = (*VentaService)(nil)
