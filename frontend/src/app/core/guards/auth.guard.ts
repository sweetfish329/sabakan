import { inject } from "@angular/core";
import { Router, type UrlTree } from "@angular/router";

import { AuthService } from "../../services/auth.service";

/**
 * Guard that protects routes requiring authentication.
 * Redirects to /login if the user is not authenticated.
 * @returns {boolean | UrlTree} True if authenticated, otherwise redirects to login
 */
export const authGuard = (): boolean | UrlTree => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (authService.isAuthenticated()) {
    return true;
  }

  return router.createUrlTree(["/login"]);
};
