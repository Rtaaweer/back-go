import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    redirectTo: '/login',
    pathMatch: 'full'
  },
  {
    path: 'login',
    loadComponent: () => import('./pages/auth/login/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'register',
    loadComponent: () => import('./pages/auth/register/register.component').then(m => m.RegisterComponent)
  },
  {
    path: 'dashboard',
    loadComponent: () => import('./pages/dashboard/dashboard.component').then(m => m.DashboardComponent)
  },
  {
    path: 'usuarios',
    loadComponent: () => import('./pages/usuarios/usuarios.component').then(m => m.UsuariosComponent)
  },
  {
    path: 'consultorios',
    loadComponent: () => import('./pages/consultorios/consultorios.component').then(m => m.ConsultoriosComponent)
  },
  {
    path: 'consultas',
    loadComponent: () => import('./pages/consultas/consultas.component').then(m => m.ConsultasComponent)
  },
  {
    path: 'expedientes',
    loadComponent: () => import('./pages/expedientes/expedientes.component').then(m => m.ExpedientesComponent)
  },
  {
    path: 'horarios',
    loadComponent: () => import('./pages/horarios/horarios.component').then(m => m.HorariosComponent)
  },
  {
    path: 'recetas',
    loadComponent: () => import('./pages/recetas/recetas.component').then(m => m.RecetasComponent)
  },
  {
    path: '**',
    redirectTo: '/login'
  }
];
