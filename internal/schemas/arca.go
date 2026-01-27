package schemas

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// ============ ENUMS Y CONSTANTES ============

type TipoComprobante int

const (
	FacturaA     TipoComprobante = 1
	NotaDebitoA  TipoComprobante = 2
	NotaCreditoA TipoComprobante = 3
	FacturaB     TipoComprobante = 6
	NotaDebitoB  TipoComprobante = 7
	NotaCreditoB TipoComprobante = 8
	FacturaC     TipoComprobante = 11
	NotaDebitoC  TipoComprobante = 12
	NotaCreditoC TipoComprobante = 13
	// FacturaCreditoA TipoComprobante = 201
	// FacturaCreditoB TipoComprobante = 206
	// FacturaCreditoC TipoComprobante = 211
)

var name_invoice_type = map[string]TipoComprobante{
	"FacturaA":     FacturaA,
	"FacturaB":     FacturaB,
	"FacturaC":     FacturaC,
	"NotaDebitoA":  NotaDebitoA,
	"NotaDebitoB":  NotaDebitoB,
	"NotaDebitoC":  NotaDebitoC,
	"NotaCreditoA": NotaCreditoA,
	"NotaCreditoB": NotaCreditoB,
	"NotaCreditoC": NotaCreditoC,
}

func GetCodeTipoComprobante(name string) (int, error) {
	code, ok := name_invoice_type[name]
	if !ok {
		return 0, fmt.Errorf("alícuota inválida: %s", name)
	}
	return int(code), nil
}

func (a *TipoComprobante) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	valor, ok := name_invoice_type[s]
	if !ok {
		return fmt.Errorf("alícuota inválida: %s", s)
	}

	*a = valor
	return nil
}

type TipoConcepto int

const (
	Productos           TipoConcepto = 1
	Servicios           TipoConcepto = 2
	ProductosYServicios TipoConcepto = 3
)

var name_concept = map[string]TipoConcepto{
	"Productos":           Productos,
	"Servicios":           Servicios,
	"ProductosYServicios": ProductosYServicios,
}

func GetCodeTipoConcepto(name string) (int, error) {
	code, ok := name_concept[name]
	if !ok {
		return 0, fmt.Errorf("alícuota inválida: %s", name)
	}
	return int(code), nil
}

func (a *TipoConcepto) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	valor, ok := name_concept[s]
	if !ok {
		return fmt.Errorf("alícuota inválida: %s", s)
	}

	*a = valor
	return nil
}

type TipoDocumento int

const (
	CUIT TipoDocumento = 80
	CUIL TipoDocumento = 86
	// CDI             TipoDocumento = 87
	// LE              TipoDocumento = 89
	// LC              TipoDocumento = 90
	// CI_Extranjera   TipoDocumento = 91
	// EnTramite       TipoDocumento = 92
	// Acta_Nacimiento TipoDocumento = 93
	// CI_Bs_As        TipoDocumento = 95
	DNI TipoDocumento = 96
	// Pasaporte       TipoDocumento = 94
	SinIdentificar TipoDocumento = 99
)

var name_document = map[string]TipoDocumento{
	"CUIT":           CUIT,
	"CUIL":           CUIL,
	"DNI":            DNI,
	"SinIdentificar": SinIdentificar,
}

func GetCodeTipoDocumento(name string) (int, error) {
	code, ok := name_document[name]
	if !ok {
		return 0, fmt.Errorf("alícuota inválida: %s", name)
	}
	return int(code), nil
}

func (a *TipoDocumento) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	valor, ok := name_document[s]
	if !ok {
		return fmt.Errorf("alícuota inválida: %s", s)
	}

	*a = valor
	return nil
}

type AlicuotaIVA int

const (
	IVA_0    AlicuotaIVA = 3
	IVA_10_5 AlicuotaIVA = 4
	IVA_21   AlicuotaIVA = 5
	IVA_27   AlicuotaIVA = 6
	IVA_5    AlicuotaIVA = 8
	IVA_2_5  AlicuotaIVA = 9
)

