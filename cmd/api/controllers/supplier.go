package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// SupplierGetByID godoc
//	@Summary		Get Supplier By ID
//	@Description	Get a supplier by its ID within a specified workplace.
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string									true	"ID of the supplier"
//	@Success		200	{object}	schemas.Response{body=schemas.Supplier}	"Supplier obtained with success"
//	@Failure		400	{object}	schemas.Response						"Bad Request"
//	@Failure		401	{object}	schemas.Response						"Auth is required"
//	@Failure		403	{object}	schemas.Response						"Not Authorized"
//	@Failure		404	{object}	schemas.Response						"Supplier not found"
//	@Failure		500	{object}	schemas.Response						"Internal server error"
//	@Router			/api/v1/supplier/{id} [get]
func (s *SupplierController) SupplierGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	supplier, err := s.SupplierService.SupplierGetByID(id)
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

	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    supplier,
		Message: "Proveedor obtenido con éxito",
	})
}

// SupplierGetAll godoc
//	@Summary		Get All Suppliers
//	@Description	Get All Suppliers
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.Supplier}	"Suppliers obtained with success"
//	@Failure		400	{object}	schemas.Response							"Bad Request"
//	@Failure		401	{object}	schemas.Response							"Auth is required"
//	@Failure		403	{object}	schemas.Response							"Not Authorized"
//	@Failure		500	{object}	schemas.Response							"Internal server error"
//	@Router			/api/v1/supplier/get_all [get]
func (s *SupplierController) SupplierGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los proveedores")
	suppliers, err := s.SupplierService.SupplierGetAll()
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

	logging.INFO("Proveedores obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    suppliers,
		Message: "Proveedores obtenidos con éxito",
	})
}

// SupplierGetByName godoc
//	@Summary		Get Supplier By Name
//	@Description	Fetches suppliers from either laundry or workshop based on the provided name and workplace.
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			name	query		string										true	"Name of the Supplier"
//	@Success		200		{object}	schemas.Response{body=[]schemas.Supplier}	"List of suppliers"
//	@Failure		400		{object}	schemas.Response							"Bad Request"
//	@Failure		401		{object}	schemas.Response							"Auth is required"
//	@Failure		403		{object}	schemas.Response							"Not Authorized"
//	@Failure		500		{object}	schemas.Response							"Internal server error"
//	@Router			/api/v1/supplier/get_by_name [get]
func (s *SupplierController) SupplierGetByName(c *fiber.Ctx) error {
	logging.INFO("Obtener proveedores por nombre")
	name := c.Query("name")
	if name == "" || len(name) < 3 {
		logging.ERROR("Error: El valor no debe de ser vacio o menor a 3 caracteres")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El valor no debe de ser vacio o menor a 3 caracteres",
		})
	}

	supplies, err := s.SupplierService.SupplierGetByName(name)
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

	logging.INFO("Proveedores obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    supplies,
		Message: "Proveedores obtenidos con éxito",
	})
}

// SupplierCreate godoc
//	@Summary		Create Supplier
//	@Description	Creates a new supplier within the specified workplace.
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			supplier	body		schemas.SupplierCreate			true	"Details of the supplier to create"
//	@Success		200			{object}	schemas.Response{body=string}	"Supplier created successfully"
//	@Failure		400			{object}	schemas.Response				"Bad Request"
//	@Failure		401			{object}	schemas.Response				"Auth is required"
//	@Failure		403			{object}	schemas.Response				"Not Authorized"
//	@Failure		422			{object}	schemas.Response				"Model is invalid"
//	@Failure		500			{object}	schemas.Response				"Internal server error"
//	@Router			/api/v1/supplier/create [post]
func (s *SupplierController) SupplierCreate(c *fiber.Ctx) error {
	logging.INFO("Crear proveedor")
	var supplierCreate schemas.SupplierCreate
	if err := c.BodyParser(&supplierCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := supplierCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := s.SupplierService.SupplierCreate(&supplierCreate)
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

	logging.INFO("Proveedor creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Proveedor creado con éxito",
	})
}

// SupplierUpdate godoc
//	@Summary		Update Supplier
//	@Description	Update a supplier's information from the specified workplace.
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			body	body		schemas.SupplierUpdate	true	"Supplier information"
//	@Success		200		{object}	schemas.Response		"Supplier updated with success"
//	@Failure		400		{object}	schemas.Response		"Bad Request"
//	@Failure		401		{object}	schemas.Response		"Auth is required"
//	@Failure		403		{object}	schemas.Response		"Not Authorized"
//	@Failure		404		{object}	schemas.Response		"Supplier not found"
//	@Failure		422		{object}	schemas.Response		"Model is invalid"
//	@Failure		500		{object}	schemas.Response		"Internal server error"
//	@Router			/api/v1/supplier/update [put]
func (s *SupplierController) SupplierUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar proveedor")
	var supplierUpdate schemas.SupplierUpdate
	if err := c.BodyParser(&supplierUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := supplierUpdate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	err := s.SupplierService.SupplierUpdate(&supplierUpdate)
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

	logging.INFO("Proveedor editado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Proveedor editado con éxito",
	})
}

// SupplierDeleteByID godoc
//	@Summary		Delete Supplier
//	@Description	Deletes a supplier based on the provided ID and workplace context.
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"ID of the supplier"
//	@Success		200	{object}	schemas.Response	"Supplier deleted with success"
//	@Failure		400	{object}	schemas.Response	"Bad Request"
//	@Failure		401	{object}	schemas.Response	"Auth is required"
//	@Failure		403	{object}	schemas.Response	"Not Authorized"
//	@Failure		404	{object}	schemas.Response	"Supplier not found"
//	@Failure		500	{object}	schemas.Response	"Internal server error"
//	@Router			/api/v1/supplier/delete/{id} [delete]
func (s *SupplierController) SupplierDeleteByID(c *fiber.Ctx) error {
	logging.INFO("Eliminar proveedor")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	err := s.SupplierService.SupplierDelete(id)
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

	logging.INFO("Proveedor eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Proveedor eliminado con éxito",
	})
}

