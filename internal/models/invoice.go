package models

import "time"

type Invoice struct {
	ID int64 `gorm:"primaryKey" json:"id"`
    // PointSaleID int64 `gorm:"not null" json:"point_sale_id"`
    // IncomeSaleID int64 `gorm:"not null" json:"income_sale_id"`
    // ClientID *int64 `gorm:"not null" json:"client_id"`
    // IsEcommerce bool `gorm:"not null;default:false" json:"is_ecommerce"`



	TypeVoucher int64 `gorm:"not null" json:"type_voucher"`
	PointSale int64 `gorm:"not null" json:"point_sale"`
    NumberVoucher int64 `gorm:"not null" json:"number_voucher"`
    DateVoucher string `gorm:"not null" json:"date_voucher"`
    Cae string `gorm:"not null" json:"cae"`
    CaeMaturity string `gorm:"not null" json:"cae_maturity"`
    Status string `gorm:"not null" json:"status"`
    IsTest bool `gorm:"not null;default:false" json:"is_test"`
    Subtotal float64 `gorm:"not null" json:"subtotal"`
    Description *string `gorm:"type:text" json:"description"`
    Discount float64 `gorm:"not null;default:0" json:"discount"`
    TypeDiscount string `gorm:"not null;default:percent" json:"type_discount" validate:"oneof=amount percent"`
    Total float64 `gorm:"not null" json:"total"`
    IsBudget bool `gorm:"not null;default:false" json:"is_budget"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}


// type FacturaElectronica struct {
// 	// Datos del Comprobante
// 	TipoComprobante int    `json:"tipo_comprobante"` // 1=A, 6=B, 11=C
// 	PuntoVenta      int    `json:"punto_venta"`      //
// 	Numero          int64  `json:"numero"`           // Correlativo
// 	Fecha           string `json:"fecha"`            // YYYY-MM-DD
// 	Concepto        int    `json:"concepto"`         // 1=Prod, 2=Serv, 3=Mixto

// 	// Emisor (Tenant)
// 	EmisorCUIT         int64  `json:"emisor_cuit"`          //
// 	EmisorNombre       string `json:"emisor_nombre"`        //
// 	RazonSocial        string `json:"razon_social"`         //
// 	IngresosBrutos     string `json:"ingresos_brutos"`      //
// 	InicioActividades  string `json:"inicio_actividades"`   //
// 	CondicionIVAEmisor string `json:"condicion_iva_emisor"` // RI o Monotributo
// 	DomicilioEmisor    string `json:"domicilio_emisor"`     //

// 	// Receptor (Cliente)
// 	ReceptorCUIT         int64  `json:"receptor_cuit"`          //
// 	ReceptorNombre       string `json:"receptor_nombre"`        //
// 	CondicionIVAReceptor string `json:"condicion_iva_receptor"` //
// 	ReceptorDomicilio    string `json:"receptor_domicilio"`     //

// 	// Totales y Desglose (Crítico para Factura A)
// 	ImporteNeto   float64 `json:"importe_neto"`          // Base imponible
// 	ImporteIVA    float64 `json:"importe_iva,omitempty"` // Total IVA
// 	ImporteExento float64 `json:"importe_exento"`        // Para artículos que no pagan IVA
// 	ImporteTotal  float64 `json:"importe_total"`         //

// 	Items []ItemIVATotal `json:"items"`

// 	// Datos de AFIP/ARCA
// 	CAE            int64  `json:"cae"`             // Código de 14 dígitos
// 	CAEVencimiento string `json:"cae_vencimiento"` //
// 	URL_QR         string `json:"url_qr"`          // URL con Base64
// }