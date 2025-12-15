package domain

import (
	"encoding/xml"

	"github.com/goccy/go-json"
)

type NilableString struct {
	Value  *string
	XsiNil string `xml:"xsi:nil,attr,omitempty"`
}

func (n NilableString) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Value == nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
		return e.EncodeElement("", start)
	}
	return e.EncodeElement(*n.Value, start)
}

func (n NilableString) MarshalJSON() ([]byte, error) {
	if n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.Value)
}

type NilableUint64 struct {
	Value  *uint64
	XsiNil string `xml:"xsi:nil,attr,omitempty"`
}

func (n NilableUint64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Value == nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
		return e.EncodeElement("", start)
	}
	return e.EncodeElement(*n.Value, start)
}

func (n NilableUint64) MarshalJSON() ([]byte, error) {
	if n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.Value)
}

type NilableFloat64 struct {
	Value  *float64
	XsiNil string `xml:"xsi:nil,attr,omitempty"`
}

func (n NilableFloat64) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Value == nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
		return e.EncodeElement("", start)
	}
	return e.EncodeElement(*n.Value, start)
}

func (n NilableFloat64) MarshalJSON() ([]byte, error) {
	if n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.Value)
}

type NilableDecimal struct {
	Value  *float64
	XsiNil string `xml:"xsi:nil,attr,omitempty"`
}

func (n NilableDecimal) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if n.Value == nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xsi:nil"}, Value: "true"})
		return e.EncodeElement("", start)
	}
	return e.EncodeElement(*n.Value, start)
}

func (n NilableDecimal) MarshalJSON() ([]byte, error) {
	if n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.Value)
}

// ------------------ NilableString ------------------
func (n *NilableString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	n.Value = &s
	return nil
}

// ------------------ NilableUint64 ------------------
func (n *NilableUint64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var u uint64
	if err := json.Unmarshal(data, &u); err != nil {
		return err
	}
	n.Value = &u
	return nil
}

// ------------------ NilableFloat64 ------------------
func (n *NilableFloat64) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	n.Value = &f
	return nil
}

// ------------------ NilableDecimal ------------------
func (n *NilableDecimal) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.Value = nil
		return nil
	}
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	n.Value = &f
	return nil
}

func (n NilableString) ValueOrNil() interface{} {
	if n.Value == nil {
		return nil
	}
	return *n.Value
}

func (n NilableUint64) ValueOrNil() interface{} {
	if n.Value == nil {
		return nil
	}
	return *n.Value
}

func (n NilableFloat64) ValueOrNil() interface{} {
	if n.Value == nil {
		return nil
	}
	return *n.Value
}

// FacturaCompraVenta modelo para emitir factura
type FacturaCompraVenta struct {
	Cabecera Cabecera  `json:"cabecera"`
	Detalle  []Detalle `json:"detalle"`
}

type Cabecera struct {
	Municipio                    string         `json:"municipio"`
	Telefono                     *string        `json:"telefono,omitempty"`
	CodigoSucursal               uint64         `json:"codigoSucursal"`
	Direccion                    string         `json:"direccion"`
	CodigoPuntoVenta             NilableUint64  `json:"codigoPuntoVenta"`
	NombreRazonSocial            NilableString  `json:"nombreRazonSocial"`
	CodigoTipoDocumentoIdentidad uint64         `json:"codigoTipoDocumentoIdentidad"`
	NumeroDocumento              string         `json:"numeroDocumento"`
	Complemento                  NilableString  `json:"complemento,omitempty"`
	CodigoCliente                string         `json:"codigoCliente"`
	EmailCliente                 string         `json:"emailCliente,omitempty"`
	CodigoMetodoPago             uint64         `json:"codigoMetodoPago"`
	NumeroTarjeta                NilableUint64  `json:"numeroTarjeta"`
	MontoTotal                   float64        `json:"montoTotal"`
	CodigoMoneda                 uint64         `json:"codigoMoneda"`
	TipoCambio                   float64        `json:"tipoCambio"`
	MontoTotalMoneda             float64        `json:"montoTotalMoneda"`
	MontoGiftCard                NilableFloat64 `json:"montoGiftCard"`
	DescuentoAdicional           NilableFloat64 `json:"descuentoAdicional"`
	CodigoExcepcion              NilableUint64  `json:"codigoExcepcion"`
	Cafc                         NilableString  `json:"cafc"`
	Leyenda                      string         `json:"leyenda"`
	Usuario                      string         `json:"usuario"`
	CodigoDocumentoSector        uint64         `json:"codigoDocumentoSector"`
	TipoFacturaDocumento         uint64         `json:"tipoFacturaDocumento"`
}

type Detalle struct {
	ActividadEconomica string         `json:"actividadEconomica"`
	CodigoProductoSin  string         `json:"codigoProductoSin"`
	CodigoProducto     string         `json:"codigoProducto"`
	Descripcion        string         `json:"descripcion"`
	Cantidad           float64        `json:"cantidad"`
	UnidadMedida       uint64         `json:"unidadMedida"`
	PrecioUnitario     float64        `json:"precioUnitario"`
	MontoDescuento     NilableFloat64 `json:"montoDescuento"`
	SubTotal           float64        `json:"subTotal"`
	NumeroSerie        NilableString  `json:"numeroSerie,omitempty"`
	NumeroImei         NilableString  `json:"numeroImei,omitempty"`
}

type AnularFacturaRequest struct {
	NumeroFactura    uint64 `json:"numerofactura"`
	CodigoSucursal   uint64 `json:"codigoSucursal"`
	CodigoPuntoVenta uint64 `json:"codigoPuntoVenta"`
	CodigoMotivo     uint64 `json:"codigoMotivo"`
}

type FacturaCompraVentaResponse struct {
	Id               uint64 `json:"id"`
	NumeroFactura    uint64 `json:"numeroFactura"`
	CodigoSucursal   uint64 `json:"codigoSucursal"`
	CodigoPuntoVenta uint64 `json:"codigoPuntoVenta"`
	Cuf              string `json:"cuf"`
	Nit              uint64 `json:"nit"`
	Url              string `json:"url"`
}

type Factura struct {
	Id               uint64 `json:"id"`
	NumeroFactura    uint64 `json:"numeroFactura"`
	CodigoSucursal   uint64 `json:"codigoSucursal"`
	CodigoPuntoVenta uint64 `json:"codigoPuntoVenta"`
	Cuf              string `json:"cuf"`
	Nit              uint64 `json:"nit"`
	Url              string `json:"url"`
	VentaId          uint64 `json:"ventaId"`
}
