package validator

import (
	"encoding/pem"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// Mantenemos tu lógica de seguridad fuerte
	if len(password) < 8 {
		return false
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)
	return hasUpper && hasNumber && hasSpecial
}

func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	return regexp.MustCompile(`^[^@]+@[^.]+$`).MatchString(username)
}

func validatePEMCertificate(fl validator.FieldLevel) bool {
	block, _ := pem.Decode([]byte(fl.Field().String()))
	return block != nil && block.Type == "CERTIFICATE"
}

func validatePEMKey(fl validator.FieldLevel) bool {
	block, _ := pem.Decode([]byte(fl.Field().String()))
	return block != nil && strings.Contains(block.Type, "KEY")
}