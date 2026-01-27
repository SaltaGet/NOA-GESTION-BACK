package controllers

import (
	"strconv"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Tenant godoc
//
//	@Summary		Tenant GetAll
//	@Description	Tenant GetAll required auth token
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.TenantResponse}	"Tenants obtenidos con éxito"
//	@Failure		400	{object}	schemas.Response								"Bad Request"
//	@Failure		401	{object}	schemas.Response								"Auth is required"
//	@Failure		403	{object}	schemas.Response								"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/tenant/get_all [get]
func (t *TenantController) GetTenants(c *fiber.Ctx) error {
	tenants, err := t.TenantService.TenantGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	if tenants == nil || len(*tenants) == 0 {
		empty := []schemas.TenantResponse{}
		tenants = &empty
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    *tenants,
		Message: "Tenants obtenidos con éxito",
	})
}

// Tenant godoc
//
//	@Summary		Tenant Create
//	@Description	Tenant Create required auth token
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			user_id			query		int64					true	"UserID"
//	@Param			TenantCreate	body		schemas.TenantCreate	true	"TenantCreate"
//	@Success		200				{object}	schemas.Response		"Tenant creado con éxito"
//	@Failure		400				{object}	schemas.Response		"Bad Request"
//	@Failure		401				{object}	schemas.Response		"Auth is required"
//	@Failure		403				{object}	schemas.Response		"Not Authorized"
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/tenant/create [post]
func (t *TenantController) TenantCreateByUserID(c *fiber.Ctx) error {
	userIDStr := c.Query("user_id", "")
	if userIDStr == "" {
		return c.Status(400).SendString("falta el parámetro user_id")
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(400).SendString("user_id debe ser un número")
	}

	var tenantCreate schemas.TenantCreate
	if err := c.BodyParser(&tenantCreate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := tenantCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*models.Admin)
	id, err := t.TenantService.TenantCreateByUserID(admin.ID, &tenantCreate, userID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Tenant creado con éxito",
	})
}

// Tenant godoc
//
//	@Summary		Tenant Create
//	@Description	Tenant Create required auth token
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			TenantUserCreate	body		schemas.TenantUserCreate	true	"TenantUserCreate"
//	@Success		200					{object}	schemas.Response			"Tenant y Usuario creados con éxito"
//	@Failure		400					{object}	schemas.Response			"Bad Request"
//	@Failure		401					{object}	schemas.Response			"Auth is required"
//	@Failure		403					{object}	schemas.Response			"Not Authorized"
//	@Failure		500					{object}	schemas.Response
//	@Router			/api/v1/tenant/create_tenant_user [post]
func (t *TenantController) TenantUserCreate(c *fiber.Ctx) error {
	var tenantUserCrate schemas.TenantUserCreate
	if err := c.BodyParser(&tenantUserCrate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := tenantUserCrate.TenantCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	if err := tenantUserCrate.UserCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*models.Admin)
	id, err := t.TenantService.TenantUserCreate(admin.ID, &tenantUserCrate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Tenant creado con éxito",
	})
}

// TenantUpdateExpiration godoc
//
//	@Summary		TenantUpdateExpiration
//	@Description	Actualizar fecha de expiración de un tenant
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			TenantUserCreate	body		schemas.TenantUpdateExpiration	true	"TenantUserCreate"
//	@Success		200					{object}	schemas.Response				"Fecha de expiración actualizada con éxito"
//	@Router			/api/v1/tenant/update_expiration [put]
func (t *TenantController) TenantUpdateExpiration(c *fiber.Ctx) error {
	var tenantUpdateExpiration schemas.TenantUpdateExpiration
	if err := c.BodyParser(&tenantUpdateExpiration); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := tenantUpdateExpiration.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*models.Admin)
	err := t.TenantService.TenantUpdateExpiration(admin.ID, &tenantUpdateExpiration)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Actualización de fecha de expiración del Tenant exitoso",
	})
}

// TenantUpdateAcceptedTerms godoc
//
//	@Summary		TenantUpdateAcceptedTerms
//	@Description	Actualizar la aceptación de los termninos y condiciones de un tenant
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response	"actualizacion con éxito"
//	@Router			/api/v1/tenant/update_temrs [put]
func (t *TenantController) TenantUpdateAcceptedTerms(c *fiber.Ctx) error {
	ip := c.IP()

	// Si estás en localhost y Fiber no detecta la IP, o devuelve el formato IPv6
	if ip == "::1" || ip == "" {
		ip = "127.0.0.1"
	}

	tenantUpdateTerms := &schemas.TenantUpdateTerms{
		IP:            ip,
		AcceptedTerms: true,
		DateAccepted:  time.Now(),
	}
	if err := tenantUpdateTerms.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	err := t.TenantService.TenantUpdateTerms(user.TenantID, tenantUpdateTerms)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Actualización de fecha de expiración del Tenant exitoso",
	})
}

// TenantGetSettings godoc
//
//	@Summary		TenantGetSettings
//	@Description	Obtener la configuración de un tenant
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=schemas.TenantSettingsResponse}	"actualizacion con éxito"
//	@Router			/api/v1/tenant/get_settings [get]
func (t *TenantController) TenantGetSettings(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.AuthenticatedUser)
	sett, err := t.TenantService.TenantGetSettings(user.TenantID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    sett,
		Message: "Configuración del Tenant obtenida exitosamente",
	})
}

// TenantUpdateSettings godoc
//
//	@Summary		TenantUpdateSettings
//	@Description	Actualizar la configuración de un tenant
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			TenantUpdateSettings	body		schemas.TenantUpdateSettings	true	"TenantUpdateSettings"
//	@Success		200						{object}	schemas.Response				"actualizacion con éxito"
//	@Router			/api/v1/tenant/update_settings [put]
func (t *TenantController) TenantUpdateSettings(c *fiber.Ctx) error {
	var tenantUpdateSettings schemas.TenantUpdateSettings
	if err := c.BodyParser(&tenantUpdateSettings); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := tenantUpdateSettings.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	err := t.TenantService.TenantUpdateSettings(user.TenantID, &tenantUpdateSettings)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Actualización de la configuración del Tenant exitoso",
	})
}

// TenantGenerateTokenToImageSetting godoc
//
//	@Summary		TenantGenerateTokenToImageSetting
//
//	@Description	### Flujo de Carga de Imágenes
//	@Description	Genera un token temporal para subir imágenes al microservicio.
//	@Description
//	@Description	**Pasos requeridos:**
//	@Description	1. Llamar a este endpoint para obtener el token.
//	@Description	2. Incluir el token en el header `x-token-tenant` del microservicio de imágenes.
//	@Description
//	@Description	**Formato del endpoint del microservicio:**
//	@Description	~~~
//	@Description	POST /ecommerce/{tenantIdentifier}/api/v1/tenant/upload_image
//	@Description	~~~
//	@Description
//	@Description	> *Nota: El token tiene una validez limitada de 30 minutos.*
//
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response
//	@Router			/api/v1/tenant/generate_token_to_image_setting [post]
func (t *TenantController) TenantGenerateTokenToImageSetting(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*schemas.AuthenticatedUser)

	token, err := utils.GenerateTokenToGrpcToSetting(user.TenantIdentifier)
	if err != nil {
		log.Err(err).Msg("Error al generar el token")
		return ctx.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Error al generar el token",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    token,
		Message: "Token generado correctamente",
	})
}
