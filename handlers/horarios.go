package handlers

import (
    "hospital-system/config"
    "hospital-system/models"
    "strconv"

    "github.com/gofiber/fiber/v2"
)

func GetHorarios(c *fiber.Ctx) error {
    db := config.GetDB()
    
    rows, err := db.Query(`
        SELECT id_horario, consultorio_id, turno, medico_id, consulta_id 
        FROM Horarios`)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al obtener horarios"})
    }
    defer rows.Close()

    var horarios []models.Horario
    for rows.Next() {
        var horario models.Horario
        err := rows.Scan(&horario.IDHorario, &horario.ConsultorioID, 
                        &horario.Turno, &horario.MedicoID, 
                        &horario.ConsultaID)
        if err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Error al escanear horario"})
        }
        horarios = append(horarios, horario)
    }

    return c.JSON(horarios)
}

func GetHorario(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    db := config.GetDB()
    var horario models.Horario
    
    err = db.QueryRow(`
        SELECT id_horario, consultorio_id, turno, medico_id, consulta_id 
        FROM Horarios WHERE id_horario = $1`, id).Scan(
        &horario.IDHorario, &horario.ConsultorioID, 
        &horario.Turno, &horario.MedicoID, 
        &horario.ConsultaID)
    
    if err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "Horario no encontrado"})
    }

    return c.JSON(horario)
}

func CreateHorario(c *fiber.Ctx) error {
    var horario models.Horario
    if err := c.BodyParser(&horario); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos: " + err.Error()})
    }

    db := config.GetDB()
    
    err := db.QueryRow(`
        INSERT INTO Horarios (consultorio_id, turno, medico_id, consulta_id) 
        VALUES ($1, $2, $3, $4) RETURNING id_horario`,
        horario.ConsultorioID, horario.Turno, 
        horario.MedicoID, horario.ConsultaID).Scan(&horario.IDHorario)
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al crear horario: " + err.Error()})
    }

    return c.Status(201).JSON(horario)
}

func UpdateHorario(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    var horario models.Horario
    if err := c.BodyParser(&horario); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
    }

    db := config.GetDB()
    
    _, err = db.Exec(`
        UPDATE Horarios 
        SET consultorio_id = $1, turno = $2, medico_id = $3, consulta_id = $4 
        WHERE id_horario = $5`,
        horario.ConsultorioID, horario.Turno, 
        horario.MedicoID, horario.ConsultaID, id)
    
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar horario"})
    }

    horario.IDHorario = id
    return c.JSON(horario)
}

func DeleteHorario(c *fiber.Ctx) error {
    id, err := strconv.Atoi(c.Params("id"))
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
    }

    db := config.GetDB()
    
    _, err = db.Exec("DELETE FROM Horarios WHERE id_horario = $1", id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar horario"})
    }

    return c.JSON(fiber.Map{"message": "Horario eliminado exitosamente"})
}