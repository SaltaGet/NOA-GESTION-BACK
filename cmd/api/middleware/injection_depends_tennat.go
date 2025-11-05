package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func InjectionDependsTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		member := c.Locals("user").(*schemas.AuthenticatedUser)

		db, err := database.GetTenantDB("", member.TenantID)
		if err != nil {
			return schemas.ErrorResponse(401, "No autenticado", err)
		}

		container := dependencies.NewTenantContainer(db)
		setupTenantControllers(c, container)

		return c.Next()
	}
}

func setupTenantControllers(c *fiber.Ctx, container *dependencies.TenantContainer) {
	controllersMap := map[string]any{
		"ClientController":       &controllers.ClientController{ClientService: container.Services.Client},
		"ExpenseController":      &controllers.ExpenseController{ExpenseService: container.Services.Expense},
		"IncomeController":       &controllers.IncomeController{IncomeService: container.Services.Income},
		"MemberController":       &controllers.MemberController{MemberService: container.Services.Member},
		"MovementTypeController": &controllers.MovementTypeController{MovementTypeService: container.Services.Movement},
		"PermissionController":   &controllers.PermissionController{PermissionService: container.Services.Permission},
		"PointSaleController":   &controllers.PointSaleController{PointSaleService: container.Services.PointSale},
		"ProductController":      &controllers.ProductController{ProductService: container.Services.Product},
		"RoleController":         &controllers.RoleController{RoleService: container.Services.Role},
		"SupplierController":     &controllers.SupplierController{SupplierService: container.Services.Supplier},
	}

	for name, ctrl := range controllersMap {
		c.Locals(name, ctrl)
	}
}
