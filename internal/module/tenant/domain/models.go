package domain


import (
	"errors"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
	modelSetting "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/setting/domain"
	modelCredential "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/credential/domain"
	modelPlan "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/plan/domain"
	modelUserTenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/user_tenant/domain"
	modelPayTenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/pay_tenant/domain"
	modelModule "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/module/domain"
)

type Tenant struct {
	ID            int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null" json:"name"`
	Identifier    string         `gorm:"type:varchar(50);not null;unique" json:"identifier"`
	Address       string         `gorm:"type:varchar(255);not null" json:"address"`
	Phone         string         `gorm:"type:varchar(20);not null" json:"phone"`
	Email         string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	CuitPdv       string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"cuit_pdv"`
	IsActive      bool           `gorm:"not null;default:true" json:"is_active"`
	PlanID        int64          `gorm:"not null" json:"plan_id"`
	Connection    string         `gorm:"type:varchar(255);not null" json:"connection"`
	Expiration    *time.Time     `gorm:"" json:"expiration"`
	AcceptedTerms bool           `gorm:"not null;default:false" json:"accepted_terms"`
	IP            *string        `gorm:"type:varchar(255);default:null" json:"ip"`
	DateAccepted  *time.Time     `gorm:"" json:"date_accepted"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Setting       modelSetting.SettingTenant  `gorm:"foreignKey:TenantID" json:"setting"`
	PayTenant     []modelPayTenant.PayTenant    `gorm:"foreignKey:TenantID" json:"pay_tenants"`
	UserTenants   []modelUserTenant.UserTenant   `gorm:"foreignKey:TenantID" json:"user_tenants"`
	Plan          modelPlan.Plan           `gorm:"foreignKey:PlanID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"plan"`
	Modules       []modelModule.TenantModule `gorm:"foreignKey:TenantID" json:"modules"`
	Credentials   modelCredential.Credential     `gorm:"foreignKey:TenantID" json:"credentials"`
}

func (t *Tenant) BeforeCreate(tx *gorm.DB) (err error) {
	t.Identifier = strings.ToLower(strings.TrimSpace(t.Identifier))

	var validSubdomain = regexp.MustCompile(`^[a-z0-9-]+$`)

	if !validSubdomain.MatchString(t.Identifier) {
		return errors.New("identifier invalid - only lowercase letters, numbers, and hyphens are allowed")
	}

	return nil
}
