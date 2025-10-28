package main

import (
	"log"
	"os"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/jobs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
	
	_ "github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/docs"
)

// @title						APP NOA Gestion
// @version					1.0
// @description				This is a api to app noa gestion
// @termsOfService				http://swagger.io/terms/
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Type "Bearer" followed by a space and the JWT token. Example: "Bearer eyJhbGciOiJIUz..."
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error al cargar el archivo .env: %v", err)
	}

	local := os.Getenv("LOCAL")
	if local == "true" {
		if err := jobs.GenerateSwagger(); err != nil {
			log.Fatalf("Error ejecutando swag init: %v", err)
		}
	}

	app := fiber.New()
	
	app.Get("/api/swagger/*", swagger.HandlerDefault)

	log.Fatal(app.Listen(":3000"))
}
