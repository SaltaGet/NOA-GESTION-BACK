package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// ProductGet godoc
//
//	@Summary		ProductGet
//	@Description	ProductGet obtener un producto por ID
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Id del producto"
//	@Success		200	{object}	schemas.Response{body=schemas.ProductFullResponse}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/product/get/{id} [get]
func (p *ProductController) ProductGetByID(ctx *fiber.Ctx) error {
	productID := ctx.Params("id")
	idint, err := validators.IdValidate(productID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	product, err := p.ProductService.ProductGetByID(idint)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Producto obtenido correctamente",
	})
}

// ProductGetByCode godoc
//
//	@Summary		ProductGetByCode
//	@Description	ProductGetByCode obtener un producto por Codigo
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			code	query		string	true	"codigo del producto"
//	@Success		200		{object}	schemas.Response{body=schemas.ProductResponse}
//	@Router			/api/v1/product/get_by_code [get]
func (p *ProductController) ProductGetByCode(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	if code == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Se necesita el codigo del producto",
		})
	}

	product, err := p.ProductService.ProductGetByCode(code)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Producto obtenido correctamente",
	})
}

// ProductGetByName godoc
//
//	@Summary		ProductGetByName
//	@Description	ProductGetByName obtener un producto por nombre
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			name	query		string	true	"nombre del producto"
//	@Success		200		{object}	schemas.Response{body=schemas.ProductFullResponse}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/product/get_by_name [get]
func (p *ProductController) ProductGetByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")
	if len(name) < 3 {
		return ctx.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El nombre debe de tener al menos 3 caracteres",
		})
	}

	products, err := p.ProductService.ProductGetByName(name)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    products,
		Message: "Productos obtenidos correctamente",
	})
}

// ProductGetByCategory godoc
//
//	@Summary		ProductGetByCategory
//	@Description	ProductGetByCategory obtener un producto por Id de categoria
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			category_id	path		string	true	"ID de la categoria"
//	@Success		200			{object}	schemas.Response{body=schemas.ProductFullResponse}
//	@Failure		400			{object}	schemas.Response
//	@Failure		401			{object}	schemas.Response
//	@Failure		422			{object}	schemas.Response
//	@Failure		404			{object}	schemas.Response
//	@Failure		500			{object}	schemas.Response
//	@Router			/api/v1/product/get_by_category/{category_id} [get]
func (p *ProductController) ProductGetByCategoryID(ctx *fiber.Ctx) error {
	categoryID := ctx.Params("category_id")
	idint, err := validators.IdValidate(categoryID)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	products, err := p.ProductService.ProductGetByCategoryID(idint)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    products,
		Message: "Productos obtenidos correctamente",
	})
}

// ProductGetAll godoc
//
//	@Summary		ProductGetAll
//	@Description	ProductGetAll obtener todos los productos
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page	query		int	false	"Número de página"				default(1)
//	@Param			limit	query		int	false	"Número de elementos por página"	default(10)
//	@Success		200		{object}	schemas.Response{body=[]schemas.ProductFullResponse}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/product/get_all [get]
func (p *ProductController) ProductGetAll(ctx *fiber.Ctx) error {
	page, err := strconv.Atoi(ctx.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	products, total, err := p.ProductService.ProductGetAll(page, limit)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": products, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Productos obtenidos correctamente",
	})
}

// ProductCreate godoc
//
//	@Summary		ProductCreate
//	@Description	ProductCreate crear nuevo producto
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			productCreate	body		schemas.ProductCreate	true	"Información del producto a crear"
//	@Success		200				{object}	schemas.Response
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		404				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/product/create [post]
func (p *ProductController) ProductCreate(ctx *fiber.Ctx) error {
	var productCreate schemas.ProductCreate
	if err := ctx.BodyParser(&productCreate); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	if err := productCreate.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	productID, err := p.ProductService.ProductCreate(&productCreate)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    productID,
		Message: "Producto creado correctamente",
	})
}

// ProductUpdate godoc
//
//	@Summary		ProductUpdate
//	@Description	ProductUpdate edita un producto ya creado
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			productUpdate	body		schemas.ProductUpdate	true	"Información del producto a editar"
//	@Success		200				{object}	schemas.Response
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		404				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/product/update [put]
func (p *ProductController) ProductUpdate(ctx *fiber.Ctx) error {
	var productUpdate schemas.ProductUpdate
	if err := ctx.BodyParser(&productUpdate); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	if err := productUpdate.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	err := p.ProductService.ProductUpdate(&productUpdate)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Producto actualizado correctamente",
	})
}

// ProductPriceUpdate godoc
//
//	@Summary		ProductPriceUpdate
//	@Description	ProductPriceUpdate edita el rpecio de una lista de productos
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			listProductUpdate	body		schemas.ListPriceUpdate	true	"Información de los productos y los precios a editar"
//	@Success		200					{object}	schemas.Response
//	@Failure		400					{object}	schemas.Response
//	@Failure		401					{object}	schemas.Response
//	@Failure		422					{object}	schemas.Response
//	@Failure		404					{object}	schemas.Response
//	@Failure		500					{object}	schemas.Response
//	@Router			/api/v1/product/list_price [put]
func (p *ProductController) ProductPriceUpdate(ctx *fiber.Ctx) error {
	var productUpdate schemas.ListPriceUpdate
	if err := ctx.BodyParser(&productUpdate); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	if err := productUpdate.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	err := p.ProductService.ProductPriceUpdate(&productUpdate)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Productos actualizados correctamente",
	})
}

// ProductDelete godoc
//
//	@Summary		ProductDelete
//	@Description	ProductDelete elimina un producto
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"ID del producto"
//	@Success		200	{object}	schemas.Response
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/product/delete/{id} [delete]
func (p *ProductController) ProductDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	err = p.ProductService.ProductDelete(idint)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Producto eliminado correctamente",
	})
}
