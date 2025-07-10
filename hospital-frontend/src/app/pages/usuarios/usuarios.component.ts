import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Usuario, TipoUsuario, CreateUsuarioRequest } from '../../models/usuario.model';
import { UsuarioService } from '../../services/usuario.service';

@Component({
  selector: 'app-usuarios',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule],
  templateUrl: './usuarios.component.html',
  styleUrls: ['./usuarios.component.css']
})
export class UsuariosComponent implements OnInit {
  usuarios: Usuario[] = [];
  usuarioForm: FormGroup;
  editingUsuario: Usuario | null = null;
  showForm = false;
  loading = false;
  tiposUsuario = Object.values(TipoUsuario);

  constructor(
    private usuarioService: UsuarioService,
    private fb: FormBuilder
  ) {
    this.usuarioForm = this.fb.group({
      nombre: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.email]],
      tipo: ['', [Validators.required]]
    });
  }

  ngOnInit(): void {
    this.loadUsuarios();
  }

  loadUsuarios(): void {
    this.loading = true;
    this.usuarioService.getUsuarios().subscribe({
      next: (usuarios) => {
        this.usuarios = usuarios;
        this.loading = false;
      },
      error: (error) => {
        console.error('Error loading usuarios:', error);
        this.loading = false;
      }
    });
  }

  openCreateForm(): void {
    this.editingUsuario = null;
    this.usuarioForm.reset();
    this.showForm = true;
  }

  openEditForm(usuario: Usuario): void {
    this.editingUsuario = usuario;
    this.usuarioForm.patchValue({
      nombre: usuario.nombre,
      email: usuario.email,
      tipo: usuario.tipo
    });
    this.showForm = true;
  }

  closeForm(): void {
    this.showForm = false;
    this.editingUsuario = null;
    this.usuarioForm.reset();
  }

  onSubmit(): void {
    if (this.usuarioForm.valid) {
      const usuarioData: CreateUsuarioRequest = this.usuarioForm.value;
      
      if (this.editingUsuario) {
        this.updateUsuario(this.editingUsuario.id_usuario!, usuarioData);
      } else {
        this.createUsuario(usuarioData);
      }
    }
  }

  createUsuario(usuarioData: CreateUsuarioRequest): void {
    this.usuarioService.createUsuario(usuarioData).subscribe({
      next: (usuario) => {
        this.usuarios.push(usuario);
        this.closeForm();
      },
      error: (error) => {
        console.error('Error creating usuario:', error);
      }
    });
  }

  updateUsuario(id: number, usuarioData: CreateUsuarioRequest): void {
    this.usuarioService.updateUsuario(id, usuarioData).subscribe({
      next: (usuario) => {
        const index = this.usuarios.findIndex(u => u.id_usuario === id);
        if (index !== -1) {
          this.usuarios[index] = usuario;
        }
        this.closeForm();
      },
      error: (error) => {
        console.error('Error updating usuario:', error);
      }
    });
  }

  deleteUsuario(id: number): void {
    if (confirm('¿Está seguro de que desea eliminar este usuario?')) {
      this.usuarioService.deleteUsuario(id).subscribe({
        next: () => {
          this.usuarios = this.usuarios.filter(u => u.id_usuario !== id);
        },
        error: (error) => {
          console.error('Error deleting usuario:', error);
        }
      });
    }
  }
}