package afip

import (
    "fmt"
    "log"
    "time"
)

func init() {
    // 1. Generar TRA
    tra, err := GenerateTRA("wsfe")
    if err != nil {
        log.Fatal(err)
    }
    
    // 2. Firmar TRA
    signature, err := SignTRA(tra, "./private.key")
    if err != nil {
        log.Fatal(err)
    }
    
    // 3. Autenticar en WSAA
    token, sign, err := LoginWSAA(tra, signature, "./cert.crt", false)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Token obtenido:", token[:50]+"...")
    
    // 4. Emitir factura
    factura := FECAEDetRequest{
        Concepto:   1,  // Productos
        DocTipo:    99, // Consumidor Final
        DocNro:     0,
        CbteDesde:  1,
        CbteHasta:  1,
        CbteFch:    time.Now().Format("20060102"),
        ImpTotal:   1210.00,
        ImpTotConc: 0,
        ImpNeto:    1000.00,
        ImpOpEx:    0,
        ImpIVA:     210.00,
        ImpTrib:    0,
        MonId:      "PES",
        MonCotiz:   1,
    }
    
    response, err := EmitirFactura(token, sign, 20123456789, factura, false)
    if err != nil {
        log.Fatal(err)
    }
    
    // 5. Procesar respuesta
    if len(response.Body.FECAESolicitarResponse.FECAESolicitarResult.FeDetResp.FECAEDetResponse) > 0 {
        det := response.Body.FECAESolicitarResponse.FECAESolicitarResult.FeDetResp.FECAEDetResponse[0]
        fmt.Printf("CAE: %s\nVencimiento: %s\n", det.CAE, det.CAEFchVto)
    }
}