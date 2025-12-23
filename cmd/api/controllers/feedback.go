package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// FeedbackGetByID godoc
//
//	@Summary		FeedbackGetByID
//	@Description	Obtener feedback por id
//	@Tags			Feedback
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string											true	"Feedback ID"
//	@Success		200	{object}	schemas.Response{body=schemas.FeedbackResponse}	"Feedback obtenido con éxito"
//	@Router			/api/v1/feedback/get/{id} [get]
func (t *FeedbackController) FeedbackGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idInt, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	plans, err := t.FeedbackService.FeedbackGetByID(idInt)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    plans,
		Message: "Feedback obtenido con éxito",
	})
}

// FeedbackGetAll godoc
//
//	@Summary		FeedbackGetAll
//	@Description	Obtener todos los feedbacks
//	@Tags			Feedback
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.FeedbackResponseDTO}	"feedbackes obtenidos con éxito"
//	@Router			/api/v1/feedback/get_all [get]
func (t *FeedbackController) FeedbackGetAll(c *fiber.Ctx) error {
	plans, err := t.FeedbackService.FeedbackGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    plans,
		Message: "feedbacks obtenidos con éxito",
	})
}

// FeedbackCreate godoc
//
//	@Summary		FeedbackCreate
//	@Description	Crear Feedback
//	@Tags			Feedback
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			new_create	body		schemas.FeedbackCreate	true	"Crear novead"
//	@Success		200			{object}	schemas.Response		"feedback creada con éxito"
//	@Router			/api/v1/feedback/create [post]
func (t *FeedbackController) FeedbackCreate(c *fiber.Ctx) error {
	var newCreate schemas.FeedbackCreate
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

	id, err := t.FeedbackService.FeedbackCreate(&newCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Feedback creado con éxito",
	})
}
