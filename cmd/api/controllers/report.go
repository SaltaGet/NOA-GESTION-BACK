package controllers

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// ReportGetReportExcel godoc
//
//	@Summary		ReportGetReportExcel
//	@Description	Obtiene un reporte en excel por fechas
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			start	query		string	true	"2006-01-02"
//	@Param			end		query		string	true	"2006-01-02"
//	@Success		200		{object}	schemas.Response
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		422		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/api/v1/report/get_excel [get]
func (c *ReportController) ReportExcelGet(ctx *fiber.Ctx) error {
	logging.INFO("Inicio ReportExcelGet")
	start := ctx.Query("start")
	if start == "" {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "El campo start no puede estar vacio", fmt.Errorf("el campo start no puede estar vacio")))
	}

	end := ctx.Query("end")
	if end == "" {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "El campo end no puede estar vacio", fmt.Errorf("el campo end no puede estar vacio")))
	}

	var dateRangeRequest = schemas.DateRangeRequest{
		FromDate: start,
		ToDate:   end,
	}

	fromDate, toDate, err := dateRangeRequest.GetParsedDates()
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	inform, err := c.ReportService.ReportExcelGet(fromDate, toDate)
	if err != nil {
		return err
	}

	ctx.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Set("Content-Disposition", "attachment; filename=inform.xlsx")
	ctx.Set("File-Name", "inform.xlsx")

	if err := inform.Write(ctx.Response().BodyWriter()); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	logging.INFO("Fin ReportExcelGet")

	return nil
}

// ReportGetByDate godoc
//
//	@Summary		ReportGetByDate
//	@Description	Obtiene un reporte por fechas tanto como ingresos, y egresos
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			form				query		string						true	"day o month"
//	@Param			dateRangeRequest	body		schemas.DateRangeRequest	true	"Rango de fechas"
//	@Success		200					{object}	schemas.Response
//	@Router			/api/v1/report/get_by_date [post]
func (c *ReportController) ReportMovementByDate(ctx *fiber.Ctx) error {
	logging.INFO("Inicio ReportMovementByDate")
	form := ctx.Query("form")
	if (form != "day" && form != "month") || form == "" {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "El campo form debe ser day o month no puede estar vacio", fmt.Errorf("form debe ser day o month")))
	}
	var dateRangeRequest schemas.DateRangeRequest
	if err := ctx.BodyParser(&dateRangeRequest); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	fromDate, toDate, err := dateRangeRequest.GetParsedDates()
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	report, err := c.ReportService.ReportMovementByDate(fromDate, toDate, form)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	logging.INFO("Fin ReportMovementByDate")
	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    report,
		Message: "Reporte mensual obtenido con exito",
	})
}

// ReportMovementByDatePointSale godoc
//
//	@Summary		ReportMovementByDatePointSale
//	@Description	Obtiene un reporte por fechas de los diferentes puntos de ventas, tanto como ingresos, ingresos por cancha y egresos
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			form				query		string						true	"day o month"
//	@Param			dateRangeRequest	body		schemas.DateRangeRequest	true	"Rango de fechas"
//	@Success		200					{object}	schemas.Response
//	@Router			/api/v1/report/get_by_date_point_sale [post]
func (c *ReportController) ReportMovementByDatePointSale(ctx *fiber.Ctx) error {
	logging.INFO("Inicio ReportMovementByDatePointSale")
	form := ctx.Query("form")
	if (form != "day" && form != "month") || form == "" {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "El campo form debe ser day o month no puede estar vacio", fmt.Errorf("form debe ser day o month")))
	}
	var dateRangeRequest schemas.DateRangeRequest
	if err := ctx.BodyParser(&dateRangeRequest); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	fromDate, toDate, err := dateRangeRequest.GetParsedDates()
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	report, err := c.ReportService.ReportMovementByDatePointSale(fromDate, toDate, form)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	logging.INFO("Fin ReportMovementByDatePointSale")
	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    report,
		Message: "Reporte mensual obtenido con exito",
	})
}

// ReportProfitableProducts godoc
//
//	@Summary		ReportGetProfitableProducts
//	@Description	Obtiene un reporte de productos por fechas de los mas vendidos, los mas rentables y un ranking de los productos
//	@Tags			Report
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			dateRangeRequest	body		schemas.DateRangeRequest	true	"Rango de fechas"
//	@Success		200					{object}	schemas.Response{body=[]schemas.ReportProfitableProducts}
//	@Failure		400					{object}	schemas.Response
//	@Failure		401					{object}	schemas.Response
//	@Failure		422					{object}	schemas.Response
//	@Failure		404					{object}	schemas.Response
//	@Failure		500					{object}	schemas.Response
//	@Router			/api/v1/report/get_profitable_products [post]
func (r *ReportController) ReportProfitableProducts(ctx *fiber.Ctx) error {
	logging.INFO("Inicio ReportProfitableProducts")
	var dateRangeRequest schemas.DateRangeRequest
	if err := ctx.BodyParser(&dateRangeRequest); err != nil {
		return schemas.HandleError(ctx, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	fromDate, toDate, err := dateRangeRequest.GetParsedDates()
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	report, err := r.ReportService.ReportProfitableProducts(fromDate, toDate)
	if err != nil {
		return schemas.HandleError(ctx, err)
	}

	logging.INFO("Fin ReportProfitableProducts")
	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    report,
		Message: "Reporte de productos obtenido con exito",
	})
}
