package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// PointSale GetAll godoc
//	@Summary		PointSale GetAll
//	@Description	PointSale GetAll required auth token
//	@Tags			PointSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PointSaleResponse}	"puntos de ventas obtenidos con éxito"
//	@Failure		400	{object}	schemas.Response									"Bad Request"
//	@Failure		401	{object}	schemas.Response									"Auth is required"
//	@Failure		403	{object}	schemas.Response									"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/point_sale/get_all [get]
func (p *PointSaleController) PointSaleGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los puntos de ventas")
	pointSale, err := p.PointSaleService.PointSaleGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Puntos de ventas obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    pointSale,
		Message: "Puntos de ventas obtenidos con éxito",
	})
}

// PointSale GetAllByMember godoc
//	@Summary		GetAllByMember GetAll
//	@Description	Obtener puntos de ventas asociados a miembro
//	@Tags			PointSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PointSaleResponse}	"puntos de ventas obtenidos con éxito"
//	@Failure		400	{object}	schemas.Response									"Bad Request"
//	@Failure		401	{object}	schemas.Response									"Auth is required"
//	@Failure		403	{object}	schemas.Response									"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/point_sale/get_all_by_member [get]
func (p *PointSaleController) PointSaleGetAllByMember(c *fiber.Ctx) error {
	member := c.Locals("user").(*schemas.AuthenticatedUser)
	logging.INFO("Obtener puntos de ventas asociados a miembro: %d", member.ID)

	permissions, err := p.PointSaleService.PointSaleGetAllByMember(member.ID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Puntos de ventas obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    permissions,
		Message: "Puntos de ventas obtenidos con éxito",
	})
}

// PointSaleCreate godoc
//	@Summary		PointSaleCreate 
//	@Description	Crear punto de venta
//	@Tags			PointSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			request	body		schemas.PointSaleCreate	true	"Crear punto de venta"
//	@Success		200		{object}	schemas.Response		"puntos de ventas creado con éxito"
//	@Failure		400		{object}	schemas.Response		"Bad Request"
//	@Failure		401		{object}	schemas.Response		"Auth is required"
//	@Failure		403		{object}	schemas.Response		"Not Authorized"
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/point_sale/create [post]
func (p *PointSaleController) PointSaleCreate(c *fiber.Ctx) error {
	logging.INFO("Crear punto de venta")
	var pointSale schemas.PointSaleCreate
	if err := c.BodyParser(&pointSale); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el modelo", err))
	}
	if err := pointSale.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	permissions, err := p.PointSaleService.PointSaleCreate(&pointSale)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Puntos de ventas obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    permissions,
		Message: "Puntos de ventas obtenidos con éxito",
	})
}
