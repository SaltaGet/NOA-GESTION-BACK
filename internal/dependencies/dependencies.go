package dependencies

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

type MainContainer struct {
	AuthController *controllers.AuthController
	UserController *controllers.UserController
	TenantController *controllers.TenantController
	// NotificationController *controllers.NotificationController
}

func NewApplication(mainDB *gorm.DB, cfg *schemas.EmailConfig) *MainContainer {
	mainRepo := &repositories.MainRepository{DB: mainDB}
	
	dialer := gomail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	emailService := &services.EmailService{Dialer: dialer}

	authServ := &services.AuthService{AuthRepository: mainRepo, UserRepository: mainRepo, TenantService: mainRepo, EmailService: emailService}
	userServ := &services.UserService{UserRepository: mainRepo}
	tenantServ := &services.TenantService{TenantRepository: mainRepo}
	// notificationServ := &services.NotificationService{NotificationRepository: mainRepo}

	return &MainContainer{
		AuthController: &controllers.AuthController{AuthService: authServ, EmailService: emailService},
		UserController: &controllers.UserController{UserService: userServ},
		TenantController: &controllers.TenantController{TenantService: tenantServ},
		// NotificationController: &controllers.NotificationController{SSEServer: sse, NotificationService: notificationServ},
	}
}

