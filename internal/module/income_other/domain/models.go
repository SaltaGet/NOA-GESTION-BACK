package domain

import (
	"time"

	domainM "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/member/domain"
	domainPS "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/point_sale/domain"
	domainTI "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/type_income/domain"
	domainCR "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/cash_register/domain"
)

type IncomeOther struct {
	ID             int64         `gorm:"primaryKey;autoIncrement" json:"id"`
	PointSaleID    *int64        `gorm:"" json:"point_sale_id"`
	PointSale      *domainPS.PointSale    `gorm:"foreignKey:PointSaleID;references:ID" json:"point_sale"`
	MemberID       *int64        `gorm:"" json:"member_id"`
	Member         *domainM.Member       `gorm:"foreignKey:MemberID;references:ID" json:"member"`
	CashRegisterID *int64        `gorm:"" json:"cash_register_id"`
	CashRegister   *domainCR.CashRegister `gorm:"foreignKey:CashRegisterID;references:ID" json:"cash_register"`
	Total          float64       `gorm:"not null" json:"total"`
	TypeIncomeID   int64         `gorm:"not null" json:"type_income_id"`
	TypeIncome     domainTI.TypeIncome    `gorm:"foreignKey:TypeIncomeID" json:"type_income"`
	Details        *string       `gorm:"type:varchar(255)" json:"details"`
	MethodIncome   string        `gorm:"not null;default:cash;type:varchar(30)" json:"method_income" validate:"oneof=cash credit card transfer"`
	CreatedAt      time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
