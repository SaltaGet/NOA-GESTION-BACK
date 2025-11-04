package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func Encrypt(plainText string) (string, error) {
	block, err := aes.NewCipher([]byte(os.Getenv("ENCRYPT_KEY")))
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al encriptar", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al encriptar", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", schemas.ErrorResponse(500, "Error al encriptar", err)
	}

	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func Decrypt(encryptedText string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al desencriptar", err)
	}

	block, err := aes.NewCipher([]byte(os.Getenv("ENCRYPT_KEY")))
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al desencriptar", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al desencriptar", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("dato inválido")
	}

	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al desencriptar", err)
	}

	return string(plainText), nil
}
