package service

import (
	"bytes"
	"context"
	"farma-santi_backend/internal/core/domain"
	"farma-santi_backend/internal/core/domain/datatype"
	"farma-santi_backend/internal/core/port"
	"farma-santi_backend/internal/core/util"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/goccy/go-json"
)

type VentaService struct {
	ventaRepository port.VentaRepository
}

func (v VentaService) ObtenerListaVentas(ctx context.Context, filtros map[string]string) (*[]domain.VentaInfo, error) {
	return v.ventaRepository.ObtenerListaVentas(ctx, filtros)
}

func (v VentaService) RegistraVenta(ctx context.Context, request *domain.VentaRequest) (*int64, error) {
	// Obtener ID de usuario desde el contexto
	val := ctx.Value(util.ContextUserIdKey)
	userIdFloat, ok := val.(int)
	if !ok {
		return nil, datatype.NewBadRequestError("ID de usuario inválido o no encontrado en el contexto")
	}
	request.UsuarioId = uint(userIdFloat)

	// Registrar venta en DB
	ventaId, err := v.ventaRepository.RegistraVenta(ctx, request)
	if err != nil {
		return nil, datatype.NewInternalServerError(err.Error())
	}
	ventaIdInt := int(*ventaId)

	// Obtener venta completa
	venta, err := v.ventaRepository.ObtenerVentaById(ctx, &ventaIdInt)
	if err != nil {
		return nil, err
	}
	if venta.Cliente.NitCi == nil {
		nitCi := uint(1)
		venta.Cliente.NitCi = &nitCi
	}

	// Construir JSON de la factura
	factura := construirFactura(venta)

	body, err := json.Marshal(factura)
	if err != nil {
		return nil, fmt.Errorf("error construyendo JSON de factura: %w", err)
	}

	// Llamada HTTP a SIAT
	url := fmt.Sprintf("%s/api/v1/facturacion/electronica", os.Getenv("URL_FACTURADOR"))
	token := os.Getenv("TOKEN_FACTURA")

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creando request: %v", err)
		return ventaId, nil
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Connection", "close")

	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error al obtener datos de SIAT: %v", err)
		return ventaId, nil // no bloqueamos la venta
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error leyendo respuesta SIAT: %v", err)
		return ventaId, nil
	}

	log.Printf("Código de estado SIAT: %d", resp.StatusCode)
	log.Printf("Datos de respuesta SIAT: %s", string(respBytes))

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error en respuesta SIAT: %s", resp.Status)
		return ventaId, nil
	}

	// Parsear respuesta
	var facturaResponse domain.FacturaCompraVentaResponse
	if err := json.Unmarshal(respBytes, &facturaResponse); err != nil {
		log.Printf("Error parseando respuesta SIAT: %v", err)
		return ventaId, nil
	}

	// Guardar factura en DB en transacción separada
	if err := v.ventaRepository.FacturarVentaById(context.Background(), &ventaIdInt, &facturaResponse); err != nil {
		log.Printf("Error al registrar factura: %v", err)
		return ventaId, nil
	}

	return ventaId, nil
}

var codigosTipoDocumentoIdentidad = map[string]uint64{
	"CI":  1,
	"NIT": 5,
}

