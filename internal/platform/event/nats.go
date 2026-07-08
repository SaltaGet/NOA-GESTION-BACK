package event

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	NatsConn *nats.Conn
)

// InitNats inicializa la conexión con NATS
func InitNats() error {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = "nats://localhost:4222"
	}

	opts := []nats.Option{
		nats.Name("NOA-GESTION-BACK"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			log.Printf("⚠️ Desconectado de NATS: %v", err)
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			log.Printf("✅ Reconectado a NATS: %s", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			log.Printf("❌ Conexión a NATS cerrada")
		}),
	}

	var err error
	NatsConn, err = nats.Connect(url, opts...)
	if err != nil {
		log.Printf("⚠️ NATS no disponible en %s: %v. La app continuará sin eventos NATS.", url, err)
		NatsConn = nil
		return err
	}

	log.Println("✅ NATS conectado correctamente")
	return nil
}

// CloseNats cierra la conexión a NATS
func CloseNats() error {
	if NatsConn != nil {
		NatsConn.Close()
	}
	return nil
}

// Publish publica un evento en NATS
func Publish(subject string, payload interface{}) error {
	if NatsConn == nil {
		return fmt.Errorf("NATS no disponible")
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("error al serializar payload: %w", err)
	}

	if err := NatsConn.Publish(subject, data); err != nil {
		return fmt.Errorf("error al publicar evento en NATS: %w", err)
	}

	return nil
}

// IsAvailable verifica si NATS está disponible
func IsAvailable() bool {
	return NatsConn != nil
}
