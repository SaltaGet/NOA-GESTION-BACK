package middleware

import (
	// "github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/cmd/api/controllers"
	"github.com/DanielChachagua/GestionCar/pkg/database"
	"github.com/DanielChachagua/GestionCar/pkg/dependencies"
	"github.com/DanielChachagua/GestionCar/pkg/key"
	"github.com/DanielChachagua/GestionCar/pkg/models"
	"github.com/DanielChachagua/GestionCar/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return unauthorized(c, "Token no proporcionado")
		}

		deps, ok := c.Locals(key.AppKey).(*dependencies.Application)
		if !ok {
			return unauthorized(c, "Dependencias no proporcionadas")
		}

		claims, err := utils.VerifyToken(token)
		if err != nil {
			return unauthorized(c, "Token inválido")
		}

		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			return unauthorized(c, "Claims inválidos")
		}

		userId, ok := mapClaims["id"].(string)
		if !ok {
			return unauthorized(c, "ID inválido en el token")
		}

		isAdmin := getBoolClaim(mapClaims, "is_admin_tenant")
		tenantID := getStringClaim(mapClaims, "tenant_id")

		if tenantID != "" {
			return handleTenantUser(c, deps, tenantID, userId, isAdmin)
		}

		return handleSuperAdmin(c, deps, userId)
	}
}

// func AuthMiddleware() fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		token := c.Get("Authorization")
// 		if token == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Token no proporcionado",
// 			})
// 		}

// 		// ctx := c.UserContext()
// 		// deps := ctx.Value(key.AppKey).(*dependencies.Application)

// 		deps, ok := c.Locals(key.AppKey).(*dependencies.Application)
// 		if !ok {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Dependencias no proporcionadas",
// 			})
// 		}

// 		claims, err := utils.VerifyToken(token)

// 		if err != nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Token inválido",
// 			})
// 		}

// 		mapClaims, ok := claims.(jwt.MapClaims)
// 		if !ok {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "Claims inválidos",
// 			})
// 		}

// 		userId, ok := mapClaims["id"].(string)
// 		if !ok {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
// 				"error": "ID inválido en el token",
// 			})
// 		}

// 		isAdmin := false
// 		if adm, ok := mapClaims["is_admin_tenant"].(bool); ok {
// 			isAdmin = adm
// 		}

// 		var tenantID string
// 		if tid, ok := mapClaims["tenant_id"].(string); ok {
// 			tenantID = tid
// 		}

// 		if tenantID != "" {
// 			tenant, err := deps.TenantController.TenantService.TenantGetByID(tenantID)
// 			if err != nil {
// 				return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
// 					Status:  false,
// 					Body:    nil,
// 					Message: err.Error(),
// 				})
// 			}

// 			connection, err := utils.Decrypt(tenant.Connection)
// 			if err != nil {
// 				return c.Status(500).JSON(models.Response{
// 					Status:  false,
// 					Body:    nil,
// 					Message: err.Error(),
// 				})
// 			}

// 			db, err := database.GetTenantDB(connection)
// 			if err != nil {
// 				return c.Status(500).JSON(models.Response{
// 					Status:  false,
// 					Body:    nil,
// 					Message: err.Error(),
// 				})
// 			}

