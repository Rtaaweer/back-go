import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private apiUrl = 'http://localhost:3001/api/v1/auth';

  constructor(private http: HttpClient) {}

  login(credentials: {email: string, password: string, totp_code?: string}): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, credentials);
  }

  register(userData: {nombre: string, email: string, tipo: string, password: string}): Observable<any> {
    return this.http.post(`${this.apiUrl}/register`, userData);
  }

  // Nuevo método para verificar si el usuario está autenticado
  isAuthenticated(): boolean {
    const token = localStorage.getItem('token');
    return !!token;
  }

  // Nuevo método para obtener el token
  getToken(): string | null {
    return localStorage.getItem('token');
  }

  // Nuevo método para logout
  logout(): void {
    localStorage.removeItem('token');
  }
}