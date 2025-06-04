package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"mime/multipart"
	"net/http"
)

type ProductoService struct {
	productoRepository port.ProductoRepository
}

func (p ProductoService) ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error) {
	return p.productoRepository.ListarFormasFarmaceuticas(ctx)
}

func (p ProductoService) ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error) {
	return p.productoRepository.ListarUnidadesMedida(ctx)
}

func (p ProductoService) ListarProductos(ctx context.Context) (*[]domain.ProductInfo, error) {
	return p.productoRepository.ListarProductos(ctx)
}

func (p ProductoService) RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error {
	for _, file := range *filesHeader {
		if !util.File.ValidarTipoArchivo(file.Filename, ".png", ".jpg", ".jpeg") {
			return &datatype.ErrorResponse{
				Code:    http.StatusBadRequest,
				Message: "Tipo de archivo no válido",
			}
		}
	}
	return p.productoRepository.RegistrarProducto(ctx, request, filesHeader)
}

func NewProductoService(productoRepository port.ProductoRepository) *ProductoService {
	return &ProductoService{productoRepository: productoRepository}
}

var _ port.ProductoService = (*ProductoService)(nil)
