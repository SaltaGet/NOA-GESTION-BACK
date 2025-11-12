package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// IncomeOtherGetByID godoc
//	@Summary		IncomeOtherGetByID
//	@Description	Obtiene ingreso de venta por ID
//	@Tags			IncomeOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string												true	"ID of the incomeOther"
//	@Success		200	{object}	schemas.Response{body=schemas.IncomeOtherResponse}	"IncomeOther details fetched successfully"
//	@Router			/api/v1/income_other/{id} [get]
func (i *IncomeOtherController) IncomeOtherGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un ingreso por ID")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	incomeOther, err := i.IncomeOtherService.IncomeOtherGetByID(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Ingreso obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    incomeOther,
		Message: "Ingreso obtenido con éxito",
	})
}

// IncomeOtherGetByDate godoc
//	@Summary		IncomeOtherGetByDate
//	@Description	Obtiene ingresos de venta por fecha
//	@Tags			IncomeOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page		query		int														false	"Page number"				default(1)
//	@Param			limit		query		int														false	"Number of items per page"	default(20)
//	@Param			fromDate	body		schemas.DateRangeRequest								true	"Fecha de inicio"
//	@Success		200			{object}	schemas.Response{body=[]schemas.IncomeOtherResponse}	"List of incomeOthers"
//	@Router			/api/v1/income_other/get_all [get]
func (i *IncomeOtherController) IncomeOtherGetByDate(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los ingresos")

	dateTime := &schemas.DateRangeRequest{}
	if err := c.QueryParser(dateTime); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "error al obtener la fecha", err))
	}
	dateFrom, dateTo, err :=dateTime.GetParsedDates()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pageParam := c.Query("page", "1")
	limitParam := c.Query("limit", "10")

	page, err := strconv.Atoi(pageParam)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit < 1 {
		limit = 20
	}

	var pointSaleID *int64

if pointID, ok := c.Locals("point_sale_id").(int64); ok {
    pointSaleID = &pointID
}

// pointSaleID ahora es un puntero a int64 O es nil.

	incomeOthers, total, err := i.IncomeOtherService.IncomeOtherGetByDate(pointSaleID, dateFrom, dateTo, page, limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Ingresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"income_dales": incomeOthers, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Egresos obtenidos con éxito",
	})
}

// IncomeOtherCreate godoc
//	@Summary		IncomeOtherCreate
//	@Description	Crea un nuevo ingreso de venta
//	@Tags			IncomeOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			incomeOtherCreate	body		schemas.IncomeOtherCreate	true	"IncomeOther information"
//	@Success		200					{object}	schemas.Response			"IncomeOther created successfully"
//	@Router			/api/v1/income_other/create [post]
func (i *IncomeOtherController) IncomeOtherCreate(c *fiber.Ctx) error {
	logging.INFO("Crear ingreso")
	var incomeOtherCreate schemas.IncomeOtherCreate
	if err := c.BodyParser(&incomeOtherCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := incomeOtherCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	pointID :=c.Locals("point_sale_id").(int64)

	id, err := i.IncomeOtherService.IncomeOtherCreate(user.ID, pointID, &incomeOtherCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Ingreso creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Ingreso creado con éxito",
	})
}

// IncomeOtherUpdate godoc
//	@Summary		IncomeOtherUpdate
//	@Description	Actualiza un ingreso de venta
//	@Tags			IncomeOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			incomeOtherUpdate	body		schemas.IncomeOtherUpdate	true	"IncomeOther data to update"
//	@Success		200					{object}	schemas.Response			"IncomeOther updated successfully"
//	@Router			/api/v1/income_other/update [put]
func (i *IncomeOtherController) IncomeOtherUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar ingreso")
	var incomeOtherUpdate schemas.IncomeOtherUpdate
	if err := c.BodyParser(&incomeOtherUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := incomeOtherUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	pointID :=c.Locals("point_sale_id").(int64)

	err := i.IncomeOtherService.IncomeOtherUpdate(user.ID, pointID, &incomeOtherUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Ingreso editado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Ingreso editado con éxito",
	})
}

// IncomeOtherDelete godoc
//	@Summary		IncomeOtherDelete
//	@Description	Elimina un ingreso de venta
//	@Tags			IncomeOther
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"ID of the incomeOther"
//	@Success		200	{object}	schemas.Response	"IncomeOther deleted successfully"
//	@Router			/api/v1/income_other/delete/{id} [delete]
func (i *IncomeOtherController) IncomeOtherDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar ingreso")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pointID :=c.Locals("point_sale_id").(int64)

	err = i.IncomeOtherService.IncomeOtherDelete(idint, pointID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Ingreso eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Ingreso eliminado con éxito",
	})
}
