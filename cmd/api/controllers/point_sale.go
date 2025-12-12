package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// PointSaleGetByID godoc
//	@Summary		PointSaleGetByID
//	@Description	Obtener punto de venta por ID
//	@Tags			PointSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string												true	"PointSale ID"
//	@Success		200	{object}	schemas.Response{body=[]schemas.PointSaleResponse}	"puntos de ventas obtenidos con éxito"
//	@Router			/api/v1/point_sale/get/{id} [get]
func (p *PointSaleController) PointSaleGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	pointSale, err := p.PointSaleService.PointSaleGetByID(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    pointSale,
		Message: "Punto de venta por ID obtenido con éxito",
	})
}

// PointSale GetAll godoc
//	@Summary		PointSale GetAll
//	@Description	PointSale GetAll required auth token
//	@Tags			PointSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PointSaleResponse}	"puntos de ventas obtenidos con éxito"
//	@Router			/api/v1/point_sale/get_all [get]
func (p *PointSaleController) PointSaleGetAll(c *fiber.Ctx) error {
	pointSale, err := p.PointSaleService.PointSaleGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

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
//	@Router			/api/v1/point_sale/get_all_by_member [get]
func (p *PointSaleController) PointSaleGetAllByMember(c *fiber.Ctx) error {
	member := c.Locals("user").(*schemas.AuthenticatedUser)

	permissions, err := p.PointSaleService.PointSaleGetAllByMember(member.ID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

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
//	@Router			/api/v1/point_sale/create [post]
func (p *PointSaleController) PointSaleCreate(c *fiber.Ctx) error {
	var pointSale schemas.PointSaleCreate
	if err := c.BodyParser(&pointSale); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el modelo", err))
	}
	if err := pointSale.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	plan := c.Locals("current_plan").(*schemas.PlanResponseDTO)
member := c.Locals("user").(*schemas.AuthenticatedUser)
	id, err := p.PointSaleService.PointSaleCreate(member.ID, &pointSale, plan)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Puntos de ventas obtenidos con éxito",
	})
}

// PointSaleUpdate godoc
//	@Summary		PointSaleUpdate 
//	@Description	Editar punto de venta
//	@Tags			PointSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			request	body		schemas.PointSaleUpdate	true	"Editar punto de venta"
//	@Success		200		{object}	schemas.Response		"puntos de ventas editado con éxito"
//	@Router			/api/v1/point_sale/update [put]
func (p *PointSaleController) PointSaleUpdate(c *fiber.Ctx) error {
	var pointSale schemas.PointSaleUpdate
	if err := c.BodyParser(&pointSale); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el modelo", err))
	}
	if err := pointSale.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)
	err := p.PointSaleService.PointSaleUpdate(member.ID, &pointSale)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Puntos de venta editado con éxito",
	})
}

// PointSaleUpdateMain godoc
//	@Summary		PointSaleUpdateMain 
//	@Description	Editar punto de venta principal
//	@Tags			PointSale
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			request	body		schemas.PointSaleUpdateMain	true	"Editar punto de venta princiapl"
//	@Success		200		{object}	schemas.Response			"puntos de venta principal editado con éxito"
//	@Router			/api/v1/point_sale/update_main [put]
func (p *PointSaleController) PointSaleUpdateMain(c *fiber.Ctx) error {
	var pointSale schemas.PointSaleUpdateMain
	if err := c.BodyParser(&pointSale); err != nil {
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el modelo", err))
	}
	if err := pointSale.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)
	err := p.PointSaleService.PointSaleUpdateMain(member.ID, &pointSale)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Puntos de venta principal editado con éxito",
	})
}
