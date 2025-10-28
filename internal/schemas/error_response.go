package schemas

import (
	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
	"github.com/gofiber/fiber/v2"
)

type ErrorStruc struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *ErrorStruc) Error() string {
	return e.Message
}

func ErrorResponse(code int, message string, err error) *ErrorStruc {
	return &ErrorStruc{
		StatusCode: code,
		Message:    message,
		Err:        err,
	}
}

func HandleError(ctx *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	if errResp, ok := err.(*ErrorStruc); ok {
		logging.ERROR("Error: %s", errResp.Err.Error())
		return ctx.Status(errResp.StatusCode).JSON(Response{
			Status:  false,
			Body:    nil,
			Message: errResp.Message,
		})
	}

	logging.ERROR("Error: %s", err.Error())
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Status:  false,
		Body:    nil,
		Message: "Error interno",
	})
}