package controllers

// import (
// 	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
// 	"github.com/gofiber/fiber/v2"
// )

// // PurchaseOrderGetByID godoc
// //	@Summary		Get Purchase Order By ID
// //	@Description	Retrieves a specific purchase order by its ID. Returns purchase order based on the tenant context.
// //	@Tags			Purchase Order
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			id	path		string										true	"ID of Purchase Order"
// //	@Success		200	{object}	schemas.Response{body=schemas.PurchaseOrder}	"Laundry order obtained successfully"
// //	@Failure		400	{object}	schemas.Response								"Bad Request"
// //	@Failure		401	{object}	schemas.Response								"Auth is required"
// //	@Failure		403	{object}	schemas.Response								"Not Authorized"
// //	@Failure		404	{object}	schemas.Response								"Purchase Order not found"
// //	@Failure		500	{object}	schemas.Response								"Internal server error"
// //	@Router			/purchase_order/{id} [get]
// func (p *PurchaseOrderController) PurchaseOrderGetByID(c *fiber.Ctx) error {
// 	logging.INFO("Obtener orden de compra por ID")
// 	id := c.Params("id")
// 	if id == "" {
// 		logging.ERROR("Error: ID is required")
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "ID is required",
// 		})
// 	}

// 	purchaseOrder, err := p.PurchaseOrderService.PurchaseOrderGetByID(id)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Orden de compra obtenida con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    purchaseOrder,
// 		Message: "Orden de compra obtenida con éxito",
// 	})
// }

// // PurchaseOrderGetAll godoc
// //	@Summary		Get All Purchase Orders
// //	@Description	Get All Purchase Orders
// //	@Tags			Purchase Order
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Success		200	{object}	schemas.Response{body=[]schemas.PurchaseOrder}	"Purchase Orders obtained with success"
// //	@Failure		400	{object}	schemas.Response									"Bad Request"
// //	@Failure		401	{object}	schemas.Response									"Auth is required"
// //	@Failure		403	{object}	schemas.Response									"Not Authorized"
// //	@Failure		500	{object}	schemas.Response									"Internal server error"
// //	@Router			/purchase_order/get_all [get]
// //	@Security		BearerAuth
// func (p *PurchaseOrderController) PurchaseOrderGetAll(c *fiber.Ctx) error {
// 	logging.INFO("Obtener todas las ordenes de compra")
// 	purchasesOrder, err := p.PurchaseOrderService.PurchaseOrderGetAll()
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Ordenes de compra obtenidas con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    purchasesOrder,
// 		Message: "Orden de compra obtenida con éxito",
// 	})
// }

// // PurchaseOrderCreate godoc
// //	@Summary		Create Purchase Order
// //	@Description	Creates a purchase order, either for laundry or workshop.
// //	@Tags			Purchase Order
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			purchaseOrderCreate	body		schemas.PurchaseOrderCreate		true	"Purchase order creation data"
// //	@Success		200					{object}	schemas.Response{body=string}	"Purchase order created successfully"
// //	@Failure		400					{object}	schemas.Response					"Bad Request"
// //	@Failure		401					{object}	schemas.Response					"Auth is required"
// //	@Failure		403					{object}	schemas.Response					"Not Authorized"
// //	@Failure		422					{object}	schemas.Response					"Model invalid"
// //	@Failure		500					{object}	schemas.Response					"Internal server error"
// //	@Router			/purchase_order/create     [post]
// //	@Security		BearerAuth
// func (p *PurchaseOrderController) PurchaseOrderCreate(c *fiber.Ctx) error {
// 	logging.INFO("Crear orden de compra")
// 	var purchaseOrderCreate schemas.PurchaseOrderCreate
// 	if err := c.BodyParser(&purchaseOrderCreate); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Invalid request" + err.Error(),
// 		})
// 	}
// 	if err := purchaseOrderCreate.Validate(); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	id, err := p.PurchaseOrderService.PurchaseOrderCreate(&purchaseOrderCreate)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Orden de compra creada con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    id,
// 		Message: "Orden de compra creada con éxito",
// 	})
// }

