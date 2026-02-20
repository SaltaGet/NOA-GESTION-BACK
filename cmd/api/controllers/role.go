package controllers

import (
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// GetRoleByID godoc
//
//	@Summary		GetRoleByID
//	@Description	Obtener role por id
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string										true	"Role ID"
//	@Success		200	{object}	schemas.Response{body=schemas.RoleResponse}	"Roles retrieved successfully"
//	@Router			/api/v1/role/get/{id} [get]
func (r *RoleController) RoleGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	intID, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	role, err := r.RoleService.RoleGetByID(intID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    role,
		Message: "Roles obtenidos exitosamente",
	})
}

// GetRoles godoc
//
//	@Summary		Retrieve roles for a user in a specific workplace
//	@Description	This function fetches roles based on the user's role and workplace identifier
//	             from the context. It requires both user and workplace information to be present
//	             in the request context.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.RoleResponse}	"Roles retrieved successfully"
//	@Router			/api/v1/role/get_all [get]
func (r *RoleController) RoleGetAll(c *fiber.Ctx) error {
	roles, err := r.RoleService.RoleGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    roles,
		Message: "Roles obtenidos exitosamente",
	})
}

// Create Role godoc
//
//	@Summary		Retrieve roles for a user in a specific workplace
//	@Description	This function fetches roles based on the user's role and workplace identifier
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			role	body		schemas.RoleCreate	true	"Role object"
//	@Success		200		{object}	schemas.Response	"Roles created successfully"
//	@Router			/api/v1/role/create [post]
func (r *RoleController) RoleCreate(c *fiber.Ctx) error {
	var roleCreate schemas.RoleCreate
	if err := validators.ValidateRequest(c, &roleCreate); err != nil {
		return schemas.HandleError(c, err)
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)
	id, err := r.RoleService.RoleCreate(member.ID, &roleCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Rol creado exitosamente",
	})
}

// UpdateRole godoc
//
//	@Summary		UpdateRole
//	@Description	Actualizar rol
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			role	body		schemas.RoleUpdate	true	"Role object"
//	@Success		200		{object}	schemas.Response	"Roles updated successfully"
//	@Router			/api/v1/role/update [put]
func (r *RoleController) RoleUpdate(c *fiber.Ctx) error {
	var roleCreate schemas.RoleUpdate
	if err := validators.ValidateRequest(c, &roleCreate); err != nil {
		return schemas.HandleError(c, err)
	}

	if roleCreate.ID == 1 {
		return schemas.HandleError(c, schemas.ErrorResponse(422, "No se puede actualizar el rol admin", fmt.Errorf("no se puede actualizar el rol admin")))
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)
	err := r.RoleService.RoleUpdate(member.ID, &roleCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Rol editado exitosamente",
	})
}
