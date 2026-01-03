import { inject } from "@angular/core";
import { Router, type UrlTree } from "@angular/router";

import { AuthService } from "../../services/auth.service";

/**
 * Guard that protects routes for guests only (unauthenticated users).
 * Redirects to / if the user is already authenticated.
 * @returns {boolean | UrlTree} True if not authenticated, otherwise redirects to home
 */
export const guestGuard = (): boolean | UrlTree => {
  const authService = inject(AuthService);
  const router = inject(Router);

  if (!authService.isAuthenticated()) {
    return true;
  }

  return router.createUrlTree(["/"]);
};
