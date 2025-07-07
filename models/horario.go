package models

type Horario struct {
    IDHorario     int  `json:"id_horario"`
    ConsultorioID int  `json:"consultorio_id"`
    Turno         string `json:"turno"`
    MedicoID      int  `json:"medico_id"`
    ConsultaID    *int `json:"consulta_id"`
}