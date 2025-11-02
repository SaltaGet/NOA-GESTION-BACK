package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, controllers *controllers.AuthController) {
	auth := app.Group("api/v1/auth")

	//  Login godoc
	//	@Summary		Login user
	//	@Description	Login user required identifier and password
	//	@Tags			Auth
	//	@Accept			json
	//	@Produce		json
	//	@Param			credentials	body		models.AuthLogin	true	"Credentials"
	//	@Success		200			{object}	models.Response
	//	@Failure		400			{object}	models.Response
	//	@Failure		401			{object}	models.Response
	//	@Failure		422			{object}	models.Response
	//	@Failure		404			{object}	models.Response
	//	@Failure		500			{object}	models.Response
	//	@Router			api/v1/auth/login [post]
	auth.Post("/login", controllers.AuthLogin)

	//  LoginTenant godoc
	//	@Summary		Login Tenant
	//	@Description	Login tenant required tenant_id
	//	@Tags			Auth
	//	@Accept			json
	//	@Produce		json
	//	@Security		BearerAuth
	//	@Param			tenant_id	path		string	true	"tenant_id"
	//	@Success		200			{object}	models.Response
	//	@Failure		400			{object}	models.Response
	//	@Failure		401			{object}	models.Response
	//	@Failure		403			{object}	models.Response
	//	@Failure		404			{object}	models.Response
	//	@Failure		422			{object}	models.Response
	//	@Failure		500			{object}	models.Response
	//	@Router			api/v1/auth/tenant_login/{tenant_id} [get]
	auth.Post("/tenant_login/:tenant_id", controllers.AuthTenant)

	auth.Post("/tenant_logout", controllers.LogoutTenant)

	auth.Post("/logout", controllers.Logout)
}
