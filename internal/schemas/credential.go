package schemas

import (
	"fmt"
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
	ResponsibilityFrontIVA *string `json:"responsibility_front_iva" example:"Responsibility Front IVA"`
	Cuit                   *string `json:"cuit" example:"20-12345678-9"`
	ArcaCertificate        *string `json:"arca_certificate" example:"key.pem"`
	ArcaKey                *string `json:"arca_key" example:"key.pem"`
	ArcaCertificateTest    *string `json:"arca_certificate_test" example:"key.pem"`
	ArcaKeyTest            *string `json:"arca_key_test" example:"key.pem"`
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
	ResponsibilityFrontIVA string `json:"responsibility_front_iva" validate:"required"`
	GrossIncome            string `json:"gross_income" validate:"required"`
	StartActivities        string `json:"start_activities" validate:"required"`
	Cuit                   string `json:"cuit" validate:"required"`
	ArcaCertificate        string `json:"arca_certificate" validate:"required"`
	ArcaKey                string `json:"arca_key" validate:"required"`
}

func (c *CredentialArcaRequest) Validate() error {
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

type MPUserResponse struct {
	ID               int64     `json:"id"`
	Nickname         string    `json:"nickname"`
	Email            string    `json:"email"`
	FirstName        string    `json:"first_name"`
	LastName         string    `json:"last_name"`
	RegistrationDate time.Time `json:"registration_date"` // 2019-01-08...
	CountryID        string    `json:"country_id"`       // "AR"
	SellerExperience string    `json:"seller_experience"` // "NEWBIE"

	// Datos de la empresa/local (Donde sale "Pizzería Metro")
	Company struct {
		BrandName      string `json:"brand_name"`      // "Pizzería Metro"
		SoftDescriptor string `json:"soft_descriptor"` // "PIZZERIAMETRO"
	} `json:"company"`

	// Estado de la cuenta
	Status struct {
		SiteStatus             string `json:"site_status"`             // "active"
		MercadoPagoAccountType string `json:"mercadopago_account_type"` // "personal"
		Billing    struct {
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
				recommendation += recommendation + " Dirección pendiente de verificación, por favor completar la información en Mercado Pago. "
			}
			recommendation += code + " "
		}
		return &recommendation
	}
	return nil
}