package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/alexandrevicenzi/go-sse"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

// NotificationAlert se suscribe a un canal SSE específico del tenant.
func (n *NotificationController) NotificationAlert(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.AuthenticatedUser)
	// Creamos un nombre de canal único para el tenant.
	channelName := fmt.Sprintf("/api/v1/notification/alert/tenant/%d", user.TenantID)
	c.Locals("sse_channel", channelName)

	// Headers SSE después de CORS
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache, no-transform")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	// Configurar el SSE handler
	handler := fasthttpadaptor.NewFastHTTPHandler(n.SSEServer)
	handler(c.Context())

	return nil
}

// NotificationAlertProduct se suscribe a un canal SSE específico del tenant para productos.
func (n *NotificationController) NotificationAlertProduct(c *fiber.Ctx) error {
	user := c.Locals("user").(*schemas.AuthenticatedUser)
	// Canal específico para notificaciones de productos del tenant.
	channelName := fmt.Sprintf("/api/v1/notification/alert_product/tenant/%d", user.TenantID)
	c.Locals("sse_channel", channelName)

	// Headers SSE después de CORS
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache, no-transform")
	c.Set("Connection", "keep-alive")
	c.Set("X-Accel-Buffering", "no")

	// Configurar el SSE handler
	handler := fasthttpadaptor.NewFastHTTPHandler(n.SSEServer)
	handler(c.Context())

	return nil
}

// SendStockNotification envía notificaciones de stock bajo para un tenant específico.
func (n *NotificationController) SendStockNotification(tenantID int64) error {
	log.Printf("Enviando notificación de stock para tenant %d", tenantID)

	products, err := n.NotificationService.NotificationStock(tenantID)
	if err != nil {
		log.Printf("Error obteniendo productos con stock bajo para tenant %d: %v", tenantID, err)
		return err
	}

	if len(products) == 0 {
		return nil
	}

	loc, err := time.LoadLocation("America/Argentina/Buenos_Aires")
	if err != nil {
		log.Printf("Error cargando zona horaria: %v", err)
		loc = time.UTC
	}

	response := schemas.Response{
		Status:  true,
		Message: "Productos con stock bajo",
		Body: map[string]any{
			"event": "alert-stock",
			"response": map[string]any{
				"data":     products,
				"count":    len(products),
				"datetime": time.Now().In(loc).Format(time.RFC3339),
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error serializando datos: %v", err)
		return err
	}

	// Enviar al canal específico del tenant
	channelName := fmt.Sprintf("/api/v1/notification/alert/tenant/%d", tenantID)
	n.SSEServer.SendMessage(channelName, sse.NewMessage("", string(data), "stock-notification"))

	log.Printf("Notificación de stock enviada para tenant %d: %d productos", tenantID, len(products))
	return nil
}

// SendNotificationRefreshProducts envía una notificación para refrescar productos a un tenant específico.
func (n *NotificationController) SendNotificationRefreshProducts(tenantID int64) error {
	log.Printf("Enviando notificación de stock")

	response := schemas.Response{
		Status:  true,
		Message: "Productos con stock bajo",
		Body: map[string]any{
			"event": "refresh-products",
			"response": "Refresca la lista de productos sr del front",
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error serializando datos: %v", err)
		return err
	}

	channelName := fmt.Sprintf("/api/v1/notification/alert_product/tenant/%d", tenantID)
	n.SSEServer.SendMessage(channelName, sse.NewMessage("", string(data), "refresh-products"))

	log.Printf("Notificación de refresco de productos enviada para tenant %d", tenantID)
	return nil
}

