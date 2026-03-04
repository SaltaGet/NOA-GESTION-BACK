package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/validator"
	"github.com/gofiber/fiber/v2"
)

// CreateUser godoc
//	@Summary		Create User
//	@Description	Creates a new user.
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			userCreate	body		schemas.UserCreate	true	"User information"
//	@Success		201			{object}	schemas.Response
//	@Failure		400			{object}	schemas.Response	"Bad Request"
//	@Failure		401			{object}	schemas.Response	"Auth is required"
//	@Failure		403			{object}	schemas.Response	"Not Authorized"
//	@Failure		500			{object}	schemas.Response
//	@Router			/user/create [post]
func (u *UserController) CreateUser(c *fiber.Ctx) error {
	var userCreate schemas.UserCreate
	if err := validator.ValidateRequest(c, &userCreate); err != nil {
		return schemas.HandleError(c, err)
	}
	
	return c.Status(fiber.StatusCreated).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "User created",
	})
}