package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func MemberRoutes(app *fiber.App){
	member := app.Group("/api/v1/member", middleware.AuthMiddleware(), middleware.AuthPointSaleMiddleware())

	member.Get("/get_all", GetController("MemberController", func(c *fiber.Ctx, ctrl *controllers.MemberController) error {
		return ctrl.MemberGetAll(c)
	}))

	member.Post("/create", GetController("MemberController", func(c *fiber.Ctx, ctrl *controllers.MemberController) error {
		return ctrl.MemberCreate(c)
	}))

	member.Get("/get/:id", GetController("MemberController", func(c *fiber.Ctx, ctrl *controllers.MemberController) error {
		return ctrl.MemberGetByID(c)
	}))
}
