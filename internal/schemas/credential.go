package schemas

import (
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type CredentialMPTokenResponse struct {
	AccessToken     *string `json:"access_token" example:"token_1234567890"`
	AccessTokenTest *string `json:"access_token_test" example:"token_1234567890"`
	TokenEmail      *string `json:"token_email" example:"token_1234567890"`
}

type CredentialArcaResponse struct {
	SocialReason           *string `json:"social_reason" example:"My Company"`
	BusinessName           *string `json:"business_name" example:"My Company"`
	Address                *string `json:"address" example:"Calle 123"`
	ResponsibilityFrontIVA *string `json:"responsibility_front_iva" example:"Responsibility Front IVA"`
	Cuit                   *string `json:"cuit" example:"20-12345678-9"`
	GrossIncome            *string `json:"gross_income" example:"1000000"`
	StartActivities        *string `json:"start_activities" example:"2022-01-01"`
	Concept                *string `json:"concept" example:"productos"`
	ArcaCertificate        *string `json:"arca_certificate" example:"key.pem"`
	ArcaKey                *string `json:"arca_key" example:"key.pem"`
}

type CredentialMPTokenRequest struct {
	AccessToken     string `json:"access_token" validate:"required"`
	AccessTokenTest string `json:"access_token_test" validate:"required"`
	// TokenEmail      string `json:"token_email" example:"token_1234567890"`
}

func (c *CredentialMPTokenRequest) Validate() error {
	validate := validator.New()

	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}

type CredentialArcaRequest struct {
	SocialReason           string `json:"social_reason" validate:"required"`
	BusinessName           string `json:"business_name" validate:"required"`
	Address                string `json:"address" validate:"required"`
	ResponsibilityFrontIVA string `json:"responsibility_front_iva" validate:"required,oneof=responsable_inscripto monotributo"`
	GrossIncome            string `json:"gross_income" validate:"required,numeric,min=7,max=14"`
	StartActivities        string `json:"start_activities" validate:"required,datetime=2006-01-02"`
	Cuit                   string `json:"cuit" validate:"required,numeric,len=11"`
	Concept                string `json:"concept" validate:"required,oneof=productos servicios productos_servicios"`
	ArcaCertificate        string `json:"arca_certificate" validate:"required"`
	ArcaKey                string `json:"arca_key" validate:"required"`
}

func (c *CredentialArcaRequest) Validate() error {
	validate := validator.New()

	// 1. Registramos la validación para el Certificado
	validate.RegisterValidation("is_pem_cert", func(fl validator.FieldLevel) bool {
		block, _ := pem.Decode([]byte(fl.Field().String()))
		return block != nil && block.Type == "CERTIFICATE"
	})

	// 2. Registramos la validación para la Clave Privada
	validate.RegisterValidation("is_pem_key", func(fl validator.FieldLevel) bool {
		block, _ := pem.Decode([]byte(fl.Field().String()))
		// Aceptamos cualquier tipo de clave (PRIVATE KEY, RSA PRIVATE KEY, etc)
		return block != nil && strings.Contains(block.Type, "KEY")
	})

	err := validate.Struct(c)
	if err == nil {
		return nil
	}

	// Manejo de errores de validación
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorResponse(500, "Error interno de validación", err)
	}

	verr := validationErrors[0]
	field := verr.Field()
	tag := verr.Tag()

	var message string
	switch tag {
	case "is_pem_cert":
		message = "El certificado no tiene un formato PEM válido (debe empezar con -----BEGIN CERTIFICATE-----)"
	case "is_pem_key":
		message = "La clave privada no tiene un formato PEM válido (debe empezar con -----BEGIN PRIVATE KEY-----)"
	case "required":
		message = fmt.Sprintf("El campo %s es obligatorio", field)
	case "numeric":
		message = fmt.Sprintf("El campo %s solo debe contener números", field)
	case "len", "min", "max":
		message = fmt.Sprintf("El campo %s tiene una longitud inválida", field)
	case "datetime":
		message = fmt.Sprintf("El campo %s debe tener formato YYYY-MM-DD", field)
	case "oneof":
		message = fmt.Sprintf("El valor del campo %s no es una opción permitida", field)
	default:
		message = fmt.Sprintf("Error en el campo %s: (%s)", field, tag)
	}

	// message := fmt.Sprintf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)

	return ErrorResponse(422, message, fmt.Errorf("%s", message))
}

type MPUserResponse struct {
	ID               int64     `json:"id"`
	Nickname         string    `json:"nickname"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	RegistrationDate time.Time `json:"registration_date"` // 2019-01-08...
	CountryID        string    `json:"country_id"`        // "AR"
	SellerExperience string    `json:"seller_experience"` // "NEWBIE"

	// Datos de la empresa/local (Donde sale "Pizzería Metro")
	Company struct {
		BrandName      string `json:"brand_name"`      // "Pizzería Metro"
		SoftDescriptor string `json:"soft_descriptor"` // "PIZZERIAMETRO"
	} `json:"company"`

	// Estado de la cuenta
	Status struct {
		SiteStatus             string `json:"site_status"`              // "active"
		MercadoPagoAccountType string `json:"mercadopago_account_type"` // "personal"
		Billing                struct {
			Allow bool     `json:"allow"`
			Codes []string `json:"codes"` // Aquí vendrá ["address_pending"]
		} `json:"billing"`
	} `json:"status"`

	// Reputación y métricas de venta
	SellerReputation struct {
		LevelID           interface{} `json:"level_id"` // Usamos interface{} porque puede ser null o string
		PowerSellerStatus interface{} `json:"power_seller_status"`
		Transactions      struct {
			Total     int    `json:"total"`
			Completed int    `json:"completed"`
			Period    string `json:"period"`
		} `json:"transactions"`
		Metrics struct {
			Sales struct {
				Completed int `json:"completed"`
			} `json:"sales"`
			Claims struct {
				Rate float64 `json:"rate"`
			} `json:"claims"`
		} `json:"metrics"`
	} `json:"seller_reputation"`
}

func (c *MPUserResponse) Recommendations() *string {
	if len(c.Status.Billing.Codes) > 0 {
		recommendation := "La cuenta tiene los siguientes problemas: "
		for _, code := range c.Status.Billing.Codes {
			if code == "address_pending" {
				recommendation += " Dirección pendiente de verificación, por favor completar la información en Mercado Pago. "
			}
			recommendation += code + " "
		}
		return &recommendation
	}
	return nil
}
