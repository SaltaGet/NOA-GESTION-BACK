package schemas

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type MemberResponse struct {
	ID         int64                `json:"id"`
	FirstName  string              `json:"first_name"`
	LastName   string              `json:"last_name"`
	Username   string              `json:"username"`
	Email      string              `json:"email"`
	Address    *string             `json:"address"`
	Phone      *string             `json:"phone"`
	IsAdmin    bool                `json:"is_admin"`
	IsActive   bool                `json:"is_active"`
	Role       RoleResponse        `json:"role"`
	PointSales []PointSaleResponse `json:"point_sales"`
}

type MemberResponseDTO struct {
	ID        int64         `json:"id"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Address   *string      `json:"address"`
	Phone     *string      `json:"phone"`
	IsActive  bool         `json:"is_active"`
	IsAdmin   bool         `json:"is_admin"`
	Role      RoleResponseDTO `json:"role"`
}

type MemberSimpleDTO struct {
	ID        int64    `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Username  string  `json:"username"`
}

type MemberCreate struct {
	FirstName    string  `json:"first_name" validate:"required" example:"John"`
	LastName     string  `json:"last_name" validate:"required" example:"Doe"`
	Username     string  `json:"username" validate:"required" example:"johndoe"`
	Email        string  `json:"email" validate:"email,required" example:"a@b.com"`
	Password     string  `json:"password" validate:"required,password" example:"Password123*"`
	Address      *string `json:"address,omitempty" example:"casita roja|null"`
	Phone        *string `json:"phone,omitempty" example:"123123123|null"`
	RoleID       int64   `json:"role_id" validate:"required" example:"1"`
	PointSaleIDs []int64 `json:"point_sales_ids" validate:"required" example:"1,2,3"`
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 8 {
		return false
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}
	if !regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password) {
		return false
	}
	return true
}

type MemberUpdate struct {
	ID           int64    `json:"id" validate:"required" example:"1"`
	FirstName    string  `json:"first_name" validate:"required" example:"John"`
	LastName     string  `json:"last_name" validate:"required" example:"Doe"`
	Username     string  `json:"username" validate:"required" example:"johndoe"`
	Email        string  `json:"email" validate:"email,required" example:"a@b.com"`
	Address      *string `json:"address,omitempty" example:"address|null"`
	Phone        *string `json:"phone,omitempty" example:"phone|null"`
	RoleID       int64    `json:"role_id" validate:"required" example:"1"`
	IsActive     *bool   `json:"is_active" validate:"required" example:"true"`
	PointSaleIDs []int64  `json:"point_sales_ids" validate:"required" example:"1,2,3"`
}

type MemberUpdatePassword struct {
	OldPassword string `json:"old_password" validate:"required,password" example:"Password123*"`
	NewPassword string `json:"new_password" validate:"required,password" example:"Password123*"`
	ConfirmPass string `json:"confirm_pass" validate:"required,eqfield=NewPassword" example:"Password123*"`
}