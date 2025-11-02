package controllers

// import (
// 	"github.com/DanielChachagua/GestionCar/pkg/models"
// 	"github.com/DanielChachagua/GestionCar/pkg/services"
// 	"github.com/gofiber/fiber/v2"
// )

// // PurchaseProductGetByID godoc
// //	@Summary		Get Purchase Product By ID
// //	@Description	Retrieves a specific purchase product by its ID for a given workplace.
// //	@Tags			Purchase Product
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			id					path		string												true	"ID of the purchase product"
// //	@Success		200					{object}	models.Response{body=models.PurchaseProductLaundry}	"Product obtained successfully"
// //	@Failure		400					{object}	models.Response										"Bad Request"
// //	@Failure		401					{object}	models.Response										"Auth is required"
// //	@Failure		403					{object}	models.Response										"Not Authorized"
// //	@Failure		404					{object}	models.Response										"Purchase Product not found"
// //	@Failure		500					{object}	models.Response										"Internal server error"
// //	@Router			/purchase_product/{id} [get]
// func (p *PurchaseProductController) PurchaseProductGetByID(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	if id == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "ID is required",
// 		})
// 	}

// 	purchaseProduct, err := p.PurchaseProductService.PurchaseProductGetByID(id)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	return c.Status(200).JSON(models.Response{
// 		Status:  true,
// 		Body:    purchaseProduct,
// 		Message: "Producto de compra obtenida con éxito",
// 	})
// }

// // PurchaseProductGetAllByPurhcaseID godoc
// //	@Summary		Get All Products From Purchase Order
// //	@Description	Get All Products From Purchase Order
// //	@Tags			Purchase Product
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			purchase_id			path		string													true	"ID of Purchase Order"
// //	@Success		200					{object}	models.Response{body=[]models.PurchaseProductLaundry}	"Products obtained with success"
// //	@Failure		400					{object}	models.Response											"Bad Request"
// //	@Failure		401					{object}	models.Response											"Auth is required"
// //	@Failure		403					{object}	models.Response											"Not Authorized"
// //	@Failure		404					{object}	models.Response											"Purchase Order not found"
// //	@Failure		500					{object}	models.Response											"Internal server error"
// //	@Router			/purchase_product/get_purchase/{purchase_id} [get]
// func (p *PurchaseProductController) PurchaseProductGetAllByPurhcaseID(c *fiber.Ctx) error {
// 	purchaseId := c.Params("purchase_id")
// 	if purchaseId == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "ID is required",
// 		})
// 	}

// 	purchaseProducts, err := p.PurchaseProductService.PurchaseProductGetAllByPurhcaseID(purchaseId)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	return c.Status(200).JSON(models.Response{
// 		Status:  true,
// 		Body:    purchaseProducts,
// 		Message: "Productos de orden de compra obtenida con éxito",
// 	})
// }

// // PurchaseProductCreate godoc
// //	@Summary		Create Purchase Product
// //	@Description	Creates a purchase product, either for laundry or workshop.
// //              Returns the ID of the created purchase product.
// //	@Tags			Purchase Product
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			purchaseProductCreate	body		models.PurchaseProductCreate	true	"Purchase product creation data"
// //	@Success		200						{object}	models.Response{body=string}	"Purchase product created successfully"
// //	@Failure		400						{object}	models.Response					"Bad Request"
// //	@Failure		401						{object}	models.Response					"Auth is required"
// //	@Failure		403						{object}	models.Response					"Not Authorized"
// //	@Failure		422						{object}	models.Response					"Model is invalid"
// //	@Failure		500						{object}	models.Response					"Internal server error"
// //	@Router			/purchase_product/create   [post]
// func (p *PurchaseProductController) PurchaseProductCreate(c *fiber.Ctx) error {
// 	var purchaseProductCreate models.PurchaseProductCreate
// 	if err := c.BodyParser(&purchaseProductCreate); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Invalid request",
// 		})
// 	}
// 	if err := purchaseProductCreate.Validate(); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	id, err := p.PurchaseProductService.PurchaseProductCreate(&purchaseProductCreate)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	return c.Status(200).JSON(models.Response{
// 		Status:  true,
// 		Body:    id,
// 		Message: "Producto deOrden de compra creado con éxito",
// 	})
// }

// // PurchaseProductUpdate godoc
// //	@Summary		Update Purchase Product
// //	@Description	Updates the given purchase product and returns a success message.
// //	@Tags			Purchase Product
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			id					path		string							true	"ID of the purchase product"
// //	@Param			product				body		models.PurchaseProductUpdate	true	"Purchase product update details"
// //	@Success		200					{object}	models.Response					"Purchase product updated successfully"
// //	@Failure		400					{object}	models.Response					"Bad Request"
// //	@Failure		401					{object}	models.Response					"Auth is required"
// //	@Failure		403					{object}	models.Response					"Not Authorized"
// //	@Failure		404					{object}	models.Response					"Purchase Product not found"
// //	@Failure		422					{object}	models.Response					"Model is invalid"
// //	@Failure		500					{object}	models.Response					"Internal server error"
// //	@Router			/purchase_product/update/{id} [put]
// func (p *PurchaseProductController) PurchaseProductUpdate(c *fiber.Ctx) error {
// 	var purchaseProductUpdate models.PurchaseProductUpdate
// 	if err := c.BodyParser(&purchaseProductUpdate); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Invalid request",
// 		})
// 	}
// 	if err := purchaseProductUpdate.Validate(); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	err := p.PurchaseProductService.PurchaseProductUpdate(&purchaseProductUpdate)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	return c.Status(200).JSON(models.Response{
// 		Status:  true,
// 		Body:    nil,
// 		Message: "Producto de Orden de compra editado con éxito",
// 	})
// }

// // PurchaseProductDelete godoc
// //	@Summary		Delete Purchase Product
// //	@Description	Deletes a specific purchase product.
// //	@Tags			Purchase Product
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			id					path		string			true	"ID of Purchase Product"
// //	@Success		200					{object}	models.Response	"Purchase product deleted successfully"
// //	@Failure		400					{object}	models.Response	"Bad Request"
// //	@Failure		401					{object}	models.Response	"Auth is required"
// //	@Failure		403					{object}	models.Response	"Not Authorized"
// //	@Failure		404					{object}	models.Response	"Purchase Product not found"
// //	@Failure		500					{object}	models.Response	"Internal server error"
// //	@Router			/purchase_product/{id} [delete]
// func (p *PurchaseProductController) PurchaseProductDelete(c *fiber.Ctx) error {
// 	id := c.Params("id")
// 	if id == "" {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "ID is required",
// 		})
// 	}

// 	err := p.PurchaseProductService.PurchaseProductDelete(id)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	return c.Status(200).JSON(models.Response{
// 		Status:  true,
// 		Body:    nil,
// 		Message: "Producto deOrden de compra eliminado con éxito",
// 	})
// }
