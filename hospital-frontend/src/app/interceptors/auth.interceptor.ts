import { HttpInterceptorFn } from '@angular/common/http';

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const token = localStorage.getItem('token');
  
  // No agregar token a las rutas de autenticación
  if (req.url.includes('/auth/')) {
    return next(req);
  }
  
  if (token && token.length > 10) {
    console.log('Agregando token a la petición:', token.substring(0, 20) + '...');
    console.log('URL de la petición:', req.url);
    
    const authReq = req.clone({
      setHeaders: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      }
    });
    return next(authReq);
  } else {
    console.log('No hay token válido disponible para la petición a:', req.url);
    console.log('Token actual:', token);
  }
  
  return next(req);
};