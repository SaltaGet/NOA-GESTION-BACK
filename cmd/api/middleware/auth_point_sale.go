package middleware

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/gofiber/fiber/v2"
)

func PointSaleMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		pointSaleID, ok := c.Locals("point_sale_id").(int64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "Se requiere autenticacion en punto de venta",
			})
		}

		if pointSaleID < 1 {
			return c.Status(401).JSON(schemas.Response{
				Status:  false,
				Body:    nil,
				Message: "Se requiere autenticacion en punto de venta",
			})

		} 

		return c.Next()
	}
}
