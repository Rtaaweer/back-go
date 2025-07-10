package schemas

import (
	"encoding/json"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// Esquemas JSON para validación de respuestas
var (
	// Esquema para respuesta de registro
	RegisterResponseSchema = `{
		"type": "object",
		"properties": {
			"user": {
				"type": "object",
				"properties": {
					"id_usuario": {"type": "integer"},
					"nombre": {"type": "string", "minLength": 1},
					"email": {"type": "string", "format": "email"},
					"tipo": {"type": "string", "enum": ["paciente", "medico", "enfermera", "admin"]},
					"mfa_enabled": {"type": "boolean"},
					"created_at": {"type": "string", "format": "date-time"},
					"updated_at": {"type": "string", "format": "date-time"}
				},
				"required": ["id_usuario", "nombre", "email", "tipo", "mfa_enabled", "created_at", "updated_at"]
			},
			"mfa_secret": {"type": "string", "minLength": 1},
			"qr_code_url": {"type": "string", "format": "uri"},
			"message": {"type": "string", "minLength": 1}
		},
		"required": ["user", "mfa_secret", "qr_code_url", "message"]
	}`

	// Esquema para respuesta de login exitoso
	LoginSuccessResponseSchema = `{
		"type": "object",
		"properties": {
			"access_token": {"type": "string", "minLength": 1},
			"refresh_token": {"type": "string", "minLength": 1},
			"user": {
				"type": "object",
				"properties": {
					"id_usuario": {"type": "integer"},
					"nombre": {"type": "string", "minLength": 1},
					"email": {"type": "string", "format": "email"},
					"tipo": {"type": "string", "enum": ["paciente", "medico", "enfermera", "admin"]},
					"mfa_enabled": {"type": "boolean"},
					"created_at": {"type": "string", "format": "date-time"},
					"updated_at": {"type": "string", "format": "date-time"}
				},
				"required": ["id_usuario", "nombre", "email", "tipo", "mfa_enabled", "created_at", "updated_at"]
			},
			"requires_mfa": {"type": "boolean"}
		},
		"required": ["access_token", "refresh_token", "user"]
	}`

	// Esquema para respuesta de MFA requerido
	LoginMFARequiredSchema = `{
		"type": "object",
		"properties": {
			"requires_mfa": {"type": "boolean", "const": true}
		},
		"required": ["requires_mfa"]
	}`

	// Esquema para respuesta de refresh token
	RefreshTokenResponseSchema = `{
		"type": "object",
		"properties": {
			"access_token": {"type": "string", "minLength": 1}
		},
		"required": ["access_token"]
	}`

	// Esquema para respuestas de error
	ErrorResponseSchema = `{
		"type": "object",
		"properties": {
			"error": {"type": "string", "minLength": 1},
			"code": {"type": "string"},
			"details": {"type": "object"}
		},
		"required": ["error"]
	}`

	// Esquema para respuesta de habilitación MFA
	EnableMFAResponseSchema = `{
		"type": "object",
		"properties": {
			"secret": {"type": "string", "minLength": 1},
			"qr_code_url": {"type": "string", "format": "uri"},
			"backup_codes": {
				"type": "array",
				"items": {"type": "string", "minLength": 1},
				"minItems": 1
			}
		},
		"required": ["secret", "qr_code_url", "backup_codes"]
	}`
)

// ValidateResponse valida una respuesta JSON contra un esquema
func ValidateResponse(responseData interface{}, schemaString string) error {
	// Convertir la respuesta a JSON
	responseJSON, err := json.Marshal(responseData)
	if err != nil {
		return fmt.Errorf("error al serializar respuesta: %v", err)
	}

	// Crear loaders para el esquema y los datos
	schemaLoader := gojsonschema.NewStringLoader(schemaString)
	documentLoader := gojsonschema.NewBytesLoader(responseJSON)

	// Validar
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("error en validación: %v", err)
	}

	// Verificar si es válido
	if !result.Valid() {
		errorMessages := make([]string, len(result.Errors()))
		for i, desc := range result.Errors() {
			errorMessages[i] = desc.String()
		}
		return fmt.Errorf("respuesta no válida: %v", errorMessages)
	}

	return nil
}

// ValidateRegisterResponse valida respuesta de registro
func ValidateRegisterResponse(response interface{}) error {
	return ValidateResponse(response, RegisterResponseSchema)
}

// ValidateLoginSuccessResponse valida respuesta de login exitoso
func ValidateLoginSuccessResponse(response interface{}) error {
	return ValidateResponse(response, LoginSuccessResponseSchema)
}

// ValidateLoginMFAResponse valida respuesta de MFA requerido
func ValidateLoginMFAResponse(response interface{}) error {
	return ValidateResponse(response, LoginMFARequiredSchema)
}

// ValidateRefreshTokenResponse valida respuesta de refresh token
func ValidateRefreshTokenResponse(response interface{}) error {
	return ValidateResponse(response, RefreshTokenResponseSchema)
}

// ValidateErrorResponse valida respuesta de error
func ValidateErrorResponse(response interface{}) error {
	return ValidateResponse(response, ErrorResponseSchema)
}

// ValidateEnableMFAResponse valida respuesta de habilitación MFA
func ValidateEnableMFAResponse(response interface{}) error {
	return ValidateResponse(response, EnableMFAResponseSchema)
}
