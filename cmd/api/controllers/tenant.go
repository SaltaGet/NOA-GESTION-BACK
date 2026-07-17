package controllers

import (
	"strconv"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models/master"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/cache"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/utils"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/platform/validator"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
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
	if err := validator.ValidateRequest(c, &tenantCreate); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*master.Admin)
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
	if err := validator.ValidateRequest(c, &tenantUserCrate); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*master.Admin)
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
	if err := validator.ValidateRequest(c, &tenantUpdateExpiration); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*master.Admin)
	err := t.TenantService.TenantUpdateExpiration(admin.ID, &tenantUpdateExpiration)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	if cache.IsAvailable() {
		_ = cache.InvalidateTenantConnection(tenantUpdateExpiration.ID)
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

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	err := t.TenantService.TenantUpdateTerms(user.TenantID, tenantUpdateTerms)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	if cache.IsAvailable() {
		_ = cache.InvalidateAuthUser(user.ID, user.TenantID)
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
	if err := validator.ValidateRequest(c, &tenantUpdateSettings); err != nil {
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

// TenantGetWithModules godoc
//
//	@Summary		TenantGetWithModules
//	@Description	Obtener información completa de un tenant con sus módulos asociados
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		AdminAuth
//	@Param			id	query		int64														true	"Tenant ID"
//	@Success		200	{object}	schemas.Response{body=schemas.TenantWithModulesResponse}	"Tenant con módulos obtenido con éxito"
//	@Failure		400	{object}	schemas.Response											"Bad Request"
//	@Failure		401	{object}	schemas.Response											"Auth is required"
//	@Failure		403	{object}	schemas.Response											"Not Authorized"
//	@Failure		404	{object}	schemas.Response											"Not Found"
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/tenant/get_with_modules [get]
func (t *TenantController) TenantGetWithModules(c *fiber.Ctx) error {
	idStr := c.Query("id", "")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Falta el parámetro id",
		})
	}

	tenantID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El parámetro id debe ser un número",
		})
	}

	result, err := t.TenantService.TenantGetWithModules(tenantID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    result,
		Message: "Tenant con módulos obtenido con éxito",
	})
}

// TenantGetAllWithModules godoc
//
//	@Summary		TenantGetAllWithModules
//	@Description	Obtener lista de tenants con información reducida y sus módulos asociados
//	@Tags			Tenant
//	@Accept			json
//	@Produce		json
//	@Security		AdminAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.TenantSimpleWithModulesResponse}	"Tenants con módulos obtenidos con éxito"
//	@Failure		401	{object}	schemas.Response													"Auth is required"
//	@Failure		403	{object}	schemas.Response													"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/tenant/get_all_with_modules [get]
func (t *TenantController) TenantGetAllWithModules(c *fiber.Ctx) error {
	result, err := t.TenantService.TenantGetAllWithModules()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	if result == nil {
		result = []schemas.TenantSimpleWithModulesResponse{}
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    result,
		Message: "Tenants con módulos obtenidos con éxito",
	})
}

