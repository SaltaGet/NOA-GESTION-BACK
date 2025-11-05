package schemas

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type PointSaleCreate struct {
	Name        string  `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description *string `gorm:"size:200" json:"description"`
	IsDeposit   bool    `gorm:"not null;default:false" json:"is_deposit"`
}

func (p *PointSaleCreate) Validate() error {
	validate := validator.New()
	err := validate.Struct(p)
	if err == nil {
		return nil
	}

	validationErr := err.(validator.ValidationErrors)[0]
	field := validationErr.Field()
	tag := validationErr.Tag()
	param := validationErr.Param()

	return fmt.Errorf("campo %s es invalido, revisar: (%s) (%s)", field, tag, param)
}

type PointSaleResponse struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	IsDeposit   bool    `json:"is_deposit"`
}
