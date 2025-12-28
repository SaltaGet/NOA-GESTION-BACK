package test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"io"
// 	"net/http/httptest"
// 	"testing"
// 	"time"

// 	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/controllers"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/test/services"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// // =========================================================
// // setupApp — Helper para configurar Fiber con middleware fake
// // =========================================================
// func setupApp(ctrl controllers.CashRegisterController) *fiber.App {
// 	app := fiber.New()

// 	// Middleware para setear user y pointSale
// 	app.Use(func(c *fiber.Ctx) error {
// 		c.Locals("user", &schemas.AuthenticatedUser{ID: 10})
// 		c.Locals("point_sale_id", int64(1))
// 		return c.Next()
// 	})

// 	// Rutas
// 	app.Post("/cash_register/open/:point_sale_id", ctrl.CashRegisterOpen)
// 	app.Post("/cash_register/close/:point_sale_id", ctrl.CashRegisterClose)
// 	app.Get("/cash_register/inform", ctrl.CashRegiterInform)
// 	app.Get("/cash_register/exist-open/:point_sale_id", ctrl.CashRegisterExistOpen)
// 	app.Get("/cash_register/:id", ctrl.CashRegisterGetByID)

// 	return app
// }

// // -----------------------------------------------------------
// // TEST: CashRegisterOpen
// // -----------------------------------------------------------
// func TestCashRegisterOpenController(t *testing.T) {
// 	mock := new(services.MockCashRegisterService)
// 	ctrl := controllers.CashRegisterController{CashRegisterService: mock}
// 	app := setupApp(ctrl)

// 	body := schemas.CashRegisterOpen{OpenAmount: 1000}

// 	mock.
// 		On("CashRegisterOpen", int64(1), int64(10), body).
// 		Return(nil)

// 	jsonBody, _ := json.Marshal(body)

// 	req := httptest.NewRequest("POST", "/cash_register/open/1", bytes.NewReader(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := app.Test(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)

// 	mock.AssertExpectations(t)
// }

// // -----------------------------------------------------------
// // TEST: CashRegisterGetByID
// // -----------------------------------------------------------
// func TestCashRegisterGetByIDController(t *testing.T) {
// 	mock := new(services.MockCashRegisterService)
// 	ctrl := controllers.CashRegisterController{CashRegisterService: mock}
// 	app := setupApp(ctrl)

// 	expected := &schemas.CashRegisterFullResponse{
// 		ID:         20,
// 		OpenAmount: 1500,
// 		IsClose:    false,
// 		CreatedAt:  time.Now(),
// 		MemberOpen: schemas.MemberSimpleDTO{
// 			ID:        1,
// 			FirstName: "test",
// 			LastName:  "test",
// 			Username:  "test",
// 		},
// 	}

// 	mock.
// 		On("CashRegisterGetByID", int64(1), int64(20)).
// 		Return(expected, nil)

// 	req := httptest.NewRequest("GET", "/cash_register/20", nil)
// 	resp, err := app.Test(req)

// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)

// 	var result schemas.CashRegisterFullResponse
// 	body, _ := io.ReadAll(resp.Body)
// 	json.Unmarshal(body, &result)

// 	assert.Equal(t, expected.ID, result.ID)
// 	mock.AssertExpectations(t)
// }

// // -----------------------------------------------------------
// // TEST: CashRegisterClose
// // -----------------------------------------------------------
// func TestCashRegisterCloseController(t *testing.T) {
// 	mock := new(services.MockCashRegisterService)
// 	ctrl := controllers.CashRegisterController{CashRegisterService: mock}
// 	app := setupApp(ctrl)

// 	body := schemas.CashRegisterClose{CloseAmount: 800}

// 	mock.
// 		On("CashRegisterClose", int64(1), int64(10), body).
// 		Return(nil)

// 	jsonBody, _ := json.Marshal(body)

// 	req := httptest.NewRequest("POST", "/cash_register/close/1", bytes.NewReader(jsonBody))
// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := app.Test(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)

// 	mock.AssertExpectations(t)
// }

// // -----------------------------------------------------------
// // TEST: CashRegisterInform
// // -----------------------------------------------------------
// func TestCashRegisterInformController(t *testing.T) {
//     mock := new(services.MockCashRegisterService)
//     ctrl := controllers.CashRegisterController{CashRegisterService: mock}

//     app := fiber.New()

//     app.Use(func(c *fiber.Ctx) error {
//         c.Locals("point_sale_id", int64(1))   // 👈 nombre correcto
//         c.Locals("user", &schemas.AuthenticatedUser{ID: 10})
//         return c.Next()
//     })

//     from := time.Now().Add(-24 * time.Hour)
//     to := time.Now()

//     expected := []*schemas.CashRegisterInformResponse{
//         {
//             ID:         1,
//             OpenAmount: 1000,
//             MemberOpen: schemas.MemberSimpleDTO{
//                 ID:        1,
//                 FirstName: "test",
//                 LastName:  "test",
//                 Username:  "test",
//             },
//             HourOpen:  time.Now(),
//             IsClose:   false,
//             CreatedAt: time.Now(),
//         },
//     }

//     // 💥 Importante: usar mock.Anything para fechas (no van a coincidir por nanosegundos)
//     mock.
//         On("CashRegisterInform", int64(1), int64(10), mock.Anything, mock.Anything).
//         Return(expected, nil)

//     app.Post("/cash_register/inform", ctrl.CashRegiterInform)

//     // Body JSON correcto según el controlador
//     body := map[string]string{
//         "from": from.Format(time.RFC3339),
//         "to":   to.Format(time.RFC3339),
//     }
//     bodyJSON, _ := json.Marshal(body)

//     req := httptest.NewRequest(
//         "POST",
//         "/cash_register/inform",
//         bytes.NewReader(bodyJSON),
//     )
//     req.Header.Set("Content-Type", "application/json")

//     resp, err := app.Test(req)
//     assert.NoError(t, err)
//     assert.Equal(t, 200, resp.StatusCode)

//     mock.AssertExpectations(t)
// }


// // -----------------------------------------------------------
// // TEST: CashRegisterExistOpen
// // -----------------------------------------------------------
// func TestCashRegisterExistOpenController(t *testing.T) {
// 	mock := new(services.MockCashRegisterService)
// 	ctrl := controllers.CashRegisterController{CashRegisterService: mock}
// 	app := setupApp(ctrl)

// 	mock.
// 		On("CashRegisterExistOpen", int64(1)).
// 		Return(true, nil)

// 	req := httptest.NewRequest("GET", "/cash_register/exist-open/1", nil)

// 	resp, err := app.Test(req)
// 	assert.NoError(t, err)
// 	assert.Equal(t, 200, resp.StatusCode)

// 	mock.AssertExpectations(t)
// }