// // // PurchaseOrderUpdate godoc
// // //	@Summary		Update Purchase Order
// // //	@Description	Updates an existing purchase order with new details.
// // //              Validates the request body and workplace context.
// // //              Returns a success message if the update is successful.
// // //	@Tags			Purchase Order
// // //	@Accept			json
// // //	@Produce		json
// // //	@Security		BearerAuth
// // //	@Param			purchaseOrderUpdate	body		schemas.PurchaseOrderUpdate	true	"Purchase order update data"
// // //	@Success		200					{object}	schemas.Response				"Purchase order updated successfully"
// // //	@Failure		400					{object}	schemas.Response				"Bad Request"
// // //	@Failure		401					{object}	schemas.Response				"Auth is required"
// // //	@Failure		403					{object}	schemas.Response				"Not Authorized"
// // //	@Failure		422					{object}	schemas.Response				"Model invalid"
// // //	@Failure		500					{object}	schemas.Response				"Internal server error"
// // //	@Router			/purchase_order/update [put]
// // func (p *PurchaseOrderController) PurchaseOrderUpdate(c *fiber.Ctx) error {
// // 	logging.INFO("Actualizar orden de compra")
// // 	var purchaseOrderUpdate schemas.PurchaseOrderUpdate
// // 	if err := c.BodyParser(&purchaseOrderUpdate); err != nil {
// // 		logging.ERROR("Error: %s", err.Error())
// // 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// // 			Status:  false,
// // 			Body:    nil,
// // 			Message: "Invalid request" + err.Error(),
// // 		})
// // 	}
// // 	if err := purchaseOrderUpdate.Validate(); err != nil {
// // 		logging.ERROR("Error: %s", err.Error())
// // 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// // 			Status:  false,
// // 			Body:    nil,
// // 			Message: err.Error(),
// // 		})
// // 	}

// // 	err := p.PurchaseOrderService.PurchaseOrderUpdate(&purchaseOrderUpdate)
// // 	if err != nil {
// // 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// // 			logging.ERROR("Error: %s", errResp.Err.Error())
// // 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// // 				Status:  false,
// // 				Body:    nil,
// // 				Message: errResp.Message,
// // 			})
// // 		}
// // 		logging.ERROR("Error: %s", err.Error())
// // 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// // 			Status:  false,
// // 			Body:    nil,
// // 			Message: "Error interno",
// // 		})
// // 	}

// // 	logging.INFO("Orden de compra editada con éxito")
// // 	return c.Status(200).JSON(schemas.Response{
// // 		Status:  true,
// // 		Body:    nil,
// // 		Message: "Orden de compra editada con éxito",
// // 	})
// // }

// // PurchaseOrderDelete godoc
// //	@Summary		Delete Purchase Order
// //	@Description	Deletes a specific purchase order by its ID.
// //	@Tags			Purchase Order
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			id	path		string			true	"ID of Purchase Order"
// //	@Success		200	{object}	schemas.Response	"Purchase order deleted successfully"
// //	@Failure		400	{object}	schemas.Response	"Bad Request"
// //	@Failure		401	{object}	schemas.Response	"Auth is required"
// //	@Failure		403	{object}	schemas.Response	"Not Authorized"
// //	@Failure		404	{object}	schemas.Response	"Purchase order not found"
// //	@Failure		500	{object}	schemas.Response	"Internal server error"
// //	@Router			/purchase_order/delete/{id} [delete]
// func (p *PurchaseOrderController) PurchaseOrderDelete(c *fiber.Ctx) error {
// 	logging.INFO("Eliminar orden de compra")
// 	id := c.Params("id")
// 	if id == "" {
// 		logging.ERROR("Error: ID is required")
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "ID is required",
// 		})
// 	}

// 	err := p.PurchaseOrderService.PurchaseOrderDelete(id)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Orden de compra eliminada con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    nil,
// 		Message: "Orden de compra eliminada con éxito",
// 	})
// }

