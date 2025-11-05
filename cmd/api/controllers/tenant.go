package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

//	Tenant godoc
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
	logging.INFO("Obtener todos los tenants")
	tenants, err := t.TenantService.TenantGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	if tenants == nil || len(*tenants) == 0 {
		empty := []schemas.TenantResponse{}
		tenants = &empty
	}

	logging.INFO("Tenants obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    *tenants,
		Message: "Tenants obtenidos con éxito",
	})
}

//	Tenant godoc
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
	logging.INFO("Crear tenant")
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
		logging.ERROR("Invalid request %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := tenantCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(422).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := t.TenantService.TenantCreateByUserID(&tenantCreate, userID)
	if err != nil {
		if errResp, ok := err.(*schemas.ErrorStruc); ok {
			logging.ERROR("Error: %s", errResp.Err.Error())
			return c.Status(errResp.StatusCode).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: errResp.Message,
			})
		}
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Error interno",
		})
	}

	logging.INFO("Tenant creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Tenant creado con éxito",
	})
}

//	Tenant godoc
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
	logging.INFO("Crear tenant y user")
	var tenantUserCrate schemas.TenantUserCreate
	if err := c.BodyParser(&tenantUserCrate); err != nil {
		logging.ERROR("Invalid request %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := tenantUserCrate.TenantCreate.Validate(); err != nil {
		logging.ERROR("Tenant Error: %s", err.Error())
		return c.Status(422).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}
	if err := tenantUserCrate.UserCreate.Validate(); err != nil {
		logging.ERROR("User Error: %s", err.Error())
		return c.Status(422).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := t.TenantService.TenantUserCreate(&tenantUserCrate)
	if err != nil {
		if errResp, ok := err.(*schemas.ErrorStruc); ok {
			logging.ERROR("Error: %s", errResp.Err.Error())
			return c.Status(errResp.StatusCode).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: errResp.Message,
			})
		}
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Error interno",
		})
	}

	logging.INFO("Tenant creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Tenant creado con éxito",
	})
}
