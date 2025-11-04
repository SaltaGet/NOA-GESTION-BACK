package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

//	Permissions GetAll godoc
//	@Summary		Permissions GetAll
//	@Description	Permissions GetAll required auth token
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PermissionResponse}	"Members obtenidos con éxito"
//	@Failure		400	{object}	schemas.Response									"Bad Request"
//	@Failure		401	{object}	schemas.Response									"Auth is required"
//	@Failure		403	{object}	schemas.Response									"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/permission/get_all [get]
func (p *PermissionController) PermissionGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los permisos")
	permissions, err := p.PermissionService.PermissionGetAll()
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

	logging.INFO("Permisos obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    permissions,
		Message: "Permisos obtenidos con éxito",
	})
}

//	Permissions GetAllToMe godoc
//
//	@Summary		Permissions GetAlltoME
//	@Description	Permissions GetAllToMe required auth token
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PermissionResponse}	"Members obtenidos con éxito"
//	@Failure		400	{object}	schemas.Response									"Bad Request"
//	@Failure		401	{object}	schemas.Response									"Auth is required"
//	@Failure		403	{object}	schemas.Response									"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/permission/get_to_me [get]
func (p *PermissionController) PermissionGetToMe(c *fiber.Ctx) error {
	logging.INFO("Obtener todos mis permisos")
	user := c.Locals("user").(*schemas.AuthenticatedUser)

	if user.IsAdminTenant {
		permissions, err := p.PermissionService.PermissionGetAll()
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

		logging.INFO("Tiene todos los permsios")
		return c.Status(fiber.StatusForbidden).JSON(schemas.Response{
			Status:  false,
			Body:    permissions,
			Message: "Tiene todos los permsios",
		})
	}

	permissions, err := p.PermissionService.PermissionGetToMe(*user.RoleID)
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

	logging.INFO("Permisos obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    permissions,
		Message: "Permisos obtenidos con éxito",
	})
}
