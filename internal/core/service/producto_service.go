package service

import (
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"github.com/google/uuid"
	"mime/multipart"
	"strings"
)

type ProductoService struct {
	productoRepository port.ProductoRepository
}

func (p ProductoService) ObtenerProductoById(ctx context.Context, id *uuid.UUID) (*domain.ProductoDetail, error) {
	return p.productoRepository.ObtenerProductoById(ctx, id)
}

func (p ProductoService) HabilitarProducto(ctx context.Context, id *uuid.UUID) error {
	return p.productoRepository.HabilitarProducto(ctx, id)
}

func (p ProductoService) DeshabilitarProducto(ctx context.Context, id *uuid.UUID) error {
	return p.productoRepository.DeshabilitarProducto(ctx, id)
}

func (p ProductoService) RegistrarProducto(ctx context.Context, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error {
	request.NombreComercial = strings.TrimSpace(request.NombreComercial)
	request.NombreComercial = strings.ToUpper(request.NombreComercial)
	if len(request.PrincipiosActivos) == 0 {
		return datatype.NewBadRequestError("Lista de principios activos vacía")
	}
	for _, file := range *filesHeader {
		if !util.File.ValidarTipoArchivo(file.Filename, ".png", ".jpg", ".jpeg") {
			return datatype.NewBadRequestError("Tipo de archivo no válido")
		}
	}
	return p.productoRepository.RegistrarProducto(ctx, request, filesHeader)
}

func (p ProductoService) ModificarProducto(ctx context.Context, id *uuid.UUID, request *domain.ProductRequest, filesHeader *[]*multipart.FileHeader) error {
	request.NombreComercial = strings.TrimSpace(request.NombreComercial)
	request.NombreComercial = strings.ToUpper(request.NombreComercial)
	for _, file := range *filesHeader {
		if !util.File.ValidarTipoArchivo(file.Filename, ".png", ".jpg", ".jpeg") {
			return datatype.NewBadRequestError("Tipo de archivo no válido")
		}
	}
	return p.productoRepository.ModificarProducto(ctx, id, request, filesHeader)
}

func (p ProductoService) ListarFormasFarmaceuticas(ctx context.Context) (*[]domain.FormaFarmaceutica, error) {
	return p.productoRepository.ListarFormasFarmaceuticas(ctx)
}

func (p ProductoService) ListarUnidadesMedida(ctx context.Context) (*[]domain.UnidadMedida, error) {
	return p.productoRepository.ListarUnidadesMedida(ctx)
}

func (p ProductoService) ListarProductos(ctx context.Context) (*[]domain.ProductoInfo, error) {
	return p.productoRepository.ListarProductos(ctx)
}

func NewProductoService(productoRepository port.ProductoRepository) *ProductoService {
	return &ProductoService{productoRepository: productoRepository}
}

var _ port.ProductoService = (*ProductoService)(nil)
