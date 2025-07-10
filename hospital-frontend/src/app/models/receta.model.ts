export interface Receta {
  id_receta?: number;
  fecha: Date;
  medico_id: number;
  medicamento?: string;
  dosis?: string;
  consultorio_id: number;
  paciente_id: number;
}

export interface CreateRecetaRequest {
  fecha: Date;
  medico_id: number;
  medicamento?: string;
  dosis?: string;
  consultorio_id: number;
  paciente_id: number;
}