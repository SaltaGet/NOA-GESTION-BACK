package afip

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"os"
	"time"
)

type LoginTicketRequest struct {
    XMLName xml.Name `xml:"loginTicketRequest"`
    Version string   `xml:"version,attr"`
    Header  struct {
        UniqueID     int64  `xml:"uniqueId"`
        GenerationTime string `xml:"generationTime"`
        ExpirationTime string `xml:"expirationTime"`
    } `xml:"header"`
    Service string `xml:"service"`
}

func GenerateTRA(service string) (string, error) {
    tra := LoginTicketRequest{
        Version: "1.0",
        Service: service, // "wsfe" para facturación
    }
    
    // Generar ID único
    tra.Header.UniqueID = time.Now().Unix()
    
    // Fechas (válido por 12 horas)
    now := time.Now()
    tra.Header.GenerationTime = now.Format("2006-01-02T15:04:05")
    tra.Header.ExpirationTime = now.Add(12 * time.Hour).Format("2006-01-02T15:04:05")
    
    // Convertir a XML
    xmlData, err := xml.MarshalIndent(tra, "", "  ")
    if err != nil {
        return "", err
    }
    
    return string(xmlData), nil
}

// Firmar el TRA con tu clave privada
func SignTRA(traXML string, privateKeyPath string) (string, error) {
    // Leer clave privada
    keyData, err := os.ReadFile(privateKeyPath)
    if err != nil {
        return "", err
    }
    
    block, _ := pem.Decode(keyData)
    privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
    if err != nil {
        return "", err
    }
    
    rsaKey := privateKey.(*rsa.PrivateKey)
    
    // Calcular hash
    hash := sha256.Sum256([]byte(traXML))
    
    // Firmar
    signature, err := rsa.SignPKCS1v15(rand.Reader, rsaKey, crypto.SHA256, hash[:])
    if err != nil {
        return "", err
    }
    
    return base64.StdEncoding.EncodeToString(signature), nil
}