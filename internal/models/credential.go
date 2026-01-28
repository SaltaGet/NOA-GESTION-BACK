package models

import "time"

type Credential struct {
	ID       int64 `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID int64 `gorm:"uniqueIndex;not null" json:"tenant_id"`

	AccessTokenMP     *string `gorm:"varchar(255)" json:"access_token_mp"`
	AccessTokenTestMP *string `gorm:"varchar(255)" json:"access_token_test_mp"`

	SocialReason           *string `gorm:"varchar(255)" json:"social_reason"`
	BusinessName           *string `gorm:"varchar(255)" json:"business_name"`
	Address                *string `gorm:"varchar(255)" json:"address"`
	ResponsibilityFrontIVA *string `gorm:"varchar(255)" json:"responsibility_front_iva"`
	GrossIncome            *string `gorm:"varchar(255)" json:"gross_income"`
	StartActivities        *string `gorm:"date" json:"start_activities"`
	Cuit                   *string `gorm:"varchar(255);unique" json:"cuit"`
	Concept                *string `gorm:"varchar(255)" json:"concept"`
	ArcaCertificate        *string `gorm:"type:LONGTEXT" json:"arca_certificate"`
	ArcaKey                *string `gorm:"type:LONGTEXT" json:"arca_key"`

	TokenArca       *string    `gorm:"varchar(255)" json:"token_arca"`
	SignArca        *string    `gorm:"varchar(255)" json:"sign_arca"`
	ExpireTokenArca *time.Time `gorm:"type:datetime" json:"expire_token_arca"`

	TokenEmail *string `gorm:"varchar(255)" json:"token_email"`
}
