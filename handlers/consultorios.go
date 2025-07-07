package handlers

import (
    "hospital-system/config"
    "hospital-system/models"
    "strconv"

    "github.com/gofiber/fiber/v2"
)

func GetConsultorios(c *fiber.Ctx) error {
    db := config.GetDB()
    
    rows, err := db.Query(`
        SELECT id_consultorio, tipo, medico_id, ubicacion, nombre 
        FROM Consultorios`)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultorios"})
    }
    defer rows.Close()

    var consultorios []models.Consultorio
    for rows.Next() {
        var consultorio models.Consultorio
        err := rows.Scan(&consultorio.IDConsultorio, &consultorio.Tipo, 
                        &consultorio.MedicoID, &consultorio.Ubicacion, 
                        &consultorio.Nombre)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al escanear consultorio"})
        }
        consultorios = append(consultorios, consultorio)
    }

    return c.JSON(consultorios)
}

func GetConsultorio(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    db := config.GetDB()
    var consultorio models.Consultorio
    
    err = db.QueryRow(`
        SELECT id_consultorio, tipo, medico_id, ubicacion, nombre 
        FROM Consultorios WHERE id_consultorio = $1`, id).Scan(
        &consultorio.IDConsultorio, &consultorio.Tipo, 
        &consultorio.MedicoID, &consultorio.Ubicacion, 
        &consultorio.Nombre)
    
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Consultorio no encontrado"})
    }

    return c.JSON(consultorio)
}

func CreateConsultorio(c *fiber.Ctx) error {
    var consultorio models.Consultorio
    if err := c.BodyParser(&consultorio); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos: " + err.Error()})
    }

    db := config.GetDB()
    
    err := db.QueryRow(`
        INSERT INTO Consultorios (tipo, medico_id, ubicacion, nombre) 
        VALUES ($1, $2, $3, $4) RETURNING id_consultorio`,
        consultorio.Tipo, consultorio.MedicoID, 
        consultorio.Ubicacion, consultorio.Nombre).Scan(&consultorio.IDConsultorio)
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al crear consultorio: " + err.Error()})
    }

    return c.Status(201).JSON(consultorio)
}

func UpdateConsultorio(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    var consultorio models.Consultorio
    if err := c.BodyParser(&consultorio); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
    }

    db := config.GetDB()
    
    _, err = db.Exec(`
        UPDATE Consultorios 
        SET tipo = $1, medico_id = $2, ubicacion = $3, nombre = $4 
        WHERE id_consultorio = $5`,
        consultorio.Tipo, consultorio.MedicoID, 
        consultorio.Ubicacion, consultorio.Nombre, id)
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar consultorio"})
    }

    consultorio.IDConsultorio = id
    return c.JSON(consultorio)
}

func DeleteConsultorio(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    db := config.GetDB()
    
    _, err = db.Exec("DELETE FROM Consultorios WHERE id_consultorio = $1", id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar consultorio"})
    }

    return c.JSON(fiber.Map{"message": "Consultorio eliminado exitosamente"})
}