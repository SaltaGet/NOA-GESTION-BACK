package controllers

import (
	"strconv"


	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

//	Member godoc
//
//	@Summary		Memeber GetAll
//	@Description	Memeber GetAll required auth token
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//
//	@Param			limit		query		int													false	"Limite por pagina, default 10"
//	@Param			page		query		int													false	"Numero de pagina, default 1"
//	@Param			first_name	query		string												false	"Nombre del miembro"
//	@Param			last_name	query		string												false	"Apellido del miembro"
//	@Param			username	query		string												false	"Username"
//	@Param			email		query		string												false	"Correo del miembro"
//
//	@Success		200			{object}	schemas.Response{body=[]schemas.MemberResponseDTO}	"Members obtenidos con éxito"
//	@Router			/api/v1/member/get_all [get]
func (m *MemberController) MemberGetAll(c *fiber.Ctx) error {
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}

	search := &map[string]string{}
	username := c.Query("identifier")
	if username != "" {
		(*search)["identifier"] = username
	}
	firstName := c.Query("first_name")
	if firstName != "" {
		(*search)["first_name"] = firstName
	}
	lastName := c.Query("last_name")
	if lastName != "" {
		(*search)["last_name"] = lastName
	}
	email := c.Query("email")
	if email != "" {
		(*search)["email"] = email
	}
	isActive := c.Query("is_active")
	if isActive != "" {
		(*search)["is_active"] = isActive
	}

	members, total, err := m.MemberService.MemberGetAll(int(limit), int(page), search)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"data": members, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Miembros obtenidos con éxito",
	})
}

// MemberGetByID godoc
//
//	@Summary		MemberGetByID
//	@Description	Obtener miembro por ID
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string											true	"Member ID"
//	@Success		200	{object}	schemas.Response{body=[]schemas.MemberResponse}	"Members obtenidos con éxito"
//	@Router			/api/v1/member/get/{id} [get]
func (m *MemberController) MemberGetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	memebers, err := m.MemberService.MemberGetByID(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    memebers,
		Message: "Miembros obtenidos con éxito",
	})
}

// MemberCreate godoc
//
//	@Summary		Memeber Create
//	@Description	Memeber Create required auth token
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			member_create	body		schemas.MemberCreate	true	"MemberCreate"
//	@Success		200				{object}	schemas.Response		"Members creado con éxito"
//	@Router			/api/v1/member/create [post]
func (m *MemberController) MemberCreate(c *fiber.Ctx) error {
	
	var memberCreate schemas.MemberCreate
	if err := c.BodyParser(&memberCreate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := memberCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	plan := c.Locals("current_plan").(*schemas.PlanResponseDTO)
member := c.Locals("user").(*schemas.AuthenticatedUser)
	id, err := m.MemberService.MemberCreate(member.ID, &memberCreate, plan)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Miembro creado con éxito",
	})
}

// MemberUpdate godoc
//
//	@Summary		MemberUpdate
//	@Description	Update member
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			member_update	body		schemas.MemberUpdate	true	"MemberUpdate"
//	@Success		200				{object}	schemas.Response		"Members actualizado con éxito"
//	@Router			/api/v1/member/update [put]
func (m *MemberController) MemberUpdate(c *fiber.Ctx) error {

	var memberUpdate schemas.MemberUpdate
	if err := c.BodyParser(&memberUpdate); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := memberUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)
	err := m.MemberService.MemberUpdate(member.ID, &memberUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Miembro editado con éxito",
	})
}

// MemberUpdatePassword godoc
//
//	@Summary		MemberUpdatePassword
//	@Description	Update member password
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			member_update_update	body		schemas.MemberUpdatePassword	true	"MemberUpdatePassword"
//	@Success		200						{object}	schemas.Response				"Members obtenidos con éxito"
//	@Router			/api/v1/member/update_password [put]
func (m *MemberController) MemberUpdatePassword(c *fiber.Ctx) error {
	var memberUpdatePassword schemas.MemberUpdatePassword
	if err := c.BodyParser(&memberUpdatePassword); err != nil {
		log.Err(err).Msg("Error al parsear el body")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := memberUpdatePassword.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)

	err := m.MemberService.MemberUpdatePassword(member.ID, &memberUpdatePassword)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Password del miembro editado con éxito",
	})
}

// MemberDelete godoc
//
//	@Summary		MemberDelete
//	@Description	Eliminar miembro por ID
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string				true	"Member ID"
//	@Success		200	{object}	schemas.Response	"Members obtenidos con éxito"
//	@Router			/api/v1/member/delete/{id} [delete]
func (m *MemberController) MemberDelete(c *fiber.Ctx) error {
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	member := c.Locals("user").(*schemas.AuthenticatedUser)
	err = m.MemberService.MemberDelete(member.ID, idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Miembros obtenidos con éxito",
	})
}
