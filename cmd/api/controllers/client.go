package controllers

import (
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/validators"
	"github.com/gofiber/fiber/v2"
)

// ClientGetByID godoc
//
//	@Summary		Get client by id
//	@Description	Get client by id
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Id del cliente"
//	@Success		200	{object}	schemas.Response{body=schemas.ClientResponse}
//	@Router			/api/v1/client/{id} [get]
func (cl *ClientController) ClientGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un cliente por ID")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil{
		return schemas.HandleError(c, err)
	}

	client, err := cl.ClientService.ClientGetByID(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Cliente obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    client,
		Message: "Cliente obtenido con éxito",
	})
}

// ClientGetByFilter godoc
//
//	@Summary		ClientGetByFilter
//	@Description	Obtener clientes por filtro nombre, apellido o idenficador
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			search	query		string	true	"Search"
//	@Success		200		{object}	schemas.Response{body=[]schemas.ClientResponseDTO}
//	@Router			/api/v1/client [get]
func (cl *ClientController) ClientGetByFilter(c *fiber.Ctx) error {
	logging.INFO("Obtener clientes")
	search := c.Query("search")
	if len(search) < 3 {
		logging.ERROR("Error: El valor no debe de ser vacio o menor a 3 caracteres")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El valor no debe de ser vacio o menor a 3 caracteres",
		})
	}

	clients, err := cl.ClientService.ClientGetByFilter(search)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Clientes obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    clients,
		Message: "Clientes obtenidos con éxito",
	})
}

// ClientGetAll godoc
//
//	@Summary		Get All Clients
//	@Description	Get All Clients
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			limit		query		int		false	"Limite por pagina, default 10"
//	@Param			page		query		int		false	"Pagina, default 1"
//	@Param			identifier	query		string	false	"Identificador del cliente"
//	@Param			first_name	query		string	false	"Nombre del cliente"
//	@Param			last_name	query		string	false	"Apellido del cliente"
//	@Param			email		query		string	false	"Correo del cliente"
//	@Success		200			{object}	schemas.Response{body=[]schemas.ClientResponseDTO}
//	@Router			/api/v1/client/get_all [get]
func (cl *ClientController) ClientGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los clientes")
	limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
	if err != nil {
		limit = 10
	}
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}

	search := &map[string]string{}
	identifier := c.Query("identifier")
	if identifier != "" {
		(*search)["identifier"] = identifier
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

	clients, total, err := cl.ClientService.ClientGetAll(limit, page, search)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	logging.INFO("Clientes obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    map[string]any{"clients": clients, "total": total, "page": page, "limit": limit, "total_pages": totalPages},
		Message: "Clientes obtenidos con éxito",
	})
}

// ClientCreate godoc
//
//	@Summary		ClientCreate
//	@Description	Crear un cliente
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			clientCreate	body		schemas.ClientCreate	true	"Información del cliente"
//	@Success		200				{object}	schemas.Response
//	@Router			/api/v1/client/create [post]
func (cl *ClientController) ClientCreate(c *fiber.Ctx) error {
	logging.INFO("Crear un cliente")

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	var clientCreate schemas.ClientCreate
	if err := c.BodyParser(&clientCreate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := clientCreate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	clientCreated, err := cl.ClientService.ClientCreate(user.ID, &clientCreate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Cliente creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    clientCreated,
		Message: "Cliente creado con éxito",
	})
}

// ClientUpdate godoc
//
//	@Summary		Actualizar un cliente
//	@Description	Actualizar un cliente
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			ClientUpdate	body		schemas.ClientUpdate	true	"Cliente a actualizar"
//	@Success		200				{object}	schemas.Response
//	@Router			/api/v1/client/update [put]
func (cl *ClientController) ClientUpdate(c *fiber.Ctx) error {
	logging.INFO("Actualizar un cliente")
	var clientUpdate schemas.ClientUpdate
	if err := c.BodyParser(&clientUpdate); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "Invalid request" + err.Error(),
		})
	}
	if err := clientUpdate.Validate(); err != nil {
		return schemas.HandleError(c, err)
	}

	err := cl.ClientService.ClientUpdate(&clientUpdate)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Cliente actualizado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Cliente actualizado con éxito",
	})
}

// ClientDelete godoc
//
//	@Summary		Delete client by ID
//	@Description	Eliminar un cliente por ID
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Param			id	path		string	true	"Id del cliente"
//	@Success		200	{object}	schemas.Response
//	@Router			/api/v1/client/delete/{id} [delete]
func (cl *ClientController) ClientDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar un cliente")
	id := c.Params("id")
	idint, err := validators.IdValidate(id)
	if err != nil{
		return schemas.HandleError(c, err)
	}

	err = cl.ClientService.ClientDelete(idint)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	logging.INFO("Cliente eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Cliente eliminado con éxito",
	})
}
