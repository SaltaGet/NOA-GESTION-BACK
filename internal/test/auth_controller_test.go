package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/test/services"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService es un mock para la interfaz ports.AuhtService

func TestAuthLogin(t *testing.T) {
	// Crear una instancia del mock del servicio
	mockAuthService := new(services.MockAuthService)

	// Crear el controlador con el servicio mockeado
	// Nota: Hay un typo en tu struct, debería ser AuthService. Usamos AuhtService para que coincida.
	authController := controllers.AuthController{
		AuthService: mockAuthService,
	}

	// Crear una nueva app de Fiber para el test
	app := fiber.New()
	app.Post("/api/v1/auth/login", authController.AuthLogin)

	t.Run("Login exitoso", func(t *testing.T) {
		// Configurar el mock para que devuelva un token y sin error
		mockAuthService.On("AuthLogin", "test@user", "password123").Return("fake-jwt-token", nil).Once()

		// Crear el cuerpo de la petición
		loginCredentials := schemas.AuthLogin{
			Username: "test@user",
			Password: "password123",
		}
		body, _ := json.Marshal(loginCredentials)

		// Crear la petición HTTP
		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		// Ejecutar la petición
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Verificar el cuerpo de la respuesta
		var responseBody schemas.Response
		respBodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBodyBytes, &responseBody)

		assert.True(t, responseBody.Status)
		assert.Equal(t, "Login exitoso", responseBody.Message)
		assert.Nil(t, responseBody.Body)

		// Verificar la cookie
		cookies := resp.Cookies()
		assert.Len(t, cookies, 1)
		assert.Equal(t, "access_token", cookies[0].Name)
		assert.Equal(t, "fake-jwt-token", cookies[0].Value)

		// Asegurarse que todas las expectativas del mock se cumplieron
		mockAuthService.AssertExpectations(t)
	})

	t.Run("Login fallido - credenciales incorrectas", func(t *testing.T) {
		// Configurar el mock para que devuelva un error
		expectedError := schemas.ErrorResponse(fiber.StatusUnauthorized, "Credenciales inválidas", errors.New("invalid credentials"))
		mockAuthService.On("AuthLogin", "wrong@user", "wrongpass").Return("", expectedError).Once()

		loginCredentials := schemas.AuthLogin{
			Username: "wrong@user",
			Password: "wrongpass",
		}
		body, _ := json.Marshal(loginCredentials)

		req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)

		var responseBody schemas.Response
		respBodyBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBodyBytes, &responseBody)

		assert.False(t, responseBody.Status)
		assert.Equal(t, "Credenciales inválidas", responseBody.Message)

		mockAuthService.AssertExpectations(t)
	})
}

func TestAuthPointSale(t *testing.T) {
	mockAuthService := new(services.MockAuthService)

	authController := controllers.AuthController{
		AuthService: mockAuthService,
	}

	app := fiber.New()

	// Simular middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &schemas.AuthenticatedUser{
			ID:       1,
			Username: "test@user",
		})
		return c.Next()
	})

	app.Post("/api/v1/auth/login_point_sale/:point_sale_id", authController.AuthPointSale)

	t.Run("Login PointSale exitoso", func(t *testing.T) {
		mockAuthService.
			On("AuthPointSale", mock.AnythingOfType("*schemas.AuthenticatedUser"), int64(5)).
			Return("token-ps", nil).Once()

		req := httptest.NewRequest("POST", "/api/v1/auth/login_point_sale/5", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var jsonResp schemas.Response
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &jsonResp)

		assert.True(t, jsonResp.Status)
		assert.Equal(t, "Login a Punto de venta exitoso, token enviado en cookie", jsonResp.Message)

		cookies := resp.Cookies()
		assert.Equal(t, "token-ps", cookies[0].Value)

		mockAuthService.AssertExpectations(t)
	})
}

func TestLogoutPointSale(t *testing.T) {
	mockAuthService := new(services.MockAuthService)

	authController := controllers.AuthController{
		AuthService: mockAuthService,
	}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &schemas.AuthenticatedUser{
			ID: 1,
		})
		return c.Next()
	})

	app.Post("/api/v1/auth/logout_point_sale", authController.LogoutPointSale)

	t.Run("Logout PointSale exitoso", func(t *testing.T) {

		mockAuthService.
			On("LogoutPointSale", mock.AnythingOfType("*schemas.AuthenticatedUser")).
			Return("newToken", nil).Once()

		req := httptest.NewRequest("POST", "/api/v1/auth/logout_point_sale", nil)

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var jsonResp schemas.Response
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &jsonResp)

		assert.True(t, jsonResp.Status)
		assert.Equal(t, "Logout de Punto de venta exitoso, token enviado en cookie", jsonResp.Message)

		cookies := resp.Cookies()
		assert.Equal(t, "newToken", cookies[0].Value)

		mockAuthService.AssertExpectations(t)
	})
}

