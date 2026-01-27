package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

// EcommerceGetByID godoc
//
//	@Summary		EcommerceGetByID
//	@Description	Obtener una compra electrónica por ID
//	@Tags			Ecommerce
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"ID"
//	@Success		200	{object}	schemas.Response{body=schemas.EcommerceResponse}
//	@Router			/api/v1/ecommerce/get/{id} [get]
func (ec *EcommerceController) EcommerceGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Message: "ID parameter is required",
		})
	}
	idInt, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	ecommerce, err := ec.EcommerceService.GetByID(idInt)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Message: "Información de compra electrónica obtenida correctamente",
		Body:    ecommerce,
	})
}

// EcommerceGetByReference godoc
//
//	@Summary		EcommerceGetByReference
//	@Description	Obtener una compra electrónica por referencia
//	@Tags			Ecommerce
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			reference	path		string	true	"Referencia"
//	@Success		200			{object}	schemas.Response{body=schemas.EcommerceResponse}
//	@Router			/api/v1/ecommerce/get_by_reference/{reference} [get]
func (ec *EcommerceController) EcommerceGetByReference(c *fiber.Ctx) error {
	reference := c.Params("reference")
	if reference == "" {
		log.Error().Msg("Reference parameter is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Message: "Reference parametro es requirido",
		})
	}
	if !validators.IsValidUUIDv4(reference) {
		log.Error().Msg("Reference must be a valid UUIDv4")
		return c.Status(fiber.StatusUnprocessableEntity).JSON(schemas.Response{
			Status:  false,
			Message: "La referencia debe ser un UUIDv4 válido",
		})
	}

	ecommerce, err := ec.EcommerceService.GetByReference(reference)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Message: "Información de compra electrónica obtenida correctamente",
		Body:    ecommerce,
	})
}

// EcommerceGetAll godoc
//
//	@Summary		EcommerceGetAll
//	@Description	Obtener todas las compras electrónicas
//	@Tags			Ecommerce
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			page	query		int		false	"pagina"	default(1)
//	@Param			limit	query		int		false	"limite"	default(10)
//	@Param			status	query		string	false	"estado"
//	@Success		200		{object}	schemas.Response{body=[]schemas.EcommerceResponseDTO}
//	@Router			/api/v1/ecommerce/get_all [get]
func (ec *EcommerceController) EcommerceGetAll(c *fiber.Ctx) error {
	page := c.Query("page", "1")
	limit := c.Query("limit", "10")
	pg, err := validators.IntValidate(page)
	if err != nil {
		return schemas.HandleError(c, err)
	}
	lt, err := validators.IntValidate(limit)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	state := c.Query("status", "")
	var status *string
	if state == "" {
		status = nil
	} else {
		status = &state
	}


	ecommerces, err := ec.EcommerceService.GetAll(pg, lt, status)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Message: "Información de compra electrónica obtenida correctamente",
		Body:    ecommerces,
	})
}

// DepositProdEcommerceUpdateStatusuctGetAll godoc
//
//	@Summary		EcommerceUpdateStatus
//	@Description	Actualizar estado de una compra electrónica
//	@Tags			Ecommerce
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response
//	@Router			/api/v1/ecommerce/update_status [put]
func (ec *EcommerceController) EcommerceUpdateStatus(c *fiber.Ctx) error {
	var statusUpdate schemas.EcommerceStatusUpdate
	if err := c.BodyParser(&statusUpdate); err != nil {
		return schemas.HandleError(c, err)
	}
	if err := statusUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := ec.EcommerceService.UpdateStatus(&statusUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Message: "Estado de compra electrónica actualizado correctamente",
		Body:    nil,
	})
}
