package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// ArcaEmitInvoice godoc
//
//	@Summary		ArcaEmitInvoice
//	@Description	### Emitir Factura
//	@Tags			Arca
//	@Accept			json
//	@Produce		json
//	@Param			datos	body		schemas.FacturaRequest	true	"datos de la factura"
//	@Success		200		{object}	schemas.Response
//	@Router			/api/v1/arca/emit_invoice [post]
func (a *ArcaController) ArcaEmitInvoice(c *fiber.Ctx) error {
	var factReq schemas.FacturaRequest
	if err := c.BodyParser(&factReq); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := factReq.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	factura, err := a.ArcaService.EmitInvoice(user, int64(1), &factReq, true)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    factura,
		Message: "Factura emitida con éxito",
	})
}