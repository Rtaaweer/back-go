import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { finalize } from 'rxjs';

// Servicio simple de loading
export class LoadingService {
  private loading = false;
  
  show() {
    this.loading = true;
    console.log('Loading started');
  }
  
  hide() {
    this.loading = false;
    console.log('Loading finished');
  }
  
  get isLoading() {
    return this.loading;
  }
}

export const loadingInterceptor: HttpInterceptorFn = (req, next) => {
  const loadingService = inject(LoadingService);
  
  loadingService.show();
  
  return next(req).pipe(
    finalize(() => {
      loadingService.hide();
    })
  );
};