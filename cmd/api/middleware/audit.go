package middleware

// import (
// 	"github.com/DanielChachagua/GestionCar/pkg/models"
// 	"github.com/DanielChachagua/GestionCar/pkg/repositories"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/google/uuid"
// )

// func AuditMiddleware() fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         err := c.Next() 

//         method := c.Method()
//         if (method == fiber.MethodPost || method == fiber.MethodPut || method == fiber.MethodDelete) &&
//             c.Response().StatusCode() >= 200 && c.Response().StatusCode() < 300 {

//             user, _ := c.Locals("user").(*models.User)
//             audit := models.AuditLog{
//                 ID:        uuid.NewString(),
//                 UserID:    "",
//                 Method:    method,
//                 Path:      c.Path(),
//             }
//             if user != nil {
//                 audit.UserID = user.ID
//             }

//             go repositories.Repo.DB.Create(&audit)
//         }

//         return err
//     }
// }






// type AuditLog struct {
// 	ID        string
// 	UserID    string
// 	Method    string
// 	Path      string
// 	Action    string
// 	TargetID  string
// 	// otros campos...
// }

// func DeleteClient(c *fiber.Ctx) error {
// 	clientID := c.Params("id")
// 	// ... lógica de borrado ...
// 	c.Locals("audit_info", map[string]interface{}{
// 			"action":    "delete_client",
// 			"client_id": clientID,
// 	})
// 	return c.SendStatus(fiber.StatusNoContent)
// }

// auditInfo, _ := c.Locals("audit_info").(map[string]interface{})
// if auditInfo != nil {
//     if v, ok := auditInfo["action"].(string); ok {
//         audit.Action = v
//     }
//     if v, ok := auditInfo["client_id"].(string); ok {
//         audit.TargetID = v
//     }
// }