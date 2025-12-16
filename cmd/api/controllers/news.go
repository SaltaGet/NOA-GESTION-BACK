package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// NewsGetByID godoc
//
//	@Summary		NewsGetByID
//	@Description	Obtene novedad por id
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string											true	"News ID"
//	@Success		200	{object}	schemas.Response{body=[]schemas.NewsResponse}	"Novedad obtenida con éxito"
//	@Router			/api/v1/news/get/{id} [get]
func (t *NewsController) NewsGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	plans, err := t.NewsService.NewsGetByID(idInt)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    plans,
		Message: "Novedad obtenida con éxito",
	})
}

// NewsGetAll godoc
//
//	@Summary		NewsGetAll
//	@Description	Obtener todos los planes
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.NewsResponseDTO}	"novedades obtenidas con éxito"
//	@Router			/api/v1/news/get_all [get]
func (t *NewsController) NewsGetAll(c *fiber.Ctx) error {
	plans, err := t.NewsService.NewsGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    plans,
		Message: "planes obtenidos con éxito",
	})
}

// NewsCreate godoc
//
//	@Summary		NewsCreate
//	@Description	Crear Novedad
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			new_create	body		schemas.NewsCreate	true	"Crear novead"
//	@Success		200			{object}	schemas.Response	"novedad creada con éxito"
//	@Router			/api/v1/news/create [post]
func (t *NewsController) NewsCreate(c *fiber.Ctx) error {
	var newCreate schemas.NewsCreate
	if err := c.BodyParser(&newCreate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := newCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*models.Admin)

	id, err := t.NewsService.NewsCreate(admin.ID, &newCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "News creado con éxito",
	})
}

// NewsUpdate godoc
//
//	@Summary		NewsUpdate
//	@Description	Actualizar novedad
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			new_update	body		schemas.NewsUpdate	true	"Editar novedad"
//	@Success		200			{object}	schemas.Response	"novedad editada con éxito"
//	@Router			/api/v1/news/update [put]
func (t *NewsController) NewsUpdate(c *fiber.Ctx) error {
	var newUpdate schemas.NewsUpdate
	if err := c.BodyParser(&newUpdate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := newUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*models.Admin)
	err := t.NewsService.NewsUpdate(admin.ID, &newUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "News actualizado con éxito",
	})
}

// NewsDelete godoc
//
//	@Summary		NewsDelete
//	@Description	eliminar novedad
//	@Tags			News
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"id novedad"
//	@Success		200	{object}	schemas.Response	"novedad eliminada con éxito"
//	@Router			/api/v1/news/delete/{id} [put]
func (t *NewsController) NewsDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	admin := c.Locals("user_admin").(*models.Admin)
	err = t.NewsService.NewsDelete(admin.ID, idInt)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "News eliminada con éxito",
	})
}
