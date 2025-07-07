package models

type Consultorio struct {
    IDConsultorio int    `json:"id_consultorio"`
    Tipo          string `json:"tipo"`
    MedicoID      *int   `json:"medico_id"`
    Ubicacion     string `json:"ubicacion"`
    Nombre        string `json:"nombre"`
}