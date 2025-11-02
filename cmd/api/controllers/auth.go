package controllers

import (
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func (a *AuthController) AuthLogin(c *fiber.Ctx) error {
	logging.INFO("Login")
	var loginRequest schemas.AuthLogin
	if err := c.BodyParser(&loginRequest); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	if err := loginRequest.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, err)
	}

	token, err := a.AuthService.AuthLogin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
	}
	
	c.Cookie(cookie)
	
	logging.INFO("Login exitoso")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login exitoso",
	})
}

func (a *AuthController) AuthTenant(c *fiber.Ctx) error {
	logging.INFO("Login tenant")
	id := c.Params("tenant_id")
	if id == "" {
		logging.ERROR("ID is required")
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID es requerido",
		})
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	token, err := a.AuthService.AuthGetTenant(user, id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	logging.INFO("Login tenant exitoso")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login a Punto de venta exitoso, token enviado en cookie",
	})
}

func (a *AuthController) LogoutTenant(c *fiber.Ctx) error {
	logging.INFO("Logout tenant")
	user := c.Locals("user").(*schemas.AuthenticatedUser)

	token, err := a.AuthService.LogoutPointSale(user.ID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	logging.INFO("Logout tenant exitoso")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Logout de Punto de venta exitoso, token enviado en cookie",
	})
}

func (a *AuthController) Logout(ctx *fiber.Ctx) error {
	logging.INFO("Logout")
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
	})

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Logout exitoso",
	})
}
