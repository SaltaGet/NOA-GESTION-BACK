package controllers

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

// ClientGetByID godoc
//	@Summary		Get client by id
//	@Description	Get client by id
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Id del cliente"
//	@Success		200	{object}	schemas.Response{body=schemas.Client}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		403	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/client/{id} [get]
func (cl *ClientController) ClientGetByID(c *fiber.Ctx) error {
	logging.INFO("Obtener un cliente por ID")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	client, err := cl.ClientService.ClientGetByID(id)
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

	logging.INFO("Cliente obtenido con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    client,
		Message: "Cliente obtenido con éxito",
	})
}

// ClientGetAll godoc
//	@Summary		Get All Clients
//	@Description	Get All Clients
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	schemas.Response{body=[]schemas.Client}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		403	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/client/get_all [get]
func (cl *ClientController) ClientGetAll(c *fiber.Ctx) error {
	logging.INFO("Obtener todos los clientes")
	clients, err := cl.ClientService.ClientGetAll()
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

	logging.INFO("Clientes obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    clients,
		Message: "Clientes obtenidos con éxito",
	})
}

// ClientGetByName godoc
//	@Summary		Get Client By Name
//	@Description	Get Client By Name
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			name	query		string	true	"Name"
//	@Success		200		{object}	schemas.Response{body=[]schemas.Client}
//	@Failure		400		{object}	schemas.Response
//	@Failure		401		{object}	schemas.Response
//	@Failure		403		{object}	schemas.Response
//	@Failure		404		{object}	schemas.Response
//	@Failure		500		{object}	schemas.Response
//	@Router			/client/get_by_name [get]
func (cl *ClientController) ClientGetByName(c *fiber.Ctx) error {
	logging.INFO("Obtener un cliente por nombre")
	name := c.Query("name")
	if name == "" || len(name) < 3 {
		logging.ERROR("Error: El valor no debe de ser vacio o menor a 3 caracteres")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "El valor no debe de ser vacio o menor a 3 caracteres",
		})
	}

	clients, err := cl.ClientService.ClientGetByName(name)
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

	logging.INFO("Clientes obtenidos con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    clients,
		Message: "Clientes obtenidos con éxito",
	})
}

// ClientUpdate actualiza un cliente
//	@Summary		Actualizar un cliente
//	@Description	Actualizar un cliente
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			ClientUpdate	body		schemas.ClientUpdate	true	"Cliente a actualizar"
//	@Success		200				{object}	schemas.Response
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		403				{object}	schemas.Response
//	@Failure		404				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/client/update [put]
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
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}
	err := cl.ClientService.ClientUpdate(&clientUpdate)
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

	logging.INFO("Cliente actualizado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Cliente actualizado con éxito",
	})
}

// ClientDelete godoc
//	@Summary		Delete client by ID
//	@Description	Delete client by ID
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			id	path		string	true	"Id del cliente"
//	@Success		200	{object}	schemas.Response{body=schemas.Client}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		403	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/client/delete/{id} [delete]
func (cl *ClientController) ClientDelete(c *fiber.Ctx) error {
	logging.INFO("Eliminar un cliente")
	id := c.Params("id")
	if id == "" {
		logging.ERROR("Error: ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID is required",
		})
	}

	err := cl.ClientService.ClientDelete(id)
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

	logging.INFO("Cliente eliminado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Cliente eliminado con éxito",
	})
}

// CreateClient godoc
//	@Summary		Create client
//	@Description	Create client
//	@Tags			Client
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			clientCreate	body		schemas.ClientCreate	true	"Información del cliente"
//	@Success		200				{object}	schemas.Response
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		403				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/client/create [post]
func (cl *ClientController) CreateClient(c *fiber.Ctx) error {
	logging.INFO("Crear un cliente")
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
		logging.ERROR("Error: %s", err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: err.Error(),
		})
	}

	clientCreated, err := cl.ClientService.ClientCreate(&clientCreate)
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

	logging.INFO("Cliente creado con éxito")
	return c.Status(200).JSON(schemas.Response{
		Status:  true,
		Body:    clientCreated,
		Message: "Cliente creado con éxito",
	})
}
