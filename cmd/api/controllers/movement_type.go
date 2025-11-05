package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// GetMovementTypeByID godoc
//
//	@Summary		Get Movement Type By ID
//	@Description	Get Movement Type By ID
//	@Tags			Movement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string										true	"ID of the movement type"
//	@Success		200	{object}	schemas.Response{body=schemas.MovementType}	"Movement type details"
//	@Failure		400	{object}	schemas.Response							"Bad Request"
//	@Failure		401	{object}	schemas.Response							"Auth is required"
//	@Failure		403	{object}	schemas.Response							"Not Authorized"
//	@Failure		404	{object}	schemas.Response							"Expense not found"
//	@Failure		500	{object}	schemas.Response							"Internal server error"
//	@Router			/api/v1/movement/{id} [get]
func (m *MovementTypeController) GetMovementTypeByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un movimiento por ID")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	movement, err := m.MovementTypeService.MovementTypeGetByID(id)
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

	logging.INFO("Movimiento obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    movement,
		Message: "Movimiento obtenido con éxito",
	})
}

// GetAllMovementTypes godoc
//	@Summary		Get all movement types
//	@Description	Get all movement types from either laundry or workshop based on the provided isIncome query parameter.
//	@Tags			Movement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			isIncome	query		bool											true	"Is income movement type"
//	@Success		200			{object}	schemas.Response{body=[]schemas.MovementType}	"List of movement types"
//	@Failure		400			{object}	schemas.Response								"Bad Request"
//	@Failure		401			{object}	schemas.Response								"Auth is required"
//	@Failure		403			{object}	schemas.Response								"Not Authorized"
//	@Failure		404			{object}	schemas.Response								"Expense not found"
//	@Failure		500			{object}	schemas.Response								"Internal server error"
//	@Router			/api/v1/movement/get_all [get]
func (m *MovementTypeController) GetAllMovementTypes(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los movimientos")
	isIncomeStr := c.Query("isIncome")
	isIncome := false
	if isIncomeStr != "" {
		var err error
		isIncome, err = strconv.ParseBool(isIncomeStr)
		if err != nil {
			logging.ERROR("Error: %s", err.Error())
			return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "Invalid value for isIncome",
			})
		}
	}

	movements, err := m.MovementTypeService.MovementTypeGetAll(isIncome)
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

	logging.INFO("Movimientos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    movements,
		Message: "Movimientos obtenidos con éxito",
	})
}

// MovementTypeCreate godoc
//
//	@Summary		Create Movement Type
//	@Description	This endpoint creates a new movement type based on the provided JSON payload.
//	@Tags			Movement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			movementType	body		schemas.MovementTypeCreate		true	"Movement Type Details"
//	@Success		200				{object}	schemas.Response{body=string}	"Movement created successfully"
//	@Failure		400				{object}	schemas.Response				"Bad Request"
//	@Failure		401				{object}	schemas.Response				"Auth is required"
//	@Failure		403				{object}	schemas.Response				"Not Authorized"
//	@Failure		404				{object}	schemas.Response				"Expense not found"
//	@Failure		422				{object}	schemas.Response				"Model invalid"
//	@Failure		500				{object}	schemas.Response				"Internal server error"
//	@Router			/api/v1/movement/create [post]
func (m *MovementTypeController) MovementTypeCreate(c *fiber.Ctx) error {
	logging.INFO("Crear un movimiento")
	var movementCreate schemas.MovementTypeCreate
	if err := c.BodyParser(&movementCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := movementCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := m.MovementTypeService.MovementTypeCreate(&movementCreate)
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

	logging.INFO("Movimiento creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Movimiento creado con éxito",
	})
}

// MovementTypeUpdate godoc
//	@Summary		Update Movement Type
//	@Description	This endpoint updates a movement type based on the provided JSON payload.
//	@Tags			Movement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			movementType	body		schemas.MovementTypeUpdate	true	"Movement Type Details"
//	@Success		200				{object}	schemas.Response			"Movement updated successfully"
//	@Failure		400				{object}	schemas.Response			"Bad Request"
//	@Failure		401				{object}	schemas.Response			"Auth is required"
//	@Failure		403				{object}	schemas.Response			"Not Authorized"
//	@Failure		404				{object}	schemas.Response			"Expense not found"
//	@Failure		422				{object}	schemas.Response			"Model invalid"
//	@Failure		500				{object}	schemas.Response			"Internal server error"
//	@Router			/api/v1/movement/update [put]
func (m *MovementTypeController) MovementTypeUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar un movimiento")
	var movementUpdate schemas.MovementTypeUpdate
	if err := c.BodyParser(&movementUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := movementUpdate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	err := m.MovementTypeService.MovementTypeUpdate(&movementUpdate)
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

	logging.INFO("Movimiento editado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Movimiento editado con éxito",
	})
}

// MovementTypeDelete godoc
//	@Summary		Delete Movement Type
//	@Description	Deletes a movement type based on its ID.
//	@Tags			Movement
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"ID of the movement type"
//	@Success		200	{object}	schemas.Response	"Movement type deleted successfully"
//	@Failure		400	{object}	schemas.Response	"Bad Request"
//	@Failure		401	{object}	schemas.Response	"Auth is required"
//	@Failure		403	{object}	schemas.Response	"Not Authorized"
//	@Failure		404	{object}	schemas.Response	"Expense not found"
//	@Failure		500	{object}	schemas.Response	"Internal server error"
//	@Router			/api/v1/movement/delete/{id} [delete]
func (m *MovementTypeController) MovementTypeDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar un movimiento")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	err := m.MovementTypeService.MovementTypeDelete(id)
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

	logging.INFO("Movimiento eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Movimiento eliminado con éxito",
	})
}
