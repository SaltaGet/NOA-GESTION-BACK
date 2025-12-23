package afip

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"os"
)

type SOAPEnvelope struct {
    XMLName xml.Name `xml:"soap:Envelope"`
    SOAP    string   `xml:"xmlns:soap,attr"`
    Body    SOAPBody `xml:"soap:Body"`
}

type SOAPBody struct {
    LoginCMS LoginCMS `xml:"loginCms"`
}

type LoginCMS struct {
    XMLName xml.Name `xml:"loginCms"`
    Xmlns   string   `xml:"xmlns,attr"`
    In0     string   `xml:"in0"` // TRA firmado en base64
}

type LoginCMSResponse struct {
    XMLName xml.Name `xml:"Envelope"`
    Body    struct {
        LoginCMSReturn struct {
            Token string `xml:"token"`
            Sign  string `xml:"sign"`
        } `xml:"loginCmsReturn"`
    } `xml:"Body"`
}

func LoginWSAA(tra, signature, certPath string, production bool) (token, sign string, err error) {
    // Leer certificado
    certData, err := os.ReadFile(certPath)
    if err != nil {
        return "", "", err
    }
    
    // Construir CMS (TRA + firma + certificado)
    cms := buildCMS(tra, signature, string(certData))
    
    // SOAP Request
    envelope := SOAPEnvelope{
        SOAP: "http://schemas.xmlsoap.org/soap/envelope/",
        Body: SOAPBody{
            LoginCMS: LoginCMS{
                Xmlns: "https://wsaa.afip.gov.ar/ws/services/LoginCms",
                In0:   cms,
            },
        },
    }
    
    xmlData, _ := xml.Marshal(envelope)
    
    // Endpoint
    url := "https://wsaahomo.afip.gov.ar/ws/services/LoginCms" // Testing
    if production {
        url = "https://wsaa.afip.gov.ar/ws/services/LoginCms"
    }
    
    // Hacer request
    resp, err := http.Post(url, "text/xml", bytes.NewBuffer(xmlData))
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    
    // Parsear respuesta
    var response LoginCMSResponse
    xml.Unmarshal(body, &response)
    
    return response.Body.LoginCMSReturn.Token, 
           response.Body.LoginCMSReturn.Sign, 
           nil
}

func buildCMS(tra, signature, cert string) string {
    // Aquí debes construir el formato CMS (PKCS#7)
    // Es complejo, considera usar crypto/pkcs7 o similar
    // Por simplicidad, muchos usan openssl via exec
    return "" // Implementación completa requiere PKCS#7
}