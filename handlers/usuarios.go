package handlers

import (
    "menchaca-api/config"
    "menchaca-api/models"
    "strconv"

    "github.com/gofiber/fiber/v2"
)

// GetUsuarios obtiene todos los usuarios
func GetUsuarios(c *fiber.Ctx) error {
    db := config.GetDB()
    
    rows, err := db.Query("SELECT id_usuario, nombre, tipo FROM Usuarios")
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "No se pudieron obtener usuarios"})
    }
    defer rows.Close() 

    var usuarios []models.Usuario
    for rows.Next() {
        var usuario models.Usuario
        err := rows.Scan(&usuario.IDUsuario, &usuario.Nombre, &usuario.Tipo)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al obtener usuario"})
        }
        usuarios = append(usuarios, usuario)
    }

    return c.JSON(usuarios)
}

// GetUsuario obtiene un usuario por ID
func GetUsuario(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    db := config.GetDB()
    var usuario models.Usuario
    
    err = db.QueryRow("SELECT id_usuario, nombre, tipo FROM Usuarios WHERE id_usuario = $1", id).Scan(
        &usuario.IDUsuario, &usuario.Nombre, &usuario.Tipo)
    
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Usuario no encontrado"})
    }

    return c.JSON(usuario)
}

// CreateUsuario crea un nuevo usuario
func CreateUsuario(c *fiber.Ctx) error {
    var usuario models.Usuario
    if err := c.BodyParser(&usuario); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
    }

    db := config.GetDB()
    
    err := db.QueryRow(
        "INSERT INTO Usuarios (nombre, tipo) VALUES ($1, $2) RETURNING id_usuario",
        usuario.Nombre, usuario.Tipo).Scan(&usuario.IDUsuario)
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al crear usuario"})
    }

    return c.Status(201).JSON(usuario)
}

// UpdateUsuario actualiza un usuario
func UpdateUsuario(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    var usuario models.Usuario
    if err := c.BodyParser(&usuario); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
    }

    db := config.GetDB()
    
    _, err = db.Exec(
        "UPDATE Usuarios SET nombre = $1, tipo = $2 WHERE id_usuario = $3",
        usuario.Nombre, usuario.Tipo, id)
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar usuario"})
    }

    usuario.IDUsuario = id
    return c.JSON(usuario)
}

// DeleteUsuario elimina un usuario
func DeleteUsuario(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    db := config.GetDB()
    
    _, err = db.Exec("DELETE FROM Usuarios WHERE id_usuario = $1", id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar usuario"})
    }

    return c.JSON(fiber.Map{"message": "Usuario eliminado correctamente"})
}