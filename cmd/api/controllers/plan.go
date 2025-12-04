package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

//	PlanGetAll godoc
//
//	@Summary		PlanGetAll
//	@Description	Obtener todos los planes
//	@Tags			Plan
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PlanResponseDTO}	"planes obtenidos con éxito"
//	@Router			/api/v1/plan/get_all [get]
func (t *PlanController) PlanGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los planes")
	plans, err := t.PlanService.PlanGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Planes obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    plans,
		Message: "planes obtenidos con éxito",
	})
}

//	PlanGetByID godoc
//
//	@Summary		PlanGetByID
//	@Description	Obtener plan por id
//	@Tags			Plan
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string										true	"Plan ID"
//	@Success		200	{object}	schemas.Response{body=schemas.PlanResponse}	"planes obtenidos con éxito"
//	@Router			/api/v1/plan/get/{id} [get]
func (t *PlanController) PlanGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener plan por id")
	id := c.Params("id")
	idInt, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	plan, err := t.PlanService.PlanGetByID(idInt)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Plan obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    plan,
		Message: "plane obtenido con éxito",
	})
}

//	PlanCreate godoc
//
//	@Summary		PlanCreate
//	@Description	Crear plan
//	@Tags			Plan
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			plan_create	body		schemas.PlanCreate	true	"Crear plan"
//	@Success		200			{object}	schemas.Response	"plan creado con éxito"
//	@Router			/api/v1/plan/create [post]
func (t *PlanController) PlanCreate(c *fiber.Ctx) error {
	logging.INFO("Crear plan")
	var planCreate schemas.PlanCreate
	if err := c.BodyParser(&planCreate); err != nil {
		logging.ERROR("Invalid request %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := planCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	id, err := t.PlanService.PlanCreate(&planCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Plan creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Plan creado con éxito",
	})
}

//	PlanUpdate godoc
//
//	@Summary		PlanUpdate
//	@Description	Actualizar plan
//	@Tags			Plan
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			plan_create	body		schemas.PlanUpdate	true	"Crear plan"
//	@Success		200			{object}	schemas.Response	"plan creado con éxito"
//	@Router			/api/v1/plan/update [put]
func (t *PlanController) PlanUpdate(c *fiber.Ctx) error {
	logging.INFO("Crear plan")
	var planUpdate schemas.PlanUpdate
	if err := c.BodyParser(&planUpdate); err != nil {
		logging.ERROR("Invalid request %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := planUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := t.PlanService.PlanUpdate(&planUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Plan actualizado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Plan actualizado con éxito",
	})
}