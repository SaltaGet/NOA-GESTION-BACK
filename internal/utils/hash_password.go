package utils

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
    hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
    if err != nil {
        return "", models.ErrorResponse(500, "Error al crear el hash de la contraseña", err)
    }
    return hash, nil
}

// Verificar una contraseña
func CheckPasswordHash(password, hash string) bool {
    match, err := argon2id.ComparePasswordAndHash(password, hash)
    return err == nil && match
}