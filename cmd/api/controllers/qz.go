package controllers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"sync"
)

// Variables compartidas — cargadas una sola vez en memoria (lazy singleton).
var (
	qzOnce       sync.Once
	qzPrivateKey *rsa.PrivateKey
	qzPublicPEM  []byte
	qzLoadErr    error
)

// QzResponse representa la respuesta del endpoint QZ Tray.
type QzResponse struct {
	Certificate string `json:"certificate"` // Contenido del qz-public.pem
	Signature   string `json:"signature"`   // Firma RSA SHA-512 en Base64
}

// loadQZKeys carga las claves desde disco una sola vez.
func loadQZKeys() error {
	qzOnce.Do(func() {
		privatePath := os.Getenv("QZ_PRIVATE_KEY_PATH")
		publicPath := os.Getenv("QZ_PUBLIC_KEY_PATH")

		if privatePath == "" {
			qzLoadErr = errors.New("QZ_PRIVATE_KEY_PATH no configurado")
			return
		}
		if publicPath == "" {
			qzLoadErr = errors.New("QZ_PUBLIC_KEY_PATH no configurado")
			return
		}

		// Leer clave privada
		privBytes, err := os.ReadFile(privatePath)
		if err != nil {
			qzLoadErr = err
			return
		}

		block, _ := pem.Decode(privBytes)
		if block == nil {
			qzLoadErr = errors.New("no se pudo decodificar el PEM de la clave privada")
			return
		}

		privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			// Fallback a PKCS1
			privKey1, err1 := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err1 != nil {
				qzLoadErr = err
				return
			}
			qzPrivateKey = privKey1
		} else {
			rsaKey, ok := privKey.(*rsa.PrivateKey)
			if !ok {
				qzLoadErr = errors.New("la clave privada no es RSA")
				return
			}
			qzPrivateKey = rsaKey
		}

		// Leer certificado público
		pubBytes, err := os.ReadFile(publicPath)
		if err != nil {
			qzLoadErr = err
			return
		}
		qzPublicPEM = pubBytes
	})

	return qzLoadErr
}
