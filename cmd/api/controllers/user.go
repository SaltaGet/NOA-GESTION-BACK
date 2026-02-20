package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
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
	if err := validators.ValidateRequest(c, &userCreate); err != nil {
		return schemas.HandleError(c, err)
	}
	// userCreated, err := u.UserService.UserCreate(&userCreate)
	// if err != nil {
	// 	if errResp, ok := err.(*schemas.ErrorStruc); ok {
	// 		logging.ERROR("Error: %s", errResp.Err.Error())
	// 		return c.Status(errResp.StatusCode).JSON(schemas.Response{
	// 			Status:  false,
	// 			Body:    nil,
	// 			Message: errResp.Message,
	// 		})
	// 	}
	// 	logging.ERROR("Error: %s", err.Error())
	// 	return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
	// 		Status:  false,
	// 		Body:    nil,
	// 		Message: "Error interno",
	// 	})
	// }
	
	return c.Status(fiber.StatusCreated).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "User created",
	})
}