// 			container := dependencies.NewTenantContainer(db)
// 			setupTenantControllers(c, container)
// 			// c.Locals("TenantContainer", container)
// 			// attendanceCtrl := &controllers.AttendanceController{
// 			// 	AttendanceService: container.Services.Attendance,
// 			// }
// 			// clientCtrl := &controllers.ClientController{
// 			// 	ClientService: container.Services.Client,
// 			// }
// 			// employeeCtrl := &controllers.EmployeeController{
// 			// 	EmployeeService: container.Services.Employee,
// 			// }
// 			// expenseCtrl := &controllers.ExpenseController{
// 			// 	ExpenseService: container.Services.Expense,
// 			// }
// 			// incomeCtrl := &controllers.IncomeController{
// 			// 	IncomeService: container.Services.Income,
// 			// }
// 			// memberCtrl := &controllers.MemberController{
// 			// 	MemberService: container.Services.Member,
// 			// }
// 			// movementCtrl := &controllers.MovementTypeController{
// 			// 	MovementTypeService: container.Services.Movement,
// 			// }
// 			// permissionCtrl := &controllers.PermissionController{
// 			// 	PermissionService: container.Services.Permission,
// 			// }
// 			// productCtrl := &controllers.ProductController{
// 			// 	ProductService: container.Services.Product,
// 			// }
// 			// purchaseCtrl := &controllers.PurchaseOrderController{
// 			// 	PurchaseOrderService: container.Services.Purchase,
// 			// }
// 			// purchaseProductCtrl := &controllers.PurchaseProductController{
// 			// 	PurchaseProductService: container.Services.PurchaseProduct,
// 			// }
// 			// resumeCtrl := &controllers.ResumeController{
// 			// 	ResumeExpenseService: container.Services.Resume,
// 			// 	ResumeIncomeService: container.Services.Resume,
// 			// }
// 			// roleCtrl := &controllers.RoleController{
// 			// 	RoleService: container.Services.Role,
// 			// }
// 			// serviceCtrl := &controllers.ServiceController{
// 			// 	ServiceService: container.Services.Service,
// 			// }
// 			// supplierCtrl := &controllers.SupplierController{
// 			// 	SupplierService: container.Services.Supplier,
// 			// }
// 			// vehicleCtrl := &controllers.VehicleController{
// 			// 	VehicleService: container.Services.Vehicle,
// 			// }

// 			// c.Locals("AttendanceController", attendanceCtrl)
// 			// c.Locals("ClientController", clientCtrl)
// 			// c.Locals("EmployeeController", employeeCtrl)
// 			// c.Locals("ExpenseController", expenseCtrl)
// 			// c.Locals("IncomeController", incomeCtrl)
// 			// c.Locals("MemberController", memberCtrl)
// 			// c.Locals("MovementTypeController", movementCtrl)
// 			// c.Locals("PermissionController", permissionCtrl)
// 			// c.Locals("ProductController", productCtrl)
// 			// c.Locals("PurchaseOrderController", purchaseCtrl)
// 			// c.Locals("PurchaseProductController", purchaseProductCtrl)
// 			// c.Locals("ResumeController", resumeCtrl)
// 			// c.Locals("RoleController", roleCtrl)
// 			// c.Locals("ServiceController", serviceCtrl)
// 			// c.Locals("SupplierController", supplierCtrl)
// 			// c.Locals("VehicleController", vehicleCtrl)

// 			userFromToken := models.AuthenticatedUser{}
// 			if !isAdmin {
// 				user, err := container.Services.Member.MemberGetByID(userId)
// 				if err != nil {
// 					return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 						Status:  false,
// 						Body:    nil,
// 						Message: "Error interno",
// 					})
// 				}

// 				permissions := func(perms []models.Permission) []string {
// 					names := make([]string, 0, len(perms))
// 					for _, p := range perms {
// 						names = append(names, p.Code)
// 					}
// 					return names
// 				}(user.Role.Permissions)

// 				userFromToken = models.AuthenticatedUser{
// 					ID:            user.ID,
// 					FirstName:     user.FirstName,
// 					LastName:      user.LastName,
// 					Username:      user.Username,
// 					IsAdminTenant: false,
// 					RoleID:        &user.Role.ID,
// 					RoleName:      &user.Role.Name,
// 					Permissions:   permissions,
// 					TenantID:      &tenantID,
// 					TenantName:    &tenant.Name,
// 					Identifier:    &tenant.Identifier,
// 				}
// 			} else {
// 				user, err := deps.AuthController.AuthService.CurrentUser(userId)

