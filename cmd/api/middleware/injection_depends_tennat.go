package middleware

import (
	// "github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/tenant"
	tenant_cache "github.com/SaltaGet/NOA-GESTION-BACK/internal/cache/tenant"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
	// "github.com/SaltaGet/NOA-GESTION-BACK/internal/key"
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

		// ⚡ Usar el container cacheado
		container := tenant_cache.GetTenantContainer(db, member.TenantID)

		// Guardarlo en el contexto
		c.Locals("tenant", container)

		return c.Next()
	}
}

// func InjectionDependsTenant() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		member := c.Locals("user").(*schemas.AuthenticatedUser)

// 			db, err := database.GetTenantDB("", member.TenantID)
// 		if err != nil {
// 			return schemas.ErrorResponse(401, "No autenticado", err)
// 		}

// 		container := dependencies.NewTenantContainer(db)
// 		setupTenantControllers(c, container)

//			return c.Next()
//		}
//	}


// func setupTenantControllers(c *fiber.Ctx, container *dependencies.TenantContainer) {
// 	deps := c.Locals(key.AppKey).(*dependencies.MainContainer)
// 	notifController := &deps.NotificationController

// 	controllersMap := map[string]any{
// 		"CashRegisterController": &controllers.CashRegisterController{CashRegisterService: container.Services.CashRegister},
// 		"CategoryController":     &controllers.CategoryController{CategoryService: container.Services.Category},
// 		"ClientController":       &controllers.ClientController{ClientService: container.Services.Client},
// 		"DepositController":      &controllers.DepositController{DepositService: container.Services.Deposit},
// 		"ExpenseBuyController":   &controllers.ExpenseBuyController{ExpenseBuyService: container.Services.ExpenseBuy},
// 		"ExpenseOtherController": &controllers.ExpenseOtherController{ExpenseOtherService: container.Services.ExpenseOther},
// 		"IncomeOtherController":  &controllers.IncomeOtherController{IncomeOtherService: container.Services.IncomeOther},
// 		"IncomeSaleController":   &controllers.IncomeSaleController{IncomeSaleService: container.Services.IncomeSale},
// 		"MemberController":       &controllers.MemberController{MemberService: container.Services.Member},
// 		"MovementStockController": &controllers.MovementStockController{
// 			MovementStockService:   container.Services.MovementStock,
// 			NotificationController: *notifController,
// 		},
// 		"PermissionController":   &controllers.PermissionController{PermissionService: container.Services.Permission},
// 		"PointSaleController":    &controllers.PointSaleController{PointSaleService: container.Services.PointSale},
// 		"ProductController":      &controllers.ProductController{ProductService: container.Services.Product},
// 		"ReportController":       &controllers.ReportController{ReportService: container.Services.Report},
// 		"RoleController":         &controllers.RoleController{RoleService: container.Services.Role},
// 		"StockController":        &controllers.StockController{StockService: container.Services.Stock},
// 		"SupplierController":     &controllers.SupplierController{SupplierService: container.Services.Supplier},
// 		"TypeMovementController": &controllers.TypeMovementController{TypeMovementService: container.Services.TypeMovement},
// 	}

// 	for name, ctrl := range controllersMap {
// 		c.Locals(name, ctrl)
// 	}
// }

// package middleware

// import (
// 	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/dependencies"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
// 	"github.com/gofiber/fiber/v2"
// )

// func InjectionDependsTenant() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		member := c.Locals("user").(*schemas.AuthenticatedUser)

// 			db, err := database.GetTenantDB("", member.TenantID)
// 		if err != nil {
// 			return schemas.ErrorResponse(401, "No autenticado", err)
// 		}

// 		container := dependencies.NewTenantContainer(db)
// 		setupTenantControllers(c, container)

// 		return c.Next()
// 	}
// }

// func setupTenantControllers(c *fiber.Ctx, container *dependencies.TenantContainer) {
// 	controllersMap := map[string]any{
// 		"CashRegisterController": &controllers.CashRegisterController{CashRegisterService: container.Services.CashRegister},
// 		"CategoryController":     &controllers.CategoryController{CategoryService: container.Services.Category},
// 		"ClientController":       &controllers.ClientController{ClientService: container.Services.Client},
// 		"DepositController":      &controllers.DepositController{DepositService: container.Services.Deposit},
// 		"ExpenseBuyController":   &controllers.ExpenseBuyController{ExpenseBuyService: container.Services.ExpenseBuy},
// 		"ExpenseOtherController": &controllers.ExpenseOtherController{ExpenseOtherService: container.Services.ExpenseOther},
// 		"IncomeOtherController":       &controllers.IncomeOtherController{IncomeOtherService: container.Services.IncomeOther},
// 		"IncomeSaleController":       &controllers.IncomeSaleController{IncomeSaleService: container.Services.IncomeSale},
// 		"MemberController":       &controllers.MemberController{MemberService: container.Services.Member},
// 		"PermissionController":   &controllers.PermissionController{PermissionService: container.Services.Permission},
// 		"PointSaleController":   &controllers.PointSaleController{PointSaleService: container.Services.PointSale},
// 		"ProductController":      &controllers.ProductController{ProductService: container.Services.Product},
// 		"ReportController":       &controllers.ReportController{ReportService: container.Services.Report},
// 		"RoleController":         &controllers.RoleController{RoleService: container.Services.Role},
// 		"StockController":        &controllers.StockController{StockService: container.Services.Stock},
// 		"SupplierController":     &controllers.SupplierController{SupplierService: container.Services.Supplier},
// 		"TypeMovementController": &controllers.TypeMovementController{TypeMovementService: container.Services.TypeMovement},
// 	}

// 	for name, ctrl := range controllersMap {
// 		c.Locals(name, ctrl)
// 	}
// }
