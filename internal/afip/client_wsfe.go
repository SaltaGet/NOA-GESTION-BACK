package afip

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
)

type FECAERequest struct {
    XMLName xml.Name `xml:"soap:Envelope"`
    SOAP    string   `xml:"xmlns:soap,attr"`
    XSI     string   `xml:"xmlns:xsi,attr"`
    XSD     string   `xml:"xmlns:xsd,attr"`
    Body    struct {
        FECAESolicitar struct {
            Xmlns string `xml:"xmlns,attr"`
            Auth  struct {
                Token string `xml:"Token"`
                Sign  string `xml:"Sign"`
                Cuit  int64  `xml:"Cuit"`
            } `xml:"Auth"`
            FeCAEReq struct {
                FeCabReq struct {
                    CantReg   int `xml:"CantReg"`
                    PtoVta    int `xml:"PtoVta"`
                    CbteTipo  int `xml:"CbteTipo"`
                } `xml:"FeCabReq"`
                FeDetReq struct {
                    FECAEDetRequest []FECAEDetRequest `xml:"FECAEDetRequest"`
                } `xml:"FeDetReq"`
            } `xml:"FeCAEReq"`
        } `xml:"FECAESolicitar"`
    } `xml:"soap:Body"`
}

type FECAEDetRequest struct {
    Concepto    int     `xml:"Concepto"`
    DocTipo     int     `xml:"DocTipo"`
    DocNro      int64   `xml:"DocNro"`
    CbteDesde   int     `xml:"CbteDesde"`
    CbteHasta   int     `xml:"CbteHasta"`
    CbteFch     string  `xml:"CbteFch"` // YYYYMMDD
    ImpTotal    float64 `xml:"ImpTotal"`
    ImpTotConc  float64 `xml:"ImpTotConc"`
    ImpNeto     float64 `xml:"ImpNeto"`
    ImpOpEx     float64 `xml:"ImpOpEx"`
    ImpIVA      float64 `xml:"ImpIVA"`
    ImpTrib     float64 `xml:"ImpTrib"`
    MonId       string  `xml:"MonId"` // PES = Pesos
    MonCotiz    float64 `xml:"MonCotiz"`
}

type FECAEResponse struct {
    XMLName xml.Name `xml:"Envelope"`
    Body    struct {
        FECAESolicitarResponse struct {
            FECAESolicitarResult struct {
                FeDetResp struct {
                    FECAEDetResponse []struct {
                        CAE       string `xml:"CAE"`
                        CAEFchVto string `xml:"CAEFchVto"`
                        CbteFch   string `xml:"CbteFch"`
                        Resultado string `xml:"Resultado"`
                    } `xml:"FECAEDetResponse"`
                } `xml:"FeDetResp"`
                Errors struct {
                    Err []struct {
                        Code int    `xml:"Code"`
                        Msg  string `xml:"Msg"`
                    } `xml:"Err"`
                } `xml:"Errors"`
            } `xml:"FECAESolicitarResult"`
        } `xml:"FECAESolicitarResponse"`
    } `xml:"Body"`
}

func EmitirFactura(token, sign string, cuit int64, factura FECAEDetRequest, production bool) (*FECAEResponse, error) {
    request := FECAERequest{
        SOAP: "http://schemas.xmlsoap.org/soap/envelope/",
        XSI:  "http://www.w3.org/2001/XMLSchema-instance",
        XSD:  "http://www.w3.org/2001/XMLSchema",
    }
    
    request.Body.FECAESolicitar.Xmlns = "http://ar.gov.afip.dif.FEV1/"
    request.Body.FECAESolicitar.Auth.Token = token
    request.Body.FECAESolicitar.Auth.Sign = sign
    request.Body.FECAESolicitar.Auth.Cuit = cuit
    
    request.Body.FECAESolicitar.FeCAEReq.FeCabReq.CantReg = 1
    request.Body.FECAESolicitar.FeCAEReq.FeCabReq.PtoVta = 1
    request.Body.FECAESolicitar.FeCAEReq.FeCabReq.CbteTipo = 11 // Factura C
    
    request.Body.FECAESolicitar.FeCAEReq.FeDetReq.FECAEDetRequest = []FECAEDetRequest{factura}
    
    xmlData, _ := xml.MarshalIndent(request, "", "  ")
    
    url := "https://wswhomo.afip.gov.ar/wsfev1/service.asmx" // Testing
    if production {
        url = "https://servicios1.afip.gov.ar/wsfev1/service.asmx"
    }
    
    resp, err := http.Post(url, "text/xml; charset=utf-8", bytes.NewBuffer(xmlData))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    var response FECAEResponse
    xml.Unmarshal(body, &response)
    
    return &response, nil
}