func TestLogout(t *testing.T) {
	authController := controllers.AuthController{}

	app := fiber.New()
	app.Post("/api/v1/auth/logout", authController.Logout)

	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var jsonResp schemas.Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &jsonResp)

	assert.True(t, jsonResp.Status)
	assert.Equal(t, "Logout exitoso", jsonResp.Message)

	cookies := resp.Cookies()
	assert.Equal(t, "", cookies[0].Value)
}

func TestCurrentUser(t *testing.T) {
	authController := controllers.AuthController{}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("user", &schemas.AuthenticatedUser{
			ID:       2,
			Username: "john@doe",
		})
		return c.Next()
	})

	app.Get("/api/v1/auth/current_user", authController.CurrentUser)

	req := httptest.NewRequest("GET", "/api/v1/auth/current_user", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var jsonResp schemas.Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &jsonResp)

	user := jsonResp.Body.(map[string]interface{})

	assert.Equal(t, "john@doe", user["username"])
}

func TestCurrentPlan(t *testing.T) {
	authController := controllers.AuthController{}

	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("current_plan", &schemas.PlanResponseDTO{
			ID: 1,
			Name: "Básico",
			PriceMounthly: 25,
			PriceYearly: 250,
			Description: "es lo que hay",
			Features: "emmm no hay nada aqui",
			AmountPointSale: 1,
			AmountMember: 5,
		})
		return c.Next()
	})

	app.Get("/api/v1/auth/current_plan", authController.CurrentPlan)

	req := httptest.NewRequest("GET", "/api/v1/auth/current_plan", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var jsonResp schemas.Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &jsonResp)

	plan := jsonResp.Body.(map[string]interface{})
	assert.Equal(t, "Básico", plan["name"])
}

func TestCurrentTenant(t *testing.T) {
	authController := controllers.AuthController{}

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("current_tenant", &schemas.TenantResponse{
			ID:   99,
			Name: "Tenant Example",
		})
		return c.Next()
	})

	app.Get("/api/v1/auth/current_tenant", authController.CurrentTenant)

	req := httptest.NewRequest("GET", "/api/v1/auth/current_tenant", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var jsonResp schemas.Response
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &jsonResp)

	tenant := jsonResp.Body.(map[string]interface{})
	assert.Equal(t, float64(99), tenant["id"])
}

func TestAuthLoginAdmin(t *testing.T) {
	mockAuthService := new(services.MockAuthService)

	authController := controllers.AuthController{
		AuthService: mockAuthService,
	}

	app := fiber.New()
	app.Post("/api/v1/auth/login_admin", authController.AuthLoginAdmin)

	t.Run("Login Admin OK", func(t *testing.T) {
		mockAuthService.
			On("AuthLoginAdmin", "admin@test", "pass123").
			Return("admin-token", nil).Once()

		body, _ := json.Marshal(schemas.AuthLoginAdmin{
			Username: "admin@test",
			Password: "pass123",
		})

		req := httptest.NewRequest("POST", "/api/v1/auth/login_admin", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)

		var result schemas.Response
		respBytes, _ := io.ReadAll(resp.Body)
		json.Unmarshal(respBytes, &result)

		assert.True(t, result.Status)
		assert.Equal(t, "Login exitoso", result.Message)

		cookies := resp.Cookies()
		assert.Equal(t, "access_token_admin", cookies[0].Name)
		assert.Equal(t, "admin-token", cookies[0].Value)

		mockAuthService.AssertExpectations(t)
	})
}

func TestAuthForgotPassword(t *testing.T) {
	mockAuthService := new(services.MockAuthService)
	authController := controllers.AuthController{AuthService: mockAuthService}

	app := fiber.New()
	app.Post("/api/v1/auth/forgot_password", authController.AuthForgotPassword)

	t.Run("Forgot password OK", func(t *testing.T) {
		mockAuthService.
			On("AuthForgotPassword", mock.AnythingOfType("*schemas.AuthForgotPassword")).
			Return(nil).Once()

		body, _ := json.Marshal(schemas.AuthForgotPassword{
			Username:         "user",
			TenantIdentifier: "tenant",
		})

		req := httptest.NewRequest("POST", "/api/v1/auth/forgot_password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

func TestAuthResetPassword(t *testing.T) {
	mockAuthService := new(services.MockAuthService)
	authController := controllers.AuthController{AuthService: mockAuthService}

	app := fiber.New()
	app.Post("/api/v1/auth/reset_password", authController.AuthResetPassword)

	t.Run("Reset password OK", func(t *testing.T) {
		mockAuthService.
			On("AuthResetPassword", mock.AnythingOfType("*schemas.AuthResetPassword")).
			Return(nil).Once()

		body, _ := json.Marshal(schemas.AuthResetPassword{
			Token:       "abcd1234",
			NewPassword: "Qwer1234*",
			ConfirmPass: "Qwer1234*",
		})

		req := httptest.NewRequest("POST", "/api/v1/auth/reset_password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})
}