var name_alicuota = map[string]AlicuotaIVA{
	"IVA_0":    IVA_0,
	"IVA_10_5": IVA_10_5,
	"IVA_21":   IVA_21,
	"IVA_27":   IVA_27,
	"IVA_5":    IVA_5,
	"IVA_2_5":  IVA_2_5,
}

func GetCodeAlicuotaIVA(name string) (int, error) {
	valor, ok := name_alicuota[name]
	if !ok {
		return 0, fmt.Errorf("alícuota inválida: %s", name)
	}
	return int(valor), nil
}

func (a *AlicuotaIVA) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	valor, ok := name_alicuota[s]
	if !ok {
		return fmt.Errorf("alícuota inválida: %s", s)
	}

	*a = valor
	return nil
}

type CondicionIVAReceptor int

const (
	IVAResponsableInscripto CondicionIVAReceptor = 1
	IVASujetoExento         CondicionIVAReceptor = 4
	ConsumidorFinal         CondicionIVAReceptor = 5
	ResponsableMonotributo  CondicionIVAReceptor = 6
	// SujetoNoCategorizado    CondicionIVAReceptor = 7
	// ProveedorExterior       CondicionIVAReceptor = 8
	// ClienteExterior         CondicionIVAReceptor = 9
	// IVALiberado             CondicionIVAReceptor = 10
	// MonotributoSocial       CondicionIVAReceptor = 13
	// IVANoAlcanzado          CondicionIVAReceptor = 15
	// MonotributoTrabajador   CondicionIVAReceptor = 16
)

var name_condition = map[string]CondicionIVAReceptor{
	"responsable_inscripto": IVAResponsableInscripto,
	"exento":         IVASujetoExento,
	"consumidor_final":         ConsumidorFinal,
	"monotributista":  ResponsableMonotributo,
	// "SujetoNoCategorizado":    SujetoNoCategorizado,
	// "ProveedorExterior":       ProveedorExterior,
	// "ClienteExterior":         ClienteExterior,
	// "IVALiberado":             IVALiberado,
	// "MonotributoSocial":       MonotributoSocial,
	// "IVANoAlcanzado":          IVANoAlcanzado,
	// "MonotributoTrabajador":   MonotributoTrabajador,
}
func GetCodeConditionIVA(name string) (int, error) {
	valor, ok := name_condition[name]
	if !ok {
		return 0, fmt.Errorf("alícuota inválida: %s", name)
	}
	return int(valor), nil
}

func (a *CondicionIVAReceptor) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	valor, ok := name_condition[s]
	if !ok {
		return fmt.Errorf("alícuota inválida: %s", s)
	}

	*a = valor
	return nil
}

// ============================================
// WSFE - FACTURACIÓN ELECTRÓNICA
// ============================================

type FECAERequest struct {
	XMLName  xml.Name `xml:"FECAESolicitar"`
	Xmlns    string   `xml:"xmlns,attr"`
	Auth     Auth     `xml:"Auth"`
	FeCAEReq FeCAEReq `xml:"FeCAEReq"`
}

type Auth struct {
	Token string `xml:"Token"`
	Sign  string `xml:"Sign"`
	Cuit  int64  `xml:"Cuit"`
}

type FeCAEReq struct {
	FeCabReq FeCabReq          `xml:"FeCabReq"`
	FeDetReq []FECAEDetRequest `xml:"FeDetReq>FECAEDetRequest"`
}

type FeCabReq struct {
	CantReg  int `xml:"CantReg"`
	PtoVta   int `xml:"PtoVta"`
	CbteTipo int `xml:"CbteTipo"`
}

