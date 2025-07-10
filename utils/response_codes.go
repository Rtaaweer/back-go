package utils

// ResponseCode estructura para códigos de respuesta estandarizados
type ResponseCode struct {
	StatusCode int         `json:"statuscode"`
	IntCode    string      `json:"intcode"`
	Data       interface{} `json:"data"`
}

// Códigos para Login
const (
	// Login Success
	LOGIN_SUCCESS      = "S01"
	LOGIN_MFA_REQUIRED = "S02"

	// Login Errors
	LOGIN_INVALID_CREDENTIALS = "E01"
	LOGIN_INVALID_MFA         = "E02"
	LOGIN_PARSE_ERROR         = "E03"
	LOGIN_SERVER_ERROR        = "E04"
	LOGIN_TOKEN_ERROR         = "E05"

	// Login Warnings
	LOGIN_ACCOUNT_LOCKED   = "W01"
	LOGIN_PASSWORD_EXPIRED = "W02"
)

// Códigos para Registro
const (
	// Register Success
	REGISTER_SUCCESS = "S03"

	// Register Errors - Usar códigos únicos
	REGISTER_PARSE_ERROR   = "E06" // Cambiar de E04 a E06
	REGISTER_WEAK_PASSWORD = "E07" // Agregar esta constante
	REGISTER_HASH_ERROR    = "E08" // Cambiar de E05 a E08
	REGISTER_MFA_ERROR     = "E09"
	REGISTER_QR_ERROR      = "E10"
	REGISTER_DB_ERROR      = "E11"

	// Register Warnings
	REGISTER_EMAIL_EXISTS = "W03"
)

// NewResponse crea una nueva respuesta estandarizada
func NewResponse(statusCode int, intCode string, data interface{}) ResponseCode {
	return ResponseCode{
		StatusCode: statusCode,
		IntCode:    intCode,
		Data:       data,
	}
}

// GetCodeDescription devuelve la descripción del código
func GetCodeDescription(code string) string {
	descriptions := map[string]string{
		// Login
		LOGIN_SUCCESS:             "Login exitoso",
		LOGIN_MFA_REQUIRED:        "MFA requerido",
		LOGIN_INVALID_CREDENTIALS: "Credenciales inválidas",
		LOGIN_INVALID_MFA:         "Código MFA inválido",
		LOGIN_PARSE_ERROR:         "Error al parsear datos de login",
		LOGIN_SERVER_ERROR:        "Error del servidor en login",
		LOGIN_TOKEN_ERROR:         "Error al generar token",
		LOGIN_ACCOUNT_LOCKED:      "Cuenta bloqueada",
		LOGIN_PASSWORD_EXPIRED:    "Contraseña expirada",

		// Register
		REGISTER_SUCCESS:       "Registro exitoso",
		REGISTER_PARSE_ERROR:   "Error al parsear datos de registro",
		REGISTER_WEAK_PASSWORD: "Contraseña débil",
		REGISTER_HASH_ERROR:    "Error al procesar contraseña",
		REGISTER_MFA_ERROR:     "Error al generar MFA",
		REGISTER_QR_ERROR:      "Error al generar código QR",
		REGISTER_DB_ERROR:      "Error al crear usuario",
		REGISTER_EMAIL_EXISTS:  "Email ya existe",
	}

	if desc, exists := descriptions[code]; exists {
		return desc
	}
	return "Código desconocido"
}
