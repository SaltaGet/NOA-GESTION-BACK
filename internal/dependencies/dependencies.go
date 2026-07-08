package dependencies

import (
	"database/sql"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/repositories"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/services"
)	

type MainContainer struct {
	AuthController *controllers.AuthController
	UserController *controllers.UserController
	TenantController *controllers.TenantController
	// NotificationController *controllers.NotificationController
	PlanController *controllers.PlanController
	NewsController *controllers.NewsController
	FeedbackController *controllers.FeedbackController
	ModuleController *controllers.ModuleController
	CredentialController *controllers.CredentialController
}

func NewApplication(mainDB *sql.DB) *MainContainer {
	mainRepo := &repositories.MainRepository{DB: mainDB}
	
	emailService := &services.EmailService{}

	authServ := &services.AuthService{AuthRepository: mainRepo, TenantService: mainRepo, EmailService: emailService, PlanRepository: mainRepo, ModuleRepository: mainRepo}
	userServ := &services.UserService{UserRepository: mainRepo}
	tenantServ := &services.TenantService{TenantRepository: mainRepo, EmailService: emailService}
	// notificationServ := &services.NotificationService{NotificationRepository: mainRepo}
	planServ := &services.PlanService{PlanRepository: mainRepo}
	newsServ := &services.NewsService{NewsRepository: mainRepo}
	feedbackServ := &services.FeedbackService{FeedbackRepository: mainRepo}
	moduleServ := &services.ModuleService{ModuleRepository: mainRepo}
	credentialServ := &services.CredentialService{CredentialRepository: mainRepo}



	return &MainContainer{
		AuthController: &controllers.AuthController{AuthService: authServ, EmailService: emailService},
		UserController: &controllers.UserController{UserService: userServ},
		TenantController: &controllers.TenantController{TenantService: tenantServ},
		// NotificationController: &controllers.NotificationController{SSEServer: sse, NotificationService: notificationServ},
		PlanController: &controllers.PlanController{PlanService: planServ},
		NewsController: &controllers.NewsController{NewsService: newsServ},
		FeedbackController: &controllers.FeedbackController{FeedbackService: feedbackServ},
		ModuleController: &controllers.ModuleController{ModuleService: moduleServ},
		CredentialController: &controllers.CredentialController{CredentialService: credentialServ},
	}
}

