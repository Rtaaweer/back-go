import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-dashboard',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  
  constructor(
    private router: Router,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    // Verificar autenticación al cargar el dashboard
    const token = this.authService.getToken();
    console.log('Token en dashboard:', token ? token.substring(0, 20) + '...' : 'No hay token');
    
    if (!token) {
      console.log('No hay token, redirigiendo a login');
      alert('Sesión expirada. Por favor, inicia sesión nuevamente.');
      this.router.navigate(['/login']);
      return;
    }
    
    // Verificar que el token no esté vacío o sea inválido
    if (token.length < 10) {
      console.log('Token inválido, redirigiendo a login');
      localStorage.removeItem('token');
      this.router.navigate(['/login']);
      return;
    }
    
    console.log('Dashboard cargado correctamente con token válido');
  }

  navigateTo(route: string): void {
    this.router.navigate([route]);
  }

  logout(): void {
    this.authService.logout();
    this.router.navigate(['/login']);
  }
}