package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// SupplierGetByID godoc
//	@Summary		Get Supplier By ID
//	@Description	Get a supplier by its ID within a specified workplace.
//	@Tags			Supplier
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string											true	"ID of the supplier"
//	@Success		200	{object}	schemas.Response{body=schemas.SupplierResponse}	"Supplier obtained with success"
//	@Failure		400	{object}	schemas.Response								"Bad Request"
//	@Failure		401	{object}	schemas.Response								"Auth is required"
//	@Failure		403	{object}	schemas.Response								"Not Authorized"
//	@Failure		404	{object}	schemas.Response								"Supplier not found"
//	@Failure		500	{object}	schemas.Response								"Internal server error"
//	@Router			/api/v1/supplier/{id} [get]
func (s *SupplierController) SupplierGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	supplier, err := s.SupplierService.SupplierGetByID(idint)
	if err != nil {
		return schemas.HandleError(c, err)
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
//	@Success		200	{object}	schemas.Response{body=[]schemas.SupplierResponseDTO}	"Suppliers obtained with success"
//	@Failure		400	{object}	schemas.Response										"Bad Request"
//	@Failure		401	{object}	schemas.Response										"Auth is required"
//	@Failure		403	{object}	schemas.Response										"Not Authorized"
//	@Failure		500	{object}	schemas.Response										"Internal server error"
//	@Router			/api/v1/supplier/get_all [get]
func (s *SupplierController) SupplierGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los proveedores")
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}

	search := &map[string]string{}
	name := c.Query("name")
	if name != "" {
		(*search)["name"] = name
	}
	identifier := c.Query("identifier")
	if identifier != "" {
		(*search)["identifier"] = identifier
	}
	companyName := c.Query("company_name")
	if companyName != "" {
		(*search)["company_name"] = companyName
	}
	email := c.Query("email")
	if email != "" {
		(*search)["email"] = email
	}
	isActive := c.Query("is_active")
	if isActive != "" {
		(*search)["is_active"] = isActive
	}

	suppliers, total, err := s.SupplierService.SupplierGetAll(int(limit), int(page), search)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Proveedores obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"suppliers": suppliers, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
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
//	@Param			supplier	body		schemas.SupplierCreate	true	"Details of the supplier to create"
//	@Success		200			{object}	schemas.Response		"Supplier created successfully"
//	@Failure		400			{object}	schemas.Response		"Bad Request"
//	@Failure		401			{object}	schemas.Response		"Auth is required"
//	@Failure		403			{object}	schemas.Response		"Not Authorized"
//	@Failure		422			{object}	schemas.Response		"Model is invalid"
//	@Failure		500			{object}	schemas.Response		"Internal server error"
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
		return schemas.HandleError(c, err)
	}
	
	id, err := s.SupplierService.SupplierCreate(&supplierCreate)
	if err != nil {
		return schemas.HandleError(c, err)
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
		return schemas.HandleError(c, err)
	}

	err := s.SupplierService.SupplierUpdate(&supplierUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
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
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	err = s.SupplierService.SupplierDelete(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Proveedor eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Proveedor eliminado con éxito",
	})
}

