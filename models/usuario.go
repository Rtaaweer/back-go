package models

import "time"

type TipoUsuario string

const (
	Paciente  TipoUsuario = "paciente"
	Medico    TipoUsuario = "medico"
	Enfermera TipoUsuario = "enfermera"
	Admin     TipoUsuario = "admin"
)

type Usuario struct {
	IDUsuario    int         `json:"id_usuario"`
	Nombre       string      `json:"nombre"`
	Email        string      `json:"email"`
	Password     string      `json:"-"` 
	Tipo         TipoUsuario `json:"tipo"`
	RefreshToken string      `json:"-"` 
	TokenExpiry  *time.Time  `json:"-"` 
	MFAEnabled   bool        `json:"mfa_enabled"`
	MFASecret    string      `json:"-"` 
	BackupCodes  []string    `json:"-"` 
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=12"`
	TOTPCode string `json:"totp_code,omitempty"` 
}

type LoginResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         Usuario `json:"user"`
	RequiresMFA  bool    `json:"requires_mfa,omitempty"` 
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type EnableMFARequest struct {
	Password string `json:"password" validate:"required"`
}

type EnableMFAResponse struct {
	Secret      string   `json:"secret"`
	QRCodeURL   string   `json:"qr_code_url"`
	BackupCodes []string `json:"backup_codes"`
}

type VerifyMFARequest struct {
	Secret   string `json:"secret" validate:"required"`
	TOTPCode string `json:"totp_code" validate:"required"`
}

type DisableMFARequest struct {
	Password string `json:"password" validate:"required"`
	TOTPCode string `json:"totp_code" validate:"required"`
}
