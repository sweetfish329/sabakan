import type { Routes } from "@angular/router";

/**
 * Application routes.
 */
export const routes: Routes = [
  { path: "", redirectTo: "containers", pathMatch: "full" },
  {
    path: "login",
    loadComponent: () =>
      import("./features/auth/login-page/login-page").then((m) => m.LoginPageComponent),
  },
  {
    path: "register",
    loadComponent: () =>
      import("./features/auth/register-page/register-page").then((m) => m.RegisterPageComponent),
  },
  {
    path: "oauth/callback",
    loadComponent: () =>
      import("./features/auth/oauth-callback/oauth-callback").then((m) => m.OAuthCallbackComponent),
  },
  {
    path: "containers",
    loadComponent: () =>
      import("./features/containers/container-list/container-list").then(
        (m) => m.ContainerListComponent,
      ),
  },
  {
    path: "containers/:id",
    loadComponent: () =>
      import("./features/containers/container-detail/container-detail").then(
        (m) => m.ContainerDetailComponent,
      ),
  },
];
