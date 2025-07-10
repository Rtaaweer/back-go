import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { RecetaService } from '../../services/receta.service';
import { Receta, CreateRecetaRequest } from '../../models/receta.model';

@Component({
  selector: 'app-recetas',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './recetas.component.html',
  styleUrls: ['./recetas.component.css']
})
export class RecetasComponent implements OnInit {
  recetas: Receta[] = [];
  recetaForm: FormGroup;
  isModalOpen = false;
  isEditing = false;
  currentRecetaId: number | null = null;
  loading = false;
  error: string | null = null;

  constructor(
    private recetaService: RecetaService,
    private fb: FormBuilder
  ) {
    this.recetaForm = this.fb.group({
      fecha: ['', [Validators.required]],
      medico_id: ['', [Validators.required, Validators.min(1)]],
      paciente_id: ['', [Validators.required, Validators.min(1)]],
      consultorio_id: ['', [Validators.required, Validators.min(1)]],
      medicamento: ['', [Validators.required, Validators.minLength(3)]],
      dosis: ['', [Validators.required, Validators.minLength(2)]]
    });
  }

  ngOnInit(): void {
    this.loadRecetas();
  }

  loadRecetas(): void {
    this.loading = true;
    this.recetaService.getRecetas().subscribe({
      next: (recetas) => {
        this.recetas = recetas;
        this.loading = false;
      },
      error: (error) => {
        this.error = 'Error al cargar las recetas';
        this.loading = false;
        console.error('Error:', error);
      }
    });
  }

  openModal(receta?: Receta): void {
    this.isModalOpen = true;
    this.isEditing = !!receta;
    this.currentRecetaId = receta?.id_receta || null;
    
    if (receta) {
      // Formatear la fecha para el input date
      const fechaFormatted = new Date(receta.fecha).toISOString().split('T')[0];
      this.recetaForm.patchValue({
        fecha: fechaFormatted,
        medico_id: receta.medico_id,
        paciente_id: receta.paciente_id,
        consultorio_id: receta.consultorio_id,
        medicamento: receta.medicamento,
        dosis: receta.dosis
      });
    } else {
      this.recetaForm.reset();
      // Establecer fecha por defecto a hoy
      const today = new Date().toISOString().split('T')[0];
      this.recetaForm.patchValue({ fecha: today });
    }
  }

  closeModal(): void {
    this.isModalOpen = false;
    this.isEditing = false;
    this.currentRecetaId = null;
    this.recetaForm.reset();
    this.error = null;
  }

  onSubmit(): void {
    if (this.recetaForm.valid) {
      this.loading = true;
      const recetaData: CreateRecetaRequest = {
        ...this.recetaForm.value,
        fecha: new Date(this.recetaForm.value.fecha)
      };

      if (this.isEditing && this.currentRecetaId) {
        this.recetaService.updateReceta(this.currentRecetaId, recetaData).subscribe({
          next: () => {
            this.loadRecetas();
            this.closeModal();
            this.loading = false;
          },
          error: (error) => {
            this.error = 'Error al actualizar la receta';
            this.loading = false;
            console.error('Error:', error);
          }
        });
      } else {
        this.recetaService.createReceta(recetaData).subscribe({
          next: () => {
            this.loadRecetas();
            this.closeModal();
            this.loading = false;
          },
          error: (error) => {
            this.error = 'Error al crear la receta';
            this.loading = false;
            console.error('Error:', error);
          }
        });
      }
    }
  }

  deleteReceta(id: number): void {
    if (confirm('¿Estás seguro de que quieres eliminar esta receta?')) {
      this.loading = true;
      this.recetaService.deleteReceta(id).subscribe({
        next: () => {
          this.loadRecetas();
          this.loading = false;
        },
        error: (error) => {
          this.error = 'Error al eliminar la receta';
          this.loading = false;
          console.error('Error:', error);
        }
      });
    }
  }

  formatDate(date: Date): string {
    return new Date(date).toLocaleDateString('es-ES');
  }
}