// 				if err != nil {
// 					if errResp, ok := err.(*models.ErrorStruc); ok {
// 						return c.Status(errResp.StatusCode).JSON(models.Response{
// 							Status:  false,
// 							Body:    nil,
// 							Message: errResp.Message,
// 						})
// 					}
// 					return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 						Status:  false,
// 						Body:    nil,
// 						Message: "Error interno",
// 					})
// 				}

// 				for _, userTenant := range user.UserTenants {
// 					if userTenant.TenantID == tenantID {
// 						userFromToken = models.AuthenticatedUser{
// 							ID:            user.ID,
// 							FirstName:     user.FirstName,
// 							LastName:      user.LastName,
// 							Username:      user.Username,
// 							IsAdminTenant: userTenant.IsAdmin,
// 							RoleID:        nil,
// 							RoleName:      nil,
// 							Permissions:   nil,
// 							TenantID:      &tenantID,
// 							TenantName:    &tenant.Name,
// 							Identifier:    &tenant.Identifier,
// 						}
// 						break
// 					}
// 				}
// 			}

// 			c.Locals("user", &userFromToken)

// 			return c.Next()
// 		}

// 		user, err := deps.AuthController.AuthService.CurrentUser(userId)

// 		if err != nil {
// 			if errResp, ok := err.(*models.ErrorStruc); ok {
// 				return c.Status(errResp.StatusCode).JSON(models.Response{
// 					Status:  false,
// 					Body:    nil,
// 					Message: errResp.Message,
// 				})
// 			}
// 			return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: "Error interno",
// 			})
// 		}

// 		userFromToken := models.AuthenticatedUser{
// 			ID:            user.ID,
// 			FirstName:     user.FirstName,
// 			LastName:      user.LastName,
// 			Username:      user.Username,
// 			IsAdminTenant: true,
// 			RoleID:        nil,
// 			RoleName:      nil,
// 			Permissions:   nil,
// 			TenantID:      nil,
// 			TenantName:    nil,
// 			Identifier:    nil,
// 		}

// 		c.Locals("user", &userFromToken)

// 		return c.Next()
// 	}
// }



func setupTenantControllers(c *fiber.Ctx, container *dependencies.TenantContainer) {
	controllersMap := map[string]interface{}{
		"AttendanceController":        &controllers.AttendanceController{AttendanceService: container.Services.Attendance},
		"ClientController":            &controllers.ClientController{ClientService: container.Services.Client},
		"EmployeeController":          &controllers.EmployeeController{EmployeeService: container.Services.Employee},
		"ExpenseController":           &controllers.ExpenseController{ExpenseService: container.Services.Expense},
		"IncomeController":            &controllers.IncomeController{IncomeService: container.Services.Income},
		"MemberController":            &controllers.MemberController{MemberService: container.Services.Member},
		"MovementTypeController":      &controllers.MovementTypeController{MovementTypeService: container.Services.Movement},
		"PermissionController":        &controllers.PermissionController{PermissionService: container.Services.Permission},
		"ProductController":           &controllers.ProductController{ProductService: container.Services.Product},
		"PurchaseOrderController":     &controllers.PurchaseOrderController{PurchaseOrderService: container.Services.Purchase},
		"PurchaseProductController":   &controllers.PurchaseProductController{PurchaseProductService: container.Services.PurchaseProduct},
		"ResumeController":            &controllers.ResumeController{ResumeExpenseService: container.Services.Resume, ResumeIncomeService: container.Services.Resume},
		"RoleController":              &controllers.RoleController{RoleService: container.Services.Role},
		"ServiceController":           &controllers.ServiceController{ServiceService: container.Services.Service},
		"SupplierController":          &controllers.SupplierController{SupplierService: container.Services.Supplier},
		"VehicleController":           &controllers.VehicleController{VehicleService: container.Services.Vehicle},
	}

	for name, ctrl := range controllersMap {
		c.Locals(name, ctrl)
	}
}

