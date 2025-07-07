package models

import "time"

type Receta struct {
    IDReceta      int       `json:"id_receta"`
    Fecha         time.Time `json:"fecha"`
    MedicoID      int       `json:"medico_id"`
    Medicamento   *string   `json:"medicamento"`
    Dosis         *string   `json:"dosis"`
    ConsultorioID int       `json:"consultorio_id"`
    PacienteID    int       `json:"paciente_id"`
}