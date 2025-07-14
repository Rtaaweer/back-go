package middleware

import (
	"fmt"
	"hospital-system/config"
	"hospital-system/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTMiddleware actualizado para cargar permisos
func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("[AUTH] Verificando autenticación para: %s %s\n", c.Method(), c.Path())

		// Obtener el header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			fmt.Printf("[AUTH] Token de autorización faltante\n")
			return c.Status(401).JSON(fiber.Map{
				"error": "Token de autorización requerido",
			})
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			fmt.Printf("[AUTH] Formato de token inválido\n")
			return c.Status(401).JSON(fiber.Map{
				"error": "Formato de token inválido",
			})
		}

		claims, err := utils.ValidateToken(tokenParts[1])
		if err != nil {
			fmt.Printf("[AUTH] Token inválido: %v\n", err)
			return c.Status(401).JSON(fiber.Map{
				"error": "Token inválido",
			})
		}

		fmt.Printf("[AUTH] Usuario autenticado: ID=%d, Email=%s, Tipo=%s\n", claims.UserID, claims.Email, claims.Tipo)

		// Si el token no tiene permisos cargados, cargarlos desde la base de datos
		if len(claims.Permissions) == 0 && claims.RoleID != nil {
			fmt.Printf("[AUTH] Cargando permisos para role_id: %d\n", *claims.RoleID)
			permissions, err := loadUserPermissions(*claims.RoleID)
			if err != nil {
				fmt.Printf("[AUTH] Error al cargar permisos: %v\n", err)
				return c.Status(500).JSON(fiber.Map{
					"error": "Error al cargar permisos",
				})
			}
			claims.Permissions = permissions
			fmt.Printf("[AUTH] Permisos cargados: %d permisos encontrados\n", len(permissions))
			for _, perm := range permissions {
				fmt.Printf("[AUTH]   - %s:%s\n", perm.Resource, perm.Action)
			}
		} else if len(claims.Permissions) > 0 {
			fmt.Printf("[AUTH] Permisos ya cargados en token: %d permisos\n", len(claims.Permissions))
			for _, perm := range claims.Permissions {
				fmt.Printf("[AUTH]   - %s:%s\n", perm.Resource, perm.Action)
			}
		} else {
			fmt.Printf("[AUTH] ⚠️ Usuario sin role_id asignado\n")
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_tipo", claims.Tipo)
		c.Locals("user_role_id", claims.RoleID)
		c.Locals("user_permissions", claims.Permissions)

		return c.Next()
	}
}

// RequirePermission middleware para verificar permisos específicos
func RequirePermission(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Printf("[PERM] Verificando permiso: %s:%s para %s %s\n", resource, action, c.Method(), c.Path())

		permissions, ok := c.Locals("user_permissions").([]utils.Permission)
		if !ok {
			fmt.Printf("[PERM] No se pudieron obtener los permisos del usuario\n")
			return c.Status(403).JSON(fiber.Map{
				"error": "No se pudieron obtener los permisos del usuario",
			})
		}

		userID := c.Locals("user_id")
		fmt.Printf("[PERM] Usuario ID: %v tiene %d permisos\n", userID, len(permissions))

		// Verificar si el usuario tiene el permiso específico
		for _, permission := range permissions {
			if permission.Resource == resource && permission.Action == action {
				fmt.Printf("[PERM] Permiso concedido: %s:%s\n", resource, action)
				return c.Next()
			}
		}

		fmt.Printf("[PERM]  Permiso denegado: %s:%s\n", resource, action)
		fmt.Printf("[PERM]  Permisos disponibles:\n")
		for _, perm := range permissions {
			fmt.Printf("[PERM]   - %s:%s\n", perm.Resource, perm.Action)
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Permisos insuficientes",
			"required": map[string]string{
				"resource": resource,
				"action":   action,
			},
		})
	}
}

// loadUserPermissions carga los permisos de un usuario desde la base de datos
func loadUserPermissions(roleID int) ([]utils.Permission, error) {
	fmt.Printf("[DB] Consultando permisos para role_id: %d\n", roleID)

	db := config.GetDB()
	var permissions []utils.Permission

	// Volver a la consulta original que solo selecciona resource y action
	rows, err := db.Query(`
		SELECT p.resource, p.action 
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`, roleID)
	if err != nil {
		fmt.Printf("[DB] Error en consulta de permisos: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var permission utils.Permission
		err := rows.Scan(&permission.Resource, &permission.Action)
		if err != nil {
			fmt.Printf("[DB] Error al escanear permiso: %v\n", err)
			return nil, err
		}
		permissions = append(permissions, permission)
		fmt.Printf("[DB] Permiso encontrado: %s:%s\n", permission.Resource, permission.Action)
		count++
	}

	// Verificar si no se encontraron permisos
	if count == 0 {
		fmt.Printf("[DB] ⚠️ No se encontraron permisos para el role_id: %d\n", roleID)
		// Agregar permisos por defecto según el rol para evitar bloqueos
		switch roleID {
		case 1: // Admin
			permissions = append(permissions, utils.Permission{Resource: "usuarios", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "usuarios", Action: "create"})
			permissions = append(permissions, utils.Permission{Resource: "usuarios", Action: "update"})
			permissions = append(permissions, utils.Permission{Resource: "usuarios", Action: "delete"})
			permissions = append(permissions, utils.Permission{Resource: "consultorios", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "consultorios", Action: "create"})
			permissions = append(permissions, utils.Permission{Resource: "consultorios", Action: "update"})
			permissions = append(permissions, utils.Permission{Resource: "consultorios", Action: "delete"})
			permissions = append(permissions, utils.Permission{Resource: "consultas", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "consultas", Action: "create"})
			permissions = append(permissions, utils.Permission{Resource: "expedientes", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "expedientes", Action: "create"})
			permissions = append(permissions, utils.Permission{Resource: "expedientes", Action: "update"})
			permissions = append(permissions, utils.Permission{Resource: "expedientes", Action: "delete"})
			permissions = append(permissions, utils.Permission{Resource: "horarios", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "horarios", Action: "create"})
			permissions = append(permissions, utils.Permission{Resource: "horarios", Action: "update"})
			permissions = append(permissions, utils.Permission{Resource: "horarios", Action: "delete"})
			permissions = append(permissions, utils.Permission{Resource: "recetas", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "recetas", Action: "create"})
			permissions = append(permissions, utils.Permission{Resource: "recetas", Action: "update"})
			permissions = append(permissions, utils.Permission{Resource: "recetas", Action: "delete"})
			fmt.Printf("[DB] Agregados permisos por defecto para Admin (role_id: 1)\n")
		case 2: // Médico
			permissions = append(permissions, utils.Permission{Resource: "usuarios", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "consultas", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "consultas", Action: "create"})
			fmt.Printf("[DB] Agregados permisos por defecto para Médico (role_id: 2)\n")
		case 3: // Enfermera
			permissions = append(permissions, utils.Permission{Resource: "usuarios", Action: "read"})
			permissions = append(permissions, utils.Permission{Resource: "consultas", Action: "read"})
			fmt.Printf("[DB] Agregados permisos por defecto para Enfermera (role_id: 3)\n")
		case 4: // Paciente
			permissions = append(permissions, utils.Permission{Resource: "consultas", Action: "read"})
			fmt.Printf("[DB] Agregados permisos por defecto para Paciente (role_id: 4)\n")
		}
	}

	fmt.Printf("[DB] Total permisos cargados: %d\n", len(permissions))
	return permissions, nil
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userTipo := c.Locals("user_tipo").(string)

		for _, role := range roles {
			if userTipo == role {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Permisos insuficientes",
		})
	}
}
