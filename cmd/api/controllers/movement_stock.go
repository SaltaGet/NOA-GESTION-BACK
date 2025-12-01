package controllers

import (
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// MovementStockGet godoc
//
//	@Summary		MovementStockGet
//	@Description	MovementStockGet Obtener un movimiento de stock por ID
//	@Tags			MovementStock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"ID del movimiento de stock"
//	@Success		200	{object}	schemas.Response{body=schemas.MovementStockResponse}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/movement_stock/get/{id} [get]
func (m *MovementStockController) MovementStockGet(c *fiber.Ctx) error {
	logging.INFO("Inicio MovementStockGet")
	id := c.Params("id")
	if id == "" {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Se necesita el id del movimiento de stock", fmt.Errorf("se necesita el id del movimiento de stock")))
	}

	idUint, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(422, "el id debe ser un número", err))
	}

	movementStock, err := m.MovementStockService.MovementStockGetByID(idUint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Fin MovementStockGet")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    movementStock,
		Message: "Movimiento de stock obtenido correctamente",
	})
}

// MovementStockGetByDate godoc
//
//	@Summary		MovementStockGetByDate
//	@Description	Obtener movimeintos de sotck por paginacion
//	@Tags			MovementStock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page		query		int							false	"Número de página"				default(1)
//	@Param			limit		query		int							false	"Número de elementos por página"	default(10)
//	@Param			fromDate	query		schemas.DateRangeRequest	true	"Fecha de inicio"
//	@Success		200			{object}	schemas.Response{body=[]schemas.MovementStockResponse}
//	@Failure		400			{object}	schemas.Response
//	@Failure		401			{object}	schemas.Response
//	@Failure		422			{object}	schemas.Response
//	@Failure		404			{object}	schemas.Response
//	@Failure		500			{object}	schemas.Response
//	@Router			/api/v1/movement_stock/get_by_date [get]
func (m *MovementStockController) MovementStockGetByDate(c *fiber.Ctx) error {
	logging.INFO("Inicio MovementStockGetByDate")
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	formDate := &schemas.DateRangeRequest{}
	formDate.FromDate = c.Query("from_date")
	formDate.ToDate = c.Query("to_date")
	
	fromDate, toDate, err := formDate.GetParsedDates()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	movementsStock, total, err := m.MovementStockService.MovementStockGetByDate(page, limit, fromDate, toDate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Fin MovementStockGetByDate")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]interface{}{"data": movementsStock, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Movimiento de stock obtenido correctamente",
	})
}

// MovementStockList godoc
//
//	@Summary		MovementStockList
//	@Description	movimiento de stock entre doposito y puntos de ventas por lista
//	@Tags			MovementStock
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			movement_stock	body		[]schemas.MovementStockList	true	"movimiento de stock"
//	@Success		200				{object}	schemas.Response
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		404				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/movement_stock/move_list [post]
func (m *MovementStockController) MoveStockList(c *fiber.Ctx) error {
	logging.INFO("Inicio MoveStockList")
	user := c.Locals("user").(*schemas.AuthenticatedUser)

	var movementStock []*schemas.MovementStockList
	if err := c.BodyParser(&movementStock); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	// Validar que no esté vacío
	if len(movementStock) == 0 {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "No hay movimientos para procesar", fmt.Errorf("lista vacía")))
	}

	// toDepositHasPointSale := false
	// toPointSale := false
	// Validar cada movimiento
	for idx, ms := range movementStock {
		// for _, msi := range ms.MovementStockItem {
		// 	if msi.ToType == "point_sale" {
		// 		toPointSale = true
		// 	}
		// 	if msi.FromType == "deposit" && msi.ToType == "point_sale" {
		// 		toDepositHasPointSale = true
		// 		break
		// 	}
		// }
		if err := ms.Validate(); err != nil {
			return schemas.HandleError(c, schemas.ErrorResponse(400,
				fmt.Sprintf("Error de validación en producto (posición %d, ID: %d)", idx+1, ms.ProductID),
				err))
		}
	}

	err := m.MovementStockService.MoveStockList(user.ID, movementStock)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	
	// if toDepositHasPointSale {
	// 	// Opción 2: Llamada directa en goroutine (MÁS SIMPLE)
	// 	logging.INFO("Intento notificacion de stock por movimiento desde deposito a punto de venta")
	// 	go func() {
	// 		if err := m.NotificationController.SendStockNotification(user.TenantID); err != nil {
	// 			log.Printf("Error enviando notificación de stock: %v", err)
	// 		}
	// 	}()
	// }

	// if toPointSale {
	// 	// Opción 2: Llamada directa en goroutine (MÁS SIMPLE)
	// 	logging.INFO("Intento notificacion de stock por movimiento desde deposito a punto de venta")
	// 	go func() {
	// 		if err := m.NotificationController.SendNotificationRefreshProducts(user.TenantID); err != nil {
	// 			log.Printf("Error enviando notificación de stock: %v", err)
	// 		}
	// 	}()
	// }

	logging.INFO("Fin MoveStockList")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Movimiento de stock realizado correctamente",
	})
}
