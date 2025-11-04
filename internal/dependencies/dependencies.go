package dependencies

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services"
	"gorm.io/gorm"
)

type MainContainer struct {
	AuthController *controllers.AuthController
	UserController *controllers.UserController
	TenantController *controllers.TenantController
}

func NewApplication(mainDB *gorm.DB) *MainContainer {
	mainRepo := &repositories.MainRepository{DB: mainDB}

	authServ := &services.AuthService{AuthRepository: mainRepo, UserRepository: mainRepo, TenantService: mainRepo}
	userServ := &services.UserService{UserRepository: mainRepo}
	tenantServ := &services.TenantService{TenantRepository: mainRepo}

	return &MainContainer{
		AuthController: &controllers.AuthController{AuthService: authServ},
		UserController: &controllers.UserController{UserService: userServ},
		TenantController: &controllers.TenantController{TenantService: tenantServ},
	}
}

