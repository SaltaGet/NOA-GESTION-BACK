package repositories

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"gorm.io/gorm"
)

func (r *ArcaRepository) GetCredentialsArca(tenantID int64) (*models.Credential, error) {
	db := database.GetMainDB()
	var credential models.Credential
	err := db.
		Select("id", "social_reason", "business_name", "address", "responsibility_front_iva", "cuit", "gross_income", "start_activities", "arca_certificate", "arca_key", "token_arca", "sign_arca", "expire_token_arca", "concept").
		Where("tenant_id = ?", tenantID).First(&credential).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Credenciales no encontradas", err)
		}
		return nil, schemas.ErrorResponse(500, "Error interno al obtener credenciales", err)
	}

	if credential.SocialReason == nil || credential.ResponsibilityFrontIVA == nil || credential.Cuit == nil || credential.ArcaCertificate == nil || credential.ArcaKey == nil {
		return nil, schemas.ErrorResponse(422, "La entidad no se puede procesar por datos incompletos, revise y complete adecuadamente las credenciales", errors.New("La entidad no se puede procesar por datos incompletos."))
	}

	return &credential, nil
}

func (r *ArcaRepository) SetTokenSignArca(v *schemas.CredentialsValidation) error {
	db := database.GetMainDB()

	result := db.Model(&models.Credential{}).
		Where("cuit = ?", v.CUIT).
		Updates(map[string]interface{}{
			"token_arca":        v.Token,
			"sign_arca":         v.Sign,
			"expire_token_arca": v.Expiration,
		})

	if result.Error != nil {
		return schemas.ErrorResponse(500, "Error interno al actualizar credenciales", result.Error)
	}

	if result.RowsAffected == 0 {
		return schemas.ErrorResponse(404, "Credenciales no encontradas para actualizar", errors.New("Credencial no encontrada"))
	}

	return nil
}
	// creds, err := loadCredentials()
	// if err != nil || time.Now().After(creds.Expiration) {
	// 	fmt.Println("📡 Solicitando nuevas credenciales al WSAA...")

	// 	wsaa, err := schemas.NewWSAA(config)
	// 	if err != nil {
	// 		log.Fatalf("❌ Error creando cliente WSAA: %v", err)
	// 	}

	// 	creds, err = wsaa.GetCredentials()
	// 	if err != nil {
	// 		log.Fatalf("❌ Error obteniendo credenciales: %v", err)
	// 	}

	// 	if err := saveCredentials(creds); err != nil {
	// 		log.Printf("⚠️  No se pudieron guardar las credenciales: %v", err)
	// 	} else {
	// 		fmt.Println("💾 Credenciales guardadas en credentials.env")
	// 	}
	// } else {
	// 	fmt.Println("✅ Usando credenciales existentes de credentials.env")
	// }

func (r *ArcaRepository) GetLastestInvoice(w *schemas.WSFEClient, pointSale, TypeInvoice int) (int64, error) {
	req := schemas.FECompUltimoAutorizadoRequest{
		Xmlns:    "http://ar.gov.afip.dif.FEV1/",
		Auth:     w.Auth,
		PtoVta:   pointSale,
		CbteTipo: TypeInvoice,
	}

	reqXML, _ := xml.Marshal(req)
	soapEnv := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
			<soap:Body>
				%s
			</soap:Body>
		</soap:Envelope>`, string(reqXML))

	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
	if err != nil {
		return 0, err
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECompUltimoAutorizado")

	resp, err := w.Client.Do(httpReq)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response schemas.FECompUltimoAutorizadoResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return 0, err
	}

	if len(response.Body.Response.Result.Errors.Err) > 0 {
		return 0, fmt.Errorf("AFIP Error: %s", response.Body.Response.Result.Errors.Err[0].Msg)
	}

	return response.Body.Response.Result.CbteNro, nil
}

func (r *ArcaRepository) GetInfoInvoice(w *schemas.WSFEClient, pointSale, typeInvoice int, numberInvoice int64) (*schemas.FECompConsultaResponse, error) {
	req := schemas.FECompConsultaRequest{
		Xmlns: "http://ar.gov.afip.dif.FEV1/",
		Auth:  w.Auth,
	}
	req.FeCompConsReq.CbteTipo = typeInvoice
	req.FeCompConsReq.PtoVta = pointSale
	req.FeCompConsReq.CbteNro = numberInvoice

	reqXML, _ := xml.MarshalIndent(req, "    ", "  ")
	soapEnv := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
			<soap:Body>
				%s
			</soap:Body>
		</soap:Envelope>`, string(reqXML))

	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECompConsultar")

	resp, err := w.Client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var response schemas.FECompConsultaResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error parseando respuesta: %v", err)
	}

	// Verificar errores de AFIP
	if len(response.Body.Response.Result.Errors.Err) > 0 {
		return nil, fmt.Errorf("AFIP Error: %s", response.Body.Response.Result.Errors.Err[0].Msg)
	}

	return &response, nil
}

