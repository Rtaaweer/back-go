export interface Consultorio {
  id_consultorio?: number;
  nombre: string;
  ubicacion?: string;
  capacidad?: number;
  equipamiento?: string;
}

export interface CreateConsultorioRequest {
  nombre: string;
  ubicacion?: string;
  capacidad?: number;
  equipamiento?: string;
}