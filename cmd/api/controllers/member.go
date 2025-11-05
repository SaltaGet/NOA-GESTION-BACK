package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

//	Member godoc
//
//	@Summary		Memeber GetAll
//	@Description	Memeber GetAll required auth token
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.MemberDTO}	"Members obtenidos con éxito"
//	@Failure		400	{object}	schemas.Response							"Bad Request"
//	@Failure		401	{object}	schemas.Response							"Auth is required"
//	@Failure		403	{object}	schemas.Response							"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/member/get_all [get]
func (m *MemberController) MemberGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los miembros")
	memebers, err := m.MemberService.MemberGetAll()
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

	logging.INFO("Miembros obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    memebers,
		Message: "Miembros obtenidos con éxito",
	})
}

//	Member godoc
//
//	@Summary		Memeber GetAll
//	@Description	Memeber GetAll required auth token
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string											true	"Member ID"
//	@Success		200	{object}	schemas.Response{body=[]schemas.MemberResponse}	"Members obtenidos con éxito"
//	@Failure		400	{object}	schemas.Response								"Bad Request"
//	@Failure		401	{object}	schemas.Response								"Auth is required"
//	@Failure		403	{object}	schemas.Response								"Not Authorized"
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/member/get/{id} [get]
func (m *MemberController) MemberGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los miembros")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	memebers, err := m.MemberService.MemberGetByID(id)
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

	logging.INFO("Miembros obtenidos con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    memebers,
		Message: "Miembros obtenidos con éxito",
	})
}

//	Member godoc
//
//	@Summary		Memeber Create
//	@Description	Memeber Create required auth token
//	@Tags			Member
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			member_create	body		schemas.MemberCreate	true	"MemberCreate"
//	@Success		200				{object}	schemas.Response		"Members obtenidos con éxito"
//	@Failure		400				{object}	schemas.Response		"Bad Request"
//	@Failure		401				{object}	schemas.Response		"Auth is required"
//	@Failure		403				{object}	schemas.Response		"Not Authorized"
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/member/create [post]
func (m *MemberController) MemberCreate(c *fiber.Ctx) error {
	logging.INFO("Crear miembro")
	user := c.Locals("user").(*schemas.AuthenticatedUser)

	var memberCreate schemas.MemberCreate
	if err := c.BodyParser(&memberCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := memberCreate.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(422).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	id, err := m.MemberService.MemberCreate(&memberCreate, user)
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

	logging.INFO("Miembro creado con éxito")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    id,
		Message: "Miembro creado con éxito",
	})
}
