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
  templateUrl: './login.component.html',  // ← Cambio aquí: de 'template' a 'templateUrl'
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
      password: ['', [Validators.required]]
    });
  }

  onSubmit() {
    if (this.loginForm.valid) {
      this.isLoading = true;
      // Aquí iría la lógica de autenticación
      console.log('Datos del formulario:', this.loginForm.value);
      
      // Simulación de login
      setTimeout(() => {
        this.isLoading = false;
        // Redirigir al dashboard o página principal
        // this.router.navigate(['/dashboard']);
      }, 2000);
    }
  }

  goToRegister(event: Event) {
    event.preventDefault();
    // this.router.navigate(['/auth/register']);
    console.log('Navegar a registro');
  }
}