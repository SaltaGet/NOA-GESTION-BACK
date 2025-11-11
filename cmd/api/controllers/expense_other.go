package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// ExpenseOtherGetByID godoc
//
//	@Summary		ExpenseOtherGetByID
//	@Description	Obtener un egreso por ID
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string											true	"ID of ExpenseOther"
//	@Success		200	{object}	schemas.Response{body=schemas.ExpenseOtherResponse}	"ExpenseOther obtained successfully"
//	@Router			/api/v1/expense_other/{id} [get]
func (e *ExpenseOtherController) ExpenseOtherGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un egreso por ID")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	expenseOther, err := e.ExpenseOtherService.ExpenseOtherGetByID(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Egreso obtenido con éxito")
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
//	@Param			fromDate	query		schemas.DateRangeRequest	true	"Fecha de inicio"
//	@Param			page	query		int											false	"Page number"				default(1)
//	@Param			limit	query		int											false	"Number of items per page"	default(20)
//	@Success		200		{object}	schemas.Response{body=[]schemas.ExpenseOtherResponseSimple}	"List of expenseOthers"
//	@Router			/api/v1/expense_other/get_all [get]
func (e *ExpenseOtherController) ExpenseOtherGetByDate(c *fiber.Ctx) error {
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

	pointID :=c.Locals("point_sale_id").(int64)

	expenseOthers, total, err := e.ExpenseOtherService.ExpenseOtherGetByDate(&pointID, fromDate, toDate, page, limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Egresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"expense_others": expenseOthers, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
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
//	@Param			expenseOtherCreate	body		schemas.ExpenseOtherCreate			true	"ExpenseOther information"
//	@Success		200				{object}	schemas.Response"ExpenseOther created successfully"
//	@Router			/api/v1/expense_other/create [post]
func (e *ExpenseOtherController) ExpenseOtherCreate(c *fiber.Ctx) error {
	logging.INFO("Crear un egreso")
	var expenseOtherCreate schemas.ExpenseOtherCreate
	if err := c.BodyParser(&expenseOtherCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
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
	pointID :=c.Locals("point_sale_id").(int64)

	id, err := e.ExpenseOtherService.ExpenseOtherCreate(user.ID, pointID, &expenseOtherCreate)
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

// ExpenseOtherUpdate godoc
//
//	@Summary		ExpenseOtherUpdate
//	@Description	Actualizar un egreso de compra
//	@Tags			ExpenseOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			expenseOtherUpdate	body		schemas.ExpenseOtherUpdate	true	"ExpenseOther data to update"
//	@Success		200				{object}	schemas.Response		"ExpenseOther updated successfully"
//	@Router			/api/v1/expense_other/update [put]
func (e *ExpenseOtherController) ExpenseOtherUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar un egreso")
	var expenseOtherUpdate schemas.ExpenseOtherUpdate
	if err := c.BodyParser(&expenseOtherUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
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

	err := e.ExpenseOtherService.ExpenseOtherUpdate(user.ID, &expenseOtherUpdate)
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
//	@Router			/api/v1/expenseOther/delete/{id} [delete]
func (e *ExpenseOtherController) ExpenseOtherDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar un egreso")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	err = e.ExpenseOtherService.ExpenseOtherDelete(idint)
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
