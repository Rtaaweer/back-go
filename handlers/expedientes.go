package handlers

import (
	"hospital-system/config"
	"hospital-system/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetExpedientes(c *fiber.Ctx) error {
	db := config.GetDB()

	rows, err := db.Query(`
        SELECT id_expediente, antecedentes, historial_clinico, paciente_id, seguro 
        FROM Expedientes`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener expedientes"})
	}
	defer rows.Close()

	var expedientes []models.Expediente
	for rows.Next() {
		var expediente models.Expediente
		err := rows.Scan(&expediente.IDExpediente, &expediente.Antecedentes,
			&expediente.HistorialClinico, &expediente.PacienteID,
			&expediente.Seguro)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al escanear expediente"})
		}
		expedientes = append(expedientes, expediente)
	}

	return c.JSON(expedientes)
}

func GetExpediente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	db := config.GetDB()
	var expediente models.Expediente

	err = db.QueryRow(`
        SELECT id_expediente, antecedentes, historial_clinico, paciente_id, seguro 
        FROM Expedientes WHERE id_expediente = $1`, id).Scan(
		&expediente.IDExpediente, &expediente.Antecedentes,
		&expediente.HistorialClinico, &expediente.PacienteID,
		&expediente.Seguro)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Expediente no encontrado"})
	}

	return c.JSON(expediente)
}

func CreateExpediente(c *fiber.Ctx) error {
	var expediente models.Expediente
	if err := c.BodyParser(&expediente); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos: " + err.Error()})
	}

	db := config.GetDB()

	err := db.QueryRow(`
        INSERT INTO Expedientes (antecedentes, historial_clinico, paciente_id, seguro) 
        VALUES ($1, $2, $3, $4) RETURNING id_expediente`,
		expediente.Antecedentes, expediente.HistorialClinico,
		expediente.PacienteID, expediente.Seguro).Scan(&expediente.IDExpediente)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear expediente: " + err.Error()})
	}

	return c.Status(201).JSON(expediente)
}

func UpdateExpediente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	var expediente models.Expediente
	if err := c.BodyParser(&expediente); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()

	_, err = db.Exec(`
        UPDATE Expedientes 
        SET antecedentes = $1, historial_clinico = $2, paciente_id = $3, seguro = $4 
        WHERE id_expediente = $5`,
		expediente.Antecedentes, expediente.HistorialClinico,
		expediente.PacienteID, expediente.Seguro, id)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al actualizar expediente"})
	}

	expediente.IDExpediente = id
	return c.JSON(expediente)
}

func DeleteExpediente(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "ID inválido"})
	}

	db := config.GetDB()

	_, err = db.Exec("DELETE FROM Expedientes WHERE id_expediente = $1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al eliminar expediente"})
	}

	return c.JSON(fiber.Map{"message": "Expediente eliminado exitosamente"})
}