type FECAEDetRequest struct {
	Concepto               int            `xml:"Concepto"`
	DocTipo                int            `xml:"DocTipo"`
	DocNro                 int64          `xml:"DocNro"`
	CbteDesde              int64          `xml:"CbteDesde"`
	CbteHasta              int64          `xml:"CbteHasta"`
	CbteFch                string         `xml:"CbteFch"`
	ImpTotal               float64        `xml:"ImpTotal"`
	ImpTotConc             float64        `xml:"ImpTotConc"`
	ImpNeto                float64        `xml:"ImpNeto"`
	ImpOpEx                float64        `xml:"ImpOpEx"`
	ImpTrib                float64        `xml:"ImpTrib"`
	ImpIVA                 float64        `xml:"ImpIVA"`
	FchServDesde           string         `xml:"FchServDesde,omitempty"`
	FchServHasta           string         `xml:"FchServHasta,omitempty"`
	FchVtoPago             string         `xml:"FchVtoPago,omitempty"`
	MonId                  string         `xml:"MonId"`
	MonCotiz               float64        `xml:"MonCotiz"`
	CondicionIVAReceptorId int            `xml:"CondicionIVAReceptorId,omitempty"` // ← NOMBRE CORRECTO
	Tributos               *TributosArray `xml:"Tributos,omitempty"`
	Iva                    *IvaArray      `xml:"Iva,omitempty"`
}

type IvaArray struct {
	AlicIva []AlicIva `xml:"AlicIva"`
}

type AlicIva struct {
	Id      int     `xml:"Id"`
	BaseImp float64 `xml:"BaseImp"`
	Importe float64 `xml:"Importe"`
}

type TributosArray struct {
	Tributo []Tributo `xml:"Tributo"`
}

type Tributo struct {
	Id      int     `xml:"Id"`
	Desc    string  `xml:"Desc"`
	BaseImp float64 `xml:"BaseImp"`
	Alic    float64 `xml:"Alic"`
	Importe float64 `xml:"Importe"`
}

type FECAEResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Response struct {
			Results struct {
				FeCabResp FeCabResp `xml:"FeCabResp"`
				FeDetResp struct {
					FECAEDetResponse []FECAEDetResponse `xml:"FECAEDetResponse"`
				} `xml:"FeDetResp"`
			} `xml:"FECAESolicitarResult"`
			Errors struct {
				Err []ErrorAfip `xml:"Err"`
			} `xml:"Errors"`
		} `xml:"FECAESolicitarResponse"`
	} `xml:"Body"`
}

type FeCabResp struct {
	Cuit       int64  `xml:"Cuit"`
	PtoVta     int    `xml:"PtoVta"`
	CbteTipo   int    `xml:"CbteTipo"`
	FchProceso string `xml:"FchProceso"`
	CantReg    int    `xml:"CantReg"`
	Resultado  string `xml:"Resultado"`
	Reproceso  string `xml:"Reproceso"`
}

type FECAEDetResponse struct {
	Concepto      int    `xml:"Concepto"`
	DocTipo       int    `xml:"DocTipo"`
	DocNro        int64  `xml:"DocNro"`
	CbteDesde     int64  `xml:"CbteDesde"`
	CbteHasta     int64  `xml:"CbteHasta"`
	CbteFch       string `xml:"CbteFch"`
	Resultado     string `xml:"Resultado"`
	CAE           string `xml:"CAE"`
	CAEFchVto     string `xml:"CAEFchVto"`
	Observaciones struct {
		Obs []Observacion `xml:"Obs"`
	} `xml:"Observaciones"`
}

type Observacion struct {
	Code int    `xml:"Code"`
	Msg  string `xml:"Msg"`
}

type ErrorAfip struct {
	Code int    `xml:"Code"`
	Msg  string `xml:"Msg"`
}

type WSFEClient struct {
	BaseURL string
	Auth    Auth
	Client  *http.Client
}

// ============================================
// MODELO DE FACTURA
// ============================================

type FacturaRequest struct {
	TipoComprobante   int         `json:"tipo_comprobante"`
	TipoDocumento     int         `json:"tipo_documento"`
	NumeroDocumento   int64 		`json:"numero_documento"`
	Domicilio         string      `json:"domicilio"`
	CondicionIVA      string 			`json:"condicion_iva" validate:"required,oneof=responsable_inscripto exento consumidor_final monotributista"`
	NombreRazonSocial string 		`json:"nombre_razon_social"`
	Items         []ItemIVATotal `json:"items"`
}

