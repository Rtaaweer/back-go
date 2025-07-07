package handlers

import (
	"hospital-system/config"
	"hospital-system/models"

	"github.com/gofiber/fiber/v2"
)

func GetConsultas(c *fiber.Ctx) error {
	db := config.GetDB()

	rows, err := db.Query(`
        SELECT id_consulta, consultorio_id, medico_id, paciente_id, 
               tipo, horario, diagnostico, costo 
        FROM Consultas`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultas"})
	}
	defer rows.Close()

	var consultas []models.Consulta
	for rows.Next() {
		var consulta models.Consulta
		err := rows.Scan(&consulta.IDConsulta, &consulta.ConsultorioID,
			&consulta.MedicoID, &consulta.PacienteID,
			&consulta.Tipo, &consulta.Horario,
			&consulta.Diagnostico, &consulta.Costo)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al escanear consulta"})
		}
		consultas = append(consultas, consulta)
	}

	return c.JSON(consultas)
}

func CreateConsulta(c *fiber.Ctx) error {
	var consulta models.Consulta
	if err := c.BodyParser(&consulta); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Error al parsear datos"})
	}

	db := config.GetDB()

	err := db.QueryRow(`
        INSERT INTO Consultas (consultorio_id, medico_id, paciente_id, tipo, horario, diagnostico, costo) 
        VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_consulta`,
		consulta.ConsultorioID, consulta.MedicoID, consulta.PacienteID,
		consulta.Tipo, consulta.Horario, consulta.Diagnostico, consulta.Costo).Scan(&consulta.IDConsulta)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al crear consulta"})
	}

	return c.Status(201).JSON(consulta)
}
