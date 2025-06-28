package models

type TipoUsuario string

const (
    Paciente TipoUsuario = "paciente"
    Medico   TipoUsuario = "medico"
    Enfermera TipoUsuario = "enfermera"
    Admin    TipoUsuario = "admin"
)

type Usuario struct {
    IDUsuario int         `json:"id_usuario"`
    Nombre    string      `json:"nombre"`
    Tipo      TipoUsuario `json:"tipo"`
}