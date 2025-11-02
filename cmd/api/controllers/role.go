package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// GetRoles godoc
//
//	@Summary		Retrieve roles for a user in a specific workplace
//	@Description	This function fetches roles based on the user's role and workplace identifier
//	             from the context. It requires both user and workplace information to be present
//	             in the request context.
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.RoleResponse}	"Roles retrieved successfully"
//	@Failure		400	{object}	schemas.Response								"Bad request if user or workplace is missing"
//	@Failure		500	{object}	schemas.Response								"Internal server error on failure"
//	@Router			/role/get_all [get]
func (r *RoleController) RoleGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los roles")
	roles, err := r.RoleService.RoleGetAll()
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

	logging.INFO("Roles obtenidos exitosamente")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    roles,
		Message: "Roles obtenidos exitosamente",
	})
}

// Create Role godoc
//	@Summary		Retrieve roles for a user in a specific workplace
//	@Description	This function fetches roles based on the user's role and workplace identifier
//	@Tags			Role
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			role	body		schemas.RoleCreate	true	"Role object"
//	@Success		200		{object}	schemas.Response		"Roles retrieved successfully"
//	@Failure		400		{object}	schemas.Response		"Bad request if user or workplace is missing"
//	@Failure		500		{object}	schemas.Response		"Internal server error on failure"
//	@Router			/role/create [post]
func (r *RoleController) RoleCreate(c *fiber.Ctx) error {
	logging.INFO("Crear rol")
	var roleCreate schemas.RoleCreate
	if err := c.BodyParser(&roleCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Bad request" + err.Error(),
		})
	}
	if err := roleCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(422).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := r.RoleService.RoleCreate(&roleCreate)
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
	
	logging.INFO("Rol creado exitosamente")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Rol creado exitosamente",
	})
}
