package schemas

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgconn" // Driver de Postgres
	"gorm.io/gorm"
)

func extractPostgresField(pgErr *pgconn.PgError) string {
	// 1. Prioridad: Si Postgres nos da la columna directamente (común en Not Null)
	if pgErr.ColumnName != "" {
		return strings.ReplaceAll(pgErr.ColumnName, "_", " ")
	}

	// 2. Violación de Único (Duplicados)
	if pgErr.Code == "23505" {
		re := regexp.MustCompile(`Key \((.*?)\)=`)
		match := re.FindStringSubmatch(pgErr.Detail)
		if len(match) > 1 {
			return strings.ReplaceAll(match[1], "_", " ")
		}
	}

	if pgErr.Code == "22001" {
		return "el largo de uno de los campos"
	}

	return "campo"
}

func isConnectionError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "io timeout") ||
		strings.Contains(msg, "network is unreachable") ||
		strings.Contains(msg, "failed to connect")
}

func HandlerErrorGorm(err error, entity string) error {
	var pgErr *pgconn.PgError

	// Mensajes adicionales para robustez
	const (
		MsgNotFound     = "El registro de %s no existe."
		MsgDuplicated   = "Ya existe un/a %s con este dato: %s."
		MsgForeignKey   = "No se puede eliminar: este/a %s tiene información asociada."
		MsgRelationErr  = "Error de relación: el dato asociado para %s no es válido."
		MsgNotNull      = "El campo %s es obligatorio."
		MsgInvalidType  = "El formato de un dato para %s es incorrecto (ej. texto en campo numérico)."
		MsgDataTooLong  = "Uno de los campos de %s supera el largo permitido (ej. más de 30 caracteres)."
		MsgDeadlock     = "El servidor está muy ocupado procesando esta información. Reintente."
		MsgInternal     = "Ocurrió un error inesperado."
		MsgDBConnection = "Problema de conexión con la base de datos."
	)

	// 1. Errores de GORM
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrorResponse(404, fmt.Sprintf(MsgNotFound, entity), err)
	}

	// 2. Errores de PostgreSQL (SQLSTATE)
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		// --- INTEGRIDAD ---
		case "23505": // unique_violation
			field := extractPostgresField(pgErr)
			return ErrorResponse(400, fmt.Sprintf(MsgDuplicated, entity, field), err)
		case "23503": // foreign_key_violation
			if strings.Contains(pgErr.Message, "is still referenced") {
				return ErrorResponse(409, fmt.Sprintf(MsgForeignKey, entity), err)
			}
			return ErrorResponse(400, fmt.Sprintf(MsgRelationErr, entity), err)
		case "23502": // not_null_violation
			return ErrorResponse(400, fmt.Sprintf(MsgNotNull, pgErr.ColumnName), err)

		case "22001": // string_data_right_truncation
			field := extractPostgresField(pgErr)
			return ErrorResponse(400, fmt.Sprintf("El valor es demasiado largo en %s.", field), err)
		case "22P02": // invalid_text_representation (Ej: enviar "letras" a un campo UUID o INT)
			return ErrorResponse(400, fmt.Sprintf(MsgInvalidType, entity), err)

		case "40P01": // deadlock_detected
			return ErrorResponse(409, MsgDeadlock, err)
		case "40001": // serialization_failure
			return ErrorResponse(409, "Conflicto de actualización. Por favor, reintente.", err)

		case "42703": // undefined_column
			return ErrorResponse(500, "Error de sistema: Columna inexistente.", err)
		case "42P01": // undefined_table
			return ErrorResponse(500, "Error de sistema: Tabla inexistente.", err)
		}
	}

	// 3. Errores de Conexión
	if isConnectionError(err) {
		return ErrorResponse(503, MsgDBConnection, err)
	}

	return ErrorResponse(500, MsgInternal, err)
}
