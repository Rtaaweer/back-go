import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Expediente, CreateExpedienteRequest } from '../../models/expediente.model';
import { ExpedienteService } from '../../services/expediente.service';
import { UsuarioService } from '../../services/usuario.service';
import { Usuario } from '../../models/usuario.model';

@Component({
  selector: 'app-expedientes',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  templateUrl: './expedientes.component.html',
  styleUrls: ['./expedientes.component.css']
})
export class ExpedientesComponent implements OnInit {
  expedientes: Expediente[] = [];
  expedienteForm: FormGroup;
  editingExpediente: Expediente | null = null;
  showForm = false;
  loading = false;
  pacientes: Usuario[] = [];

  constructor(
    private expedienteService: ExpedienteService,
    private usuarioService: UsuarioService,
    private fb: FormBuilder
  ) {
    this.expedienteForm = this.fb.group({
      paciente_id: ['', [Validators.required]],
      antecedentes: [''],
      historial_clinico: [''],
      seguro: ['']
    });
  }

  ngOnInit(): void {
    this.loadExpedientes();
    this.loadPacientes();
  }

  loadExpedientes(): void {
    this.loading = true;
    this.expedienteService.getExpedientes().subscribe({
      next: (expedientes) => {
        this.expedientes = expedientes;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading expedientes:', error);
        this.loading = false;
      }
    });
  }

  loadPacientes(): void {
    this.usuarioService.getUsuarios().subscribe({
      next: (usuarios) => {
        this.pacientes = usuarios.filter(u => u.tipo === 'paciente');
      },
      error: (error) => {
        console.error('Error loading pacientes:', error);
      }
    });
  }

  openCreateForm(): void {
    this.editingExpediente = null;
    this.expedienteForm.reset();
    this.showForm = true;
  }

  openEditForm(expediente: Expediente): void {
    this.editingExpediente = expediente;
    this.expedienteForm.patchValue({
      paciente_id: expediente.paciente_id,
      antecedentes: expediente.antecedentes,
      historial_clinico: expediente.historial_clinico,
      seguro: expediente.seguro
    });
    this.showForm = true;
  }

  closeForm(): void {
    this.showForm = false;
    this.editingExpediente = null;
    this.expedienteForm.reset();
  }

  onSubmit(): void {
    if (this.expedienteForm.valid) {
      const expedienteData: CreateExpedienteRequest = this.expedienteForm.value;
      
      if (this.editingExpediente) {
        this.updateExpediente(this.editingExpediente.id_expediente!, expedienteData);
      } else {
        this.createExpediente(expedienteData);
      }
    }
  }

  createExpediente(expedienteData: CreateExpedienteRequest): void {
    this.expedienteService.createExpediente(expedienteData).subscribe({
      next: (expediente) => {
        this.expedientes.push(expediente);
        this.closeForm();
      },
      error: (error) => {
        console.error('Error creating expediente:', error);
      }
    });
  }

  updateExpediente(id: number, expedienteData: CreateExpedienteRequest): void {
    this.expedienteService.updateExpediente(id, expedienteData).subscribe({
      next: (expediente) => {
        const index = this.expedientes.findIndex(e => e.id_expediente === id);
        if (index !== -1) {
          this.expedientes[index] = expediente;
        }
        this.closeForm();
      },
      error: (error) => {
        console.error('Error updating expediente:', error);
      }
    });
  }

  deleteExpediente(id: number): void {
    if (confirm('¿Está seguro de que desea eliminar este expediente?')) {
      this.expedienteService.deleteExpediente(id).subscribe({
        next: () => {
          this.expedientes = this.expedientes.filter(e => e.id_expediente !== id);
        },
        error: (error) => {
          console.error('Error deleting expediente:', error);
        }
      });
    }
  }

  getPacienteNombre(pacienteId: number): string {
    const paciente = this.pacientes.find(p => p.id_usuario === pacienteId);
    return paciente ? paciente.nombre : 'N/A';
  }
}