package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

//	Permissions GetAll godoc
//	@Summary		Permissions GetAll
//	@Description	Permissions GetAll required auth token
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PermissionResponse}	"Permisos obtenidos con éxito"
//	@Router			/api/v1/permission/get_all [get]
func (p *PermissionController) PermissionGetAll(c *fiber.Ctx) error {
	permissions, err := p.PermissionService.PermissionGetAll()
	if err != nil {
		return schemas.HandleError(c, err)
	}

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
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.PermissionResponse}	"Permisos obtenidos con éxito"
//	@Router			/api/v1/permission/get_to_me [get]
func (p *PermissionController) PermissionGetToMe(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.AuthenticatedUser)

	if user.IsAdmin {
		permissions, err := p.PermissionService.PermissionGetAll()
		if err != nil {
			return schemas.HandleError(c, err)
		}

		return c.Status(200).JSON(schemas.Response{
			Status:  false,
			Body:    permissions,
			Message: "Tiene todos los permsios",
		})
	}

	permissions, err := p.PermissionService.PermissionGetToMe(user.RoleID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    permissions,
		Message: "Permisos obtenidos con éxito",
	})
}

//	PermissionUpdateAll godoc
//
//	@Summary		PermissionUpdateAll
//	@Description	update all
//	@Tags			Permission
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response	"Permisos obtenidos con éxito"
//	@Router			/api/v1/permission/update_all [put]
func (p *PermissionController) PermissionUpdateAll(c *fiber.Ctx) error {
		err := p.PermissionService.PermissionUpdateAll()
		if err != nil {
			return schemas.HandleError(c, err)
		}

		return c.Status(200).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Permsios Actualizados",
		})
}
