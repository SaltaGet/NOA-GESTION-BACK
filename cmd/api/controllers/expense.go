package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// GetExpenseByID godoc
//	@Summary		Get Expense By ID
//	@Description	Get Expense By ID
//	@Tags			Expense
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string											true	"ID of Expense"
//	@Success		200	{object}	schemas.Response{body=schemas.ExpenseResponse}	"Expense obtained successfully"
//	@Failure		400	{object}	schemas.Response								"Bad Request"
//	@Failure		401	{object}	schemas.Response								"Auth is required"
//	@Failure		403	{object}	schemas.Response								"Not Authorized"
//	@Failure		404	{object}	schemas.Response								"Expense not found"
//	@Failure		500	{object}	schemas.Response
//	@Router			/expense/{id} [get]
func (e *ExpenseController) GetExpenseByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un egreso por ID")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	expense, err := e.ExpenseService.ExpenseGetByID(id)
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

	logging.INFO("Egreso obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    expense,
		Message: "Egreso obtenido con éxito",
	})
}

// GetAllExpenses godoc
//	@Summary		Get all expenses
//	@Description	Fetches all expenses from the specified tenant.
//	@Tags			Expense
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int											false	"Page number"				default(1)
//	@Param			limit	query		int											false	"Number of items per page"	default(20)
//	@Success		200		{object}	schemas.Response{body=[]schemas.ExpenseDTO}	"List of expenses"
//	@Failure		400		{object}	schemas.Response							"Bad Request"
//	@Failure		401		{object}	schemas.Response							"Auth is required"
//	@Failure		403		{object}	schemas.Response							"Not Authorized"
//	@Failure		500		{object}	schemas.Response							"Internal server error"
//	@Router			/expense/get_all [get]
func (e *ExpenseController) GetAllExpenses(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los egresos")
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

	expenses, err := e.ExpenseService.ExpenseGetAll(page, limit)
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

	logging.INFO("Egresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    expenses,
		Message: "Egresos obtenidos con éxito",
	})
}

// GetExpenseToday godoc
//	@Summary		Get expense today
//	@Description	Fetches all expenses from the specified tenant, on the current day.
//	@Tags			Expense
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			page	query		int											false	"Page number"				default(1)
//	@Param			limit	query		int											false	"Number of items per page"	default(20)
//	@Success		200		{object}	schemas.Response{body=[]schemas.ExpenseDTO}	"List of expenses"
//	@Failure		400		{object}	schemas.Response							"Bad Request"
//	@Failure		401		{object}	schemas.Response							"Auth is required"
//	@Failure		403		{object}	schemas.Response							"Not Authorized"
//	@Failure		500		{object}	schemas.Response							"Internal server error"
//	@Router			/expense/get_today [get]
func (e *ExpenseController) GetExpenseToday(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los egresos de hoy")
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
	
	expenses, err := e.ExpenseService.ExpenseGetToday(page, limit)
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

	logging.INFO("Egresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    expenses,
		Message: "Egresos obtenidos con éxito",
	})
}

// CreateExpense godoc
//	@Summary		Create Expense
//	@Description	Parses the request body to create a new expense entry.
//	@Tags			Expense
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			expenseCreate	body		schemas.ExpenseCreate			true	"Expense information"
//	@Success		200				{object}	schemas.Response{body=string}	"Expense created successfully"
//	@Failure		400				{object}	schemas.Response				"Bad Request"
//	@Failure		401				{object}	schemas.Response				"Auth is required"
//	@Failure		403				{object}	schemas.Response				"Not Authorized"
//	@Failure		422				{object}	schemas.Response				"Model Invalid"
//	@Failure		500				{object}	schemas.Response				"Internal server error"
//	@Router			/expense/create [post]
func (e *ExpenseController) CreateExpense(c *fiber.Ctx) error {
	logging.INFO("Crear un egreso")
	var expenseCreate schemas.ExpenseCreate
	if err := c.BodyParser(&expenseCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := e.ExpenseService.ExpenseCreate(&expenseCreate)
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

	logging.INFO("Egreso creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Egreso creado con éxito",
	})
}

// UpdateExpense godoc
//	@Summary		Update Expense
//	@Description	Updates the details of an expense based on the provided data.
//	@Tags			Expense
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			expenseUpdate	body		schemas.ExpenseUpdate	true	"Expense data to update"
//	@Success		200				{object}	schemas.Response		"Expense updated successfully"
//	@Failure		400				{object}	schemas.Response		"Bad Request"
//	@Failure		401				{object}	schemas.Response		"Auth is required"
//	@Failure		403				{object}	schemas.Response		"Not Authorized"
//	@Failure		422				{object}	schemas.Response		"Model Invalid"
//	@Failure		500				{object}	schemas.Response		"Internal server error"
//	@Router			/expense/update [put]
func (e *ExpenseController) UpdateExpense(c *fiber.Ctx) error {
	logging.INFO("Actualizar un egreso")
	var expenseUpdate schemas.ExpenseUpdate
	if err := c.BodyParser(&expenseUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseUpdate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	err := e.ExpenseService.ExpenseUpdate(&expenseUpdate)
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

	logging.INFO("Egreso editado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso editado con éxito",
	})
}

// DeleteExpense deletes an expense by its ID from the specified workplace.
//	@Summary		Delete Expense
//	@Description	Deletes an expense based on the provided ID and workplace context.
//	@Tags			Expense
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"ID of the expense"
//	@Success		200	{object}	schemas.Response	"Expense deleted successfully"
//	@Failure		400	{object}	schemas.Response	"Bad Request"
//	@Failure		401	{object}	schemas.Response	"Auth is required"
//	@Failure		403	{object}	schemas.Response	"Not Authorized"
//	@Failure		500	{object}	schemas.Response	"Internal server error"
//	@Router			/expense/delete/{id} [delete]
func (e *ExpenseController) DeleteExpense(c *fiber.Ctx) error {
	logging.INFO("Eliminar un egreso")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	err := e.ExpenseService.ExpenseDelete(id)
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

	logging.INFO("Egreso eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso eliminado con éxito",
	})
}

