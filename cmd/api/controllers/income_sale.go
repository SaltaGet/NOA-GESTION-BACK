package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// IncomeSaleGetByID godoc
//	@Summary		IncomeSaleGetByID
//	@Description	Obtiene ingreso de venta por ID
//	@Tags			IncomeSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string												true	"ID of the incomeSale"
//	@Success		200	{object}	schemas.Response{body=schemas.IncomeSaleResponse}	"IncomeSale details fetched successfully"
//	@Router			/api/v1/income_sale/{id} [get]
func (i *IncomeSaleController) IncomeSaleGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un ingreso por ID")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pointID := c.Locals("point_sale_id").(int64)

	incomeSale, err := i.IncomeSaleService.IncomeSaleGetByID(pointID, idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Ingreso obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    incomeSale,
		Message: "Ingreso obtenido con éxito",
	})
}

// IncomeSaleGetByDate godoc
//	@Summary		IncomeSaleGetByDate
//	@Description	Obtiene ingresos de venta por fecha
//	@Tags			IncomeSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page		query		int														false	"Page number"				default(1)
//	@Param			limit		query		int														false	"Number of items per page"	default(20)
//	@Param			fromDate	query		schemas.DateRangeRequest								true	"Fecha de inicio"
//	@Success		200			{object}	schemas.Response{body=[]schemas.IncomeSaleResponseDTO}	"List of incomeSales"
//	@Router			/api/v1/income_sale/get_by_date [get]
func (i *IncomeSaleController) IncomeSaleGetByDate(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los ingresos")

	formDate := &schemas.DateRangeRequest{}
	formDate.FromDate = c.Query("from_date")
	formDate.ToDate = c.Query("to_date")
	
	fromDate, toDate, err := formDate.GetParsedDates()
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

	pointID := c.Locals("point_sale_id").(int64)

	incomeSales, total, err := i.IncomeSaleService.IncomeSaleGetByDate(pointID, fromDate, toDate, page, limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Ingresos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": incomeSales, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Egresos obtenidos con éxito",
	})
}

// IncomeSaleCreate godoc
//	@Summary		IncomeSaleCreate
//	@Description	Crea un nuevo ingreso de venta
//	@Tags			IncomeSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			incomeSaleCreate	body		schemas.IncomeSaleCreate	true	"IncomeSale information"
//	@Success		200					{object}	schemas.Response			"IncomeSale created successfully"
//	@Router			/api/v1/income_sale/create [post]
func (i *IncomeSaleController) IncomeSaleCreate(c *fiber.Ctx) error {
	logging.INFO("Crear ingreso")
	var incomeSaleCreate schemas.IncomeSaleCreate
	if err := c.BodyParser(&incomeSaleCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := incomeSaleCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	pointID :=c.Locals("point_sale_id").(int64)

	id, err := i.IncomeSaleService.IncomeSaleCreate(user.ID, pointID, &incomeSaleCreate)
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

// IncomeSaleUpdate godoc
//	@Summary		IncomeSaleUpdate
//	@Description	Actualiza un ingreso de venta
//	@Tags			IncomeSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			incomeSaleUpdate	body		schemas.IncomeSaleUpdate	true	"IncomeSale data to update"
//	@Success		200					{object}	schemas.Response			"IncomeSale updated successfully"
//	@Router			/api/v1/income_sale/update [put]
func (i *IncomeSaleController) IncomeSaleUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar ingreso")
	var incomeSaleUpdate schemas.IncomeSaleUpdate
	if err := c.BodyParser(&incomeSaleUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := incomeSaleUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)
	pointID :=c.Locals("point_sale_id").(int64)

	err := i.IncomeSaleService.IncomeSaleUpdate(user.ID, pointID, &incomeSaleUpdate)
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

// IncomeSaleDelete godoc
//	@Summary		IncomeSaleDelete
//	@Description	Elimina un ingreso de venta
//	@Tags			IncomeSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"ID of the incomeSale"
//	@Success		200	{object}	schemas.Response	"IncomeSale deleted successfully"
//	@Router			/api/v1/income_sale/delete/{id} [delete]
func (i *IncomeSaleController) IncomeSaleDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar ingreso")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pointID :=c.Locals("point_sale_id").(int64)

	err = i.IncomeSaleService.IncomeSaleDelete(idint, pointID)
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