func (r *ArcaRepository) SendToWSAA(w *schemas.WSAA, cms string) ([]byte, error) {
	url := "https://wsaahomo.afip.gov.ar/ws/services/LoginCms"
	if !w.Config.Homologacion {
		url = "https://wsaa.afip.gov.ar/ws/services/LoginCms"
	}

	soapRequest := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:wsaa="http://wsaa.view.sua.dvadac.desein.afip.gov">
				<soapenv:Header/>
				<soapenv:Body>
						<wsaa:loginCms>
								<wsaa:in0>%s</wsaa:in0>
						</wsaa:loginCms>
				</soapenv:Body>
		</soapenv:Envelope>`, cms)

	req, err := http.NewRequest("POST", url, strings.NewReader(soapRequest))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "")

	resp, err := w.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WSAA error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// func (r *ArcaRepository) EmitInvoice(factura *schemas.Factura) error {
// 	return nil
// }

func (r *ArcaRepository) EmitInvoice(w *schemas.WSFEClient, factura *schemas.Factura) (*schemas.FECAEDetResponse, error) {
	req := schemas.FECAERequest{
		Xmlns: "http://ar.gov.afip.dif.FEV1/",
		Auth:  w.Auth,
		FeCAEReq: schemas.FeCAEReq{
			FeCabReq: schemas.FeCabReq{
				CantReg:  1,
				PtoVta:   factura.PuntoVenta,
				CbteTipo: factura.TipoComprobante,
			},
			FeDetReq: []schemas.FECAEDetRequest{factura.ToFECAEDetRequest()},
		},
	}

	soapEnv := w.BuildSOAPEnvelope(req)

	httpReq, err := http.NewRequest("POST", w.BaseURL, strings.NewReader(soapEnv))
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al realizar solicitud a ARCA", err)
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")
	httpReq.Header.Set("SOAPAction", "http://ar.gov.afip.dif.FEV1/FECAESolicitar")

	resp, err := w.Client.Do(httpReq)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error al realizar solicitud a ARCA", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, schemas.ErrorResponse(resp.StatusCode, "WSFE error en la solicitud", errors.New(string(body)))
	}

	var response schemas.FECAEResponse
	if err := xml.Unmarshal(body, &response); err != nil {
		return nil, schemas.ErrorResponse(500, "Error al parsear la respuesta de ARCA", err)
	}

	if len(response.Body.Response.Errors.Err) > 0 {
		errMsg := ""
		for _, e := range response.Body.Response.Errors.Err {
			errMsg += fmt.Sprintf("Error %d: %s\n", e.Code, e.Msg)
		}
		return nil, schemas.ErrorResponse(500, "Error en la respuesta de ARCA", fmt.Errorf("errores de AFIP:\n%s", errMsg))
	}

	if len(response.Body.Response.Results.FeDetResp.FECAEDetResponse) == 0 {
		return nil, schemas.ErrorResponse(500, "Error en la respuesta de ARCA", fmt.Errorf("sin resultados en la respuesta"))
	}

	detResp := response.Body.Response.Results.FeDetResp.FECAEDetResponse[0]

	if detResp.Resultado != "A" {
		obsMsg := ""
		for _, obs := range detResp.Observaciones.Obs {
			obsMsg += fmt.Sprintf("Obs %d: %s\n", obs.Code, obs.Msg)
		}
		return nil, schemas.ErrorResponse(500, "Factura no autorizada por ARCA", fmt.Errorf("resultado: %s\nobservaciones:\n%s", detResp.Resultado, obsMsg))
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("✅ FACTURA AUTORIZADA")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("\n🎫 CAE: %s\n", detResp.CAE)
	fmt.Printf("📅 Vencimiento CAE: %s\n", detResp.CAEFchVto)
	fmt.Printf("📋 Comprobante: %04d-%08d-%08d\n",
		factura.PuntoVenta,
		factura.TipoComprobante,
		detResp.CbteDesde)
	fmt.Printf("💰 Total: $%.2f\n", factura.ImporteTotal)

	if len(detResp.Observaciones.Obs) > 0 {
		fmt.Println("\n⚠️  Observaciones:")
		for _, obs := range detResp.Observaciones.Obs {
			fmt.Printf("   - [%d] %s\n", obs.Code, obs.Msg)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎉 PROCESO COMPLETADO")
	fmt.Println(strings.Repeat("=", 60))

	return &detResp, nil
}