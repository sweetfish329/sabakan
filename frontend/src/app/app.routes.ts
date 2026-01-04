import type { Routes } from "@angular/router";

import { authGuard } from "./core/guards/auth.guard";
import { guestGuard } from "./core/guards/guest.guard";

/**
 * Application routes.
 */
export const routes: Routes = [
  { path: "", redirectTo: "containers", pathMatch: "full" },
  {
    path: "login",
    canActivate: [guestGuard],
    /**
     * @returns {Promise<any>} Login page component
     */
    loadComponent: () =>
      import("./features/auth/login-page/login-page").then((mod) => mod.LoginPageComponent),
  },
  {
    path: "register",
    canActivate: [guestGuard],
    /**
     * @returns {Promise<any>} Register page component
     */
    loadComponent: () =>
      import("./features/auth/register-page/register-page").then(
        (mod) => mod.RegisterPageComponent,
      ),
  },
  {
    path: "oauth/callback",
    /**
     * @returns {Promise<any>} OAuth callback component
     */
    loadComponent: () =>
      import("./features/auth/oauth-callback/oauth-callback").then(
        (mod) => mod.OAuthCallbackComponent,
      ),
  },
  {
    path: "containers",
    canActivate: [authGuard],
    /**
     * @returns {Promise<any>} Container list component
     */
    loadComponent: () =>
      import("./features/containers/container-list/container-list").then(
        (mod) => mod.ContainerListComponent,
      ),
  },
  {
    path: "containers/:id",
    canActivate: [authGuard],
    /**
     * @returns {Promise<any>} Container detail component
     */
    loadComponent: () =>
      import("./features/containers/container-detail/container-detail").then(
        (mod) => mod.ContainerDetailComponent,
      ),
  },
  {
    path: "mods",
    canActivate: [authGuard],
    /**
     * @returns {Promise<any>} Mod list component
     */
    loadComponent: () =>
      import("./features/mods/mod-list/mod-list").then((mod) => mod.ModListComponent),
  },
  {
    path: "game-servers",
    canActivate: [authGuard],
    /**
     * @returns {Promise<any>} Game server list component
     */
    loadComponent: () =>
      import("./features/game-servers/game-server-list/game-server-list").then(
        (mod) => mod.GameServerListComponent,
      ),
  },
  {
    path: "game-servers/:slug",
    canActivate: [authGuard],
    /**
     * @returns {Promise<any>} Game server detail component
     */
    loadComponent: () =>
      import("./features/game-servers/game-server-list/game-server-list").then(
        (mod) => mod.GameServerListComponent,
      ),
  },
];
