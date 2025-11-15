package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

//	TypeMovementGetAll godoc
//
//	@Summary		TypeMovementGetAll
//	@Description	Retorna los tipos de movimientos de ingresos o egresos
//	@Tags			TypeMovement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			type_movement	query		string													true	"TypeMovement"	enum(income,expense)
//	@Success		200				{object}	schemas.Response{body=[]schemas.TypeMovementResponse}	"Movimientos obtenidos con éxito"
//	@Router			/api/v1/type_movement/get_all [get]
func (t *TypeMovementController) TypeMovementGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los tipos de movimientos")
	typeMovement := c.Query("type_movement")
	if typeMovement != "income" && typeMovement != "expense" {
		logging.ERROR("Error: type_movement es requerido")
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Error: type_movement es requerido",
		})
	}

	movement, err := t.TypeMovementService.TypeMovementGetAll(typeMovement)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Tipos de movimientos obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    movement,
		Message: "Tipos de movimientos obtenidos con éxito",
	})
}

//	TypeMovementCreate godoc
//
//	@Summary		TypeMovementCreate
//	@Description	Crea un nuevo tipo de movimiento
//	@Tags			TypeMovement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			movement_create	body		schemas.TypeMovementCreate	true	"MovementCreate"
//	@Success		200				{object}	schemas.Response
//	@Router			/api/v1/type_movement/create [post]
func (t *TypeMovementController) TypeMovementCreate(c *fiber.Ctx) error {
	movementCreate := schemas.TypeMovementCreate{}
	if err := c.BodyParser(&movementCreate); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el modelo", err))
	}
	if err := movementCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := t.TypeMovementService.TypeMovementCreate(movementCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Movimiento creado con éxito",
	})
}

//	TypeMovementUpdate godoc
//
//	@Summary		TypeMovementUpdate
//	@Description	Actualiza un tipo de movimiento
//	@Tags			TypeMovement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			movement_update	body		schemas.TypeMovementUpdate	true	"MovementUpdate"
//	@Success		200				{object}	schemas.Response
//	@Router			/api/v1/type_movement/update [put]
func (t *TypeMovementController) TypeMovementUpdate(c *fiber.Ctx) error {
	movementUpdate := schemas.TypeMovementUpdate{}
	if err := c.BodyParser(&movementUpdate); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el modelo", err))
	}
	if err := movementUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := t.TypeMovementService.TypeMovementUpdate(movementUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Movimiento actualizado con éxito",
	})
}
