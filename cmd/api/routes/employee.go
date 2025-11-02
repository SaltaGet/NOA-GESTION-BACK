package routes

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func EmployeeRoutes(app *fiber.App){
	emp := app.Group("/employee", middleware.AuthMiddleware(), middleware.TenantMiddleware())

	emp.Get("/get_all", GetController("EmployeeController", func(c *fiber.Ctx, ctrl *controllers.EmployeeController) error {
		return ctrl.GetAllEmployees(c)
	}))

	emp.Get("/get_by_name", GetController("EmployeeController", func(c *fiber.Ctx, ctrl *controllers.EmployeeController) error {
		return ctrl.GetEmployeeByName(c)
	}))

	emp.Post("/create", GetController("EmployeeController", func(c *fiber.Ctx, ctrl *controllers.EmployeeController) error {
		return ctrl.CreateEmployee(c)
	}))

	emp.Put("/update", GetController("EmployeeController", func(c *fiber.Ctx, ctrl *controllers.EmployeeController) error {
		return ctrl.UpdateEmployee(c)
	}))

	emp.Delete("/delete/:id", GetController("EmployeeController", func(c *fiber.Ctx, ctrl *controllers.EmployeeController) error {
		return ctrl.DeleteEmployee(c)
	}))

	emp.Get("/:id", GetController("EmployeeController", func(c *fiber.Ctx, ctrl *controllers.EmployeeController) error {
		return ctrl.GetEmployeeByID(c)
	}))

}
