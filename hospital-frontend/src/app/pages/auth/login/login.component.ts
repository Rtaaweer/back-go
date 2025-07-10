import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { ButtonModule } from 'primeng/button';
import { InputTextModule } from 'primeng/inputtext';
import { PasswordModule } from 'primeng/password';
import { CardModule } from 'primeng/card';
import { MessageModule } from 'primeng/message';
import { ToastModule } from 'primeng/toast'; // Agregar ToastModule
import { MessageService } from 'primeng/api'; // Agregar MessageService
import { AuthService } from '../../../services/auth.service';

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
    MessageModule,
    ToastModule // Agregar ToastModule a los imports
  ],
  providers: [MessageService], // Agregar MessageService como provider
  templateUrl: './login.component.html',  // ← Cambio aquí: de 'template' a 'templateUrl'
  styles: []
})
export class LoginComponent {
  loginForm: FormGroup;
  isLoading = false;
  showMFAInput = true; // Cambiar a true para que aparezca siempre

  constructor(
    private fb: FormBuilder,
    private router: Router,
    private authService: AuthService,
    private messageService: MessageService // Inyectar MessageService
  ) {
    this.loginForm = this.fb.group({
      email: ['', [Validators.required, Validators.email]],
      password: ['', [Validators.required]],
      totp_code: [''] // Campo opcional para MFA, siempre visible
    });
  }

  onSubmit() {
    if (this.loginForm.valid) {
      this.isLoading = true;
      
      const loginData = {
        email: this.loginForm.value.email,
        password: this.loginForm.value.password,
        ...(this.loginForm.value.totp_code && { totp_code: this.loginForm.value.totp_code })
      };
      
      this.authService.login(loginData).subscribe({
        next: (response) => {
          console.log('Respuesta completa del servidor:', response);
          
          // Si el servidor requiere MFA pero no se proporcionó
          if (response.intcode === 'S02' || response.data?.requires_mfa) {
            // El campo ya está visible, solo mostrar mensaje
            this.messageService.add({
              severity: 'warn',
              summary: 'MFA Requerido',
              detail: 'Por favor ingresa tu código de autenticación de 6 dígitos'
            });
            // Hacer el campo MFA obligatorio
            this.loginForm.get('totp_code')?.setValidators([Validators.required]);
            this.loginForm.get('totp_code')?.updateValueAndValidity();
            return;
          }
          
          // Verificar si es login exitoso
          if (response.intcode === 'S01') {
            console.log('Login exitoso:', response);
            
            const token = response.data?.access_token;
            
            if (token) {
              localStorage.setItem('token', token);
              console.log('Token guardado exitosamente:', token.substring(0, 20) + '...');
              this.router.navigate(['/dashboard']);
            } else {
              console.error('No se recibió token del servidor. Respuesta:', response);
              alert('Error: No se recibió token de autenticación del servidor');
            }
          } else {
            const errorMessage = response.data?.error || response.data?.message || 'Error desconocido';
            alert(`Error (${response.intcode}): ${errorMessage}`);
          }
          
          this.isLoading = false;
        },
        error: (error) => {
          console.error('Error completo en login:', error);
          
          if (error.status === 401) {
            if (this.loginForm.get('totp_code')?.value) {
              alert('Código MFA incorrecto. Intenta nuevamente.');
              this.loginForm.get('totp_code')?.setValue('');
            } else {
              alert('Credenciales incorrectas');
            }
          } else if (error.status === 0) {
            alert('No se puede conectar al servidor. Verifica que el backend esté corriendo en el puerto 3001.');
          } else {
            alert(`Error en el servidor (${error.status}): ${error.message || 'Error desconocido'}`);
          }
          this.isLoading = false;
        }
      });
    }
  }

  goToRegister(event: Event) {
    event.preventDefault();
    this.router.navigate(['/register']);
  }
}