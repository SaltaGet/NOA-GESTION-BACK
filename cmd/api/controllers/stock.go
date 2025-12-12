package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// StockGet godoc
//
//	@Summary		StockGet
//	@Description	StockGet obtener un producto-stock por ID
//	@Tags			Stock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Id del producto"
//	@Success		200	{object}	schemas.Response{body=schemas.ProductStockFullResponse}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/stock/get/{id} [get]
func (p *StockController) StockGetByID(ctx *fiber.Ctx) error {
	stockID := ctx.Params("id")
	idint, err := validators.IdValidate(stockID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	pointID := ctx.Locals("point_sale_id").(int64)

	stock, err := p.StockService.StockGetByID(idint, pointID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    stock,
		Message: "Stocko obtenido correctamente",
	})
}

// StockGetByCode godoc
//
//	@Summary		StockGetByCode
//	@Description	obtener un producto-stock por Codigo
//	@Tags			Stock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			code	query		string	true	"codigo del stocko"
//	@Success		200		{object}	schemas.Response{body=schemas.ProductResponse}
//	@Router			/api/v1/stock/get_by_code [get]
func (p *StockController) StockGetByCode(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Se necesita el codigo del producto",
		})
	}

	pointID := ctx.Locals("point_sale_id").(int64)

	stock, err := p.StockService.StockGetByCode(code, pointID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    stock,
		Message: "Stock obtenido correctamente",
	})
}

// StockGetByName godoc
//
//	@Summary		StockGetByName
//	@Description	StockGetByName obtener un producto-stock por nombre
//	@Tags			Stock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			name	query		string	true	"nombre del stocko"
//	@Success		200		{object}	schemas.Response{body=schemas.ProductStockFullResponse}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/stock/get_by_name [get]
func (p *StockController) StockGetByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	if len(name) < 3 {
		return ctx.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El nombre debe de tener al menos 3 caracteres",
		})
	}

	pointID := ctx.Locals("point_sale_id").(int64)

	stocks, err := p.StockService.StockGetByName(name, pointID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    stocks,
		Message: "Stock obtenidos correctamente",
	})
}

// StockGetByCategory godoc
//
//	@Summary		StockGetByCategory
//	@Description	obtener un producto-stock por Id de categoria
//	@Tags			Stock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			category_id	path		string	true	"ID de la categoria"
//	@Success		200			{object}	schemas.Response{body=schemas.ProductStockFullResponse}
//	@Failure		400			{object}	schemas.Response
//	@Failure		401			{object}	schemas.Response
//	@Failure		422			{object}	schemas.Response
//	@Failure		404			{object}	schemas.Response
//	@Failure		500			{object}	schemas.Response
//	@Router			/api/v1/stock/get_by_category/{category_id} [get]
func (p *StockController) StockGetByCategoryID(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("category_id")
	idint, err := validators.IdValidate(categoryID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	pointID := ctx.Locals("point_sale_id").(int64)

	stocks, err := p.StockService.StockGetByCategoryID(idint, pointID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    stocks,
		Message: "Stockos obtenidos correctamente",
	})
}

// StockGetAll godoc
//
//	@Summary		StockGetAll
//	@Description	StockGetAll obtener todos los producto-stock
//	@Tags			Stock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page	query		int	false	"Número de página"				default(1)
//	@Param			limit	query		int	false	"Número de elementos por página"	default(10)
//	@Success		200		{object}	schemas.Response{body=[]schemas.ProductStockFullResponse}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/stock/get_all [get]
func (p *StockController) StockGetAll(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	pointID := ctx.Locals("point_sale_id").(int64)

	stocks, total, err := p.StockService.StockGetAll(page, limit, pointID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": stocks, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Stockos obtenidos correctamente",
	})
}
