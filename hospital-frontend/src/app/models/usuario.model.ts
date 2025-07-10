export enum TipoUsuario {
  PACIENTE = 'paciente',
  MEDICO = 'medico',
  ENFERMERA = 'enfermera',
  ADMIN = 'admin'
}

export interface Usuario {
  id_usuario?: number;
  nombre: string;
  email?: string;
  tipo: TipoUsuario;
  role_id?: number;
  mfa_enabled?: boolean;
  created_at?: Date;
  updated_at?: Date;
}

export interface CreateUsuarioRequest {
  nombre: string;
  email?: string;
  tipo: TipoUsuario;
  role_id?: number;
}