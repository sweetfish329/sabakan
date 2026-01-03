import { Component, ChangeDetectionStrategy, inject, signal } from "@angular/core";
import { Router, RouterLink } from "@angular/router";
import { MatCardModule } from "@angular/material/card";

import {
  RegisterFormComponent,
  type RegisterFormData,
} from "../../../shared/components/register-form/register-form";
import {
  SocialLoginButtonComponent,
  type OAuthProvider,
} from "../../../shared/components/social-login-button/social-login-button";
import { AuthService } from "../../../services/auth.service";

/**
 * Registration page component.
 */
@Component({
  selector: "app-register-page",
  standalone: true,
  imports: [RouterLink, MatCardModule, RegisterFormComponent, SocialLoginButtonComponent],
  template: `
    <div class="register-page">
      <div class="register-page__container">
        <!-- Logo -->
        <div class="register-page__logo">
          <img src="assets/images/SABAKAN-LOGO.png" alt="Sabakan" class="register-page__logo-image" />
          <h1 class="register-page__logo-text">Sabakan</h1>
        </div>

        <!-- Card -->
        <mat-card class="register-page__card">
          <mat-card-header>
            <mat-card-title>Create account</mat-card-title>
            <mat-card-subtitle>Join Sabakan today</mat-card-subtitle>
          </mat-card-header>

          <mat-card-content>
            <app-register-form
              [loading]="loading()"
              [errorMessage]="errorMessage()"
              (submitted)="onRegister($event)"
            />

            <div class="register-page__divider">
              <div class="register-page__divider-line"></div>
              <span>or</span>
              <div class="register-page__divider-line"></div>
            </div>

            <div class="register-page__social">
              <div class="register-page__social-item">
                <app-social-login-button
                  provider="google"
                  [loading]="oauthLoading() === 'google'"
                  (clicked)="onOAuthLogin($event)"
                />
              </div>
              <div class="register-page__social-item">
                <app-social-login-button
                  provider="discord"
                  [loading]="oauthLoading() === 'discord'"
                  (clicked)="onOAuthLogin($event)"
                />
              </div>
            </div>
          </mat-card-content>

          <mat-card-actions>
            <p class="register-page__footer">
              Already have an account?
              <a routerLink="/login" class="register-page__link">Sign in</a>
            </p>
          </mat-card-actions>
        </mat-card>
      </div>

      <!-- Background decorations -->
      <div class="register-page__bg-decoration register-page__bg-decoration--1"></div>
      <div class="register-page__bg-decoration register-page__bg-decoration--2"></div>
    </div>
  `,
  styleUrl: "./register-page.scss",
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegisterPageComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  readonly loading = signal(false);
  readonly errorMessage = signal("");
  // eslint-disable-next-line unicorn/no-null
  readonly oauthLoading = signal<OAuthProvider | null>(null);

  onRegister(data: RegisterFormData): void {
    this.loading.set(true);
    this.errorMessage.set("");

    this.authService.register(data).subscribe({
      next: () => {
        this.loading.set(false);
        // Auto-login after registration
        this.authService.login({ username: data.username, password: data.password }).subscribe({
          next: () => {
            void this.router.navigate(["/"]);
          },
          error: () => {
            void this.router.navigate(["/login"]);
          },
        });
      },
      error: (err: { error?: { message?: string } }) => {
        this.loading.set(false);
        this.errorMessage.set(err.error?.message ?? "Registration failed. Please try again.");
      },
    });
  }

  onOAuthLogin(provider: OAuthProvider): void {
    this.oauthLoading.set(provider);
    globalThis.location.href = this.authService.getOAuthUrl(provider);
  }
}
