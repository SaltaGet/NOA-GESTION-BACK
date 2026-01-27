package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (c *CredentialService) CredentialGetMPToken(tenantID int64) (*schemas.CredentialMPTokenResponse, error) {
	return c.CredentialRepository.CredentialGetMPToken(tenantID)
}

func (c *CredentialService) CredentialSetMPToken(tenantID int64, request *schemas.CredentialMPTokenRequest) (*string, error) {
	client := &http.Client{Timeout: 5 * time.Second}

	req, _ := http.NewRequest("GET", "https://api.mercadopago.com/users/me", nil)
	req.Header.Set("Authorization", "Bearer "+request.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error de red de mercadopago, intentelo nuevamente", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, schemas.ErrorResponse(resp.StatusCode, resp.Status, errors.New("El access token esta expirado"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, schemas.ErrorResponse(resp.StatusCode, "error al leer respuesta", err)
	}

	var dataJson map[string]interface{}
	if err := json.Unmarshal(body, &dataJson); err != nil {
		return nil, schemas.ErrorResponse(500, "Error al leer respuesta", err)
	}
	fmt.Println(dataJson)

	var data schemas.MPUserResponse
	if err := json.Unmarshal(body, &data); err != nil {
	// if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, schemas.ErrorResponse(500, "Error al leer respuesta", err)
	}

	// --- LÓGICA DE VALIDACIÓN ---

	// 1. Validar que la cuenta esté activa (Sin deudas bloqueantes o problemas legales)
	if data.Status.SiteStatus != "active" {
		return nil, schemas.ErrorResponse(401, "la cuenta de Mercado Pago no está activa", fmt.Errorf("la cuenta de Mercado Pago no está activa"))
	}

	// 2. Validar tasa de reclamos (Opcional: Ejemplo más del 10%)
	if data.SellerReputation.Metrics.Claims.Rate > 0.10 {
		return nil, schemas.ErrorResponse(401, "el vendedor tiene una tasa de reclamos demasiado alta", fmt.Errorf("el vendedor tiene una tasa de reclamos demasiado alta (%.2f%%)", data.SellerReputation.Metrics.Claims.Rate*100))
	}

	if data.Status.Billing.Allow == false {
		return nil, schemas.ErrorResponse(401, "la cuenta tiene restricciones de facturación", fmt.Errorf("la cuenta tiene restricciones de facturación"))
	}

	message := data.Recommendations() 

	return message, c.CredentialRepository.CredentialSetMPToken(tenantID, request)
}

func (c *CredentialService) CredentialGetArca(tenantID int64) (*schemas.CredentialArcaResponse, error) {
	return c.CredentialRepository.CredentialGetArca(tenantID)
}

func (c *CredentialService) CredentialSetArca(tenantID int64, request *schemas.CredentialArcaRequest) error {
	return c.CredentialRepository.CredentialSetArca(tenantID, request)
}