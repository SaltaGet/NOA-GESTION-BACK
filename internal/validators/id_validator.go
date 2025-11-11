package validators

import (
	"fmt"
	"strconv"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
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