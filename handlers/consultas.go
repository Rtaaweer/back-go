package handlers

import (
	"hospital-system/config"
	"hospital-system/models"

	"github.com/gofiber/fiber/v2"
)

func GetConsultas(c *fiber.Ctx) error {
	db := config.GetDB()

	rows, err := db.Query(`
        SELECT c.id_consulta, c.consultorio_id, c.medico_id, c.paciente_id, 
               c.tipo, c.horario, c.diagnostico, c.costo,
               m.nombre as medico_nombre,
               p.nombre as paciente_nombre,
               con.nombre as consultorio_nombre
        FROM Consultas c
        LEFT JOIN Usuarios m ON c.medico_id = m.id_usuario AND m.tipo = 'medico'
        LEFT JOIN Usuarios p ON c.paciente_id = p.id_usuario AND p.tipo = 'paciente'
        LEFT JOIN Consultorios con ON c.consultorio_id = con.id_consultorio`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultas"})
	}
	defer rows.Close()

	var consultas []map[string]interface{}
	for rows.Next() {
		var consulta models.Consulta
		var medicoNombre, pacienteNombre, consultorioNombre *string
		
		err := rows.Scan(&consulta.IDConsulta, &consulta.ConsultorioID,
			&consulta.MedicoID, &consulta.PacienteID,
			&consulta.Tipo, &consulta.Horario,
			&consulta.Diagnostico, &consulta.Costo,
			&medicoNombre, &pacienteNombre, &consultorioNombre)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Error al escanear consulta"})
		}
		
		// Crear un mapa con toda la información
		consultaCompleta := map[string]interface{}{
			"id_consulta":       consulta.IDConsulta,
			"consultorio_id":    consulta.ConsultorioID,
			"medico_id":         consulta.MedicoID,
			"paciente_id":       consulta.PacienteID,
			"tipo":              consulta.Tipo,
			"horario":           consulta.Horario,
			"diagnostico":       consulta.Diagnostico,
			"costo":             consulta.Costo,
			"medico_nombre":     getValue(medicoNombre),
			"paciente_nombre":   getValue(pacienteNombre),
			"consultorio_nombre": getValue(consultorioNombre),
		}
		
		consultas = append(consultas, consultaCompleta)
	}

	return c.JSON(consultas)
}

// Función auxiliar para manejar valores nulos
func getValue(s *string) string {
	if s == nil {
		return "N/A"
	}
	return *s
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
