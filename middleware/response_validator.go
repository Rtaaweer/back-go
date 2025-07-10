package middleware

import (
	"encoding/json"
	"hospital-system/schemas"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ResponseValidator middleware para validar respuestas JSON (solo logging)
func ResponseValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continuar con el handler
		err := c.Next()
		if err != nil {
			return err
		}

		// Solo validar respuestas JSON
		contentType := c.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return nil
		}

		// Obtener el cuerpo de la respuesta
		responseBody := c.Response().Body()
		if len(responseBody) == 0 {
			return nil
		}

		// Parsear JSON
		var responseData interface{}
		if err := json.Unmarshal(responseBody, &responseData); err != nil {
			log.Printf("W - Error al parsear respuesta JSON: %v", err)
			return nil
		}

		// Determinar qué esquema usar basado en la ruta y código de estado
		path := c.Path()
		statusCode := c.Response().StatusCode()

		var validationErr error

		// Validar según la ruta y código de estado
		switch {
		case path == "/api/v1/auth/register" && statusCode == 201:
			validationErr = schemas.ValidateRegisterResponse(responseData)
		case path == "/api/v1/auth/login" && statusCode == 200:
			// Verificar si requiere MFA o es login exitoso
			if responseMap, ok := responseData.(map[string]interface{}); ok {
				if requiresMFA, exists := responseMap["requires_mfa"]; exists && requiresMFA == true {
					validationErr = schemas.ValidateLoginMFAResponse(responseData)
				} else {
					validationErr = schemas.ValidateLoginSuccessResponse(responseData)
				}
			}
		case path == "/api/v1/auth/refresh" && statusCode == 200:
			validationErr = schemas.ValidateRefreshTokenResponse(responseData)
		case strings.HasPrefix(path, "/api/v1/auth/mfa/enable") && statusCode == 200:
			validationErr = schemas.ValidateEnableMFAResponse(responseData)
		case statusCode >= 400:
			validationErr = schemas.ValidateErrorResponse(responseData)
		}

		// Log errores de validación pero no interrumpir la respuesta
		if validationErr != nil {
			log.Printf("F - Respuesta inválida para %s [%d]: %v", path, statusCode, validationErr)
		} else {
			log.Printf("S - Respuesta válida para %s [%d]", path, statusCode)
		}

		return nil
	}
}

// SimpleResponseValidator versión simplificada que funciona con Fiber v2
func SimpleResponseValidator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ejecutar el handler original
		err := c.Next()
		if err != nil {
			return err
		}

		// Validar solo después de que la respuesta se haya enviado
		// Esta validación es principalmente para logging y debugging
		go func() {
			// Solo validar respuestas JSON
			contentType := c.Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				return
			}

			// Obtener información de la respuesta
			path := c.Path()
			statusCode := c.Response().StatusCode()
			responseBody := c.Response().Body()

			if len(responseBody) == 0 {
				return
			}

			// Parsear JSON
			var responseData interface{}
			if err := json.Unmarshal(responseBody, &responseData); err != nil {
				log.Printf("W - Error al parsear respuesta JSON para validación: %v", err)
				return
			}

			// Validar según la ruta y código de estado
			var validationErr error
			switch {
			case path == "/api/v1/auth/register" && statusCode == 201:
				validationErr = schemas.ValidateRegisterResponse(responseData)
			case path == "/api/v1/auth/login" && statusCode == 200:
				if responseMap, ok := responseData.(map[string]interface{}); ok {
					if requiresMFA, exists := responseMap["requires_mfa"]; exists && requiresMFA == true {
						validationErr = schemas.ValidateLoginMFAResponse(responseData)
					} else {
						validationErr = schemas.ValidateLoginSuccessResponse(responseData)
					}
				}
			case path == "/api/v1/auth/refresh" && statusCode == 200:
				validationErr = schemas.ValidateRefreshTokenResponse(responseData)
			case strings.HasPrefix(path, "/api/v1/auth/mfa/enable") && statusCode == 200:
				validationErr = schemas.ValidateEnableMFAResponse(responseData)
			case statusCode >= 400:
				validationErr = schemas.ValidateErrorResponse(responseData)
			}

			// Log errores de validación
			if validationErr != nil {
				log.Printf("F - Respuesta inválida para %s [%d]: %v", path, statusCode, validationErr)
			} else {
				log.Printf("S - Respuesta válida para %s [%d]", path, statusCode)
			}
		}()

		return nil
	}
}

// ValidateResponseData función helper para validar datos de respuesta manualmente
func ValidateResponseData(path string, statusCode int, responseData interface{}) error {
	switch {
	case path == "/api/v1/auth/register" && statusCode == 201:
		return schemas.ValidateRegisterResponse(responseData)
	case path == "/api/v1/auth/login" && statusCode == 200:
		if responseMap, ok := responseData.(map[string]interface{}); ok {
			if requiresMFA, exists := responseMap["requires_mfa"]; exists && requiresMFA == true {
				return schemas.ValidateLoginMFAResponse(responseData)
			} else {
				return schemas.ValidateLoginSuccessResponse(responseData)
			}
		}
	case path == "/api/v1/auth/refresh" && statusCode == 200:
		return schemas.ValidateRefreshTokenResponse(responseData)
	case strings.HasPrefix(path, "/api/v1/auth/mfa/enable") && statusCode == 200:
		return schemas.ValidateEnableMFAResponse(responseData)
	case statusCode >= 400:
		return schemas.ValidateErrorResponse(responseData)
	}
	return nil
}
