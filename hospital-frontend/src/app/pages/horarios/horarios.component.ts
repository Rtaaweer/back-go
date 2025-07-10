import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Horario, CreateHorarioRequest } from '../../models/horario.model';
import { HorarioService } from '../../services/horario.service';
import { UsuarioService } from '../../services/usuario.service';
import { ConsultorioService } from '../../services/consultorio.service';
import { Usuario } from '../../models/usuario.model';
import { Consultorio } from '../../models/consultorio.model';

@Component({
  selector: 'app-horarios',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  templateUrl: './horarios.component.html',
  styleUrls: ['./horarios.component.css']
})
export class HorariosComponent implements OnInit {
  horarios: Horario[] = [];
  horarioForm: FormGroup;
  editingHorario: Horario | null = null;
  showForm = false;
  loading = false;
  medicos: Usuario[] = [];
  consultorios: Consultorio[] = [];

  constructor(
    private horarioService: HorarioService,
    private usuarioService: UsuarioService,
    private consultorioService: ConsultorioService,
    private fb: FormBuilder
  ) {
    this.horarioForm = this.fb.group({
      consultorio_id: ['', [Validators.required]],
      medico_id: ['', [Validators.required]],
      turno: ['', [Validators.required]],
      consulta_id: ['']
    });
  }

  ngOnInit(): void {
    this.loadHorarios();
    this.loadMedicos();
    this.loadConsultorios();
  }

  loadHorarios(): void {
    this.loading = true;
    this.horarioService.getHorarios().subscribe({
      next: (horarios) => {
        this.horarios = horarios;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading horarios:', error);
        this.loading = false;
      }
    });
  }

  loadMedicos(): void {
    this.usuarioService.getUsuarios().subscribe({
      next: (usuarios) => {
        this.medicos = usuarios.filter(u => u.tipo === 'medico');
      },
      error: (error) => {
        console.error('Error loading medicos:', error);
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
    this.editingHorario = null;
    this.horarioForm.reset();
    this.showForm = true;
  }

  openEditForm(horario: Horario): void {
    this.editingHorario = horario;
    this.horarioForm.patchValue({
      consultorio_id: horario.consultorio_id,
      medico_id: horario.medico_id,
      turno: horario.turno,
      consulta_id: horario.consulta_id
    });
    this.showForm = true;
  }

  closeForm(): void {
    this.showForm = false;
    this.editingHorario = null;
    this.horarioForm.reset();
  }

  onSubmit(): void {
    if (this.horarioForm.valid) {
      const horarioData: CreateHorarioRequest = this.horarioForm.value;
      
      if (this.editingHorario) {
        this.updateHorario(this.editingHorario.id_horario!, horarioData);
      } else {
        this.createHorario(horarioData);
      }
    }
  }

  createHorario(horarioData: CreateHorarioRequest): void {
    this.horarioService.createHorario(horarioData).subscribe({
      next: (horario) => {
        this.horarios.push(horario);
        this.closeForm();
      },
      error: (error) => {
        console.error('Error creating horario:', error);
      }
    });
  }

  updateHorario(id: number, horarioData: CreateHorarioRequest): void {
    this.horarioService.updateHorario(id, horarioData).subscribe({
      next: (horario) => {
        const index = this.horarios.findIndex(h => h.id_horario === id);
        if (index !== -1) {
          this.horarios[index] = horario;
        }
        this.closeForm();
      },
      error: (error) => {
        console.error('Error updating horario:', error);
      }
    });
  }

  deleteHorario(id: number): void {
    if (confirm('¿Está seguro de que desea eliminar este horario?')) {
      this.horarioService.deleteHorario(id).subscribe({
        next: () => {
          this.horarios = this.horarios.filter(h => h.id_horario !== id);
        },
        error: (error) => {
          console.error('Error deleting horario:', error);
        }
      });
    }
  }

  getMedicoNombre(medicoId: number): string {
    const medico = this.medicos.find(m => m.id_usuario === medicoId);
    return medico ? medico.nombre : 'N/A';
  }

  getConsultorioNombre(consultorioId: number): string {
    const consultorio = this.consultorios.find(c => c.id_consultorio === consultorioId);
    return consultorio ? consultorio.nombre : 'N/A';
  }
}