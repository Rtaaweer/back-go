package utils

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// IsStrongPassword verifica si la contraseña cumple con los requisitos de seguridad
// Retorna true si es fuerte, false si no
func IsStrongPassword(password string) bool {
	return ValidatePasswordStrength(password) == nil
}

// ValidatePasswordStrength valida que la contraseña cumpla con los requisitos de seguridad
func ValidatePasswordStrength(password string) error {
	// Verificar longitud mínima de 12 caracteres
	if len(password) < 12 {
		return errors.New("la contraseña debe tener al menos 12 caracteres")
	}

	// Verificar que contenga al menos un número
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("la contraseña debe contener al menos un número")
	}

	// Verificar que contenga al menos una letra minúscula
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLower {
		return errors.New("la contraseña debe contener al menos una letra minúscula")
	}

	// Verificar que contenga al menos una letra mayúscula
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpper {
		return errors.New("la contraseña debe contener al menos una letra mayúscula")
	}

	// Verificar que contenga al menos un símbolo especial
	hasSymbol := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?~` + "`" + `]`).MatchString(password)
	if !hasSymbol {
		return errors.New("la contraseña debe contener al menos un símbolo especial (!@#$%^&*()_+-=[]{}|;':,.<>?~)")
	}

	return nil
}

// HashPasswordWithValidation valida la fortaleza de la contraseña antes de hashearla
func HashPasswordWithValidation(password string) (string, error) {
	// Validar la fortaleza de la contraseña primero
	if err := ValidatePasswordStrength(password); err != nil {
		return "", err
	}

	// Si la contraseña es válida, proceder con el hash
	return HashPassword(password)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
