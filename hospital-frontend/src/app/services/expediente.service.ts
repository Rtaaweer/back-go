import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Expediente, CreateExpedienteRequest } from '../models/expediente.model';

@Injectable({
  providedIn: 'root'
})
export class ExpedienteService {
  private apiUrl = 'http://localhost:3001/api/v1/expedientes';

  constructor(private http: HttpClient) {}

  getExpedientes(): Observable<Expediente[]> {
    return this.http.get<Expediente[]>(this.apiUrl);
  }

  getExpediente(id: number): Observable<Expediente> {
    return this.http.get<Expediente>(`${this.apiUrl}/${id}`);
  }

  createExpediente(expediente: CreateExpedienteRequest): Observable<Expediente> {
    return this.http.post<Expediente>(this.apiUrl, expediente);
  }

  updateExpediente(id: number, expediente: CreateExpedienteRequest): Observable<Expediente> {
    return this.http.put<Expediente>(`${this.apiUrl}/${id}`, expediente);
  }

  deleteExpediente(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/${id}`);
  }
}