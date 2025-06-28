package models

type Expediente struct {
    IDExpediente     int     `json:"id_expediente"`
    Antecedentes     *string `json:"antecedentes"`
    HistorialClinico *string `json:"historial_clinico"`
    PacienteID       int     `json:"paciente_id"`
    Seguro           *string `json:"seguro"`
}