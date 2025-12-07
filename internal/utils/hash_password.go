package utils

import (
	"runtime"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	var params = &argon2id.Params{
		Memory:      16 * 1024,
		Iterations:  2,
		Parallelism: uint8(runtime.NumCPU()),
		SaltLength:  16,
		KeyLength:   32,
	}

	hash, err := argon2id.CreateHash(password, params)
	if err != nil {
		return "", schemas.ErrorResponse(500, "Error al crear el hash de la contraseña", err)
	}
	return hash, nil
}

// Verificar una contraseña
func CheckPasswordHash(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	return err == nil && match
}
