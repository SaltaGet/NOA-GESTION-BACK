package routes

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func CategoryRoutes(app *fiber.App) {
	category := app.Group("/api/v1/category",  middleware.AuthMiddleware(), middleware.InjectionDependsTenant(), middleware.AuthPointSaleMiddleware())

	category.Post("/create", GetController("CategoryController", func(c *fiber.Ctx, ctrl *controllers.CategoryController) error {
		return ctrl.CategoryCreate(c)
	}))

	category.Get("/get_all", GetController("CategoryController", func(c *fiber.Ctx, ctrl *controllers.CategoryController) error {
		return ctrl.CategoryGetAll(c)
	}))

	category.Put("/update", GetController("CategoryController", func(c *fiber.Ctx, ctrl *controllers.CategoryController) error {
		return ctrl.CategoryUpdate(c)
	}))

	category.Get("/get/:id", GetController("CategoryController", func(c *fiber.Ctx, ctrl *controllers.CategoryController) error {
		return ctrl.CategoryGet(c)
	}))

	category.Delete("/delete/:id", GetController("CategoryController", func(c *fiber.Ctx, ctrl *controllers.CategoryController) error {
		return ctrl.CategoryDelete(c)
	}))
}