func getTenantMemberUser(c *fiber.Ctx, container *dependencies.TenantContainer, tenantID string, tenant *models.Tenant, userId string) (models.AuthenticatedUser, error) {
	member, err := container.Services.Member.MemberGetByID(userId)
	if err != nil {
		return models.AuthenticatedUser{}, internalError(c, "Error interno")
	}

	permissions := make([]string, len(member.Role.Permissions))
	for i, p := range member.Role.Permissions {
		permissions[i] = p.Code
	}

	return models.AuthenticatedUser{
		ID:            member.ID,
		FirstName:     member.FirstName,
		LastName:      member.LastName,
		Username:      member.Username,
		IsAdminTenant: false,
		RoleID:        &member.Role.ID,
		RoleName:      &member.Role.Name,
		Permissions:   permissions,
		TenantID:      &tenantID,
		TenantName:    &tenant.Name,
		Identifier:    &tenant.Identifier,
	}, nil
}

func getAdminUser(c *fiber.Ctx, deps *dependencies.Application, tenantID string, tenant *models.Tenant, userId string) (models.AuthenticatedUser, error) {
	user, err := deps.AuthController.AuthService.CurrentUser(userId)
	if err != nil {
		if errResp, ok := err.(*models.ErrorStruc); ok {
			return models.AuthenticatedUser{}, c.Status(errResp.StatusCode).JSON(models.Response{
				Status:  false,
				Body:    nil,
				Message: errResp.Message,
			})
		}
		return models.AuthenticatedUser{}, internalError(c, "Error interno")
	}

	for _, t := range user.UserTenants {
		if t.TenantID == tenantID {
			return models.AuthenticatedUser{
				ID:            user.ID,
				FirstName:     user.FirstName,
				LastName:      user.LastName,
				Username:      user.Username,
				IsAdmin: user.IsAdmin,
				IsAdminTenant: t.IsAdmin,
				TenantID:      &tenantID,
				TenantName:    &tenant.Name,
				Identifier:    &tenant.Identifier,
			}, nil
		}
	}

	return models.AuthenticatedUser{}, unauthorized(c, "No autorizado en este tenant")
}

func handleTenantUser(c *fiber.Ctx, deps *dependencies.Application, tenantID, userId string, isAdmin bool) error {
	tenant, err := deps.TenantController.TenantService.TenantGetByID(tenantID)
	if err != nil {
		return unauthorized(c, err.Error())
	}

	connection, err := utils.Decrypt(tenant.Connection)
	if err != nil {
		return internalError(c, err.Error())
	}

	db, err := database.GetTenantDB(connection)
	if err != nil {
		return internalError(c, err.Error())
	}

	container := dependencies.NewTenantContainer(db)
	setupTenantControllers(c, container)

	var user models.AuthenticatedUser
	if isAdmin {
		user, err = getAdminUser(c, deps, tenantID, tenant, userId)
	} else {
		user, err = getTenantMemberUser(c, container, tenantID, tenant, userId)
	}

	if err != nil {
		return err
	}

	c.Locals("user", &user)
	return c.Next()
}

func handleSuperAdmin(c *fiber.Ctx, deps *dependencies.Application, userId string) error {
	user, err := deps.AuthController.AuthService.CurrentUser(userId)
	if err != nil {
		if errResp, ok := err.(*models.ErrorStruc); ok {
			return c.Status(errResp.StatusCode).JSON(models.Response{
				Status:  false,
				Body:    nil,
				Message: errResp.Message,
			})
		}
		return internalError(c, "Error interno")
	}

	authUser := models.AuthenticatedUser{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Username:      user.Username,
		IsAdmin:       user.IsAdmin,
		IsAdminTenant: true,
	}

	c.Locals("user", &authUser)
	return c.Next()
}

func unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": message,
	})
}

func internalError(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
		Status:  false,
		Body:    nil,
		Message: message,
	})
}

func getBoolClaim(claims jwt.MapClaims, key string) bool {
	val, ok := claims[key].(bool)
	return ok && val
}

func getStringClaim(claims jwt.MapClaims, key string) string {
	val, ok := claims[key].(string)
	if ok {
		return val
	}
	return ""
}