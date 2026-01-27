package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	// "github.com/rs/zerolog/log"
)

func (s *ArcaService) EmitInvoice(user *schemas.AuthenticatedUser, pointSaleID int64, req *schemas.FacturaRequest, isHomo bool) (*schemas.FacturaElectronica, error) {
	credentials, err := s.ArcaRepository.GetCredentialsArca(user.TenantID)
	if err != nil {
		return nil, err
	}

	cuit, err := strconv.ParseInt(*credentials.Cuit, 10, 64)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al parsear CUIT, revisa que el CUIT sea correcto", err)
	}

	wsaa, err := schemas.NewWSAA(schemas.WSAAConfig{
		Homologacion: isHomo,
		CertFile:     *credentials.ArcaCertificate,
		KeyFile:      *credentials.ArcaKey,
		CUIT:         cuit,
		Service:      "wsfe",
	})

	cred, err := schemas.LoadCredentials(*credentials.TokenArca, *credentials.SignArca, *credentials.Cuit, *credentials.ExpireTokenArca)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al cargar credenciales", err)
	}

	if time.Now().After(cred.Expiration) {

		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error interno al crear cliente WSAA", err)
		}

		ticketXML, err := wsaa.CreateTicketXML()
		if err != nil {
			return nil, schemas.ErrorResponse(500, "error creando XML ticket", err)
		}

		signedCMS, err := wsaa.SignWithOpenSSL(ticketXML)
		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error firmando CMS", err)
		}

		response, err := s.ArcaRepository.SendToWSAA(wsaa, signedCMS)
		if err != nil {
			return nil, err
		}

		newCred, err := wsaa.ParseResponse(response)
		if err != nil {
			return nil, schemas.ErrorResponse(500, "Error parseando respuesta del WSAA", err)
		}

		err = s.ArcaRepository.SetTokenSignArca(&schemas.CredentialsValidation{
			CUIT:       cuit,
			Token:      newCred.Token,
			Sign:       newCred.Sign,
			Expiration: newCred.Expiration,
		})
		if err != nil {
			return nil, err
		}

		cred.Token = newCred.Token
		cred.Sign = newCred.Sign
		cred.Expiration = newCred.Expiration
	}

	var wsfe *schemas.WSFEClient
	wsfe = schemas.NewWSFEClient(cred.Token, cred.Sign, cuit, isHomo)

	factura := &schemas.Factura{
        PuntoVenta:       int(pointSaleID),
        TipoDocumento:    req.TipoDocumento,
        NumeroDocumento:  req.NumeroDocumento,
        Concepto:         req.Concepto,
        Fecha:            time.Now().Format("20060102"),
        ImporteNeto:      req.ImporteNeto,
        ImporteNoGravado: req.ImporteNoGravado,
        ImporteExento:    req.ImporteExento,
        ImporteIVA:       req.ImporteIVA,
        ImporteTributos:  req.ImporteTributos,
        ImporteTotal:     req.ImporteTotal,
        Alicuotas:        req.Alicuotas,
        Tributos:         req.Tributos,
        MonedaId:         req.MonedaId,
        MonedaCotiz:      req.MonedaCotiz,
    }

	// responsability := *credentials.ResponsibilityFrontIVA
	// if responsability == "monotributo" {
	// 	factura.TipoComprobante = 11 // Monotributo
	// } else if responsability == "responsable_inscripto" {
	// 	if factura.CondicionIVA == 5 {
	// 		factura.TipoComprobante = 6 // Factura B
	// 	} else {
	// 		factura.TipoComprobante = 1 // Factura A
	// 	}
	// } else {
	// 	return nil, schemas.ErrorResponse(400, "Responsabilidad frente al IVA no valida", errors.New("Responsabilidad frente al IVA no valida"))
	// }
	emisorResponsabilidad := *credentials.ResponsibilityFrontIVA // Debería venir de tu DB
    
    // Suponemos que req.CondicionIVA es la condición fiscal del RECEPTOR (Cliente)
    // 1: RI, 4: Monotributista, 5: Consumidor Final, 6: Exento
    switch emisorResponsabilidad {
    case "monotributo":
        // El Monotributista SIEMPRE emite Factura C
        factura.TipoComprobante = 11 
        // Limpieza de seguridad: El monotributista no debe enviar IVA
        factura.ImporteIVA = 0
        factura.Alicuotas = nil 

    case "responsable_inscripto":
        if req.CondicionIVAReceptor == 1 || req.CondicionIVAReceptor == 4 {
            // RI a RI o RI a Monotributista = Factura A
            factura.TipoComprobante = 1
        } else if req.CondicionIVAReceptor == 5 || req.CondicionIVAReceptor == 6 {
            // RI a Consumidor Final o Exento = Factura B
            factura.TipoComprobante = 6
        } else {
            return nil, schemas.ErrorResponse(400, "Condición de IVA del receptor no soportada", nil)
        }

    default:
        return nil, schemas.ErrorResponse(400, "Responsabilidad del emisor no válida", nil)
    }

	lastInvoice, err := s.ArcaRepository.GetLastestInvoice(wsfe, int(pointSaleID), factura.TipoComprobante)
	if err != nil {
		return nil, err
	}

	factura.NumeroDesde = lastInvoice + 1
	factura.NumeroHasta = lastInvoice + 1

	fecae, err := s.ArcaRepository.EmitInvoice(wsfe, factura)
	if err != nil {
		return nil, err
	}

	return fecae, nil
}

