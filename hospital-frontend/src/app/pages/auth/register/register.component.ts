import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { PasswordModule } from 'primeng/password';
import { DropdownModule } from 'primeng/dropdown';
import { CardModule } from 'primeng/card';
import { MessageModule } from 'primeng/message';

@Component({
  selector: 'app-register',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    ButtonModule,
    InputTextModule,
    PasswordModule,
    DropdownModule,
    CardModule,
    MessageModule
  ],
  template: `
    <div class="auth-form-container">
      <p-card>
        <div class="text-center">
          <h1>Registro</h1>
          <p>Crear nueva cuenta en el sistema</p>
        </div>
        
        <form [formGroup]="registerForm" (ngSubmit)="onSubmit()">
          <div class="p-field">
            <label for="nombre">Nombre Completo</label>
            <input 
              pInputText 
              id="nombre" 
              formControlName="nombre" 
              placeholder="Ingresa tu nombre completo"
              [class.ng-invalid]="registerForm.get('nombre')?.invalid && registerForm.get('nombre')?.touched"
            />
            <small 
              class="p-error" 
              *ngIf="registerForm.get('nombre')?.invalid && registerForm.get('nombre')?.touched"
            >
              El nombre es requerido
            </small>
          </div>
          
          <div class="p-field">
            <label for="email">Correo Electrónico</label>
            <input 
              pInputText 
              id="email" 
              formControlName="email" 
              placeholder="ejemplo@hospital.com"
              [class.ng-invalid]="registerForm.get('email')?.invalid && registerForm.get('email')?.touched"
            />
            <small 
              class="p-error" 
              *ngIf="registerForm.get('email')?.invalid && registerForm.get('email')?.touched"
            >
              El correo electrónico es requerido y debe ser válido
            </small>
          </div>
          
          <div class="p-field">
            <label for="tipo">Tipo de Usuario</label>
            <p-dropdown 
              formControlName="tipo" 
              [options]="tiposUsuario" 
              placeholder="Selecciona el tipo de usuario"
              optionLabel="label" 
              optionValue="value"
            ></p-dropdown>
            <small 
              class="p-error" 
              *ngIf="registerForm.get('tipo')?.invalid && registerForm.get('tipo')?.touched"
            >
              Debes seleccionar un tipo de usuario
            </small>
          </div>
          
          <div class="p-field">
            <label for="password">Contraseña</label>
            <p-password 
              formControlName="password" 
              placeholder="Mínimo 6 caracteres"
              [toggleMask]="true"
              [feedback]="true"
            ></p-password>
            <small 
              class="p-error" 
              *ngIf="registerForm.get('password')?.invalid && registerForm.get('password')?.touched"
            >
              La contraseña debe tener al menos 6 caracteres
            </small>
          </div>
          
          <div class="p-field">
            <label for="confirmPassword">Confirmar Contraseña</label>
            <p-password 
              formControlName="confirmPassword" 
              placeholder="Repite tu contraseña"
              [toggleMask]="true"
              [feedback]="false"
            ></p-password>
            <small 
              class="p-error" 
              *ngIf="registerForm.hasError('passwordMismatch') && registerForm.get('confirmPassword')?.touched"
            >
              Las contraseñas no coinciden
            </small>
          </div>
          
          <div class="p-field">
            <p-button 
              label="Registrarse" 
              type="submit" 
              [disabled]="registerForm.invalid"
              [loading]="isLoading"
            ></p-button>
          </div>
        </form>
        
        <div class="text-center mt-3">
          <p>¿Ya tienes cuenta? <a routerLink="/login">Inicia sesión aquí</a></p>
        </div>
      </p-card>
    </div>
  `,
  styles: []
})
export class RegisterComponent {
  registerForm: FormGroup;
  isLoading = false;
  
  tiposUsuario = [
    { label: 'Médico', value: 'medico' },
    { label: 'Enfermero/a', value: 'enfermero' },
    { label: 'Administrador', value: 'admin' },
    { label: 'Recepcionista', value: 'recepcionista' }
  ];

  constructor(
    private fb: FormBuilder,
    private router: Router
  ) {
    this.registerForm = this.fb.group({
      nombre: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      tipo: ['', [Validators.required]],
      password: ['', [Validators.required, Validators.minLength(6)]],
      confirmPassword: ['', [Validators.required]]
    }, { validators: this.passwordMatchValidator });
  }

  passwordMatchValidator(form: FormGroup) {
    const password = form.get('password');
    const confirmPassword = form.get('confirmPassword');
    
    if (password && confirmPassword && password.value !== confirmPassword.value) {
      return { passwordMismatch: true };
    }
    return null;
  }

  onSubmit() {
    if (this.registerForm.valid) {
      this.isLoading = true;
      
      // Simulación de registro (sin backend)
      setTimeout(() => {
        console.log('Registro simulado:', this.registerForm.value);
        alert('Registro exitoso! Redirigiendo al login...');
        this.isLoading = false;
        this.router.navigate(['/auth/login']);
      }, 2000);
    } else {
      // Marcar todos los campos como tocados para mostrar errores
      Object.keys(this.registerForm.controls).forEach(key => {
        this.registerForm.get(key)?.markAsTouched();
      });
    }
  }
}