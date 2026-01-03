import { inject } from "@angular/core";
import type { HttpInterceptorFn, HttpErrorResponse } from "@angular/common/http";
import { catchError, throwError } from "rxjs";

import { AuthService } from "../../services/auth.service";

/**
 * HTTP interceptor that handles authentication.
 * - Adds Authorization header with JWT token to API requests
 * - Excludes auth endpoints from token injection
 * - Clears tokens on 401 errors for auth endpoints
 * @param {HttpRequest<unknown>} req - The outgoing HTTP request
 * @param {HttpHandlerFn} next - The next handler in the chain
 * @returns {Observable<HttpEvent<unknown>>} The HTTP response observable
 */
export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);

  // Skip auth endpoints
  const isAuthEndpoint = req.url.startsWith("/auth");

  // Add Authorization header if token exists and not an auth endpoint
  let modifiedReq = req;
  if (!isAuthEndpoint) {
    const token = authService.getAccessToken();
    if (token !== undefined) {
      modifiedReq = req.clone({
        setHeaders: {
          Authorization: `Bearer ${token}`,
        },
      });
    }
  }

  return next(modifiedReq).pipe(
    catchError((error: HttpErrorResponse) => {
      // Clear tokens on 401 for auth endpoints (e.g., refresh failed)
      if (error.status === 401 && isAuthEndpoint) {
        authService.clearTokens();
      }
      return throwError(() => error);
    }),
  );
};