type Factura struct {
	PuntoVenta       int
	TipoComprobante  int
	NumeroDesde      int64
	NumeroHasta      int64
	TipoDocumento    int
	NumeroDocumento  int64
	CondicionIVA     int
	Concepto         int
	Fecha            string
	FechaServDesde   string
	FechaServHasta   string
	FechaVtoPago     string
	ImporteNeto      float64
	ImporteNoGravado float64
	ImporteExento    float64
	ImporteIVA       float64
	ImporteTributos  float64
	ImporteTotal     float64
	Alicuotas        []ItemIVA
	Tributos         []ItemTributo
	MonedaId         string
	MonedaCotiz      float64
}

type ItemIVA struct {
	Codigo        int
	BaseImponible float64
	Importe       float64
}

type ItemIVATotal struct {
	ProductID 	int64 `json:"product_id"`
	Product string `json:"product"`
	Quantity    float64 `json:"quantity"`
	Codigo        string `json:"codigo_iva"`
	BaseImponible float64 `json:"base_imponible"`
	Importe       float64 `json:"importe_iva"`
}

type ItemTributo struct {
	Codigo        int
	Descripcion   string
	BaseImponible float64
	Alicuota      float64
	Importe       float64
}

func (f *Factura) ToFECAEDetRequest() FECAEDetRequest {
	det := FECAEDetRequest{
		Concepto:               f.Concepto,
		DocTipo:                f.TipoDocumento,
		DocNro:                 f.NumeroDocumento,
		CbteDesde:              f.NumeroDesde,
		CbteHasta:              f.NumeroHasta,
		CbteFch:                f.Fecha,
		ImpTotal:               f.ImporteTotal,
		ImpTotConc:             f.ImporteNoGravado,
		ImpNeto:                f.ImporteNeto,
		ImpOpEx:                f.ImporteExento,
		ImpTrib:                f.ImporteTributos,
		ImpIVA:                 f.ImporteIVA,
		MonId:                  f.MonedaId,
		MonCotiz:               f.MonedaCotiz,
		FchServDesde:           f.FechaServDesde,
		FchServHasta:           f.FechaServHasta,
		FchVtoPago:             f.FechaVtoPago,
		CondicionIVAReceptorId: f.CondicionIVA, // ← MAPEO CORRECTO
	}

	if len(f.Alicuotas) > 0 {
		det.Iva = &IvaArray{}
		for _, iva := range f.Alicuotas {
			det.Iva.AlicIva = append(det.Iva.AlicIva, AlicIva{
				Id:      iva.Codigo,
				BaseImp: iva.BaseImponible,
				Importe: iva.Importe,
			})
		}
	}

	if len(f.Tributos) > 0 {
		det.Tributos = &TributosArray{}
		for _, trib := range f.Tributos {
			det.Tributos.Tributo = append(det.Tributos.Tributo, Tributo{
				Id:      trib.Codigo,
				Desc:    trib.Descripcion,
				BaseImp: trib.BaseImponible,
				Alic:    trib.Alicuota,
				Importe: trib.Importe,
			})
		}
	}

	return det
}

// ============================================
// WSAA - AUTENTICACIÓN
// ============================================

type LoginCMSResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		LoginCmsReturn struct {
			Return string `xml:"loginCmsReturn"`
		} `xml:"loginCmsResponse"`
	} `xml:"Body"`
}

type TicketResponse struct {
	XMLName xml.Name `xml:"loginTicketResponse"`
	Header  struct {
		Source      string `xml:"source"`
		Destination string `xml:"destination"`
		UniqueID    int64  `xml:"uniqueId"`
		Generation  string `xml:"generationTime"`
		Expiration  string `xml:"expirationTime"`
	} `xml:"header"`
	Credentials struct {
		Token string `xml:"token"`
		Sign  string `xml:"sign"`
	} `xml:"credentials"`
}

