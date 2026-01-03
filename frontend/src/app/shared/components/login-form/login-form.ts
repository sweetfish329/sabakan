import { Component, input, output, ChangeDetectionStrategy, inject } from "@angular/core";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatInputModule } from "@angular/material/input";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { trigger, transition, style, animate, state } from "@angular/animations";

/**
 * Login form submission data.
 */
export interface LoginFormData {
  username: string;
  password: string;
}

/**
 * Login form component with validation and animations.
 *
 * @example
 * ```html
 * <app-login-form
 *   [loading]="isLoading"
 *   [errorMessage]="error"
 *   (submitted)="onLogin($event)"
 * />
 * ```
 */
@Component({
  selector: "app-login-form",
  standalone: true,
  imports: [
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
  ],
  animations: [
    trigger("fadeSlideIn", [
      transition(":enter", [
        style({ opacity: 0, transform: "translateY(-10px)" }),
        animate("300ms ease-out", style({ opacity: 1, transform: "translateY(0)" })),
      ]),
    ]),
    trigger("shake", [
      state("idle", style({ transform: "translateX(0)" })),
      state("error", style({ transform: "translateX(0)" })),
      transition("idle => error", [
        animate("100ms", style({ transform: "translateX(-5px)" })),
        animate("100ms", style({ transform: "translateX(5px)" })),
        animate("100ms", style({ transform: "translateX(-5px)" })),
        animate("100ms", style({ transform: "translateX(5px)" })),
        animate("100ms", style({ transform: "translateX(0)" })),
      ]),
    ]),
  ],
  template: `
		<form [formGroup]="form" (ngSubmit)="onSubmit()" class="login-form" [@fadeSlideIn]>
			@if (errorMessage()) {
				<div class="login-form__error" [@fadeSlideIn] [@shake]="errorMessage() ? 'error' : 'idle'">
					<mat-icon>error_outline</mat-icon>
					<span>{{ errorMessage() }}</span>
				</div>
			}

			<mat-form-field appearance="outline" class="login-form__field">
				<mat-label>Username</mat-label>
				<mat-icon matPrefix class="login-form__icon">person</mat-icon>
				<input
					matInput
					formControlName="username"
					placeholder="Enter your username"
					autocomplete="username"
				/>
				@if (form.controls.username.hasError('required')) {
					<mat-error>Username is required</mat-error>
				}
			</mat-form-field>

			<mat-form-field appearance="outline" class="login-form__field">
				<mat-label>Password</mat-label>
				<mat-icon matPrefix class="login-form__icon">lock</mat-icon>
				<input
					matInput
					[type]="showPassword ? 'text' : 'password'"
					formControlName="password"
					placeholder="Enter your password"
					autocomplete="current-password"
				/>
				<button
					mat-icon-button
					matSuffix
					type="button"
					(click)="showPassword = !showPassword"
					class="login-form__toggle"
				>
					<mat-icon>{{ showPassword ? 'visibility_off' : 'visibility' }}</mat-icon>
				</button>
				@if (form.controls.password.hasError('required')) {
					<mat-error>Password is required</mat-error>
				}
			</mat-form-field>

			<button
				mat-flat-button
				type="submit"
				class="login-form__submit"
				[disabled]="loading() || form.invalid"
			>
				@if (loading()) {
					<mat-spinner diameter="24" />
				} @else {
					<span class="login-form__submit-text">Sign In</span>
					<mat-icon class="login-form__submit-icon">arrow_forward</mat-icon>
				}
			</button>
		</form>
	`,
  styles: `
		.login-form {
			display: flex;
			flex-direction: column;
			gap: 20px;

			&__field {
				width: 100%;

				// Override Material form field background
				::ng-deep {
					.mdc-text-field--outlined {
						background: rgba(15, 23, 42, 0.6) !important;
						border-radius: 12px !important;

						.mdc-notched-outline__leading,
						.mdc-notched-outline__notch,
						.mdc-notched-outline__trailing {
							border-color: rgba(148, 163, 184, 0.3) !important;
						}

						&:hover .mdc-notched-outline__leading,
						&:hover .mdc-notched-outline__notch,
						&:hover .mdc-notched-outline__trailing {
							border-color: rgba(148, 163, 184, 0.5) !important;
						}

						&.mdc-text-field--focused .mdc-notched-outline__leading,
						&.mdc-text-field--focused .mdc-notched-outline__notch,
						&.mdc-text-field--focused .mdc-notched-outline__trailing {
							border-color: #8b5cf6 !important;
							border-width: 2px !important;
						}
					}

					.mat-mdc-form-field-focus-overlay {
						background: transparent !important;
					}

					.mdc-floating-label {
						color: #94a3b8 !important;
					}

					.mdc-floating-label--float-above {
						color: #a5b4fc !important;
					}

					.mat-mdc-input-element {
						color: #f1f5f9 !important;
						font-size: 15px;
						caret-color: #8b5cf6;

						&::placeholder {
							color: #64748b !important;
						}
					}
				}
			}

			&__icon {
				color: #64748b;
				margin-right: 8px;
			}

			&__toggle {
				color: #64748b;
				transition: color 0.2s ease;

				&:hover {
					color: #94a3b8;
				}
			}

			&__error {
				display: flex;
				align-items: center;
				gap: 10px;
				padding: 14px 18px;
				background: rgba(239, 68, 68, 0.15);
				border: 1px solid rgba(239, 68, 68, 0.4);
				border-radius: 12px;
				color: #fca5a5;
				font-size: 14px;
				font-weight: 500;

				mat-icon {
					font-size: 22px;
					width: 22px;
					height: 22px;
					color: #f87171;
				}
			}

			&__submit {
				height: 56px;
				font-size: 16px;
				font-weight: 600;
				border-radius: 14px;
				margin-top: 8px;
				background: linear-gradient(135deg, #8b5cf6 0%, #7c3aed 100%) !important;
				color: #fff !important;
				transition: all 0.3s ease;
				display: flex;
				align-items: center;
				justify-content: center;
				gap: 8px;
				box-shadow: 0 4px 20px rgba(139, 92, 246, 0.3);

				&:hover:not(:disabled) {
					background: linear-gradient(135deg, #a78bfa 0%, #8b5cf6 100%) !important;
					box-shadow: 0 6px 25px rgba(139, 92, 246, 0.5);
					transform: translateY(-2px);
				}

				&:active:not(:disabled) {
					transform: translateY(0);
				}

				&:disabled {
					background: linear-gradient(135deg, #475569 0%, #334155 100%) !important;
					color: #94a3b8 !important;
					box-shadow: none;
					cursor: not-allowed;
				}

				&-text {
					letter-spacing: 0.5px;
				}

				&-icon {
					font-size: 20px;
					width: 20px;
					height: 20px;
					transition: transform 0.3s ease;
				}

				&:hover:not(:disabled) &-icon {
					transform: translateX(4px);
				}
			}
		}
	`,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginFormComponent {
  private readonly fb = inject(FormBuilder);

  /** Whether the form is in loading state. */
  readonly loading = input(false);

  /** Error message to display. */
  readonly errorMessage = input("");

  /** Emits when the form is submitted with valid data. */
  readonly submitted = output<LoginFormData>();

  protected showPassword = false;

  protected readonly form = this.fb.nonNullable.group({
    username: ["", [Validators.required]],
    password: ["", [Validators.required]],
  });

  protected onSubmit(): void {
    if (this.form.valid) {
      this.submitted.emit(this.form.getRawValue());
    }
  }
}
