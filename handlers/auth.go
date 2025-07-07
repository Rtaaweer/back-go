package handlers

import (
	"database/sql"
	"fmt"
	"hospital-system/config"
	"hospital-system/models"
	"hospital-system/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

func Register(c *fiber.Ctx) error {
	var req struct {
		Nombre   string             `json:"nombre" validate:"required"`
		Email    string             `json:"email" validate:"required,email"`
		Password string             `json:"password" validate:"required,min=12"`
		Tipo     models.TipoUsuario `json:"tipo" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	if len(req.Password) < 12 {
		return c.Status(400).JSON(fiber.Map{"error": "La contraseña debe tener al menos 12 caracteres"})
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al procesar contraseña"})
	}
//genera lo del mfa
	mfaSecret, err := utils.GenerateMFASecret()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar secreto MFA"})
	}

	qrCodeURL, err := utils.GenerateQRCode(req.Email, mfaSecret, "Hospital System")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar código QR"})
	}

	db := config.GetDB()
	var usuario models.Usuario

	err = db.QueryRow(`
		INSERT INTO Usuarios (nombre, email, password, tipo, mfa_enabled, mfa_secret, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
		RETURNING id_usuario, nombre, email, tipo, created_at, updated_at`,
		req.Nombre, req.Email, hashedPassword, req.Tipo, true, mfaSecret, time.Now(), time.Now()).Scan(
		&usuario.IDUsuario, &usuario.Nombre, &usuario.Email, &usuario.Tipo, &usuario.CreatedAt, &usuario.UpdatedAt)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear usuario"})
	}

	return c.Status(201).JSON(fiber.Map{
		"user":        usuario,
		"mfa_secret":  mfaSecret,
		"qr_code_url": qrCodeURL,
		"message":     "Usuario registrado. Configura Microsoft Authenticator con el código QR o secreto.",
	})
}

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
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
			return c.Status(401).JSON(fiber.Map{"error": "Credenciales inválidas"})
		}
		return c.Status(500).JSON(fiber.Map{"error": "Error del servidor"})
	}

	
	if !utils.CheckPassword(req.Password, hashedPassword) {
		logLoginAttempt(db, req.Email, clientIP, userAgent, false)
		return c.Status(401).JSON(fiber.Map{"error": "Credenciales inválidas"})
	}
//valida el codigo del mfa
	if usuario.MFAEnabled {
		if req.TOTPCode == "" {
			return c.Status(200).JSON(models.LoginResponse{
				RequiresMFA: true,
			})
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
			return c.Status(401).JSON(fiber.Map{"error": "Código MFA inválido"})
		}
	}

	logLoginAttempt(db, req.Email, clientIP, userAgent, true)

	accessToken, err := utils.GenerateAccessToken(usuario.IDUsuario, usuario.Email, string(usuario.Tipo))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar token"})
	}

	refreshToken, err := utils.GenerateRefreshToken(usuario.IDUsuario)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al generar refresh token"})
	}

	tokenExpiry := time.Now().Add(7 * 24 * time.Hour)
	_, err = db.Exec(`
		UPDATE Usuarios SET refresh_token = $1, token_expiry = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id_usuario = $3`,
		refreshToken, tokenExpiry, usuario.IDUsuario)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al guardar refresh token"})
	}

	response := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: usuario,
	}

	return c.JSON(response)
}

func logLoginAttempt(db *sql.DB, email, ip, userAgent string, success bool) {
	status := "FAILED"
	if success {
		status = "SUCCESS"
	}

	fmt.Printf("[LOGIN %s] Email: %s, IP: %s, Time: %s\n", status, email, ip, time.Now().Format(time.RFC3339))

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


