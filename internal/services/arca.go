package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (s *ArcaService) EmitInvoice(user *schemas.AuthenticatedUser, pointSaleID int64, req *schemas.FacturaRequest, isHomo bool) (*schemas.FacturaElectronica, error) {
	isHomo = true
	
	credentials, err := s.ArcaRepository.GetCredentialsArca(user.TenantID)
	if err != nil {
		return nil, err
	}

	emisorResponsabilidad := *credentials.ResponsibilityFrontIVA

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
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al crear cliente WSAA", err)
	}

	invalid := credentials.TokenArca == nil || credentials.SignArca == nil || credentials.ExpireTokenArca == nil
	if invalid {
		vacio := ""
		old := time.Now().Add(-24 * time.Hour)
		credentials.TokenArca = &vacio
		credentials.SignArca = &vacio
		credentials.ExpireTokenArca = &old
	}

	cred, err := schemas.LoadCredentials(*credentials.TokenArca, *credentials.SignArca, *credentials.Cuit, *credentials.ExpireTokenArca)
	if err != nil {
		return nil, schemas.ErrorResponse(500, "Error interno al cargar credenciales", err)
	}

	if time.Now().Add(5 * time.Minute).After(cred.Expiration) {

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

	concept, err := schemas.GetCodeTipoConcepto(*credentials.Concept)
	if err != nil {
		return nil, schemas.ErrorResponse(400, "código de concepto inválido", errors.New("código de concepto inválido"))
	}
	fecha := time.Now().Format("20060102")

	alicuotas := []schemas.ItemIVA{}
	totalNeto := 0.0
	totalIVA := 0.0
	totalFinal := 0.0
	if emisorResponsabilidad == "responsable_inscripto" {
		for _, alicuota := range req.Items {
			code, err := schemas.GetCodeAlicuotaIVA(alicuota.Codigo)
			if err != nil {
				return nil, schemas.ErrorResponse(400, "código de alicuta inválido", errors.New("código de alicuta inválido"))
			}
			alicuotas = append(alicuotas, schemas.ItemIVA{
				Codigo:        code,
				BaseImponible: alicuota.BaseImponible,
				Importe:       alicuota.Importe,
			})

			totalNeto += alicuota.BaseImponible
			totalIVA += alicuota.Importe
		}
		totalFinal = totalNeto + totalIVA
	} else if emisorResponsabilidad == "monotributo" {
		for _, alicuota := range req.Items {
			totalNeto += alicuota.Importe
		}
		totalFinal = totalNeto + totalIVA
	}

	totalFinal = math.Round(totalFinal*100) / 100
	totalNeto = math.Round(totalNeto*100) / 100
	totalIVA = math.Round(totalIVA*100) / 100

	tipoDoc, err := schemas.GetCodeTipoDocumento(req.TipoDocumento)
	if err != nil {
		return nil, schemas.ErrorResponse(400, "código de tipo de documento inválido", errors.New("código de tipo de documento inválido"))
	}

	condition, err := schemas.GetCodeConditionIVA(req.CondicionIVA)
	if err != nil {
		return nil, schemas.ErrorResponse(400, "código de condición fiscal inválido", errors.New("código de condición fiscal inválido"))
	}

	factura := &schemas.Factura{
		PuntoVenta:      int(pointSaleID),
		TipoDocumento:   tipoDoc,
		NumeroDocumento: req.NumeroDocumento,
		Concepto:        concept,
		Fecha:           fecha,
		ImporteNeto:     totalNeto,
		CondicionIVA: condition,
		ImporteIVA: totalIVA,
		ImporteTotal: totalFinal,
		Alicuotas:    alicuotas,
		MonedaId:    "PES",
		MonedaCotiz: 1,
	}

	switch emisorResponsabilidad {
	case "monotributo":
		factura.TipoComprobante = 11
		factura.ImporteIVA = 0
		factura.Alicuotas = nil

	case "responsable_inscripto":
		condition, _ := schemas.GetCodeConditionIVA(req.CondicionIVA)
		if condition == 1 || condition == 4 {
			factura.TipoComprobante = 1
		} else if condition == 5 || condition == 6 {
			factura.TipoComprobante = 6
		} else {
			return nil, schemas.ErrorResponse(400, "Condición de IVA del receptor no soportada", errors.New("Condición de IVA del receptor no soportada"))
		}

	default:
		return nil, schemas.ErrorResponse(400, "Responsabilidad del emisor no válida", errors.New("Responsabilidad del emisor no válida"))
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

	cae, err := strconv.ParseInt(fecae.CAE, 10, 64)
	urlQR := UrlQR(schemas.DatosQR{
			Ver:        1,
			Fecha:      time.Now().Format("2006-01-02"),
			Cuit:       cuit,
			PtoVta:     factura.PuntoVenta,
			TipoCmp:    factura.TipoComprobante,
			NroCmp:     int(factura.NumeroDesde),
			Importe:    factura.ImporteTotal,
			Moneda:     factura.MonedaId,
			Ctz:        1,
			TipoDocRec: factura.TipoDocumento,
			NroDocRec:  factura.NumeroDocumento,
			TipoCodAut: "E",
			CodAut:     cae,
		})
	if urlQR == nil {
		return nil, schemas.ErrorResponse(500, "Error al generar QR", errors.New("Error al generar QR"))
	}

	facturaResponse := &schemas.FacturaElectronica{
		TipoComprobante: factura.TipoComprobante,
		PuntoVenta:      factura.PuntoVenta,
		Numero:          factura.NumeroDesde,
		Fecha:           fecha,
		Concepto:        concept,

		EmisorCUIT:         cuit,
		EmisorNombre:       *credentials.BusinessName,
		RazonSocial:        *credentials.SocialReason,
		IngresosBrutos:     *credentials.GrossIncome,
		InicioActividades:  *credentials.StartActivities,
		CondicionIVAEmisor: *credentials.ResponsibilityFrontIVA,
		DomicilioEmisor:    *credentials.Address,

		ReceptorCUIT:         req.NumeroDocumento,
		ReceptorNombre:       req.Nombre,
		CondicionIVAReceptor: req.CondicionIVA,
		ReceptorDomicilio:    req.Domicilio,

		ImporteNeto:   factura.ImporteNeto,
		ImporteIVA:    factura.ImporteIVA,
		ImporteExento: factura.ImporteExento,

		Items: req.Items,

		CAE:            cae,
		CAEVencimiento: fecae.CAEFchVto,
		ImporteTotal:   factura.ImporteTotal,
		URL_QR: *urlQR,
	}

	return facturaResponse, nil
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

