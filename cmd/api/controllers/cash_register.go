package controllers

import (
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// CashRegisterExistOpen godoc
//
//	@Summary		CashRegisterExistOpen
//	@Description	Verifica si existe apertura de caja
//	@Tags			CashRegister
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/cash_register/exist_open [get]
func (r *CashRegisterController) CashRegisterExistOpen(ctx *fiber.Ctx) error {
	pointaSale := ctx.Locals("point_sale_id").(int64)
	
	existOpen, err := r.CashRegisterService.CashRegisterExistOpen(pointaSale)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	var message string

	if existOpen {
		message = "Existe apertura de caja"
	} else {
		message = "No existe apertura de caja"
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    existOpen,
		Message: message,
	})
}

// CashRegisterGetByID godoc
//
//	@Summary		CashRegisterGetByID
//	@Description	obtener caja por id
//	@Tags			CashRegister
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		int	true	"id de la caja"
//	@Success		200	{object}	schemas.Response{body=schemas.CashRegisterFullResponse}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/cash_register/get/{id} [get]
func (r *CashRegisterController) CashRegisterGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "se necesita el id de la caja", fmt.Errorf("se necesita el id de la caja")))
	}

	idint, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(422, "el id ser un número", err))
	}

	pointaSale := ctx.Locals("point_sale_id").(int64)
	
	register, err := r.CashRegisterService.CashRegisterGetByID(pointaSale, idint)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    register,
		Message: "Caja obtenida con éxito",
	})
}

// CashRegisterOpen godoc
//
//	@Summary		CashRegisterOpen
//	@Description	Apertura de caja
//	@Tags			CashRegister
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			amount_open	body		schemas.CashRegisterOpen	true	"Monto de apertura de caja"
//	@Success		200			{object}	schemas.Response
//	@Failure		400			{object}	schemas.Response
//	@Failure		401			{object}	schemas.Response
//	@Failure		422			{object}	schemas.Response
//	@Failure		404			{object}	schemas.Response
//	@Failure		500			{object}	schemas.Response
//	@Router			/api/v1/cash_register/open [post]
func (r *CashRegisterController) CashRegisterOpen(ctx *fiber.Ctx) error {
	var amountOpen schemas.CashRegisterOpen
	if err := ctx.BodyParser(&amountOpen); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}
	if err := amountOpen.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	pointaSale := ctx.Locals("point_sale_id").(int64)
	user := ctx.Locals("user").(*schemas.AuthenticatedUser)

	err := r.CashRegisterService.CashRegisterOpen(pointaSale, user.ID, amountOpen)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}	

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Apertura de caja realizada con exito",
	})
}

// CashRegisterClose godoc
//
//	@Summary		CashRegisterClose
//	@Description	Cierre de caja
//	@Tags			CashRegister
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			amount_close	body		schemas.CashRegisterClose	true	"Monto de cierre de caja"
//	@Success		200				{object}	schemas.Response
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		404				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/cash_register/close [post]
func (r *CashRegisterController) CashRegisterClose(ctx *fiber.Ctx) error {
	var amountClose schemas.CashRegisterClose
	if err := ctx.BodyParser(&amountClose); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}
	if err := amountClose.Validate(); err != nil {
		return schemas.HandleError(ctx, err)
	}

	pointaSale := ctx.Locals("point_sale_id").(int64)
	user := ctx.Locals("user").(*schemas.AuthenticatedUser)

	err := r.CashRegisterService.CashRegisterClose(pointaSale, user.ID, amountClose)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Cierre de caja realizado con exito",
	})
}

// CashRegisterInform godoc
//
//	@Summary		CashRegisterInform
//	@Description	Informes de caja
//	@Tags			CashRegister
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			register_request	query		schemas.DateRangeRequest	true	"Fechas de solicitud de informe"
//	@Success		200					{object}	schemas.Response
//	@Failure		400					{object}	schemas.Response
//	@Failure		401					{object}	schemas.Response
//	@Failure		422					{object}	schemas.Response
//	@Failure		404					{object}	schemas.Response
//	@Failure		500					{object}	schemas.Response
//	@Router			/api/v1/cash_register/inform [get]
func (r *CashRegisterController) CashRegiterInform(ctx *fiber.Ctx) error {
	formDate := &schemas.DateRangeRequest{}
	formDate.FromDate = ctx.Query("from_date")
	formDate.ToDate = ctx.Query("to_date")
	
	fromDate, toDate, err := formDate.GetParsedDates()
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	pointaSale := ctx.Locals("point_sale_id").(int64)
	user := ctx.Locals("user").(*schemas.AuthenticatedUser)

	informs, err := r.CashRegisterService.CashRegisterInform(pointaSale, user.ID, fromDate, toDate)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    informs,
		Message: "informes obtenidos con exito",
	})
}
