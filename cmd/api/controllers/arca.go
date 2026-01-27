package controllers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ---------------------------------------------------------
// Constantes de Configuración
// ---------------------------------------------------------

// URL de Homologación (Testing). Para producción usar: https://servicios1.afip.gov.ar/wsfev1/service.asmx
const URL_WSFE = "https://wswhomo.afip.gov.ar/wsfev1/service.asmx"
const SOAP_ACTION = "http://ar.gov.afip.dif.FEV1/FECAESolicitar"

// ---------------------------------------------------------
// Estructuras SOAP (Request)
// ---------------------------------------------------------

// Envelope SOAP genérico
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ soap:Envelope"`
	XSI     string   `xml:"xmlns:xsi,attr"`
	XSD     string   `xml:"xmlns:xsd,attr"`
	Body    SOAPBody
}

type SOAPBody struct {
	XMLName        xml.Name `xml:"soap:Body"`
	FECAESolicitar FECAESolicitarWrapper
}

// Wrapper para el método específico
type FECAESolicitarWrapper struct {
	XMLName xml.Name `xml:"https://wswhomo.afip.gov.ar/wsfev1/service.asmx?op= FECAESolicitar"`
	Auth    FEAuthRequest
	FeCAEReq FeCAERequest
}

// Datos de Autenticación (Lo que ya obtuviste del WSAA)
type FEAuthRequest struct {
	Token string `xml:"Token"`
	Sign  string `xml:"Sign"`
	Cuit  int64  `xml:"Cuit"`
}

// Estructura principal de la solicitud
type FeCAERequest struct {
	FeCabReq FeCabRequest // Cabecera del lote
	FeDetReq FeDetRequest // Array de facturas (detalle)
}

type FeCabRequest struct {
	CantReg int `xml:"CantReg"`
	PtoVta  int `xml:"PtoVta"`
	CbteTipo int `xml:"CbteTipo"` // 1=Factura A, 6=Factura B, 11=Factura C
}

// Usamos un slice para soportar múltiples facturas, aunque usualmente se envía 1
type FeDetRequest struct {
	FECAEDetRequest []FECAEDetRequest `xml:"FECAEDetRequest"`
}

type FECAEDetRequest struct {
	Concepto   int     `xml:"Concepto"` // 1=Productos, 2=Servicios, 3=Ambos
	DocTipo    int     `xml:"DocTipo"`  // 80=CUIT, 96=DNI, 99=Consumidor Final
	DocNro     int64   `xml:"DocNro"`
	CbteDesde  int64   `xml:"CbteDesde"`
	CbteHasta  int64   `xml:"CbteHasta"`
	CbteFch    string  `xml:"CbteFch"` // Formato YYYYMMDD
	ImpTotal   float64 `xml:"ImpTotal"`
	ImpTotConc float64 `xml:"ImpTotConc"` // No gravado
	ImpNeto    float64 `xml:"ImpNeto"`
	ImpOpEx    float64 `xml:"ImpOpEx"`    // Exento
	ImpTrib    float64 `xml:"ImpTrib"`
	ImpIVA     float64 `xml:"ImpIVA"`
	FchServDesde string `xml:"FchServDesde,omitempty"` // Solo si Concepto > 1
	FchServHasta string `xml:"FchServHasta,omitempty"` // Solo si Concepto > 1
	FchVtoPago   string `xml:"FchVtoPago,omitempty"`   // Solo si Concepto > 1
	MonId      string  `xml:"MonId"`    // "PES" para pesos
	MonCotiz   float64 `xml:"MonCotiz"` // 1 para pesos
	
	// Sub-arrays opcionales
	Iva      *ALIvaArray `xml:"Iva,omitempty"`
}

type ALIvaArray struct {
	AlicIva []AlicIva `xml:"AlicIva"`
}

type AlicIva struct {
	Id      int     `xml:"Id"` // 5=21%, 4=10.5%, 3=0%
	BaseImp float64 `xml:"BaseImp"`
	Importe float64 `xml:"Importe"`
}

// ---------------------------------------------------------
// Estructuras SOAP (Response)
// ---------------------------------------------------------

type SOAPResponseEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SOAPResponseBody
}

type SOAPResponseBody struct {
	FECAESolicitarResponse FECAESolicitarResultWrapper `xml:"FECAESolicitarResponse"`
}

type FECAESolicitarResultWrapper struct {
	Result FECAESolicitarResult `xml:"FECAESolicitarResult"`
}

type FECAESolicitarResult struct {
	FeCabResp FeCabResp
	FeDetResp FeDetResp
	Errors    *Errors `xml:"Errors,omitempty"` // Puntero para manejar nil si no hay errores
}

type FeCabResp struct {
	Cuit      int64  `xml:"Cuit"`
	PtoVta    int    `xml:"PtoVta"`
	CbteTipo  int    `xml:"CbteTipo"`
	FchProceso string `xml:"FchProceso"`
	CantReg   int    `xml:"CantReg"`
	Resultado string `xml:"Resultado"` // A=Aprobado, R=Rechazado, P=Parcial
	Reproceso string `xml:"Reproceso"`
}

type FeDetResp struct {
	FECAEDetResponse []FECAEDetResponse `xml:"FECAEDetResponse"`
}

