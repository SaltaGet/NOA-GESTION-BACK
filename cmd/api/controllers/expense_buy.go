package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// ExpenseBuyGetByID godoc
//
//	@Summary		ExpenseBuyGetByID
//	@Description	Get ExpenseBuy By ID
//	@Tags			ExpenseBuy
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string												true	"ID of ExpenseBuy"
//	@Success		200	{object}	schemas.Response{body=schemas.ExpenseBuyResponse}	"ExpenseBuy obtained successfully"
//	@Router			/api/v1/expense_buy/{id} [get]
func (e *ExpenseBuyController) ExpenseBuyGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un egreso por ID")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	expenseBuy, err := e.ExpenseBuyService.ExpenseBuyGetByID(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Egreso obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    expenseBuy,
		Message: "Egreso obtenido con éxito",
	})
}

// ExpenseBuyGetByDate godoc
//
//	@Summary		ExpenseBuyGetByDate
//	@Description	Obtiene los egresos de compras por fecha
//	@Tags			ExpenseBuy
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			fromDate	query		schemas.DateRangeRequest									true	"Fecha de inicio"
//	@Param			page		query		int															false	"Page number"				default(1)
//	@Param			limit		query		int															false	"Number of items per page"	default(20)
//	@Success		200			{object}	schemas.Response{body=[]schemas.ExpenseBuyResponseSimple}	"List of expenseBuys"
//	@Failure		400			{object}	schemas.Response											"Bad Request"
//	@Failure		401			{object}	schemas.Response											"Auth is required"
//	@Failure		403			{object}	schemas.Response											"Not Authorized"
//	@Failure		500			{object}	schemas.Response											"Internal server error"
//	@Router			/api/v1/expense_buy/get_all [get]
func (e *ExpenseBuyController) ExpenseBuyGetByDate(c *fiber.Ctx) error {
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

	formDate := &schemas.DateRangeRequest{}
	if err := c.QueryParser(formDate); err != nil {
		return schemas.HandleError(c, err)
	}
	fromDate, toDate, err := formDate.GetParsedDates()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	expenseBuys, total, err := e.ExpenseBuyService.ExpenseBuyGetByDate(fromDate, toDate, page, limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Egresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"expense_buys": expenseBuys, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Egresos obtenidos con éxito",
	})
}

// ExpenseBuyCreate godoc
//
//	@Summary		ExpenseBuyCreate
//	@Description	Crear un egreso de compra
//	@Tags			ExpenseBuy
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			expenseBuyCreate	body		schemas.ExpenseBuyCreate	true	"ExpenseBuy information"
//	@Success		200					{object}	schemas.Response"ExpenseBuy created successfully"
//	@Router			/api/v1/expense_buy/create [post]
func (e *ExpenseBuyController) ExpenseBuyCreate(c *fiber.Ctx) error {
	logging.INFO("Crear un egreso")
	var expenseBuyCreate schemas.ExpenseBuyCreate
	if err := c.BodyParser(&expenseBuyCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseBuyCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	id, err := e.ExpenseBuyService.ExpenseBuyCreate(user.ID, &expenseBuyCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Egreso creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Egreso creado con éxito",
	})
}

// ExpenseBuyUpdate godoc
//
//	@Summary		ExpenseBuyUpdate
//	@Description	Actualizar un egreso de compra
//	@Tags			ExpenseBuy
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			expenseBuyUpdate	body		schemas.ExpenseBuyUpdate	true	"ExpenseBuy data to update"
//	@Success		200					{object}	schemas.Response			"ExpenseBuy updated successfully"
//	@Router			/api/v1/expense_buy/update [put]
func (e *ExpenseBuyController) ExpenseBuyUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar un egreso")
	var expenseBuyUpdate schemas.ExpenseBuyUpdate
	if err := c.BodyParser(&expenseBuyUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseBuyUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	err := e.ExpenseBuyService.ExpenseBuyUpdate(user.ID, &expenseBuyUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Egreso editado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso editado con éxito",
	})
}

// ExpenseBuyDelete godoc
//
//	@Summary		ExpenseBuyDelete
//	@Description	Eliminar un egreso de compra
//	@Tags			ExpenseBuy
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"ID of the expenseBuy"
//	@Success		200	{object}	schemas.Response	"ExpenseBuy deleted successfully"
//	@Router			/api/v1/expenseBuy/delete/{id} [delete]
func (e *ExpenseBuyController) ExpenseBuyDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar un egreso")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	err = e.ExpenseBuyService.ExpenseBuyDelete(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Egreso eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso eliminado con éxito",
	})
}