type CredentialsValidation struct {
	Token      string
	Sign       string
	Expiration time.Time
	CUIT       int64
}

type WSAAConfig struct {
	Homologacion bool
	CertFile     string
	KeyFile      string
	CUIT         int64
	Service      string
}

type WSAA struct {
	Config WSAAConfig
	Client *http.Client
}

func NewWSAA(config WSAAConfig) (*WSAA, error) {
	return &WSAA{
		Config: config,
		Client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.Homologacion,
				},
			},
		},
	}, nil
}

func SendToWSAA(w *WSAA, signedCMS string) ([]byte, error) {
	response, err := w.Client.Post("https://wsaahomo.afip.gov.ar/ws/services/LoginCms", "text/xml; charset=utf-8", strings.NewReader(signedCMS))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func (w *WSAA) CreateTicketXML() ([]byte, error) {
	now := time.Now()

	ticket := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<loginTicketRequest version="1.0">
<header>
    <uniqueId>%d</uniqueId>
    <generationTime>%s</generationTime>
    <expirationTime>%s</expirationTime>
</header>
<service>%s</service>
</loginTicketRequest>`,
		now.Unix(),
		now.Add(-10*time.Minute).Format("2006-01-02T15:04:05"),
		now.Add(12*time.Hour).Format("2006-01-02T15:04:05"),
		w.Config.Service,
	)

	return []byte(ticket), nil
}

func (w *WSAA) validateKeys() error {
	keyData, err := os.ReadFile(w.Config.KeyFile)
	if err != nil {
		return fmt.Errorf("no se puede leer la clave privada: %v", err)
	}

	certData, err := os.ReadFile(w.Config.CertFile)
	if err != nil {
		return fmt.Errorf("no se puede leer el certificado: %v", err)
	}

	if !bytes.Contains(keyData, []byte("BEGIN")) {
		return fmt.Errorf("la clave privada no está en formato PEM")
	}

	if bytes.Contains(keyData, []byte("ENCRYPTED")) {
		return fmt.Errorf("la clave privada está encriptada. Debe desencriptarla primero")
	}

	if !bytes.Contains(certData, []byte("BEGIN")) {
		return fmt.Errorf("el certificado no está en formato PEM")
	}

	return nil
}

func (w *WSAA) SignWithOpenSSL(data []byte) (string, error) {
	// 1. Validaciones previas (sin cambios, son correctas)
	if err := w.validateKeys(); err != nil {
		return "", err
	}

	// 2. Configurar el comando para usar STDIN y STDOUT
	// Al usar "-" en -in y no poner -out, openssl usa los pipes de Go
	cmd := exec.Command("openssl", "cms", "-sign",
		"-in", "-", // Lee desde Stdin
		"-signer", w.Config.CertFile,
		"-inkey", w.Config.KeyFile,
		"-outform", "DER",
		"-nodetach",
		"-binary",
		"-md", "sha256",
	)

	// Inyectar los datos del XML directamente al proceso
	cmd.Stdin = bytes.NewReader(data)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("⚠️ Fallo cms, reintentando con smime en memoria...")

		// Fallback a smime
		cmd = exec.Command("openssl", "smime", "-sign",
			"-in", "-",
			"-signer", w.Config.CertFile,
			"-inkey", w.Config.KeyFile,
			"-outform", "DER",
			"-nodetach",
			"-binary",
		)
		cmd.Stdin = bytes.NewReader(data)
		stdout.Reset()
		stderr.Reset()
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("error openssl: %v\nStderr: %s", err, stderr.String())
		}
	}

	log.Info().Msgf("✅ CMS generado: %d bytes (en memoria)", stdout.Len())

	// 3. Codificar directamente el buffer de salida
	return base64.StdEncoding.EncodeToString(stdout.Bytes()), nil
}

func (w *WSAA) ParseResponse(data []byte) (*CredentialsValidation, error) {
	var loginResp LoginCMSResponse
	if err := xml.Unmarshal(data, &loginResp); err != nil {
		return nil, fmt.Errorf("error parseando SOAP: %v\nRespuesta: %s", err, string(data))
	}

	ticketXML := loginResp.Body.LoginCmsReturn.Return
	if ticketXML == "" {
		return nil, fmt.Errorf("respuesta vacía del WSAA")
	}

	var ticketResp TicketResponse
	if err := xml.Unmarshal([]byte(ticketXML), &ticketResp); err != nil {
		return nil, fmt.Errorf("error parseando ticket: %v", err)
	}

	expTime, err := time.Parse("2006-01-02T15:04:05.000-07:00", ticketResp.Header.Expiration)
	if err != nil {
		expTime, err = time.Parse("2006-01-02T15:04:05", ticketResp.Header.Expiration)
		if err != nil {
			expTime = time.Now().Add(12 * time.Hour)
		}
	}

	return &CredentialsValidation{
		Token:      ticketResp.Credentials.Token,
		Sign:       ticketResp.Credentials.Sign,
		Expiration: expTime,
		CUIT:       w.Config.CUIT,
	}, nil
}

func (w *WSFEClient) BuildSOAPEnvelope(req FECAERequest) string {
	reqXML, _ := xml.MarshalIndent(req, "    ", "  ")

	soap := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
  <soap:Body>
    %s
  </soap:Body>
</soap:Envelope>`, string(reqXML))

	return soap
}

