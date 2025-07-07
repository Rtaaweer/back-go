import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { PasswordModule } from 'primeng/password';
import { CardModule } from 'primeng/card';
import { MessageModule } from 'primeng/message';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterModule,
    ButtonModule,
    InputTextModule,
    PasswordModule,
    CardModule,
    MessageModule
  ],
  template: `
    <div class="auth-form-container">
      <p-card class="p-card">
        <ng-template pTemplate="header">
          <div class="text-center">
            <h1>Iniciar Sesión</h1>
            <p>Sistema Hospitalario</p>
          </div>
        </ng-template>
        
        <form [formGroup]="loginForm" (ngSubmit)="onSubmit()">
          <div class="p-field">
            <label for="email">Correo Electrónico</label>
            <input 
              pInputText 
              id="email" 
              formControlName="email" 
              placeholder="ejemplo@hospital.com"
              class="p-inputtext"
              [class.ng-invalid]="loginForm.get('email')?.invalid && loginForm.get('email')?.touched"
            />
            <small 
              class="p-error" 
              *ngIf="loginForm.get('email')?.invalid && loginForm.get('email')?.touched"
            >
              El correo electrónico es requerido y debe ser válido
            </small>
          </div>
          
          <div class="p-field">
            <label for="password">Contraseña</label>
            <p-password 
              formControlName="password" 
              placeholder="Ingresa tu contraseña"
              [toggleMask]="true"
              [feedback]="false"
              class="p-password"
            ></p-password>
            <small 
              class="p-error" 
              *ngIf="loginForm.get('password')?.invalid && loginForm.get('password')?.touched"
            >
              La contraseña es requerida
            </small>
          </div>
          
          <div class="p-field">
            <p-button 
              label="Iniciar Sesión" 
              type="submit" 
              [disabled]="loginForm.invalid"
              class="p-button"
              [loading]="isLoading"
            ></p-button>
          </div>
        </form>
        
        <ng-template pTemplate="footer">
          <div class="text-center mt-3">
            <p>¿No tienes cuenta? <a href="#" (click)="goToRegister($event)">Regístrate aquí</a></p>
          </div>
        </ng-template>
      </p-card>
    </div>
  `,
  styles: []
})
export class LoginComponent {
  loginForm: FormGroup;
  isLoading = false;

  constructor(
    private fb: FormBuilder,
    private router: Router
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required, Validators.minLength(6)]]
    });
  }

  onSubmit() {
    if (this.loginForm.valid) {
      this.isLoading = true;
      
      setTimeout(() => {
        console.log('Login simulado:', this.loginForm.value);
        alert('Login exitoso! (Simulación)');
        this.isLoading = false;
      }, 1500);
    } else {
      Object.keys(this.loginForm.controls).forEach(key => {
        this.loginForm.get(key)?.markAsTouched();
      });
    }
  }

  goToRegister(event: Event) {
    event.preventDefault();
    this.router.navigate(['/register']);
  }
}