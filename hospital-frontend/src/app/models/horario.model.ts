export interface Horario {
  id_horario?: number;
  consultorio_id: number;
  turno: string;
  medico_id: number;
  consulta_id?: number;
}

export interface CreateHorarioRequest {
  consultorio_id: number;
  turno: string;
  medico_id: number;
  consulta_id?: number;
}