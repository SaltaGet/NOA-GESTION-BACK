package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ExpenseOtherGetByID godoc
//
//	@Summary		ExpenseOtherGetByID
//	@Description	Obtener un egreso por ID
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string												true	"ID of ExpenseOther"
//	@Success		200	{object}	schemas.Response{body=schemas.ExpenseOtherResponse}	"ExpenseOther obtained successfully"
//	@Router			/api/v1/expense_other/get/{id} [get]
func (e *ExpenseOtherController) ExpenseOtherGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	expenseOther, err := e.ExpenseOtherService.ExpenseOtherGetByID(idint, nil)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    expenseOther,
		Message: "Egreso obtenido con éxito",
	})
}

// ExpenseOtherGetByDate godoc
//
//	@Summary		ExpenseOtherGetByDate
//	@Description	Obtiene los egresos por fecha
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			fromDate	query		schemas.DateRangeRequest									true	"Fecha de inicio"
//	@Param			page		query		int															false	"Page number"				default(1)
//	@Param			limit		query		int															false	"Number of items per page"	default(20)
//	@Success		200			{object}	schemas.Response{body=[]schemas.ExpenseOtherResponseDTO}	"List of expenseOthers"
//	@Router			/api/v1/expense_other/get_by_date [get]
func (e *ExpenseOtherController) ExpenseOtherGetByDate(c *fiber.Ctx) error {
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

	dateTime := &schemas.DateRangeRequest{}
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	dateTime.FromDate = fromDate
	dateTime.ToDate = toDate

	dateFrom, dateTo, err := dateTime.GetParsedDates()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	expenseOthers, total, err := e.ExpenseOtherService.ExpenseOtherGetByDate(nil, dateFrom, dateTo, page, limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": expenseOthers, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Egresos obtenidos con éxito",
	})
}

// ExpenseOtherCreate godoc
//
//	@Summary		ExpenseOtherCreate
//	@Description	Crear un egreso de compra
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			expenseOtherCreate	body		schemas.ExpenseOtherCreate	true	"ExpenseOther information"
//	@Success		200					{object}	schemas.Response"ExpenseOther created successfully"
//	@Router			/api/v1/expense_other/create [post]
func (e *ExpenseOtherController) ExpenseOtherCreate(c *fiber.Ctx) error {
	var expenseOtherCreate schemas.ExpenseOtherCreate
	if err := c.BodyParser(&expenseOtherCreate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseOtherCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	id, err := e.ExpenseOtherService.ExpenseOtherCreate(user.ID, nil, &expenseOtherCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Egreso creado con éxito",
	})
}

// ExpenseOtherUpdate godoc
//
//	@Summary		ExpenseOtherUpdate
//	@Description	Actualizar un egreso de compra
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			expenseOtherUpdate	body		schemas.ExpenseOtherUpdate	true	"ExpenseOther data to update"
//	@Success		200					{object}	schemas.Response			"ExpenseOther updated successfully"
//	@Router			/api/v1/expense_other/update [put]
func (e *ExpenseOtherController) ExpenseOtherUpdate(c *fiber.Ctx) error {
	var expenseOtherUpdate schemas.ExpenseOtherUpdate
	if err := c.BodyParser(&expenseOtherUpdate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseOtherUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	err := e.ExpenseOtherService.ExpenseOtherUpdate(user.ID, nil, &expenseOtherUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso editado con éxito",
	})
}

// ExpenseOtherDelete godoc
//
//	@Summary		ExpenseOtherDelete
//	@Description	Eliminar un egreso de compra
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"ID of the expenseOther"
//	@Success		200	{object}	schemas.Response	"ExpenseOther deleted successfully"
//	@Router			/api/v1/expense_other/delete/{id} [delete]
func (e *ExpenseOtherController) ExpenseOtherDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)
	err = e.ExpenseOtherService.ExpenseOtherDelete(member.ID, idint, nil)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso eliminado con éxito",
	})
}

// POINTSALE

