package middleware

import (
	// "context"

	"github.com/DanielChachagua/GestionCar/pkg/dependencies"
	"github.com/DanielChachagua/GestionCar/pkg/key"
	"github.com/gofiber/fiber/v2"
)

// func InjectApp(app *dependencies.Application) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		ctx := c.UserContext()
// 		ctx = context.WithValue(ctx, key.AppKey, app)
// 		// ctx = context.WithValue(ctx, key.TenantDBKey, tenantApp)
// 		c.SetUserContext(ctx)
// 		return c.Next()
// 	}
// }

func InjectApp(app *dependencies.Application) fiber.Handler {
    return func(c *fiber.Ctx) error {
        c.Locals(key.AppKey, app)
        return c.Next()
    }
}