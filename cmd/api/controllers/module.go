package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"fmt"
)


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