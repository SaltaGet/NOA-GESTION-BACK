package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// CredentialGetMPToken godoc
//
//	@Summary		CredentialGetMPToken
//	@Description	Obtener tokens de mercado pago
//	@Tags			Credential
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=schemas.CredentialMPTokenResponse}
//	@Router			/api/v1/credential/get_mercado_pago_token [get]
func (c *CredentialController) CredentialGetMPToken(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schemas.AuthenticatedUser)
	tokens, err := c.CredentialService.CredentialGetMPToken(user.TenantID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    tokens,
		Message: "Tokens obtenidos con exito!",
	})
}

// CredentialSetMPToken godoc
//
//	@Summary		CredentialSetMPToken
//	@Description	Actualizar token de mercado pago
//	@Tags			Credential
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		schemas.CredentialMPTokenRequest	true	"token de mercado pago"
//	@Success		200		{object}	schemas.Response
//	@Router			/api/v1/credential/set_mercado_pago_token [put]
func (c *CredentialController) CredentialSetMPToken(ctx *fiber.Ctx) error {
	var request schemas.CredentialMPTokenRequest
	if err := ctx.BodyParser(&request); err != nil {
		return schemas.HandleError(ctx, err)
	}
	if err := request.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	user := ctx.Locals("user").(*schemas.AuthenticatedUser)
	message, err := c.CredentialService.CredentialSetMPToken(user.TenantID, &request)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    message,
		Message: "Token actualizado con exito!",
	})
}

// CredentialGetArca godoc
//
//	@Summary		CredentialGetArca
//	@Description	Obtener credenciales de ARCA
//	@Tags			Credential
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=schemas.CredentialArcaResponse}
//	@Router			/api/v1/credential/get_arca_token [get]
func (c *CredentialController) CredentialGetArca(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schemas.AuthenticatedUser)
	cred, err := c.CredentialService.CredentialGetArca(user.TenantID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    cred,
		Message: "Credenciales obtenidas con exito!",
	})
}

// CredentialSetArca godoc
//
//	@Summary		CredentialSetArca
//	@Description	Actualizar credenciales de ARCA
//	@Tags			Credential
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		schemas.CredentialArcaRequest	true	"token de mercado pago"
//	@Success		200		{object}	schemas.Response
//	@Router			/api/v1/credential/set_arca_token [put]
func (c *CredentialController) CredentialSetArca(ctx *fiber.Ctx) error {
	var request schemas.CredentialArcaRequest
	if err := ctx.BodyParser(&request); err != nil {
		return schemas.HandleError(ctx, err)
	}
	if err := request.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	user := ctx.Locals("user").(*schemas.AuthenticatedUser)
	err := c.CredentialService.CredentialSetArca(user.TenantID, &request)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Token actualizado con exito!",
	})
}