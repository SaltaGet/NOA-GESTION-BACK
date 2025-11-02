package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/DanielChachagua/GestionCar/pkg/models"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!*"

func GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", models.ErrorResponse(500, "Error al generar contraseña", err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}