package initial

import (
	"log"
	"os"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func InitEmail() (*schemas.EmailConfig) {
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	username := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	if host == "" || port == "" || username == "" || password == "" {
		log.Fatalf("Error: configuración de email incompleta en las variables de entorno")
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatalf("Error al convertir el puerto a entero: %v", err)
	}


	emailCfg := &schemas.EmailConfig{
		Host:     host,
		Port:     portInt,
		Username: username,
		Password: password,
	}

	// Aquí podrías agregar lógica para validar la configuración del email si es necesario

	return emailCfg
}