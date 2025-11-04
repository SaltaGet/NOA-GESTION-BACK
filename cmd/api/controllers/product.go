package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// ProductGetByID godoc
//
//	@Summary		Get Product By ID
//	@Description	Get a product or part by its ID within a specified workplace.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string									true	"ID of the product"
//	@Success		200	{object}	schemas.Response{body=schemas.Product}	"Product obtained with success"
//	@Failure		400	{object}	schemas.Response						"Bad Request"
//	@Failure		401	{object}	schemas.Response						"Auth is required"
//	@Failure		403	{object}	schemas.Response						"Not Authorized"
//	@Failure		404	{object}	schemas.Response						"Expense not found"
//	@Failure		500	{object}	schemas.Response						"Internal server error"
//	@Router			/product/{id} [get]
func (p *ProductController) ProductGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un producto por ID")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	product, err := p.ProductService.ProductGetByID(id)
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

	logging.INFO("Producto obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    product,
		Message: "Parte obtenida con éxito",
	})
}

// ProductGetAll godoc
//
//	@Summary		Get All Products
//	@Description	Get All Products
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.Product}	"Products obtained with success"
//	@Failure		400	{object}	schemas.Response							"Bad Request"
//	@Failure		401	{object}	schemas.Response							"Auth is required"
//	@Failure		403	{object}	schemas.Response							"Not Authorized"
//	@Failure		500	{object}	schemas.Response							"Internal server error"
//	@Router			/product/get_all [get]
func (p *ProductController) ProductGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los productos")
	products, err := p.ProductService.ProductGetAll()
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

	logging.INFO("Productos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    products,
		Message: "Partes obtenidas con éxito",
	})
}

// ProductGetByName godoc
//
//	@Summary		Get Product By Name
//	@Description	Fetches products from either laundry or workshop based on the provided name and workplace.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name	query		string										true	"Name of the Product"
//	@Success		200		{object}	schemas.Response{body=[]schemas.Product}	"List of products"
//	@Failure		400		{object}	schemas.Response							"Bad Request"
//	@Failure		401		{object}	schemas.Response							"Auth is required"
//	@Failure		403		{object}	schemas.Response							"Not Authorized"
//	@Failure		500		{object}	schemas.Response							"Internal server error"
//	@Router			/product/get_by_name [get]
func (p *ProductController) ProductGetByName(c *fiber.Ctx) error {
	logging.INFO("Obtener productos por nombre")
	name := c.Query("name")
	if name == "" || len(name) < 3 {
		logging.ERROR("Error: El valor no debe de ser vacio o menor a 3 caracteres")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El valor no debe de ser vacio o menor a 3 caracteres",
		})
	}

	products, err := p.ProductService.ProductGetByName(name)
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

	logging.INFO("Productos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    products,
		Message: "Partes obtenidas con éxito",
	})
}

// ProductGetByIdentifier godoc
//
//	@Summary		Get Products by identifier
//	@Description	Get Products by identifier
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			identifier	query		string										true	"Identifier of product"
//	@Success		200			{object}	schemas.Response{body=[]schemas.Product}	"Products obtained with success"
//	@Failure		400			{object}	schemas.Response							"Bad Request"
//	@Failure		401			{object}	schemas.Response							"Auth is required"
//	@Failure		403			{object}	schemas.Response							"Not Authorized"
//	@Failure		500			{object}	schemas.Response							"Internal server error"
//	@Router			/product/get_by_identifier [get]
func (p *ProductController) ProductGetByIdentifier(c *fiber.Ctx) error {
	logging.INFO("Obtener productos por identificador")
	identifire := c.Query("identifier")
	if identifire == "" || len(identifire) < 3 {
		logging.ERROR("Error: El valor no debe de ser vacio o menor a 3 caracteres")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El valor no debe de ser vacio o menor a 3 caracteres",
		})
	}

	products, err := p.ProductService.ProductGetByIdentifier(identifire)
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

	logging.INFO("Productos obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    products,
		Message: "Partes obtenidas con éxito",
	})
}