// ExpenseOtherGetByIDPointSale godoc
//
//	@Summary		ExpenseOtherGetByIDPointSale
//	@Description	Obtener un egreso por ID de punto de venta
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string												true	"ID of ExpenseOther"
//	@Success		200	{object}	schemas.Response{body=schemas.ExpenseOtherResponse}	"ExpenseOther obtained successfully"
//	@Router			/api/v1/expense_other/get_point_sale/{id} [get]
func (e *ExpenseOtherController) ExpenseOtherGetByIDPointSale(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pointID := c.Locals("point_sale_id").(int64)

	expenseOther, err := e.ExpenseOtherService.ExpenseOtherGetByID(idint, &pointID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    expenseOther,
		Message: "Egreso obtenido con éxito",
	})
}

// ExpenseOtherGetByDatePointSale godoc
//
//	@Summary		ExpenseOtherGetByDatePointSale
//	@Description	Obtiene los egresos por fecha de punto de venta
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			fromDate	query		schemas.DateRangeRequest									true	"Fecha de inicio"
//	@Param			page		query		int															false	"Page number"				default(1)
//	@Param			limit		query		int															false	"Number of items per page"	default(20)
//	@Success		200			{object}	schemas.Response{body=[]schemas.ExpenseOtherResponseDTO}	"List of expenseOthers"
//	@Router			/api/v1/expense_other/get_by_date_point_sale [get]
func (e *ExpenseOtherController) ExpenseOtherGetByDatePointSale(c *fiber.Ctx) error {
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

	dateTime := &schemas.DateRangeRequest{}
	fromDate := c.Query("from_date")
	toDate := c.Query("to_date")
	dateTime.FromDate = fromDate
	dateTime.ToDate = toDate

	dateFrom, dateTo, err := dateTime.GetParsedDates()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pointID := c.Locals("point_sale_id").(int64)

	expenseOthers, total, err := e.ExpenseOtherService.ExpenseOtherGetByDate(&pointID, dateFrom, dateTo, page, limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": expenseOthers, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Egresos obtenidos con éxito",
	})
}

// ExpenseOtherCreatePointSale godoc
//
//	@Summary		ExpenseOtherCreatePointSale
//	@Description	Crear un egreso de compra de punto de venta
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			expenseOtherCreate	body		schemas.ExpenseOtherCreate	true	"ExpenseOther information"
//	@Success		200					{object}	schemas.Response"ExpenseOther created successfully"
//	@Router			/api/v1/expense_other/create_point_sale [post]
func (e *ExpenseOtherController) ExpenseOtherCreatePointSale(c *fiber.Ctx) error {
	var expenseOtherCreate schemas.ExpenseOtherCreate
	if err := c.BodyParser(&expenseOtherCreate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseOtherCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	pointID := c.Locals("point_sale_id").(int64)

	id, err := e.ExpenseOtherService.ExpenseOtherCreate(user.ID, &pointID, &expenseOtherCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Egreso creado con éxito",
	})
}

// ExpenseOtherUpdatePointSale godoc
//
//	@Summary		ExpenseOtherUpdatePointSale
//	@Description	Actualizar un egreso de compra de punto de venta
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			expenseOtherUpdate	body		schemas.ExpenseOtherUpdate	true	"ExpenseOther data to update"
//	@Success		200					{object}	schemas.Response			"ExpenseOther updated successfully"
//	@Router			/api/v1/expense_other/update_point_sale [put]
func (e *ExpenseOtherController) ExpenseOtherUpdatePointSale(c *fiber.Ctx) error {
	var expenseOtherUpdate schemas.ExpenseOtherUpdate
	if err := c.BodyParser(&expenseOtherUpdate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := expenseOtherUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	pointID := c.Locals("point_sale_id").(int64)

	err := e.ExpenseOtherService.ExpenseOtherUpdate(user.ID, &pointID, &expenseOtherUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso editado con éxito",
	})
}

// ExpenseOtherDeletePointSale godoc
//
//	@Summary		ExpenseOtherDeletePointSale
//	@Description	Eliminar un egreso de compra de punto de venta
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"ID of the expenseOther"
//	@Success		200	{object}	schemas.Response	"ExpenseOther deleted successfully"
//	@Router			/api/v1/expense_other/delete_point_sale/{id} [delete]
func (e *ExpenseOtherController) ExpenseOtherDeletePointSale(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pointID := c.Locals("point_sale_id").(int64)
	member := c.Locals("user").(*schemas.AuthenticatedUser)
	err = e.ExpenseOtherService.ExpenseOtherDelete(member.ID, idint, &pointID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Egreso eliminado con éxito",
	})
}