// Función helper para construir la factura desde la venta
func construirFactura(venta *domain.VentaDetail) domain.FacturaCompraVenta {
	telefono := "74425055"
	var detalles []domain.Detalle

	for _, d := range venta.Detalles {
		var cantidad float64
		for _, l := range d.Lotes {
			cantidad += float64(l.Cantidad)
		}
		detalles = append(detalles, domain.Detalle{
			ActividadEconomica: "477300",
			CodigoProductoSin:  "622539",
			CodigoProducto:     d.Producto.Id.String(),
			Descripcion:        d.Producto.NombreComercial,
			Cantidad:           cantidad,
			UnidadMedida:       57,
			PrecioUnitario:     d.Precio,
			MontoDescuento:     domain.NilableFloat64{},
			SubTotal:           d.Total,
			NumeroSerie:        domain.NilableString{Value: nil},
			NumeroImei:         domain.NilableString{Value: nil},
		})
	}
	var numeroDocumento = "1"
	if venta.Cliente.NitCi != nil {
		numeroDocumento = fmt.Sprintf("%d", *venta.Cliente.NitCi)
	}
	var codigoTipoDocumentoIdentidad = codigosTipoDocumentoIdentidad[venta.Cliente.Tipo]

	return domain.FacturaCompraVenta{
		Cabecera: domain.Cabecera{
			Municipio:                    "Tarija",
			Telefono:                     &telefono,
			CodigoSucursal:               0,
			Direccion:                    "ESQUINA AVENIDA LA PAZ CASA DE CUATRO PISOS CON FACHADA DE COLOR AMARILLO CON PERSIANAS DE COLOR CREMA LOS MEMBRILLOS Nro.: S/N",
			CodigoPuntoVenta:             domain.NilableUint64{Value: nil},
			NombreRazonSocial:            domain.NilableString{Value: &venta.Cliente.RazonSocial},
			CodigoTipoDocumentoIdentidad: codigoTipoDocumentoIdentidad,
			NumeroDocumento:              numeroDocumento,
			Complemento:                  domain.NilableString{Value: &venta.Cliente.Complemento.String},
			CodigoCliente:                "1",
			EmailCliente:                 "",
			CodigoMetodoPago:             1,
			NumeroTarjeta:                domain.NilableUint64{Value: nil},
			MontoTotal:                   venta.Total - venta.Descuento,
			CodigoMoneda:                 1,
			TipoCambio:                   1,
			MontoTotalMoneda:             venta.Total - venta.Descuento,
			MontoGiftCard:                domain.NilableFloat64{Value: nil},
			DescuentoAdicional:           domain.NilableFloat64{Value: &venta.Descuento},
			CodigoExcepcion:              domain.NilableUint64{Value: nil},
			Cafc:                         domain.NilableString{Value: nil},
			Leyenda:                      "Ley N° 453: Tienes derecho a recibir información sobre las características y contenidos de los productos que consumes.",
			CodigoDocumentoSector:        1,
			TipoFacturaDocumento:         1,
		},
		Detalle: detalles,
	}
}

func (v VentaService) ObtenerVentaById(ctx context.Context, id *int) (*domain.VentaDetail, error) {
	return v.ventaRepository.ObtenerVentaById(ctx, id)
}

func (v VentaService) AnularVentaById(ctx context.Context, id *int) error {
	factura, err := v.ventaRepository.ObtenerFacturaByVentaId(ctx, id)
	if err != nil {
		log.Printf("Error obteniendo factura: %v", err)
	}

	if factura != nil {
		anulacionFactura := domain.AnularFacturaRequest{
			NumeroFactura:    factura.NumeroFactura,
			CodigoSucursal:   factura.CodigoSucursal,
			CodigoPuntoVenta: factura.CodigoPuntoVenta,
			CodigoMotivo:     1,
		}

		body, err := json.Marshal(anulacionFactura)
		if err != nil {
			log.Printf("Error construyendo JSON de factura: %v", err)
			return err
		}

		url := fmt.Sprintf("%s/api/v1/facturacion/electronica/anular", os.Getenv("URL_FACTURADOR"))
		token := os.Getenv("TOKEN_FACTURA")

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error creando request: %v", err)
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		req.Header.Set("Connection", "close")

		client := &http.Client{Timeout: 15 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error al enviar solicitud a SIAT: %v", err)
		} else {
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(resp.Body)
			respBytes, _ := io.ReadAll(resp.Body)
			log.Printf("Respuesta SIAT (status %d): %s", resp.StatusCode, string(respBytes))
		}
	}

	// Anular venta en DB (aunque falle SIAT)
	err = v.ventaRepository.AnularVentaById(context.Background(), id)
	if err != nil {
		log.Printf("Error anulando venta en DB: %v", err)
	}
	return err
}

func NewVentaService(ventaRepository port.VentaRepository) *VentaService {
	return &VentaService{ventaRepository: ventaRepository}
}

var _ port.VentaService = (*VentaService)(nil)
