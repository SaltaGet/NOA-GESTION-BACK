package schemas

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"database/sql"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

const (
	Create = "crear"
	Update = "actualizar"
	Delete = "eliminar"
	Read   = "consultar"
)

func extractPostgresField(pgErr *pgconn.PgError) string {
	var rawField string

	switch pgErr.Code {
	case "23505": // Único duplicado
		re := regexp.MustCompile(`Key \((.*?)\)=`)
		match := re.FindStringSubmatch(pgErr.Detail)
		if len(match) > 1 {
			rawField = match[1]
		}
	case "23502": // Not Null
		rawField = pgErr.ColumnName
	case "22001": // Truncamiento (Varchar)
		// Postgres no siempre da el nombre de columna aquí, pero si lo hace:
		rawField = pgErr.ColumnName
	}

	if rawField != "" {
		return rawField
	}
	return "un campo"
}

func isConnectionError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "io timeout") ||
		strings.Contains(msg, "network is unreachable") ||
		strings.Contains(msg, "failed to connect")
}

// HandlerErrorDB centraliza el manejo de errores de DB para Postgres
func HandlerErrorDB(err error, entity string, action string) error {
	var pgErr *pgconn.PgError

	// Mensajes dinámicos según la acción
	msgPrefix := fmt.Sprintf("No se pudo %s el/la %s", action, entity)

	// 1. Errores de GORM (Agnósticos)
	if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
		return ErrorResponse(404, fmt.Sprintf("%s: el registro no existe.", msgPrefix), err)
	}

	// 2. Errores específicos de PostgreSQL
	if errors.As(err, &pgErr) {
		field := extractPostgresField(pgErr)

		switch pgErr.Code {
		case "23505": // unique_violation
			return ErrorResponse(400, fmt.Sprintf("No se pudo %s: ya existe un/a %s con ese %s.", action, entity, field), err)

		case "23503": // foreign_key_violation
			if strings.Contains(pgErr.Message, "is still referenced") {
				return ErrorResponse(409, fmt.Sprintf("No se pudo eliminar: este/a %s tiene información asociada.", entity), err)
			}
			return ErrorResponse(400, fmt.Sprintf("Error de relación: el dato asociado para %s no es válido.", entity), err)

		case "23502": // not_null_violation
			return ErrorResponse(400, fmt.Sprintf("El campo %s es obligatorio para %s.", field, entity), err)

		case "22001": // value_too_long
			if field == "un campo" {
				field = "uno de los campos"
			}
			return ErrorResponse(400, fmt.Sprintf("El valor es demasiado largo en %s.", field), err)

		case "22P02": // invalid_type
			return ErrorResponse(400, fmt.Sprintf("El formato de un dato para %s es incorrecto.", entity), err)

		case "40P01", "40001": // Deadlock / Serialization
			return ErrorResponse(409, "El servidor está ocupado procesando otros cambios. Por favor, reintente.", err)

		case "42703", "42P01": // Error de esquema
			return ErrorResponse(500, "Error interno de estructura de base de datos.", err)
		}
	}

	// 3. Errores de Red
	if isConnectionError(err) {
		return ErrorResponse(503, "Estamos teniendo problemas de conexión con la base de datos.", err)
	}

	// 4. Fallback: Error Interno
	return ErrorResponse(500, fmt.Sprintf("Ocurrió un error inesperado al %s.", action), err)
}
