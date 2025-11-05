package controllers

import (
	"strconv"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

//	Login godoc
//
//	@Summary		Login user
//	@Description	Login user required identifier and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		schemas.AuthLogin	true	"Credentials"
//	@Success		200			{object}	schemas.Response
//	@Failure		400			{object}	schemas.Response
//	@Failure		401			{object}	schemas.Response
//	@Failure		422			{object}	schemas.Response
//	@Failure		404			{object}	schemas.Response
//	@Failure		500			{object}	schemas.Response
//	@Router			/api/v1/auth/login [post]
func (a *AuthController) AuthLogin(c *fiber.Ctx) error {
	logging.INFO("Login")
	var loginRequest schemas.AuthLogin
	if err := c.BodyParser(&loginRequest); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	if err := loginRequest.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, err)
	}

	token, err := a.AuthService.AuthLogin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
	}

	c.Cookie(cookie)

	logging.INFO("Login exitoso")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login exitoso",
	})
}

//	LoginPointSale godoc
//
//	@Summary		LoginPointSale
//	@Description	Login al punto de venta
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//
//	@Security		CookieAuth
//
//	@Param			point_sale_id	path		string	true	"id del punto de venta"
//	@Success		200				{object}	schemas.Response
//	@Failure		400				{object}	schemas.Response
//	@Failure		401				{object}	schemas.Response
//	@Failure		403				{object}	schemas.Response
//	@Failure		404				{object}	schemas.Response
//	@Failure		422				{object}	schemas.Response
//	@Failure		500				{object}	schemas.Response
//	@Router			/api/v1/auth/login_point_sale/{point_sale_id} [post]
func (a *AuthController) AuthPointSale(c *fiber.Ctx) error {
	logging.INFO("Login tenant")
	id := c.Params("point_sale_id")
	if id == "" {
		logging.ERROR("ID is required")
		return c.Status(400).JSON(schemas.Response{
			Status:  false,
			Body:    nil,
			Message: "ID es requerido",
		})
	}

	pointSaleID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logging.ERROR("tenant_id must be a number")
		return c.Status(400).JSON(schemas.ErrorResponse(400, "tenant_id debe ser un número", err))
	}

	user := c.Locals("user").(*schemas.AuthenticatedUser)

	token, err := a.AuthService.AuthPointSale(user, pointSaleID)
	if err != nil {
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	logging.INFO("Login tenant exitoso")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login a Punto de venta exitoso, token enviado en cookie",
	})
}

//	LogoutPointSale godoc
//
//	@Summary		LogoutPointSale
//	@Description	Logout del punto de venta
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//
//	@Security		CookieAuth
//
//	@Success		200	{object}	schemas.Response
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		403	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/auth/logout_point_sale [post]
func (a *AuthController) LogoutPointSale(c *fiber.Ctx) error {
	logging.INFO("Logout tenant")
	member := c.Locals("user").(*schemas.AuthenticatedUser)


	token, err := a.AuthService.LogoutPointSale(member)
	if err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
		Expires:  time.Now().AddDate(1, 0, 0),
	}

	c.Cookie(cookie)

	logging.INFO("Logout Punto de venta exitoso")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Logout de Punto de venta exitoso, token enviado en cookie",
	})
}

//	Logout godoc
//
//	@Summary		Logout user
//	@Description	Logout user
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/auth/logout [post]
func (a *AuthController) Logout(ctx *fiber.Ctx) error {
	logging.INFO("Logout")
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
	})

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Logout exitoso",
	})
}

//	CurrentUser godoc
//
//	@Summary		CurrentUser user
//	@Description	Obtener usuario actual
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response{body=schemas.AuthenticatedUser}
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/auth/current_user [get]
func (a *AuthController) CurrentUser(c *fiber.Ctx) error {
	logging.INFO("Obtener usuario actual")
	user := c.Locals("user").(*schemas.AuthenticatedUser)

	logging.INFO("Usuario actual obtenido")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    user,
		Message: "Usuario actual obtenido",
	})
}

//	LoginAdmin godoc
//
//	@Summary		Login Admin user
//	@Description	Required identifier and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		schemas.AuthLoginAdmin	true	"Credentials"
//	@Success		200			{object}	schemas.Response
//	@Failure		400			{object}	schemas.Response
//	@Failure		401			{object}	schemas.Response
//	@Failure		422			{object}	schemas.Response
//	@Failure		404			{object}	schemas.Response
//	@Failure		500			{object}	schemas.Response
//	@Router			/api/v1/auth/login_admin [post]
func (a *AuthController) AuthLoginAdmin(c *fiber.Ctx) error {
	logging.INFO("Login")
	var loginRequest schemas.AuthLoginAdmin
	if err := c.BodyParser(&loginRequest); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, schemas.ErrorResponse(400, "Error al parsear el cuerpo de la solicitud", err))
	}

	if err := loginRequest.Validate(); err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, err)
	}

	token, err := a.AuthService.AuthLoginAdmin(loginRequest.Username, loginRequest.Password)
	if err != nil {
		logging.ERROR("Error: %s", err.Error())
		return schemas.HandleError(c, err)
	}

	cookie := &fiber.Cookie{
		Name:     "access_token_admin",
		Value:    token,
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
	}

	c.Cookie(cookie)

	logging.INFO("Login exitoso")
	return c.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Login exitoso",
	})
}

//	LogoutAdmin godoc
//
//	@Summary		Logout Admin user
//	@Description	logout user admin
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200	{object}	schemas.Response
//	@Failure		400	{object}	schemas.Response
//	@Failure		401	{object}	schemas.Response
//	@Failure		422	{object}	schemas.Response
//	@Failure		404	{object}	schemas.Response
//	@Failure		500	{object}	schemas.Response
//	@Router			/api/v1/auth/logout_admin [post]
func (a *AuthController) LogoutAdmin(ctx *fiber.Ctx) error {
	logging.INFO("Logout")
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token_admin",
		Value:    "",
		HTTPOnly: true,
		Secure:   false,  // poner en true para prod
		SameSite: "None", // para prod : "Strict",
	})

	return ctx.Status(fiber.StatusOK).JSON(schemas.Response{
		Status:  true,
		Body:    nil,
		Message: "Logout exitoso",
	})
}
