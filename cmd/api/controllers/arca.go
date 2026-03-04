package controllers

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ArcaEmitInvoice godoc
//
//	@Summary		ArcaEmitInvoice
//	@Description	### Emitir Factura
//	@Tags			Arca
//	@Accept			json
//	@Produce		json
//	@Param			income_sale_id	path		int						true	"ID de la venta"
//	@Param			homo			query		bool					false	"Homologación, poner en true para test no mandar nada o false para produccion"
//	@Param			datos			body		schemas.FacturaRequest	true	"datos de la factura"
//	@Success		200				{object}	schemas.Response{body=schemas.FacturaElectronica}
//	@Router			/api/v1/arca/emit_invoice/{income_sale_id} [post]
func (a *ArcaController) ArcaEmitInvoice(c *fiber.Ctx) error {
	incomeSaleID := c.Params("income_sale_id")
	idint, err := validator.IdValidate(incomeSaleID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	var homo bool
	if c.Query("homo") == "true" {
		homo = true
	} else {
		homo = false
	}

	var factReq schemas.FacturaRequest
	if err := validator.ValidateRequest(c, &factReq); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	pointSaeID := c.Locals("point_sae_id").(int64)

	factura, err := a.ArcaService.EmitInvoice(user, pointSaeID, idint, &factReq, homo)
	if err != nil {
		if errResp, ok := err.(*schemas.ErrorStruc); ok {
			if errResp.StatusCode == 206 {
				log.Err(err).Msgf("Error: %s", errResp.Err.Error())
				return c.Status(errResp.StatusCode).JSON(schemas.Response{
					Status:  true,
					Body:    factura,
					Message: errResp.Message,
				})
			}
		}
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    factura,
		Message: "Factura emitida con éxito",
	})
}

// ArcaGenerateKey godoc
//
//	@Summary		ArcaGenerateKey
//	@Description	Generar Key para arca
//	@Tags			Arca
//	@Accept			json
//	@Produce		json
//	@Param			datos	body		schemas.KeyRequest	true	"datos de la factura"
//	@Success		200		{object}	schemas.Response{body=schemas.KeyResponse}
//	@Router			/api/v1/arca/generate_key [post]
func (a *ArcaController) ArcaGenerateKey(c *fiber.Ctx) error {
	var keyReq schemas.KeyRequest
	if err := validator.ValidateRequest(c, &keyReq); err != nil {
		return schemas.HandleError(c, err)
	}

	// privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	// if err != nil {
	// 	log.Err(err).Msg("Error al generar la clave privada")
	// 	return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
	// 		Status:  false,
	// 		Body:    nil,
	// 		Message: "Error al generar la clave privada",
	// 	})
	// }

	// // 3. Codificar Clave Privada a formato PEM (en memoria)
	// privKeyBuf := new(bytes.Buffer)
	// pem.Encode(privKeyBuf, &pem.Block{
	// 	Type:  "RSA PRIVATE KEY",
	// 	Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	// })
	privKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Err(err).Msg("Error al generar la clave privada")
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Error al generar la clave privada",
		})
	}

	// 3. Codificar Clave Privada a formato PEM (PKCS#8)
	privKeyBuf := new(bytes.Buffer)

	// Cambio clave: MarshalPKCS8PrivateKey devuelve la clave en formato genérico PKCS#8
	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		log.Err(err).Msg("Error al serializar clave a PKCS#8")
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
			Status:  false,
			Message: "Error al procesar formato de clave",
		})
	}

	pem.Encode(privKeyBuf, &pem.Block{
		Type:  "PRIVATE KEY", // Nota: Sin la palabra "RSA"
		Bytes: pkcs8Bytes,
	})

	// 4. Crear la solicitud del CSR
	// El serialNumber en el subject requiere manejo específico en x509
	subj := pkix.Name{
		Country:      []string{"AR"},
		Organization: []string{keyReq.BusinessName},
		CommonName:   "noagestion",
		SerialNumber: "CUIT " + keyReq.Cuit,
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	// 5. Generar el CSR
	csrBytes, err := x509.CreateCertificateRequest(rand.Reader, &template, privKey)
	if err != nil {
		log.Err(err).Msg("Error al generar el CSR")
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Error al generar el CSR",
		})
	}

	// 6. Codificar CSR a formato PEM (en memoria)
	csrBuf := new(bytes.Buffer)
	pem.Encode(csrBuf, &pem.Block{
		Type:  "CERTIFICATE REQUEST",
		Bytes: csrBytes,
	})

	// --- RESULTADOS ---
	// Ahora tienes todo en buffers, puedes enviarlos por API, guardarlos en BD, etc.
	fmt.Println("--- CLAVE PRIVADA (PEM) ---")
	fmt.Println(privKeyBuf.String())

	fmt.Println("--- CSR (PEM) ---")
	fmt.Println(csrBuf.String())

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status: true,
		Body: &schemas.KeyResponse{
			Key:         privKeyBuf.String(),
			Certificate: csrBuf.String(),
		},
		Message: "Keys emitidos con éxito",
	})
}
