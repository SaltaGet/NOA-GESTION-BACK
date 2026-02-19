package schemas

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Constantes para estandarizar las acciones en los controladores
const (
	Create = "crear"
	Update = "actualizar"
	Delete = "eliminar"
	Read   = "consultar"
)

// translateFieldName convierte nombres técnicos (cuit, first_name) en legibles (CUIT, Nombre)
// func translateFieldName(field string) string {
// 	translations := map[string]string{
// 		"id":                       "identificador",
// 		"first_name":               "nombre",
// 		"last_name":                "apellido",
// 		"cuit":                     "CUIT/CUIL",
// 		"email":                    "correo electrónico",
// 		"phone":                    "teléfono",
// 		"address":                  "dirección",
// 		"responsibility_front_iva": "condición de IVA",
// 		"tenant_id":                "ID de comercio",
// 		"password":                 "contraseña",
// 		"price":                    "precio",
// 		"stock":                    "existencias",
// 	}

// 	if val, ok := translations[strings.ToLower(field)]; ok {
// 		return val
// 	}
// 	// Si no está en el mapa, limpia el snake_case: "product_code" -> "product code"
// 	return strings.ReplaceAll(field, "_", " ")
// }

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

// HandlerErrorGorm centraliza el manejo de errores de DB para Postgres
func HandlerErrorGorm(err error, entity string, action string) error {
	var pgErr *pgconn.PgError

	// Mensajes dinámicos según la acción
	msgPrefix := fmt.Sprintf("No se pudo %s el/la %s", action, entity)

	// 1. Errores de GORM (Agnósticos)
	if errors.Is(err, gorm.ErrRecordNotFound) {
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
			if field == "un campo" { field = "uno de los campos" }
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

// package schemas

// import (
// 	"errors"
// 	"fmt"
// 	"regexp"
// 	"strings"

// 	"github.com/jackc/pgx/v5/pgconn" // Driver de Postgres
// 	"gorm.io/gorm"
// )

// func extractPostgresField(pgErr *pgconn.PgError) string {
// 	// 1. Prioridad: Si Postgres nos da la columna directamente (común en Not Null)
// 	if pgErr.ColumnName != "" {
// 		return strings.ReplaceAll(pgErr.ColumnName, "_", " ")
// 	}

// 	// 2. Violación de Único (Duplicados)
// 	if pgErr.Code == "23505" {
// 		re := regexp.MustCompile(`Key \((.*?)\)=`)
// 		match := re.FindStringSubmatch(pgErr.Detail)
// 		if len(match) > 1 {
// 			return strings.ReplaceAll(match[1], "_", " ")
// 		}
// 	}

// 	if pgErr.Code == "22001" {
// 		return "el largo de uno de los campos"
// 	}

// 	return "campo"
// }

// func isConnectionError(err error) bool {
// 	msg := strings.ToLower(err.Error())
// 	return strings.Contains(msg, "connection refused") ||
// 		strings.Contains(msg, "io timeout") ||
// 		strings.Contains(msg, "network is unreachable") ||
// 		strings.Contains(msg, "failed to connect")
// }

// func HandlerErrorGorm(err error, entity string) error {
// 	var pgErr *pgconn.PgError

// 	// Mensajes adicionales para robustez
// 	const (
// 		MsgNotFound     = "El registro de %s no existe."
// 		MsgDuplicated   = "Ya existe un/a %s con este dato: %s."
// 		MsgForeignKey   = "No se puede eliminar: este/a %s tiene información asociada."
// 		MsgRelationErr  = "Error de relación: el dato asociado para %s no es válido."
// 		MsgNotNull      = "El campo %s es obligatorio."
// 		MsgInvalidType  = "El formato de un dato para %s es incorrecto (ej. texto en campo numérico)."
// 		MsgDataTooLong  = "Uno de los campos de %s supera el largo permitido (ej. más de 30 caracteres)."
// 		MsgDeadlock     = "El servidor está muy ocupado procesando esta información. Reintente."
// 		MsgInternal     = "Ocurrió un error inesperado."
// 		MsgDBConnection = "Problema de conexión con la base de datos."
// 	)

// 	// 1. Errores de GORM
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		return ErrorResponse(404, fmt.Sprintf(MsgNotFound, entity), err)
// 	}

// 	// 2. Errores de PostgreSQL (SQLSTATE)
// 	if errors.As(err, &pgErr) {
// 		switch pgErr.Code {
// 		// --- INTEGRIDAD ---
// 		case "23505": // unique_violation
// 			field := extractPostgresField(pgErr)
// 			return ErrorResponse(400, fmt.Sprintf(MsgDuplicated, entity, field), err)
// 		case "23503": // foreign_key_violation
// 			if strings.Contains(pgErr.Message, "is still referenced") {
// 				return ErrorResponse(409, fmt.Sprintf(MsgForeignKey, entity), err)
// 			}
// 			return ErrorResponse(400, fmt.Sprintf(MsgRelationErr, entity), err)
// 		case "23502": // not_null_violation
// 			return ErrorResponse(400, fmt.Sprintf(MsgNotNull, pgErr.ColumnName), err)

// 		case "22001": // string_data_right_truncation
// 			field := extractPostgresField(pgErr)
// 			return ErrorResponse(400, fmt.Sprintf("El valor es demasiado largo en %s.", field), err)
// 		case "22P02": // invalid_text_representation (Ej: enviar "letras" a un campo UUID o INT)
// 			return ErrorResponse(400, fmt.Sprintf(MsgInvalidType, entity), err)

// 		case "40P01": // deadlock_detected
// 			return ErrorResponse(409, MsgDeadlock, err)
// 		case "40001": // serialization_failure
// 			return ErrorResponse(409, "Conflicto de actualización. Por favor, reintente.", err)

// 		case "42703": // undefined_column
// 			return ErrorResponse(500, "Error de sistema: Columna inexistente.", err)
// 		case "42P01": // undefined_table
// 			return ErrorResponse(500, "Error de sistema: Tabla inexistente.", err)
// 		}
// 	}

// 	// 3. Errores de Conexión
// 	if isConnectionError(err) {
// 		return ErrorResponse(503, MsgDBConnection, err)
// 	}

// 	return ErrorResponse(500, MsgInternal, err)
// }
