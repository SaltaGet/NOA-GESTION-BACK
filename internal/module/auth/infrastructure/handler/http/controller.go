package http

import (
	"os"
	"strconv"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/auth/domain"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/module/auth/infrastructure/repository"
	repositoryTenant "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/tenant/infrastructure/repository"
	repositoryPlan "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/plan/infrastructure/repository"
	domainEmail "github.com/SaltaGet/NOA-GESTION-BACK/internal/module/email/domain"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

type AuthController struct {
	AuthService domain.AuthService
	EmailService domainEmail.EmailService
}

// Login godoc
//
// @Summary		Login user
// @Description	Login user required identifier and password
// @Tags			Auth
// @Accept			json
// @Produce		json
// @Param			credentials	body		repository.AuthLogin	true	"Credentials"
// @Success		200			{object}	utils.Response
// @Router			/api/v1/auth/login [post]
func (a *AuthController) AuthLogin(c *fiber.Ctx) error {
	var loginRequest repository.AuthLogin
	if err := validator.ValidateRequest(c, &loginRequest); err != nil {
		return utils.HandleError(c, err)
	}

	token, err := a.AuthService.AuthLogin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return utils.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   utils.Ternary(os.Getenv("ENV") == "dev", false, true),     // poner en true para prod
		SameSite: utils.Ternary(os.Getenv("ENV") == "dev", "None", "None",), // para prod : "Strict",
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Success		200				{object}	utils.Response
// @Router			/api/v1/auth/login_point_sale/{point_sale_id} [post]
func (a *AuthController) AuthPointSale(c *fiber.Ctx) error {
	id := c.Params("point_sale_id")
	if id == "" {
		log.Err(nil).Msg("ID is required")
		return c.Status(400).JSON(utils.Response{
			Status:  false,
			Body:    nil,
			Message: "ID es requerido",
		})
	}

	pointSaleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorResponse(400, "tenant_id debe ser un número", err))
	}

	user := c.Locals("user").(*repository.AuthenticatedUser)

	token, err := a.AuthService.AuthPointSale(user, pointSaleID)
	if err != nil {
		return utils.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   utils.Ternary(os.Getenv("ENV") == "dev", false, true),     // poner en true para prod
		SameSite: utils.Ternary(os.Getenv("ENV") == "dev", "None", "None",), // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Success		200	{object}	utils.Response
// @Router			/api/v1/auth/logout_point_sale [post]
func (a *AuthController) LogoutPointSale(c *fiber.Ctx) error {
	member := c.Locals("user").(*repository.AuthenticatedUser)

	token, err := a.AuthService.LogoutPointSale(member)
	if err != nil {
		return utils.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   utils.Ternary(os.Getenv("ENV") == "dev", false, true),     // poner en true para prod
		SameSite: utils.Ternary(os.Getenv("ENV") == "dev", "None", "None",), // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Success		200	{object}	utils.Response
// @Router			/api/v1/auth/logout [post]
func (a *AuthController) Logout(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   utils.Ternary(os.Getenv("ENV") == "dev", false, true),     // poner en true para prod
		SameSite: utils.Ternary(os.Getenv("ENV") == "dev", "None", "None",), // para prod : "Strict",
	})

	return ctx.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Success		200	{object}	utils.Response{body=repository.AuthenticatedUser}
// @Router			/api/v1/auth/current_user [get]
func (a *AuthController) CurrentUser(c *fiber.Ctx) error {
	user := c.Locals("user").(*repository.AuthenticatedUser)

	return c.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Success		200	{object}	utils.Response{body=repository.PlanResponseDTO}
// @Router			/api/v1/auth/current_plan [get]
func (a *AuthController) CurrentPlan(c *fiber.Ctx) error {
	plan := c.Locals("current_plan").(*repositoryPlan.PlanResponseDTO)

	return c.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Success		200	{object}	utils.Response{body=repository.TenantResponse}
// @Router			/api/v1/auth/current_tenant [get]
func (a *AuthController) CurrentTenant(c *fiber.Ctx) error {
	user := c.Locals("current_tenant").(*repositoryTenant.TenantResponse)

	return c.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Param			credentials	body		repository.AuthLoginAdmin	true	"Credentials"
// @Success		200			{object}	utils.Response
// @Router			/api/v1/auth/login_admin [post]
func (a *AuthController) AuthLoginAdmin(c *fiber.Ctx) error {
	var loginRequest repository.AuthLoginAdmin
	if err := validator.ValidateRequest(c, &loginRequest); err != nil {
		return utils.HandleError(c, err)
	}

	token, err := a.AuthService.AuthLoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		return utils.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token_admin",
		Value:    token,
		HTTPOnly: true,
		Secure:   utils.Ternary(os.Getenv("ENV") == "dev", false, true),     // poner en true para prod
		SameSite: utils.Ternary(os.Getenv("ENV") == "dev", "None", "None",), // para prod : "Strict",
	}

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Success		200	{object}	utils.Response
// @Router			/api/v1/auth/logout_admin [post]
func (a *AuthController) LogoutAdmin(ctx *fiber.Ctx) error {
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token_admin",
		Value:    "",
		HTTPOnly: true,
		Secure:   utils.Ternary(os.Getenv("ENV") == "dev", false, true),     // poner en true para prod
		SameSite: utils.Ternary(os.Getenv("ENV") == "dev", "None", "None",), // para prod : "Strict",
	})

	return ctx.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Param			forgot_password	body		repository.AuthForgotPassword	true "field to send email"
// @Success		200				{object}	utils.Response
// @Router			/api/v1/auth/forgot_password [post]
func (a *AuthController) AuthForgotPassword(ctx *fiber.Ctx) error {
	var authForgotPassword repository.AuthForgotPassword
	if err := validator.ValidateRequest(ctx, &authForgotPassword); err != nil {
		return utils.HandleError(ctx, err)
	}

	err := a.AuthService.AuthForgotPassword(&authForgotPassword)
	if err != nil {
		return utils.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.Response{
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
// @Param			reset_password	body		repository.AuthResetPassword	true "new password"
// @Success		200				{object}	utils.Response
// @Router			/api/v1/auth/reset_password [post]
func (a *AuthController) AuthResetPassword(ctx *fiber.Ctx) error {
	var authResetPassword repository.AuthResetPassword
	if err := validator.ValidateRequest(ctx, &authResetPassword); err != nil {
		return utils.HandleError(ctx, err)
	}

	err := a.AuthService.AuthResetPassword(&authResetPassword)
	if err != nil {
		return utils.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(utils.Response{
		Status:  true,
		Body:    nil,
		Message: "Contraseña actualizada con exito!",
	})
}