// ProductUpdateStock godoc
//
//	@Summary		Update Product Stock
//	@Description	Updates the stock of a product based on the given method (add, subtract, update).
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			stock	body		schemas.StockUpdate	true	"Stock update details"
//	@Success		200		{object}	schemas.Response	"Product stock updated successfully"
//	@Failure		400		{object}	schemas.Response	"Bad Request"
//	@Failure		401		{object}	schemas.Response	"Auth is required"
//	@Failure		403		{object}	schemas.Response	"Not Authorized"
//	@Failure		404		{object}	schemas.Response	"Product not found"
//	@Failure		422		{object}	schemas.Response	"Model invalid"
//	@Failure		500		{object}	schemas.Response	"Internal server error"
//	@Router			/product/update_stock [put]
func (p *ProductController) ProductUpdateStock(c *fiber.Ctx) error {
	logging.INFO("Actualizar stock de producto")
	var stockUpdate schemas.StockUpdate
	if err := c.BodyParser(&stockUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := stockUpdate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	err := p.ProductService.ProductUpdateStock(&stockUpdate)
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

	logging.INFO("Producto actualizado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Producto actualizado con éxito",
	})
}

// ProductUpdate godoc
//
//	@Summary		Update Product
//	@Description	Updates the given product and returns the updated product.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			product	body		schemas.ProductUpdate	true	"Product update details"
//	@Success		200		{object}	schemas.Response		"Product updated successfully"
//	@Failure		400		{object}	schemas.Response		"Bad Request"
//	@Failure		401		{object}	schemas.Response		"Auth is required"
//	@Failure		403		{object}	schemas.Response		"Not Authorized"
//	@Failure		404		{object}	schemas.Response		"Product not found"
//	@Failure		422		{object}	schemas.Response		"Model invalid"
//	@Failure		500		{object}	schemas.Response		"Internal server error"
//	@Router			/product/update [put]
func (p *ProductController) ProductUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar producto")
	var productUpdate schemas.ProductUpdate
	if err := c.BodyParser(&productUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := productUpdate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	err := p.ProductService.ProductUpdate(&productUpdate)
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

	logging.INFO("Producto actualizado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Producto actualizado con éxito",
	})
}

// ProductDelete godoc
//
//	@Summary		Delete Product
//	@Description	Deletes the given product with the given id.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string				true	"ID of the product"
//	@Success		200	{object}	schemas.Response	"Product deleted with success"
//	@Failure		400	{object}	schemas.Response	"Bad Request"
//	@Failure		401	{object}	schemas.Response	"Auth is required"
//	@Failure		403	{object}	schemas.Response	"Not Authorized"
//	@Failure		404	{object}	schemas.Response	"Product not found"
//	@Failure		500	{object}	schemas.Response	"Internal server error"
//	@Router			/product/delete/{id} [delete]
func (p *ProductController) ProductDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar producto")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	err := p.ProductService.ProductDelete(id)
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

	logging.INFO("Producto eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Producto eliminado con éxito",
	})
}

// ProductCreate godoc
//
//	@Summary		Create Product
//	@Description	Creates a new product in the specified workplace.
//	@Tags			Product
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			product	body		schemas.ProductCreate	true	"Details of the product to create"
//	@Success		200		{object}	schemas.Response		"Product created successfully"
//	@Failure		400		{object}	schemas.Response		"Bad Request"
//	@Failure		401		{object}	schemas.Response		"Auth is required"
//	@Failure		403		{object}	schemas.Response		"Not Authorized"
//	@Failure		422		{object}	schemas.Response		"Model invalid"
//	@Failure		500		{object}	schemas.Response		"Internal server error"
//	@Router			/product/create [post]
func (p *ProductController) ProductCreate(c *fiber.Ctx) error {
	logging.INFO("Crear producto")
	var productCreate schemas.ProductCreate
	if err := c.BodyParser(&productCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := productCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	productCreated, err := p.ProductService.ProductCreate(&productCreate)
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

	logging.INFO("Producto creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    productCreated,
		Message: "Producto creado con éxito",
	})
}
