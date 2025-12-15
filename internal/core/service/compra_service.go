package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
)

type CompraService struct {
	compraRepository port.CompraRepository
}

func (c CompraService) ObtenerListaCompras(ctx context.Context, filtros map[string]string) (*[]domain.CompraInfo, error) {
	return c.compraRepository.ObtenerListaCompras(ctx, filtros)
}

func (c CompraService) ObtenerCompraById(ctx context.Context, id *int) (*domain.CompraDetail, error) {
	return c.compraRepository.ObtenerCompraById(ctx, id)
}

func (c CompraService) RegistrarOrdenCompra(ctx context.Context, request *domain.CompraRequest) (*uint, error) {
	val := ctx.Value(util.ContextUserIdKey)
	userIdFloat, ok := val.(int)
	if !ok {
		return nil, datatype.NewBadRequestError("ID de usuario inválido o no encontrado en el contexto")
	}

	request.UsuarioId = uint(userIdFloat)

	return c.compraRepository.RegistrarOrdenCompra(ctx, request)
}

func (c CompraService) ModificarOrdenCompra(ctx context.Context, id *int, request *domain.CompraRequest) error {
	val := ctx.Value(util.ContextUserIdKey)
	userIdFloat, ok := val.(int)
	if !ok {
		return datatype.NewBadRequestError("ID de usuario inválido o no encontrado en el contexto")
	}

	request.UsuarioId = uint(userIdFloat)
	return c.compraRepository.ModificarOrdenCompra(ctx, id, request)
}

func (c CompraService) AnularOrdenCompra(ctx context.Context, id *int) error {
	return c.compraRepository.AnularOrdenCompra(ctx, id)
}

func (c CompraService) RegistrarCompra(ctx context.Context, id *int) error {
	return c.compraRepository.RegistrarCompra(ctx, id)
}

func NewCompraService(compraRepository port.CompraRepository) *CompraService {
	return &CompraService{compraRepository: compraRepository}
}

var _ port.CompraService = (*CompraService)(nil)