type FECAEDetResponse struct {
	Concepto   int    `xml:"Concepto"`
	DocTipo    int    `xml:"DocTipo"`
	DocNro     int64  `xml:"DocNro"`
	CbteDesde  int64  `xml:"CbteDesde"`
	CbteHasta  int64  `xml:"CbteHasta"`
	CbteFch    string `xml:"CbteFch"`
	Resultado  string `xml:"Resultado"`
	CAE        string `xml:"CAE"`        // DATO CLAVE
	CAEFchVto  string `xml:"CAEFchVto"`  // DATO CLAVE
	Observaciones *Observaciones `xml:"Observaciones,omitempty"`
}

type Errors struct {
	Err []ErrorType `xml:"Err"`
}

type ErrorType struct {
	Code int    `xml:"Code"`
	Msg  string `xml:"Msg"`
}

type Observaciones struct {
	Obs []ObsType `xml:"Obs"`
}

type ObsType struct {
	Code int    `xml:"Code"`
	Msg  string `xml:"Msg"`
}

// ---------------------------------------------------------
// Función Principal para Solicitar CAE
// ---------------------------------------------------------

func SolicitarCAE(token, sign string, cuitEmisor int64, factura FECAEDetRequest, ptoVta int, cbteTipo int) (*FECAESolicitarResult, error) {
	
	// 1. Construir el Request
	req := SOAPEnvelope{
		XSI: "http://www.w3.org/2001/XMLSchema-instance",
		XSD: "http://www.w3.org/2001/XMLSchema",
		Body: SOAPBody{
			FECAESolicitar: FECAESolicitarWrapper{
				Auth: FEAuthRequest{
					Token: token,
					Sign:  sign,
					Cuit:  cuitEmisor,
				},
				FeCAEReq: FeCAERequest{
					FeCabReq: FeCabRequest{
						CantReg:  1,
						PtoVta:   ptoVta,
						CbteTipo: cbteTipo,
					},
					FeDetReq: FeDetRequest{
						FECAEDetRequest: []FECAEDetRequest{factura},
					},
				},
			},
		},
	}

	// 2. Serializar a XML
	xmlPayload, err := xml.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error al serializar XML: %v", err)
	}

	// 3. Crear Petición HTTP POST
	client := &http.Client{Timeout: 30 * time.Second}
	httpReq, err := http.NewRequest("POST", URL_WSFE, bytes.NewBuffer(xmlPayload))
	if err != nil {
		return nil, fmt.Errorf("error creando request HTTP: %v", err)
	}

	// Headers OBLIGATORIOS para SOAP
	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", SOAP_ACTION)

	// 4. Enviar Petición
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error enviando request a AFIP: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	bodyBytes, _ := io.ReadAll(resp.Body)
	
	// 5. Parsear Respuesta
	var soapResp SOAPResponseEnvelope
	err = xml.Unmarshal(bodyBytes, &soapResp)
	if err != nil {
		return nil, fmt.Errorf("error parseando respuesta XML: %v", err)
	}

	return &soapResp.Body.FECAESolicitarResponse.Result, nil
}

// ---------------------------------------------------------
// Ejemplo de Uso (Main)
// ---------------------------------------------------------
func main() {
	// Datos ficticios (Reemplazar con tus credenciales y datos reales)
	miToken := "PD94... (tu token real) ...>"
	miSign := "MZI... (tu sign real) ...="
	miCuit := int64(20123456789)

	// Crear una factura de ejemplo (Factura B de 121 pesos con IVA 21%)
	nuevaFactura := FECAEDetRequest{
		Concepto:   1,            // Productos
		DocTipo:    99,           // Consumidor Final
		DocNro:     0,            // 0 para CF montos bajos
		CbteDesde:  1,            // Número de comprobante (debe ser el siguiente libre)
		CbteHasta:  1,
		CbteFch:    time.Now().Format("20060102"), // AAAAMMDD
		ImpTotal:   121.00,
		ImpTotConc: 0,
		ImpNeto:    100.00,
		ImpOpEx:    0,
		ImpTrib:    0,
		ImpIVA:     21.00,
		MonId:      "PES",
		MonCotiz:   1,
		Iva: &ALIvaArray{
			AlicIva: []AlicIva{
				{Id: 5, BaseImp: 100.00, Importe: 21.00}, // Id 5 = 21%
			},
		},
	}

	// Llamar al servicio
	resultado, err := SolicitarCAE(miToken, miSign, miCuit, nuevaFactura, 1, 6) // PtoVta 1, Tipo 6 (Factura B)

	if err != nil {
		fmt.Printf("Error técnico: %v\n", err)
		return
	}

	// Analizar resultado de negocio
	fmt.Printf("Resultado Global: %s\n", resultado.FeCabResp.Resultado)

	if resultado.FeCabResp.Resultado == "A" {
		det := resultado.FeDetResp.FECAEDetResponse[0]
		fmt.Printf("✅ FACTURA APROBADA\n")
		fmt.Printf("CAE: %s\n", det.CAE)
		fmt.Printf("Vencimiento CAE: %s\n", det.CAEFchVto)
	} else {
		fmt.Printf("❌ FACTURA RECHAZADA\n")
		// Mostrar errores globales
		if resultado.Errors != nil {
			for _, e := range resultado.Errors.Err {
				fmt.Printf("Error Global %d: %s\n", e.Code, e.Msg)
			}
		}
		// Mostrar observaciones específicas de la factura
		if len(resultado.FeDetResp.FECAEDetResponse) > 0 {
			obs := resultado.FeDetResp.FECAEDetResponse[0].Observaciones
			if obs != nil {
				for _, o := range obs.Obs {
					fmt.Printf("Observación %d: %s\n", o.Code, o.Msg)
				}
			}
		}
	}
}