type FECompConsultaRequest struct {
	XMLName       xml.Name `xml:"FECompConsultar"`
	Xmlns         string   `xml:"xmlns,attr"`
	Auth          Auth     `xml:"Auth"`
	FeCompConsReq struct {
		CbteTipo int   `xml:"CbteTipo"`
		CbteNro  int64 `xml:"CbteNro"`
		PtoVta   int   `xml:"PtoVta"`
	} `xml:"FeCompConsReq"`
}

type FECompConsultaResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Response struct {
			Result struct {
				ResultGet struct {
					CbteDesde int64   `xml:"CbteDesde"`
					CbteHasta int64   `xml:"CbteHasta"`
					CbteFch   string  `xml:"CbteFch"`
					ImpTotal  float64 `xml:"ImpTotal"`
					CodAut    string  `xml:"CodAut"` // Este es el CAE
					FchVto    string  `xml:"FchVto"` // Vto del CAE
					Resultado string  `xml:"Resultado"`
					DocTipo   int     `xml:"DocTipo"`
					DocNro    int64   `xml:"DocNro"`
					// Puedes agregar más campos según necesites del XSD
				} `xml:"ResultGet"`
				Errors struct {
					Err []ErrorAfip `xml:"Err"`
				} `xml:"Errors"`
			} `xml:"FECompConsultarResult"`
		} `xml:"FECompConsultarResponse"`
	} `xml:"Body"`
}

type FECompUltimoAutorizadoRequest struct {
	XMLName  xml.Name `xml:"FECompUltimoAutorizado"`
	Xmlns    string   `xml:"xmlns,attr"`
	Auth     Auth     `xml:"Auth"`
	PtoVta   int      `xml:"PtoVta"`
	CbteTipo int      `xml:"CbteTipo"`
}

type FECompUltimoAutorizadoResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		Response struct {
			Result struct {
				CbteNro  int64 `xml:"CbteNro"`
				PtoVta   int   `xml:"PtoVta"`
				CbteTipo int   `xml:"CbteTipo"`
				Errors   struct {
					Err []ErrorAfip `xml:"Err"`
				} `xml:"Errors"`
			} `xml:"FECompUltimoAutorizadoResult"`
		} `xml:"FECompUltimoAutorizadoResponse"`
	} `xml:"Body"`
}

