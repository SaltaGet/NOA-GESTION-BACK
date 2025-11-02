package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func MemberRoutes(app *fiber.App){
	member := app.Group("/member", middleware.AuthMiddleware(), middleware.TenantMiddleware())

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
