package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// GetIncomeByID godoc
//	@Summary		Get Income By ID
//	@Description	Fetches income details from based on the provided ID and tenant context.
//	@Tags			Income
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string											true	"ID of the income"
//	@Success		200	{object}	schemas.Response{body=schemas.IncomeResponse}	"Income details fetched successfully"
//	@Failure		400	{object}	schemas.Response								"Bad Request"
//	@Failure		401	{object}	schemas.Response								"Auth is required"
//	@Failure		403	{object}	schemas.Response								"Not Authorized"
//	@Failure		404	{object}	schemas.Response								"Expense not found"
//	@Failure		500	{object}	schemas.Response								"Internal server error"
//	@Router			/income/{id} [get]
func (i *IncomeController) GetIncomeByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un ingreso por ID")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	income, err := i.IncomeService.IncomeGetByID(id)
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

	logging.INFO("Ingreso obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    income,
		Message: "Ingreso obtenido con éxito",
	})
}

// GetAllIncomes godoc
//	@Summary		Get all incomes
//	@Description	Fetches all incomes from the specified tenant.
//	@Tags			Income
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int											false	"Page number"				default(1)
//	@Param			limit	query		int											false	"Number of items per page"	default(20)
//	@Success		200		{object}	schemas.Response{body=[]schemas.IncomeDTO}	"List of incomes"
//	@Failure		400		{object}	schemas.Response							"Bad Request"
//	@Failure		401		{object}	schemas.Response							"Auth is required"
//	@Failure		403		{object}	schemas.Response							"Not Authorized"
//	@Failure		404		{object}	schemas.Response							"Expense not found"
//	@Failure		500		{object}	schemas.Response							"Internal server error"
//	@Router			/income/get_all [get]
func (i *IncomeController) GetAllIncomes(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los ingresos")
	pageParam := c.Query("page", "1")
	limitParam := c.Query("limit", "20")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 20
	}

	incomes, err := i.IncomeService.IncomeGetAll(page, limit)
	if err != nil {
		if errResp, ok := err.(*schemas.ErrorStruc); ok {
			return c.Status(errResp.StatusCode).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: errResp.Message,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Error interno",
		})
	}

	logging.INFO("Ingresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    incomes,
		Message: "Ingresos obtenidos con éxito",
	})
}

// GetIncomeToday godoc
//	@Summary		Get Income Today
//	@Description	Fetches all incomes from the specified tenant, on the current day.
//	@Tags			Income
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.IncomeDTO}	"List of all incomes"
//	@Failure		400	{object}	schemas.Response							"Bad Request"
//	@Failure		401	{object}	schemas.Response							"Auth is required"
//	@Failure		403	{object}	schemas.Response							"Not Authorized"
//	@Failure		404	{object}	schemas.Response							"Expense not found"
//	@Failure		500	{object}	schemas.Response							"Internal server error"
//	@Router			/income/get_today [get]
func (i *IncomeController) GetIncomeToday(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los ingresos de hoy")
	pageParam := c.Query("page", "1")
	limitParam := c.Query("limit", "20")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 20
	}

	incomes, err := i.IncomeService.IncomeGetToday(page, limit)
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

	logging.INFO("Ingresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    incomes,
		Message: "Ingresos obtenidos con éxito",
	})
}

// CreateIncome godoc
//	@Summary		Create Income
//	@Description	Parses the request body to create a new income entry for either laundry or workshop.
//	@Tags			Income
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			incomeCreate	body		schemas.IncomeCreate			true	"Income information"
//	@Success		200				{object}	schemas.Response{body=string}	"Income created successfully"
//	@Failure		400				{object}	schemas.Response				"Bad Request"
//	@Failure		401				{object}	schemas.Response				"Auth is required"
//	@Failure		403				{object}	schemas.Response				"Not Authorized"
//	@Failure		404				{object}	schemas.Response				"Expense not found"
//	@Failure		422				{object}	schemas.Response				"Model Invalid"
//	@Failure		500				{object}	schemas.Response				"Internal server error"
//	@Router			/income/create [post]
func (i *IncomeController) CreateIncome(c *fiber.Ctx) error {
	logging.INFO("Crear ingreso")
	var incomeCreate schemas.IncomeCreate
	if err := c.BodyParser(&incomeCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := incomeCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := i.IncomeService.IncomeCreate(&incomeCreate)
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

	logging.INFO("Ingreso creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Ingreso creado con éxito",
	})
}

// UpdateIncome godoc
//	@Summary		Update Income
//	@Description	Updates the details of an income based on the provided data.
//	@Tags			Income
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			incomeUpdate	body		schemas.IncomeUpdate	true	"Income data to update"
//	@Success		200				{object}	schemas.Response		"Income updated successfully"
//	@Failure		400				{object}	schemas.Response		"Bad Request"
//	@Failure		401				{object}	schemas.Response		"Auth is required"
//	@Failure		403				{object}	schemas.Response		"Not Authorized"
//	@Failure		404				{object}	schemas.Response		"Expense not found"
//	@Failure		422				{object}	schemas.Response		"Model Invalid"
//	@Failure		500				{object}	schemas.Response		"Internal server error"
//	@Router			/income/update [put]
func (i *IncomeController) UpdateIncome(c *fiber.Ctx) error {
	logging.INFO("Actualizar ingreso")
	var incomeUpdate schemas.IncomeUpdate
	if err := c.BodyParser(&incomeUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := incomeUpdate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	err := i.IncomeService.IncomeUpdate(&incomeUpdate)
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

	logging.INFO("Ingreso editado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Ingreso editado con éxito",
	})
}

// DeleteIncome godoc
//	@Summary		Delete Income
//	@Description	Deletes an income entry based on the provided ID and workplace context.
//	@Tags			Income
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"ID of the income"
//	@Success		200	{object}	schemas.Response	"Income deleted successfully"
//	@Failure		400	{object}	schemas.Response	"Bad Request"
//	@Failure		401	{object}	schemas.Response	"Auth is required"
//	@Failure		403	{object}	schemas.Response	"Not Authorized"
//	@Failure		404	{object}	schemas.Response	"Expense not found"
//	@Failure		500	{object}	schemas.Response	"Error interno"
//	@Router			/income/delete/{id} [delete]
func (i *IncomeController) DeleteIncome(c *fiber.Ctx) error {
	logging.INFO("Eliminar ingreso")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	err := i.IncomeService.IncomeDelete(id)
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

	logging.INFO("Ingreso eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Ingreso eliminado con éxito",
	})
}
