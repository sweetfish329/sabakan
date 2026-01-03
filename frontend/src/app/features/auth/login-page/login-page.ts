import { Component, ChangeDetectionStrategy, inject, signal } from "@angular/core";
import { Router, RouterLink } from "@angular/router";
import { CommonModule } from "@angular/common";
import { ReactiveFormsModule, FormBuilder, Validators } from "@angular/forms";
import { MatCardModule } from "@angular/material/card";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";

import { AuthService } from "../../../services/auth.service";
import {
  SocialLoginButtonComponent,
  type OAuthProvider,
} from "../../../shared/components/social-login-button/social-login-button";

interface LoginFormData {
  email: string;
  password: string;
}

@Component({
  selector: "app-login-page",
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    RouterLink,
    MatCardModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressBarModule,
    MatProgressSpinnerModule,
    SocialLoginButtonComponent,
  ],
  template: `
    <div class="login-page">
      <div class="login-page__bg-decoration login-page__bg-decoration--1"></div>
      <div class="login-page__bg-decoration login-page__bg-decoration--2"></div>

      <div class="login-page__container">
        <div class="login-page__logo">
          <img src="assets/images/SABAKAN-LOGO.png" alt="Sabakan Logo" class="login-page__logo-image" />
          <h1 class="login-page__logo-text">Sabakan</h1>
        </div>

        <mat-card class="login-page__card">
          <mat-card-header>
            <mat-card-title>Welcome back</mat-card-title>
            <mat-card-subtitle>Sign in to your account to continue</mat-card-subtitle>
          </mat-card-header>

          <mat-card-content>
            @if (errorMessage()) {
              <div class="error-message mb-4 p-3 bg-red-500/10 border border-red-500/20 rounded-lg text-red-400 text-sm flex items-center gap-2">
                <mat-icon class="text-red-400 text-base w-4 h-4">error_outline</mat-icon>
                {{ errorMessage() }}
              </div>
            }

            <div class="login-page__social">
              <div class="login-page__social-item">
                <app-social-login-button
                  provider="google"
                  [loading]="oauthLoading() === 'google'"
                  [disabled]="loading() || (oauthLoading() !== null && oauthLoading() !== 'google')"
                  (clicked)="onOAuthLogin($event)"
                />
              </div>
              <div class="login-page__social-item">
                <app-social-login-button
                  provider="discord"
                  [loading]="oauthLoading() === 'discord'"
                  [disabled]="loading() || (oauthLoading() !== null && oauthLoading() !== 'discord')"
                  (clicked)="onOAuthLogin($event)"
                />
              </div>
            </div>

            <div class="login-page__divider">
              <div class="login-page__divider-line"></div>
              <span>Or continue with</span>
              <div class="login-page__divider-line"></div>
            </div>

            <form [formGroup]="loginForm" (ngSubmit)="onSubmit()">
              <div class="flex flex-col gap-4">
                <mat-form-field appearance="outline" class="w-full">
                  <mat-label>Email Address</mat-label>
                  <input matInput formControlName="email" type="email" placeholder="name@example.com" />
                  <mat-icon matPrefix class="mr-2 text-slate-400">email</mat-icon>
                  @if (loginForm.get('email')?.invalid && loginForm.get('email')?.touched) {
                    <mat-error>Please enter a valid email address</mat-error>
                  }
                </mat-form-field>

                <mat-form-field appearance="outline" class="w-full">
                  <mat-label>Password</mat-label>
                  <input
                    matInput
                    [formControlName]="'password'"
                    [type]="hidePassword() ? 'password' : 'text'"
                    placeholder="Enter your password"
                  />
                  <mat-icon matPrefix class="mr-2 text-slate-400">lock</mat-icon>
                  <button
                    mat-icon-button
                    matSuffix
                    type="button"
                    (click)="hidePassword.set(!hidePassword())"
                    [attr.aria-label]="'Hide password'"
                    [attr.aria-pressed]="hidePassword()"
                  >
                    <mat-icon>{{ hidePassword() ? "visibility_off" : "visibility" }}</mat-icon>
                  </button>
                  @if (loginForm.get('password')?.invalid && loginForm.get('password')?.touched) {
                    <mat-error>Password is required</mat-error>
                  }
                </mat-form-field>

                <button
                  mat-flat-button
                  color="primary"
                  type="submit"
                  class="w-full h-12 !rounded-xl !text-base !font-semibold mt-2"
                  [disabled]="loginForm.invalid || loading() || oauthLoading() !== null"
                >
                  @if (loading()) {
                    <div class="flex items-center gap-2">
                      <mat-spinner diameter="20" class="mr-2"></mat-spinner>
                      Signing in...
                    </div>
                  } @else {
                    Sign in with Email
                  }
                </button>
              </div>
            </form>
          </mat-card-content>

          <mat-card-actions>
            <p class="login-page__footer">
              Don't have an account?
              <a routerLink="/register" class="login-page__link">Sign up now</a>
            </p>
          </mat-card-actions>
          @if (loading() || oauthLoading() !== null) {
            <mat-progress-bar mode="indeterminate" class="absolute bottom-0 left-0 w-full"></mat-progress-bar>
          }
        </mat-card>
      </div>
    </div>
  `,
  styleUrl: "./login-page.scss",
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginPageComponent {
  private readonly fb = inject(FormBuilder);
  private readonly router = inject(Router);
  private readonly authService = inject(AuthService);

  readonly loginForm = this.fb.nonNullable.group({
    email: ["", [Validators.required, Validators.email]],
    password: ["", [Validators.required]],
  });

  readonly hidePassword = signal(true);

  readonly loading = signal(false);
  readonly errorMessage = signal("");
  // eslint-disable-next-line unicorn/no-null
  readonly oauthLoading = signal<OAuthProvider | null>(null);

  onLogin(data: LoginFormData): void {
    this.loading.set(true);
    this.errorMessage.set("");

    this.authService.login({ username: data.email, password: data.password }).subscribe({
      next: () => {
        void this.router.navigate(["/"]);
      },
      error: (error: { message?: string }) => {
        console.error("Login failed", error);
        this.loading.set(false);
        this.errorMessage.set(error.message ?? "Invalid email or password");
      },
    });
  }

  onSubmit(): void {
    if (this.loginForm.valid) {
      const { email, password } = this.loginForm.getRawValue();
      if (email !== "" && password !== "") {
        this.onLogin({ email, password });
      }
    } else {
      this.loginForm.markAllAsTouched();
    }
  }

  onOAuthLogin(provider: OAuthProvider): void {
    this.oauthLoading.set(provider);
    this.errorMessage.set("");

    const targetUrl = this.authService.getOAuthUrl(provider);
    globalThis.location.href = targetUrl;
  }
}
