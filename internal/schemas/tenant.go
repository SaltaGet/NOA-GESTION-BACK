package schemas

import (
	"time"
)

type TenantCreate struct {
	Name       string `json:"name" validate:"required"`
	Identifier string `json:"identifier" validate:"required"`
	Address    string `json:"address" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	CuitPdv    string `json:"cuit_pdv" validate:"required"`
	PlanID     int64  `json:"plan_id" validate:"required"`
}

type TenantUpdate struct {
	ID      string `json:"id" validate:"required"`
	Name    string `json:"name" validate:"required"`
	Address string `json:"address" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
}

type TenantResponse struct {
	ID                     int64     `json:"id" example:"1"`
	Name                   string    `json:"name" example:"Mi tienda"`
	Address                string    `json:"address" example:"mi casa 123"`
	Phone                  string    `json:"phone" example:"3884123456"`
	Email                  string    `json:"email" example:"mitienda@gmail.com"`
	IsActive               bool      `json:"is_active" example:"true"`
	Expiration             time.Time `json:"expiration" example:"2023-01-01"`
	CuitPDV                string    `json:"cuit" example:"12345678909"`
	ResponsabilityFrontIVA *string   `json:"responsability_front_iva" example:"responsable_inscripto | monotributo | null"`
	CreatedAt              time.Time `json:"created_at" example:"2023-01-01 12:00:00"`
	UpdatedAt              time.Time `json:"updated_at" example:"2023-01-01 12:00:00"`
}

type TenantUserCreate struct {
	TenantCreate TenantCreate `json:"tenant_create" validate:"required"`
	UserCreate   UserCreate   `json:"user_create" validate:"required"`
}

type TenantUpdateExpiration struct {
	ID         int64  `json:"id" validate:"required"`
	Expiration string `json:"expiration" validate:"required,datetime=2006-01-02" example:"2023-01-01"`
}

type TenantUpdateTerms struct {
	AcceptedTerms bool      `json:"accepted_terms" validate:"required"`
	IP            string    `json:"ip" validate:"required"`
	DateAccepted  time.Time `json:"date_aceept" validate:"required"`
}

type TenantUpdateSettings struct {
	Title          *string `json:"title,omitempty" example:"Mi tienda"`
	Slogan         *string `json:"slogan,omitempty" example:"Mi tienda"`
	PrimaryColor   *string `json:"primary_color,omitempty" example:"#FF0000"`
	SecondaryColor *string `json:"secondary_color,omitempty" example:"#FF0000"`
	Phone          *string `json:"phone,omitempty" example:"+54 11 1234-5678"`
}

type TenantUpdateSettingsWithTenant struct {
	Title          *string `json:"title,omitempty" example:"Mi tienda"`
	Slogan         *string `json:"slogan,omitempty" example:"Mi tienda"`
	PrimaryColor   *string `json:"primary_color,omitempty" example:"#FF0000"`
	SecondaryColor *string `json:"secondary_color,omitempty" example:"#FF0000"`
	TenantID       int64   `json:"tenant_id" validate:"required"`
}

type TenantSettingsResponse struct {
	Logo           *string `json:"logo" example:"logo_uuid"`
	FrontPage      *string `json:"front_page" example:"front_page_uuid"`
	Title          *string `json:"title" example:"Mi tienda"`
	Slogan         *string `json:"slogan" example:"Mi tienda"`
	PrimaryColor   *string `json:"primary_color" example:"#FF0000"`
	SecondaryColor *string `json:"secondary_color" example:"#FF0000"`
	Phone          *string `json:"phone" example:"555-5555"`
}
