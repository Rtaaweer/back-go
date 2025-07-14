package middleware

import (
	"encoding/json"
	"hospital-system/config"
	"os"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LoggerMiddleware middleware para registrar todas las peticiones
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Procesar la petición
		err := c.Next()

		// Calcular tiempo de respuesta
		responseTime := int(time.Since(start).Milliseconds())

		// CAPTURAR TODOS LOS VALORES ANTES DE LA GOROUTINE
		method := c.Method()
		path := c.Path()
		protocol := c.Protocol()
		statusCode := c.Response().StatusCode()
		userAgent := c.Get("User-Agent")
		ip := c.IP()
		originalURL := c.OriginalURL()

		// Obtener información del usuario autenticado (si existe)
		var email, username, role *string
		if userEmail := c.Locals("user_email"); userEmail != nil {
			emailStr := userEmail.(string)
			email = &emailStr
		}
		if userName := c.Locals("username"); userName != nil {
			usernameStr := userName.(string)
			username = &usernameStr
		}
		if userRole := c.Locals("user_role"); userRole != nil {
			roleStr := userRole.(string)
			role = &roleStr
		}

		// Preparar body como JSON string
		var bodyStr *string
		if len(c.Body()) > 0 {
			body := string(c.Body())
			bodyStr = &body
		}

		// Preparar params como JSON string
		var paramsStr *string
		if len(c.AllParams()) > 0 {
			if paramsJSON, marshalErr := json.Marshal(c.AllParams()); marshalErr == nil {
				params := string(paramsJSON)
				paramsStr = &params
			}
		}

		// Preparar query como string
		var queryStr *string
		if c.Request().URI().QueryString() != nil {
			query := string(c.Request().URI().QueryString())
			queryStr = &query
		}

		// Determinar log level basado en status code
		logLevel := "info"
		if statusCode >= 400 && statusCode < 500 {
			logLevel = "warning"
		} else if statusCode >= 500 {
			logLevel = "error"
		}

		// Obtener información del sistema
		hostname, _ := os.Hostname()
		environment := os.Getenv("ENVIRONMENT")
		if environment == "" {
			environment = "development"
		}
		nodeVersion := runtime.Version()
		pid := os.Getpid()

		// Insertar log en la base de datos usando valores capturados

go func() {
    db := config.GetDB()
    if db != nil {
        _, dbErr := db.Exec(`
            INSERT INTO Logs (method, path, protocol, status_code, response_time, user_agent, 
                             ip, hostname, body, params, query, email, username, role, 
                             log_level, environment, node_version, pid, url, timestamp) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)`,
            method,
            path,
            protocol,
            statusCode,
            responseTime,
            userAgent,
            ip,
            hostname,
            bodyStr,
            paramsStr,
            queryStr,
            email,
            username,
            role,
            logLevel,
            environment,
            nodeVersion,
            pid,
            originalURL,
            time.Now(),
        )
        if dbErr != nil {
            println("Error logging to database:", dbErr.Error())
        }
    }
}()


		return err
	}
}
