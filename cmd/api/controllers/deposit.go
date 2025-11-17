package controllers

import (
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// DepositProductGetByID godoc
//
//	@Summary		DepositProductGetByID
//	@Description	DepositProductGetByID obtener un producto por ID del producto
//	@Tags			Deposit
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Id del producto"
//	@Success		200	{object}	schemas.Response{body=schemas.DepositResponse}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/deposit/get/{id} [get]
func (d *DepositController) DepositGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener producto del deposito")
	id := c.Params("id")
	idUint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	product, err := d.DepositService.DepositGetByID(idUint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Producto del deposito obtenido con exito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Producto del deposito obtenido con exito",
	})
}

// DepositProductGetByCode godoc
//
//	@Summary		DepositProductGetByCode
//	@Description	DepositProductGetByCode obtener un producto por codigo del deposito
//	@Tags			Deposit
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			code	query		string	true	"codigo del producto"
//	@Success		200		{object}	schemas.Response{body=schemas.DepositResponse}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/deposit/get_by_code [get]
func (d *DepositController) DepositGetByCode(c *fiber.Ctx) error {
	logging.INFO("Obtener producto del deposito por codigo")
	code := c.Query("code")
	if code == "" {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Se necesita el codigo del producto", fmt.Errorf("se necesita el codigo del producto")))
	}

	product, err := d.DepositService.DepositGetByCode(code)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Producto obtenido correctamente")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Producto obtenido correctamente",
	})
}

// DepositProductGetByName godoc
//
//	@Summary		DepositProductGetByName
//	@Description	DepositProductGetByName obtener productos por por similitud de nombre
//	@Tags			Deposit
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			name	query		string	true	"nombre por aproximacion del producto"
//	@Success		200		{object}	schemas.Response{body=schemas.DepositResponse}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/deposit/get_by_name [get]
func (d *DepositController) DepositGetByName(c *fiber.Ctx) error {
	logging.INFO("Obtener productos del deposito por nombre")
	name := c.Query("name")
	if len(name) < 3 {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "El nombre debe de tener al menos 3 caracteres", fmt.Errorf("el nombre debe de tener al menos 3 caracteres")))
	}

	products, err := d.DepositService.DepositGetByName(name)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Productos obtenidos correctamente")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    products,
		Message: "Productos obtenidos correctamente",
	})
}

// DepositProductGetAll godoc
//
//	@Summary		DepositProductGetAll
//	@Description	DepositProductGetAll obtener productos por paginacion
//	@Tags			Deposit
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page	query		int	false	"pagina" default(1)
//	@Param			limit	query		int	false	"limite" default(10)
//	@Success		200		{object}	schemas.Response{body=[]schemas.DepositResponse}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/deposit/get_all [get]
func (d *DepositController) DepositGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener productos del deposito por paginacion")
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	products, total, err := d.DepositService.DepositGetAll(page, limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Productos obtenidos correctamente")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": products, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Productos obtenidos correctamente",
	})
}

// DepositProductUpdateStock godoc
//
//	@Summary		DepositProductUpdateStock
//	@Description	DepositProductUpdateStock alctualizar stock de un producto del deposito
//	@Tags			Deposit
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			stock_update	body		schemas.DepositUpdateStock	true	"nombre por aproximacion del producto"
//	@Success		200				{object}	schemas.Response{body=[]schemas.DepositResponse}
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		404				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/deposit/update_stock [put]
func (d *DepositController) DepositUpdateStock(c *fiber.Ctx) error {
	logging.INFO("Actualizar stock de un producto del deposito")
	var stockUpdate schemas.DepositUpdateStock
	if err := c.BodyParser(&stockUpdate); err != nil {
		return schemas.HandleError(c, err)
	}

	if err := stockUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := d.DepositService.DepositUpdateStock(stockUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Producto actualizado correctamente")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Producto actualizado correctamente",
	})
}
