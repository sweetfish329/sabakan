import { Component, input, output, ChangeDetectionStrategy, inject } from "@angular/core";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";

/**
 * Registration form submission data.
 */
export interface RegisterFormData {
  username: string;
  email: string;
  password: string;
}

/**
 * Registration form component with validation and animations.
 *
 * @example
 * ```html
 * <app-register-form
 *   [loading]="isLoading"
 *   [errorMessage]="error"
 *   (submitted)="onRegister($event)"
 * />
 * ```
 */
@Component({
  selector: "app-register-form",
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
  ],

  template: `
		<form [formGroup]="form" (ngSubmit)="onSubmit()" class="register-form">
			@if (errorMessage()) {
				<div class="register-form__error">
					<mat-icon>error_outline</mat-icon>
					<span>{{ errorMessage() }}</span>
				</div>
			}

			<mat-form-field appearance="outline" class="register-form__field">
				<mat-label>Username</mat-label>
				<mat-icon matPrefix class="register-form__icon">person</mat-icon>
				<input
					matInput
					formControlName="username"
					placeholder="Choose a username"
					autocomplete="username"
				/>
				@if (form.controls.username.hasError('required')) {
					<mat-error>Username is required</mat-error>
				}
				@if (form.controls.username.hasError('minlength')) {
					<mat-error>Username must be at least 3 characters</mat-error>
				}
			</mat-form-field>

			<mat-form-field appearance="outline" class="register-form__field">
				<mat-label>Email</mat-label>
				<mat-icon matPrefix class="register-form__icon">email</mat-icon>
				<input
					matInput
					type="email"
					formControlName="email"
					placeholder="Enter your email"
					autocomplete="email"
				/>
				@if (form.controls.email.hasError('required')) {
					<mat-error>Email is required</mat-error>
				}
				@if (form.controls.email.hasError('email')) {
					<mat-error>Please enter a valid email</mat-error>
				}
			</mat-form-field>

			<mat-form-field appearance="outline" class="register-form__field">
				<mat-label>Password</mat-label>
				<mat-icon matPrefix class="register-form__icon">lock</mat-icon>
				<input
					matInput
					[type]="showPassword ? 'text' : 'password'"
					formControlName="password"
					placeholder="Create a password"
					autocomplete="new-password"
				/>
				<button
					mat-icon-button
					matSuffix
					type="button"
					(click)="showPassword = !showPassword"
					class="register-form__toggle"
				>
					<mat-icon>{{ showPassword ? 'visibility_off' : 'visibility' }}</mat-icon>
				</button>
				@if (form.controls.password.hasError('required')) {
					<mat-error>Password is required</mat-error>
				}
				@if (form.controls.password.hasError('minlength')) {
					<mat-error>Password must be at least 8 characters</mat-error>
				}
			</mat-form-field>

			<mat-form-field appearance="outline" class="register-form__field">
				<mat-label>Confirm Password</mat-label>
				<mat-icon matPrefix class="register-form__icon">lock_outline</mat-icon>
				<input
					matInput
					[type]="showPassword ? 'text' : 'password'"
					formControlName="confirmPassword"
					placeholder="Confirm your password"
					autocomplete="new-password"
				/>
				@if (form.controls.confirmPassword.hasError('required')) {
					<mat-error>Please confirm your password</mat-error>
				}
				@if (form.hasError('passwordMismatch')) {
					<mat-error>Passwords do not match</mat-error>
				}
			</mat-form-field>

			<button
				mat-flat-button
				type="submit"
				class="register-form__submit"
				[disabled]="loading() || form.invalid"
			>
				@if (loading()) {
					<mat-spinner diameter="24" />
				} @else {
					<ng-container>
						<span class="register-form__submit-text">Create Account</span>
						<mat-icon class="register-form__submit-icon" iconPositionEnd>arrow_forward</mat-icon>
					</ng-container>
				}
			</button>
		</form>
	`,
  styleUrl: "./register-form.scss",
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegisterFormComponent {
  private readonly fb = inject(FormBuilder);

  /** Whether the form is in loading state. */
  readonly loading = input(false);

  /** Error message to display. */
  readonly errorMessage = input("");

  /** Emits when the form is submitted with valid data. */
  readonly submitted = output<RegisterFormData>();

  protected showPassword = false;

  protected readonly form = this.fb.nonNullable.group(
    {
      username: ["", [Validators.required, Validators.minLength(3)]],
      email: ["", [Validators.required, Validators.email]],
      password: ["", [Validators.required, Validators.minLength(8)]],
      confirmPassword: ["", [Validators.required]],
    },
    {
      validators: (control) => {
        const password = String(control.get("password")?.value || "");
        const confirmPassword = String(control.get("confirmPassword")?.value || "");
        if (password === confirmPassword) {
          // eslint-disable-next-line unicorn/no-null
          return null;
        }
        return { passwordMismatch: true };
      },
    },
  );

  protected onSubmit(): void {
    if (this.form.valid) {
      const { username, email, password } = this.form.getRawValue();
      this.submitted.emit({ username, email, password });
    }
  }
}
