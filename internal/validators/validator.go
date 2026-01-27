package validators

import (
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
)

func IdValidate(id string) (int64, error) {
	idint, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, schemas.ErrorResponse(422, "el id ser un número", err)
	}
	if idint <= 0 {
		return 0, schemas.ErrorResponse(422, "el id debe ser mayor a 0", fmt.Errorf("el id debe de ser mayor a 0"))
	}

	return idint, nil
}

func IntValidate(value string) (int, error) {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, schemas.ErrorResponse(422, "el valor debe ser un número entero", err)
	}
	if intValue < 0 {
		return 0, schemas.ErrorResponse(422, "el valor debe ser mayor o igual a 0", fmt.Errorf("el valor debe ser mayor o igual a 0"))
	}

	return intValue, nil
}

func IsValidUUIDv4(id string) bool {
    parsed, err := uuid.Parse(id)
    if err != nil {
        return false
    }

    return parsed.Version() == 4
}