func LoadCredentials(token, sign, cuitCred string, expires time.Time) (*CredentialsValidation, error) {
	creds := &CredentialsValidation{}

	cuit, err := strconv.ParseInt(cuitCred, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("error parseando CUIT: %v", err)
	}

	creds.Token = token
	creds.Sign = sign
	creds.CUIT = cuit
	creds.Expiration = expires

	return creds, nil
}

func NewWSFEClient(token, sign string, cuit int64, homologacion bool) *WSFEClient {
	url := "https://servicios1.afip.gov.ar/wsfev1/service.asmx"
	if homologacion {
		url = "https://wswhomo.afip.gov.ar/wsfev1/service.asmx"
	}

	return &WSFEClient{
		BaseURL: url,
		Auth: Auth{
			Token: token,
			Sign:  sign,
			Cuit:  cuit,
		},
		Client: &http.Client{Timeout: 30 * time.Second},
	}
}

type FacturaElectronica struct {
	// Datos del Comprobante
	TipoComprobante int    `json:"tipo_comprobante"` // 1=A, 6=B, 11=C
	PuntoVenta      int    `json:"punto_venta"`      //
	Numero          int    `json:"numero"`           // Correlativo
	Fecha           string `json:"fecha"`            // YYYY-MM-DD
	Concepto        int    `json:"concepto"`         // 1=Prod, 2=Serv, 3=Mixto

	// Emisor (Tenant)
	EmisorCUIT         int64  `json:"emisor_cuit"`          //
	EmisorNombre       string `json:"emisor_nombre"`        //
	RazonSocial        string `json:"razon_social"`         //
	IngresosBrutos     string `json:"ingresos_brutos"`      //
	InicioActividades  string `json:"inicio_actividades"`   //
	CondicionIVAEmisor string `json:"condicion_iva_emisor"` // RI o Monotributo
	DomicilioEmisor    string `json:"domicilio_emisor"`     //

	// Receptor (Cliente)
	ReceptorCUIT         int64  `json:"receptor_cuit"`          //
	ReceptorNombre       string `json:"receptor_nombre"`        //
	CondicionIVAReceptor string `json:"condicion_iva_receptor"` //
	ReceptorDomicilio    string `json:"receptor_domicilio"`     //

	// Totales y Desglose (Crítico para Factura A)
	ImporteNeto   float64 `json:"importe_neto"`   // Base imponible
	ImporteIVA    float64 `json:"importe_iva"`    // Total IVA
	ImporteExento float64 `json:"importe_exento"` // Para artículos que no pagan IVA
	ImporteTotal  float64 `json:"importe_total"`  //

	// Datos de AFIP/ARCA
	CAE            int64  `json:"cae"`             // Código de 14 dígitos
	CAEVencimiento string `json:"cae_vencimiento"` //
	URL_QR         string `json:"url_qr"`          // URL con Base64
}

type DatosQR struct {
	Ver        int     `json:"ver"`        // Versión del formato (actualmente 1)
	Fecha      string  `json:"fecha"`      // Fecha de emisión (AAAA-MM-DD)
	Cuit       int64   `json:"cuit"`       // CUIT del emisor (sin guiones)
	PtoVta     int     `json:"ptoVta"`     // Punto de venta
	TipoCmp    int     `json:"tipoCmp"`    // Tipo de comprobante (ej: 1 para Factura A)
	NroCmp     int     `json:"nroCmp"`     // Número de comprobante
	Importe    float64 `json:"importe"`    // Importe total
	Moneda     string  `json:"moneda"`     // Moneda (ej: "PES")
	Ctz        int     `json:"ctz"`        // Cotización (1 para pesos)
	TipoDocRec int     `json:"tipoDocRec"` // Tipo doc receptor (ej: 80 para CUIT)
	NroDocRec  int64   `json:"nroDocRec"`  // Número de documento receptor
	TipoCodAut string  `json:"tipoCodAut"` // Tipo de autorización ("E" para CAE)
	CodAut     int64   `json:"codAut"`     // El número de CAE
}
