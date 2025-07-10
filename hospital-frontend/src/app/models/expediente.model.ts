export interface Expediente {
  id_expediente?: number;
  antecedentes?: string;
  historial_clinico?: string;
  paciente_id: number;
  seguro?: string;
}

export interface CreateExpedienteRequest {
  antecedentes?: string;
  historial_clinico?: string;
  paciente_id: number;
  seguro?: string;
}