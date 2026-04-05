package validator

import (
	"fmt"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func init() {
	// Registro de validaciones personalizadas
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("username", validateUsername)
	validate.RegisterValidation("is_pem_cert", validatePEMCertificate)
	validate.RegisterValidation("is_pem_key", validatePEMKey)
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("campo %s falló por validación: %s", e.FailedField, e.Tag)
}

func ValidateStruct(s any) error {
	err := validate.Struct(s)
	if err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			firstErr := validationErrors[0]

			msgForUser := getFriendlyErrorMessage(firstErr)

			technicalErr := ErrorResponse{
				FailedField: firstErr.StructNamespace(),
				Tag:         firstErr.Tag(),
				Value:       firstErr.Param(),
			}

			return schemas.ErrorResponse(422, msgForUser, technicalErr)
		}
		return schemas.ErrorResponse(500, "Error interno de validación", err)
	}

	return nil
}

func getFriendlyErrorMessage(err validator.FieldError) string {
	field := err.Field()
	if strings.Contains(err.Namespace(), "[") {
        parts := strings.Split(err.Namespace(), ".")
        field = parts[len(parts)-1] 
    }
	
	param := err.Param()

	switch err.Tag() {
	case "is_pem_cert":
		return fmt.Sprintf("El certificado en '%s' no es un formato PEM válido.", field)
	case "is_pem_key":
		return fmt.Sprintf("La clave privada en '%s' no es válida.", field)
	case "datetime":
		friendlyFormat := param
		if param == "2006-01-02" {
			friendlyFormat = "AAAA-MM-DD"
		} else if param == "02-01-2006" {
			friendlyFormat = "DD-MM-AAAA"
		}
		return fmt.Sprintf("El campo '%s' no tiene un formato de fecha válido. Se espera: %s.", field, friendlyFormat)
	case "required":
		return fmt.Sprintf("El campo '%s' es obligatorio.", field)
	case "email":
		return fmt.Sprintf("'%s' no es un correo electrónico válido.", field)
	case "min":
		if err.Kind().String() == "string" {
			return fmt.Sprintf("'%s' debe tener al menos %s caracteres.", field, param)
		}
		return fmt.Sprintf("'%s' debe ser como mínimo %s.", field, param)
	case "max":
		if err.Kind().String() == "string" {
			return fmt.Sprintf("'%s' no puede tener más de %s caracteres.", field, param)
		}
		return fmt.Sprintf("'%s' no puede ser mayor a %s.", field, param)
	case "oneof":
		options := strings.ReplaceAll(param, " ", ", ")
		return fmt.Sprintf("'%s' debe ser uno de los siguientes valores: [%s].", field, options)
	case "eqfield":
		return fmt.Sprintf("'%s' debe ser igual al campo '%s'.", field, param)
	case "len":
		return fmt.Sprintf("'%s' debe tener exactamente %s caracteres.", field, param)
	case "numeric":
		return fmt.Sprintf("'%s' debe contener solo números.", field)
	case "alphanum":
		return fmt.Sprintf("'%s' solo puede contener letras y números.", field)
	case "url":
		return fmt.Sprintf("'%s' debe ser una URL válida.", field)
	case "password":
		return fmt.Sprintf("La contraseña en '%s' es muy débil.", field)
	case "gte":
		if err.Kind().String() == "string" {
			return fmt.Sprintf("El campo '%s' debe tener al menos %s caracteres.", field, param)
		}
		if err.Kind().String() == "slice" || err.Kind().String() == "array" {
			return fmt.Sprintf("Debe seleccionar al menos %s elementos en '%s'.", param, field)
		}
		return fmt.Sprintf("El valor de '%s' debe ser mayor o igual a %s.", field, param)
	case "gt":
		return fmt.Sprintf("El valor de '%s' debe ser mayor a %s.", field, param)
	case "lte":
		return fmt.Sprintf("El valor de '%s' no puede ser mayor a %s.", field, param)
	case "nefield":
		return fmt.Sprintf("El campo '%s' no puede ser igual al campo '%s'.", field, param)
	case "unique": // Muy útil con dive en slices de strings o IDs
		return fmt.Sprintf("Los elementos en '%s' deben ser únicos.", field)
	case "username":
		return fmt.Sprintf("El campo '%s' no es un nombre de usuario válido. tiene que tener el siguiente formato user@dominio", field)
	default:
		return fmt.Sprintf("El campo '%s' no supera la validación de tipo '%s'.", field, err.Tag())
	}
}

func ValidateRequest(c *fiber.Ctx, payload any) error {
	if err := c.BodyParser(payload); err != nil {
		return schemas.ErrorResponse(422, "Cuerpo de petición inválido: "+err.Error(), err)
	}

	return ValidateStruct(payload)
}
