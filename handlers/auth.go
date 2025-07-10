package handlers

import (
	"database/sql"
	"fmt"
	"hospital-system/config"
	"hospital-system/models"

	// "hospital-system/schemas" <- Eliminar esta línea
	"hospital-system/utils"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

// sendValidatedResponse envía una respuesta validada
func sendValidatedResponse(c *fiber.Ctx, statusCode int, data interface{}, validator func(interface{}) error) error {
	// Validar la respuesta antes de enviarla
	if validator != nil {
		if err := validator(data); err != nil {
			log.Printf("Error de validación de respuesta: %v", err)
			// En desarrollo, podrías querer mostrar el error
			// En producción, enviar error genérico
			return c.Status(500).JSON(utils.NewResponse(500, "E99", fiber.Map{
				"error": "Error interno del servidor",
				"code":  "INTERNAL_ERROR",
			}))
		}
	}

	return c.Status(statusCode).JSON(data)
}

// sendStandardResponse envía una respuesta con el formato estandarizado
func sendStandardResponse(c *fiber.Ctx, statusCode int, intCode string, data interface{}) error {
	response := utils.NewResponse(statusCode, intCode, data)
	return c.Status(statusCode).JSON(response)
}

func Register(c *fiber.Ctx) error {
	var req models.RegisterRequest // Cambiar de LoginRequest a RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return sendStandardResponse(c, 400, utils.REGISTER_PARSE_ERROR, fiber.Map{
			"error":       "Error al parsear datos",
			"description": utils.GetCodeDescription(utils.REGISTER_PARSE_ERROR),
		})
	}

	// Validar fortaleza de contraseña
	if !utils.IsStrongPassword(req.Password) {
		return sendStandardResponse(c, 400, utils.REGISTER_WEAK_PASSWORD, fiber.Map{
			"error":       "La contraseña debe tener al menos 12 caracteres, incluir mayúsculas, minúsculas, números y símbolos",
			"description": utils.GetCodeDescription(utils.REGISTER_WEAK_PASSWORD),
		})
	}

	// Hash de la contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_HASH_ERROR, fiber.Map{
			"error":       "Error al procesar contraseña",
			"description": utils.GetCodeDescription(utils.REGISTER_HASH_ERROR),
		})
	}

	// Generar MFA
	mfaSecret, err := utils.GenerateMFASecret()
	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_MFA_ERROR, fiber.Map{
			"error":       "Error al generar secreto MFA",
			"description": utils.GetCodeDescription(utils.REGISTER_MFA_ERROR),
		})
	}

	// Generar QR Code
	qrCodeURL, err := utils.GenerateQRCode(req.Email, mfaSecret, "Hospital System")
	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_MFA_ERROR, fiber.Map{
			"error":       "Error al generar código QR",
			"description": utils.GetCodeDescription(utils.REGISTER_MFA_ERROR),
		})
	}

	db := config.GetDB()
	var usuario models.Usuario

	// Determinar role_id basado en tipo si no se proporciona
	var roleID *int
	if req.RoleID != nil {
		roleID = req.RoleID
	} else {
		// Mapear tipo a role_id
		roleMap := map[string]int{
			"admin":     1,
			"medico":    2,
			"enfermera": 3,
			"paciente":  4,
		}
		if id, exists := roleMap[string(req.Tipo)]; exists {
			roleID = &id
		}
	}

	// Query actualizada
	err = db.QueryRow(`
		INSERT INTO Usuarios (nombre, email, password, tipo, role_id, mfa_enabled, mfa_secret, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) 
		RETURNING id_usuario, nombre, email, tipo, created_at, updated_at`,
		req.Nombre, req.Email, hashedPassword, req.Tipo, roleID, true, mfaSecret, time.Now(), time.Now()).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &usuario.Tipo, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return sendStandardResponse(c, 500, utils.REGISTER_DB_ERROR, fiber.Map{
			"error":       "Error al crear usuario",
			"description": utils.GetCodeDescription(utils.REGISTER_DB_ERROR),
		})
	}

	// Respuesta exitosa
	responseData := fiber.Map{
		"description": utils.GetCodeDescription(utils.REGISTER_SUCCESS),
		"message":     "Usuario registrado. Configura Microsoft Authenticator con el código QR o secreto.",
		"secret_key":  mfaSecret,  // Cambiar de mfa_secret a secret_key
		"qr_code_url": qrCodeURL,  // Agregar QR code URL
	}

	return sendStandardResponse(c, 201, utils.REGISTER_SUCCESS, responseData)
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return sendStandardResponse(c, 400, utils.LOGIN_PARSE_ERROR, fiber.Map{
			"error":       "Error al parsear datos",
			"description": utils.GetCodeDescription(utils.LOGIN_PARSE_ERROR),
		})
	}

	clientIP := c.IP()
	userAgent := c.Get("User-Agent")

	db := config.GetDB()
	var usuario models.Usuario
	var hashedPassword, mfaSecret string
	var backupCodes pq.StringArray

	err := db.QueryRow(`
		SELECT id_usuario, nombre, email, password, tipo, mfa_enabled, mfa_secret, backup_codes, created_at, updated_at 
		FROM Usuarios WHERE email = $1`, req.Email).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &hashedPassword, &usuario.Tipo,
		&usuario.MFAEnabled, &mfaSecret, &backupCodes, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			logLoginAttempt(db, req.Email, clientIP, userAgent, false)
			return sendStandardResponse(c, 401, utils.LOGIN_INVALID_CREDENTIALS, fiber.Map{
				"error":       "Credenciales inválidas",
				"description": utils.GetCodeDescription(utils.LOGIN_INVALID_CREDENTIALS),
			})
		}
		return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
			"error":       "Error del servidor",
			"description": utils.GetCodeDescription(utils.LOGIN_SERVER_ERROR),
		})
	}

	if !utils.CheckPassword(req.Password, hashedPassword) {
		logLoginAttempt(db, req.Email, clientIP, userAgent, false)
		return sendStandardResponse(c, 401, utils.LOGIN_INVALID_CREDENTIALS, fiber.Map{
			"error":       "Credenciales inválidas",
			"description": utils.GetCodeDescription(utils.LOGIN_INVALID_CREDENTIALS),
		})
	}

	// Validar MFA
	if usuario.MFAEnabled {
		if req.TOTPCode == "" {
			responseData := fiber.Map{
				"requires_mfa": true,
				"message":      "Se requiere código MFA",
				"description":  utils.GetCodeDescription(utils.LOGIN_MFA_REQUIRED),
			}
			return sendStandardResponse(c, 200, utils.LOGIN_MFA_REQUIRED, responseData)
		}

		validMFA := false
		if utils.ValidateTOTP(mfaSecret, req.TOTPCode) {
			validMFA = true
		} else {
			newBackupCodes, isBackupCode := utils.ValidateBackupCode(backupCodes, req.TOTPCode)
			if isBackupCode {
				validMFA = true
				_, _ = db.Exec("UPDATE Usuarios SET backup_codes = $1 WHERE id_usuario = $2",
					pq.Array(newBackupCodes), usuario.IDUsuario)
			}
		}

		if !validMFA {
			logLoginAttempt(db, req.Email, clientIP, userAgent, false)
			return sendStandardResponse(c, 401, utils.LOGIN_INVALID_MFA, fiber.Map{
				"error":       "Código MFA inválido",
				"description": utils.GetCodeDescription(utils.LOGIN_INVALID_MFA),
			})
		}
	}

	logLoginAttempt(db, req.Email, clientIP, userAgent, true)

	accessToken, err := utils.GenerateAccessToken(usuario.IDUsuario, usuario.Email, string(usuario.Tipo))
	if err != nil {
		return sendStandardResponse(c, 500, utils.LOGIN_TOKEN_ERROR, fiber.Map{
			"error":       "Error al generar token",
			"description": utils.GetCodeDescription(utils.LOGIN_TOKEN_ERROR),
		})
	}

	refreshToken, err := utils.GenerateRefreshToken(usuario.IDUsuario)
	if err != nil {
		return sendStandardResponse(c, 500, utils.LOGIN_TOKEN_ERROR, fiber.Map{
			"error":       "Error al generar refresh token",
			"description": utils.GetCodeDescription(utils.LOGIN_TOKEN_ERROR),
		})
	}

	tokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	_, err = db.Exec(`
		UPDATE Usuarios SET refresh_token = $1, token_expiry = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id_usuario = $3`,
		refreshToken, tokenExpiry, usuario.IDUsuario)

	if err != nil {
		return sendStandardResponse(c, 500, utils.LOGIN_SERVER_ERROR, fiber.Map{
			"error":       "Error al guardar refresh token",
			"description": utils.GetCodeDescription(utils.LOGIN_SERVER_ERROR),
		})
	}

	// Respuesta exitosa
	responseData := fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user":          usuario,
		"message":       "Login exitoso",
		"description":   utils.GetCodeDescription(utils.LOGIN_SUCCESS),
	}

	return sendStandardResponse(c, 200, utils.LOGIN_SUCCESS, responseData)
}

func RefreshToken(c *fiber.Ctx) error {
	var req models.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()
	var usuario models.Usuario
	var tokenExpiry time.Time

	err := db.QueryRow(`
		SELECT id_usuario, nombre, email, tipo, token_expiry 
		FROM Usuarios WHERE refresh_token = $1`, req.RefreshToken).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &usuario.Tipo, &tokenExpiry)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Refresh token inválido"})
	}

	if time.Now().After(tokenExpiry) {
		return c.Status(401).JSON(fiber.Map{"error": "Refresh token expirado"})
	}

	accessToken, err := utils.GenerateAccessToken(usuario.IDUsuario, usuario.Email, string(usuario.Tipo))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar token"})
	}

	return c.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func logLoginAttempt(db *sql.DB, email, ip, userAgent string, success bool) {
	status := "FAILED"
	if success {
		status = "SUCCESS"
	}

	fmt.Printf("[LOGIN %s] Email: %s, IP: %s, Time: %s\n", status, email, ip, time.Now().Format(time.RFC3339))
}
