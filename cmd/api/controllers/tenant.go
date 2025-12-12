package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// Tenant godoc
//
// @Summary		Tenant GetAll
// @Description	Tenant GetAll required auth token
// @Tags			Tenant
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Success		200	{object}	schemas.Response{body=[]schemas.TenantResponse}	"Tenants obtenidos con éxito"
// @Failure		400	{object}	schemas.Response								"Bad Request"
// @Failure		401	{object}	schemas.Response								"Auth is required"
// @Failure		403	{object}	schemas.Response								"Not Authorized"
// @Failure		500	{object}	schemas.Response
// @Router			/api/v1/tenant/get_all [get]
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
// @Summary		Tenant Create
// @Description	Tenant Create required auth token
// @Tags			Tenant
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Param			user_id			query		int64					true	"UserID"
// @Param			TenantCreate	body		schemas.TenantCreate	true	"TenantCreate"
// @Success		200				{object}	schemas.Response		"Tenant creado con éxito"
// @Failure		400				{object}	schemas.Response		"Bad Request"
// @Failure		401				{object}	schemas.Response		"Auth is required"
// @Failure		403				{object}	schemas.Response		"Not Authorized"
// @Failure		500				{object}	schemas.Response
// @Router			/api/v1/tenant/create [post]
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

	id, err := t.TenantService.TenantCreateByUserID(&tenantCreate, userID)
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
// @Summary		Tenant Create
// @Description	Tenant Create required auth token
// @Tags			Tenant
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Param			TenantUserCreate	body		schemas.TenantUserCreate	true	"TenantUserCreate"
// @Success		200					{object}	schemas.Response			"Tenant y Usuario creados con éxito"
// @Failure		400					{object}	schemas.Response			"Bad Request"
// @Failure		401					{object}	schemas.Response			"Auth is required"
// @Failure		403					{object}	schemas.Response			"Not Authorized"
// @Failure		500					{object}	schemas.Response
// @Router			/api/v1/tenant/create_tenant_user [post]
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

	id, err := t.TenantService.TenantUserCreate(&tenantUserCrate)
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
// @Summary		TenantUpdateExpiration
// @Description	Actualizar fecha de expiración de un tenant
// @Tags			Tenant
// @Accept			json
// @Produce		json
// @Security		CookieAuth
// @Param			TenantUserCreate	body		schemas.TenantUpdateExpiration	true	"TenantUserCreate"
// @Success		200					{object}	schemas.Response				"Fecha de expiración actualizada con éxito"
// @Router			/api/v1/tenant/update_expiration [put]
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

	err := t.TenantService.TenantUpdateExpiration(&tenantUpdateExpiration)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Actualización de fecha de expiración del Tenant exitoso",
	})
}
