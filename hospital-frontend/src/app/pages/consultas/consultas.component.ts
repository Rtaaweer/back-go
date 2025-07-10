import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Consulta, CreateConsultaRequest } from '../../models/consulta.model';
import { ConsultaService } from '../../services/consulta.service';
import { UsuarioService } from '../../services/usuario.service';
import { ConsultorioService } from '../../services/consultorio.service';
import { Usuario } from '../../models/usuario.model';
import { Consultorio } from '../../models/consultorio.model';

@Component({
  selector: 'app-consultas',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  templateUrl: './consultas.component.html',
  styleUrls: ['./consultas.component.css']
})
export class ConsultasComponent implements OnInit {
  consultas: Consulta[] = [];
  consultaForm: FormGroup;
  editingConsulta: Consulta | null = null;
  showForm = false;
  loading = false;
  medicos: Usuario[] = [];
  pacientes: Usuario[] = [];
  consultorios: Consultorio[] = [];

  constructor(
    private consultaService: ConsultaService,
    private usuarioService: UsuarioService,
    private consultorioService: ConsultorioService,
    private fb: FormBuilder
  ) {
    this.consultaForm = this.fb.group({
      consultorio_id: ['', [Validators.required]],
      medico_id: ['', [Validators.required]],
      paciente_id: ['', [Validators.required]],
      tipo: ['', [Validators.required]],
      horario: ['', [Validators.required]],
      diagnostico: ['']
    });
  }

  ngOnInit(): void {
    this.loadConsultas();
    this.loadUsuarios();
    this.loadConsultorios();
  }

  loadConsultas(): void {
    this.loading = true;
    this.consultaService.getConsultas().subscribe({
      next: (consultas) => {
        this.consultas = consultas;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading consultas:', error);
        this.loading = false;
      }
    });
  }

  loadUsuarios(): void {
    this.usuarioService.getUsuarios().subscribe({
      next: (usuarios) => {
        this.medicos = usuarios.filter(u => u.tipo === 'medico');
        this.pacientes = usuarios.filter(u => u.tipo === 'paciente');
      },
      error: (error) => {
        console.error('Error loading usuarios:', error);
      }
    });
  }

  loadConsultorios(): void {
    this.consultorioService.getConsultorios().subscribe({
      next: (consultorios) => {
        this.consultorios = consultorios;
      },
      error: (error) => {
        console.error('Error loading consultorios:', error);
      }
    });
  }

  openCreateForm(): void {
    this.editingConsulta = null;
    this.consultaForm.reset();
    this.showForm = true;
  }

  openEditForm(consulta: Consulta): void {
    this.editingConsulta = consulta;
    const horarioFormatted = new Date(consulta.horario).toISOString().slice(0, 16);
    this.consultaForm.patchValue({
      consultorio_id: consulta.consultorio_id,
      medico_id: consulta.medico_id,
      paciente_id: consulta.paciente_id,
      tipo: consulta.tipo,
      horario: horarioFormatted,
      diagnostico: consulta.diagnostico
    });
    this.showForm = true;
  }

  closeForm(): void {
    this.showForm = false;
    this.editingConsulta = null;
    this.consultaForm.reset();
  }

  onSubmit(): void {
    if (this.consultaForm.valid) {
      const consultaData: CreateConsultaRequest = {
        ...this.consultaForm.value,
        horario: new Date(this.consultaForm.value.horario)
      };
      
      if (this.editingConsulta) {
        this.updateConsulta(this.editingConsulta.id_consulta!, consultaData);
      } else {
        this.createConsulta(consultaData);
      }
    }
  }

  createConsulta(consultaData: CreateConsultaRequest): void {
    this.consultaService.createConsulta(consultaData).subscribe({
      next: (consulta) => {
        this.consultas.push(consulta);
        this.closeForm();
      },
      error: (error) => {
        console.error('Error creating consulta:', error);
      }
    });
  }

  updateConsulta(id: number, consultaData: CreateConsultaRequest): void {
    this.consultaService.updateConsulta(id, consultaData).subscribe({
      next: (consulta) => {
        const index = this.consultas.findIndex(c => c.id_consulta === id);
        if (index !== -1) {
          this.consultas[index] = consulta;
        }
        this.closeForm();
      },
      error: (error) => {
        console.error('Error updating consulta:', error);
      }
    });
  }

  deleteConsulta(id: number): void {
    if (confirm('¿Está seguro de que desea eliminar esta consulta?')) {
      this.consultaService.deleteConsulta(id).subscribe({
        next: () => {
          this.consultas = this.consultas.filter(c => c.id_consulta !== id);
        },
        error: (error) => {
          console.error('Error deleting consulta:', error);
        }
      });
    }
  }

  getMedicoNombre(medicoId: number): string {
    const medico = this.medicos.find(m => m.id_usuario === medicoId);
    return medico ? medico.nombre : 'N/A';
  }

  getPacienteNombre(pacienteId: number): string {
    const paciente = this.pacientes.find(p => p.id_usuario === pacienteId);
    return paciente ? paciente.nombre : 'N/A';
  }

  getConsultorioNombre(consultorioId: number): string {
    const consultorio = this.consultorios.find(c => c.id_consultorio === consultorioId);
    return consultorio ? consultorio.nombre : 'N/A';
  }
}