func UrlQR(data schemas.DatosQR) *string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error generando JSON:", err)
		return nil
	}

	encodedData := base64.StdEncoding.EncodeToString(jsonData)

	urlQR := "https://www.afip.gob.ar/fe/qr/?p=" + encodedData

	return &urlQR
	}
	// data := DatosQR{
	// 	Ver:        1,
	// 	Fecha:      "2026-01-26",
	// 	Cuit:       20233073572,
	// 	PtoVta:     5,
	// 	TipoCmp:    1,
	// 	NroCmp:     147774,
	// 	Importe:    4212.00,
	// 	Moneda:     "PES",
	// 	Ctz:        1,
	// 	TipoDocRec: 80,
	// 	NroDocRec:  30714150536,
	// 	TipoCodAut: "E",
	// 	CodAut:     74195192801456,
	// }

	// // 2. Convertir a JSON
	// jsonData, err := json.Marshal(data)
	// if err != nil {
	// 	fmt.Println("Error generando JSON:", err)
	// 	return
	// }

	// // 3. Codificar en Base64
	// encodedData := base64.StdEncoding.EncodeToString(jsonData)

	// // 4. Armar la URL final
	// urlQR := "https://www.afip.gob.ar/fe/qr/?p=" + encodedData

	// fmt.Println("URL para el QR:")
	// fmt.Println(urlQR)

// 	factura := &schemas.Factura{
// 		PuntoVenta:       1,
// 		TipoComprobante:  6,
// 		NumeroDesde:      5,
// 		NumeroHasta:      5,
// 		TipoDocumento:    99,
// 		NumeroDocumento:  0,
// 		CondicionIVA:     5, // Consumidor Final
// 		Concepto:         1,
// 		Fecha:            time.Now().Format("20060102"),
// 		ImporteNeto:      1000.00,
// 		ImporteNoGravado: 0.00,
// 		ImporteExento:    0.00,
// 		ImporteIVA:       210.00,
// 		ImporteTributos:  0.00,
// 		ImporteTotal:     1210.00,
// 		Alicuotas: []schemas.ItemIVA{
// 			{
// 				Codigo:        5,
// 				BaseImponible: 1000.00,
// 				Importe:       210.00,
// 			},
// 		},
// 		MonedaId:    "PES",
// 		MonedaCotiz: 1,
// 	}

// 	fmt.Printf("📄 Factura B #%d-%08d\n", factura.PuntoVenta, factura.NumeroDesde)
// 	fmt.Printf("   Cliente: Consumidor Final\n")
// 	fmt.Printf("   Neto: $%.2f\n", factura.ImporteNeto)
// 	fmt.Printf("   IVA 21%%: $%.2f\n", factura.ImporteIVA)
// 	fmt.Printf("   Total: $%.2f\n\n", factura.ImporteTotal)

// 	fmt.Println("PASO 3: Solicitando CAE a AFIP...")
// 	fmt.Println(strings.Repeat("-", 60))

// 	wsfe := NewWSFEClient(creds.Token, creds.Sign, config.CUIT, config.Homologacion)

// 	resultado, err := wsfe.SolicitarCAE(factura)
// 	if err != nil {
// 		log.Fatalf("❌ Error: %v", err)
// 	}

// 	fmt.Println("\n" + strings.Repeat("=", 60))
// 	fmt.Println("✅ FACTURA AUTORIZADA")
// 	fmt.Println(strings.Repeat("=", 60))
// 	fmt.Printf("\n🎫 CAE: %s\n", resultado.CAE)
// 	fmt.Printf("📅 Vencimiento CAE: %s\n", resultado.CAEFchVto)
// 	fmt.Printf("📋 Comprobante: %04d-%08d-%08d\n",
// 		factura.PuntoVenta,
// 		factura.TipoComprobante,
// 		resultado.CbteDesde)
// 	fmt.Printf("💰 Total: $%.2f\n", factura.ImporteTotal)

// 	if len(resultado.Observaciones.Obs) > 0 {
// 		fmt.Println("\n⚠️  Observaciones:")
// 		for _, obs := range resultado.Observaciones.Obs {
// 			fmt.Printf("   - [%d] %s\n", obs.Code, obs.Msg)
// 		}
// 	}

// 	fmt.Println("\n" + strings.Repeat("=", 60))
// 	fmt.Println("🎉 PROCESO COMPLETADO")
// 	fmt.Println(strings.Repeat("=", 60))

// 	fmt.Println("\n🔍 Consultando factura recién emitida...")
// 	info, err := wsfe.ConsultarFactura(factura.TipoComprobante, factura.PuntoVenta, factura.NumeroDesde)
// 	if err != nil {
// 		log.Info().Msgf("❌ No se pudo consultar: %v", err)
// 	} else {
// 		res := info.Body.Response.Result.ResultGet
// 		fmt.Printf("✅ Información recuperada:\n")
// 		fmt.Printf("   CAE: %s\n", res.CodAut)
// 		fmt.Printf("   Fecha: %s\n", res.CbteFch)
// 		fmt.Printf("   Importe: $%.2f\n", res.ImpTotal)
// 		fmt.Printf("   Resultado en AFIP: %s\n", res.Resultado)
// 	}

// 	fmt.Println("🔍 Consultando último número autorizado...")
// 	ultimo, err := wsfe.GetUltimoComprobante(factura.PuntoVenta, factura.TipoComprobante)
// 	if err != nil {
// 		log.Fatalf("❌ Error consultando último número: %v", err)
// 	}

// 	fmt.Printf("✅ Último comprobante autorizado para PtoVta %d, Tipo %d: %d\n", factura.PuntoVenta, factura.TipoComprobante, ultimo)
// }
