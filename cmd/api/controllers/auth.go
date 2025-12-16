package controllers

import (
	"strconv"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Login godoc
//
// @Summary		Login user
// @Description	Login user required identifier and password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		schemas.AuthLogin	true	"Credentials"
// @Success		200			{object}	schemas.Response
// @Router			/api/v1/auth/login [post]
func (a *AuthController) AuthLogin(c *fiber.Ctx) error {
	var loginRequest schemas.AuthLogin
	if err := c.BodyParser(&loginRequest); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	if err := loginRequest.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	token, err := a.AuthService.AuthLogin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,     // poner en true para prod
		SameSite: "Strict", // para prod : "Strict",
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login exitoso",
	})
}

// LoginPointSale godoc
//
// @Summary		LoginPointSale
// @Description	Login al punto de venta
// @Tags			Auth
// @Accept			json
// @Produce		json
//
// @Security		CookieAuth
//
// @Param			point_sale_id	path		string	true	"id del punto de venta"
// @Success		200				{object}	schemas.Response
// @Router			/api/v1/auth/login_point_sale/{point_sale_id} [post]
func (a *AuthController) AuthPointSale(c *fiber.Ctx) error {
	id := c.Params("point_sale_id")
	if id == "" {
		log.Err(nil).Msg("ID is required")
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID es requerido",
		})
	}

	pointSaleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).JSON(schemas.ErrorResponse(400, "tenant_id debe ser un número", err))
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	token, err := a.AuthService.AuthPointSale(user, pointSaleID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,     // poner en true para prod
		SameSite: "Strict", // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login a Punto de venta exitoso, token enviado en cookie",
	})
}

// LogoutPointSale godoc
//
// @Summary		LogoutPointSale
// @Description	Logout del punto de venta
// @Tags			Auth
// @Accept			json
// @Produce		json
//
// @Security		CookieAuth
//
// @Success		200	{object}	schemas.Response
// @Router			/api/v1/auth/logout_point_sale [post]
func (a *AuthController) LogoutPointSale(c *fiber.Ctx) error {
	member := c.Locals("user").(*schemas.AuthenticatedUser)

	token, err := a.AuthService.LogoutPointSale(member)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,     // poner en true para prod
		SameSite: "Strict", // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Logout de Punto de venta exitoso, token enviado en cookie",
	})
}

// Logout godoc
//
// @Summary		Logout user
// @Description	Logout user
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Success		200	{object}	schemas.Response
// @Router			/api/v1/auth/logout [post]
func (a *AuthController) Logout(ctx *fiber.Ctx) error {
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

// CurrentUser godoc
//
// @Summary		CurrentUser user
// @Description	Obtener usuario actual
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Success		200	{object}	schemas.Response{body=schemas.AuthenticatedUser}
// @Router			/api/v1/auth/current_user [get]
func (a *AuthController) CurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.AuthenticatedUser)

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    user,
		Message: "Usuario actual obtenido",
	})
}

// CurrentPlan godoc
//
// @Summary		CurrentPlan
// @Description	Obtener plan actual
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Success		200	{object}	schemas.Response{body=schemas.PlanResponseDTO}
// @Router			/api/v1/auth/current_plan [get]
func (a *AuthController) CurrentPlan(c *fiber.Ctx) error {
	plan := c.Locals("current_plan").(*schemas.PlanResponseDTO)

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    plan,
		Message: "plan actual obtenido",
	})
}

// CurrentTenant godoc
//
// @Summary		CurrentTenant
// @Description	Obtener tenant actual
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Success		200	{object}	schemas.Response{body=schemas.TenantResponse}
// @Router			/api/v1/auth/current_tenant [get]
func (a *AuthController) CurrentTenant(c *fiber.Ctx) error {
	user := c.Locals("current_tenant").(*schemas.TenantResponse)

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    user,
		Message: "Tenant actual obtenido",
	})
}

// LoginAdmin godoc
//
// @Summary		Login Admin user
// @Description	Required identifier and password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		schemas.AuthLoginAdmin	true	"Credentials"
// @Success		200			{object}	schemas.Response
// @Router			/api/v1/auth/login_admin [post]
func (a *AuthController) AuthLoginAdmin(c *fiber.Ctx) error {
	var loginRequest schemas.AuthLoginAdmin
	if err := c.BodyParser(&loginRequest); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	if err := loginRequest.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	token, err := a.AuthService.AuthLoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token_admin",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,  // poner en true para prod
		SameSite: "Strict", // para prod : "Strict",
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login exitoso",
	})
}

// LogoutAdmin godoc
//
// @Summary		Logout Admin user
// @Description	logout user admin
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Success		200	{object}	schemas.Response
// @Router			/api/v1/auth/logout_admin [post]
func (a *AuthController) LogoutAdmin(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token_admin",
		Value:    "",
		HTTPOnly: true,
		Secure:   true,  // poner en true para prod
		SameSite: "Strict", // para prod : "Strict",
	})

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Logout exitoso",
	})
}

// ForgotPassword godoc
//
// @Summary		ForgotPassword
// @Description	recuperar contraseña por email
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			forgot_password	body		schemas.AuthForgotPassword	true "field to send email"
// @Success		200				{object}	schemas.Response
// @Router			/api/v1/auth/forgot_password [post]
func (a *AuthController) AuthForgotPassword(ctx *fiber.Ctx) error {
	var authForgotPassword schemas.AuthForgotPassword
	if err := ctx.BodyParser(&authForgotPassword); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}
	if err := authForgotPassword.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	err := a.AuthService.AuthForgotPassword(&authForgotPassword)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Email enviado con exito!",
	})
}

// ResetPassword godoc
//
// @Summary		ResetPassword
// @Description	recuperar contraseña por email
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			reset_password	body		schemas.AuthResetPassword	true "new password"
// @Success		200				{object}	schemas.Response
// @Router			/api/v1/auth/reset_password [post]
func (a *AuthController) AuthResetPassword(ctx *fiber.Ctx) error {
	var authResetPassword schemas.AuthResetPassword
	if err := ctx.BodyParser(&authResetPassword); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}
	if err := authResetPassword.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	err := a.AuthService.AuthResetPassword(&authResetPassword)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Contraseña actualizada con exito!",
	})
}
