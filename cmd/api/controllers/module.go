package controllers

import (
	"fmt"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ModuleGet godoc
//
//	@Summary		ModuleGet
//	@Description	Obtener un modulo por ID
//	@Tags			Module
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string											true	"ID of ExpenseOther"
//	@Success		200	{object}	schemas.Response{body=schemas.ModuleResponse}	"ExpenseOther obtained successfully"
//	@Router			/api/v1/module/get/{id} [get]
func (m *ModuleController) ModuleGet(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Se necesita el id del modulo", fmt.Errorf("se necesita el id del modulo")))
	}
	uuid, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	module, err := m.ModuleService.ModuleGet(uuid)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    module,
		Message: "Modulo obtenido con éxito",
	})
}

// ModuleGetAll godoc
//
//	@Summary		ModuleGetAll
//	@Description	Obtener todos los modulos
//	@Tags			Module
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.Response{body=[]schemas.ModuleResponse}	"ExpenseOther obtained successfully"
//	@Router			/api/v1/module/get_all [get]
func (m *ModuleController) ModuleGetAll(c *fiber.Ctx) error {
	modules, err := m.ModuleService.ModuleGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    modules,
		Message: "Modulos obtenidos con éxito",
	})
}

// ModuleCreate godoc
//
//	@Summary		ModuleCreate
//	@Description	Crear nuevo modulo
//	@Tags			Module
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		schemas.ModuleCreate	true	"Modulo"
//	@Success		200		{object}	schemas.Response		"ExpenseOther obtained successfully"
//	@Router			/api/v1/module/create [post]
func (m *ModuleController) ModuleCreate(c *fiber.Ctx) error {
	var moduleCrate schemas.ModuleCreate
	if err := c.BodyParser(&moduleCrate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := moduleCrate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	id, err := m.ModuleService.ModuleCreate(&moduleCrate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Modulo creado con éxito",
	})
}

// ModuleUpdate godoc
//
//	@Summary		ModuleUpdate
//	@Description	Editar modulo
//	@Tags			Module
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		schemas.ModuleUpdate	true	"Modulo"
//	@Success		200		{object}	schemas.Response		"ExpenseOther obtained successfully"
//	@Router			/api/v1/module/update [put]
func (m *ModuleController) ModuleUpdate(c *fiber.Ctx) error {
	var moduleUpdate schemas.ModuleUpdate
	if err := c.BodyParser(&moduleUpdate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := moduleUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := m.ModuleService.ModuleUpdate(&moduleUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Modulo actualizado con éxito",
	})
}

func (m *ModuleController) ModuleDelete(c *fiber.Ctx) error {
	return nil
}

// ModuleAddTenant godoc
//
//	@Summary		ModuleAddTenant
//	@Description	Agregar un modulo a tenant
//	@Tags			Module
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		schemas.ModuleAddTenant	true	"Modulo"
//	@Success		200		{object}	schemas.Response		"ExpenseOther obtained successfully"
//	@Router			/api/v1/module/add_tenant_expiration [put]
func (m *ModuleController) ModuleAddTenant(c *fiber.Ctx) error {
	var moduleAddTenant schemas.ModuleAddTenant
	if err := c.BodyParser(&moduleAddTenant); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := moduleAddTenant.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := m.ModuleService.ModuleAddTenant(&moduleAddTenant)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Modulo actualizado con éxito",
	})
}
