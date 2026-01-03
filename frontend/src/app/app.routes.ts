import type { Routes } from "@angular/router";

/**
 * Application routes.
 */
export const routes: Routes = [
  { path: "", redirectTo: "containers", pathMatch: "full" },
  {
    path: "login",
    /**
     * @returns {Promise<any>} Login page component
     */
    loadComponent: () =>
      import("./features/auth/login-page/login-page").then((mod) => mod.LoginPageComponent),
  },
  {
    path: "register",
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
    /**
     * @returns {Promise<any>} Container detail component
     */
    loadComponent: () =>
      import("./features/containers/container-detail/container-detail").then(
        (mod) => mod.ContainerDetailComponent,
      ),
  },
];
