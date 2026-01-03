import { Component, inject, signal, ChangeDetectionStrategy } from "@angular/core";
import { Router, RouterLink } from "@angular/router";
import { MatCardModule } from "@angular/material/card";
import { trigger, transition, style, animate, stagger, query } from "@angular/animations";
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
 * Registration page component with animations.
 */
@Component({
  selector: "app-register-page",
  standalone: true,
  imports: [RouterLink, MatCardModule, RegisterFormComponent, SocialLoginButtonComponent],
  animations: [
    trigger("pageAnimation", [
      transition(":enter", [
        query(
          ".register-page__logo, .register-page__card, .register-page__social-item",
          [
            style({ opacity: 0, transform: "translateY(20px)" }),
            stagger(100, [
              animate("400ms ease-out", style({ opacity: 1, transform: "translateY(0)" })),
            ]),
          ],
          { optional: true },
        ),
      ]),
    ]),
    trigger("floatAnimation", [
      transition(":enter", [
        style({ opacity: 0, transform: "scale(0.8)" }),
        animate(
          "500ms cubic-bezier(0.34, 1.56, 0.64, 1)",
          style({ opacity: 1, transform: "scale(1)" }),
        ),
      ]),
    ]),
  ],
  template: `
		<div class="register-page" [@pageAnimation]>
			<div class="register-page__container">
				<!-- Logo -->
				<div class="register-page__logo" [@floatAnimation]>
					<img
						src="assets/images/SABAKAN-LOGO.png"
						alt="Sabakan"
						class="register-page__logo-image"
					/>
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
  styles: `
		@keyframes float {
			0%,
			100% {
				transform: translateY(0);
			}
			50% {
				transform: translateY(-8px);
			}
		}

		@keyframes pulse {
			0%,
			100% {
				opacity: 0.3;
				transform: scale(1);
			}
			50% {
				opacity: 0.5;
				transform: scale(1.05);
			}
		}

		.register-page {
			min-height: 100vh;
			display: flex;
			align-items: center;
			justify-content: center;
			padding: 24px;
			background: linear-gradient(145deg, #0a0e1a 0%, #111827 50%, #0a0e1a 100%);
			position: relative;
			overflow: hidden;

			&__bg-decoration {
				position: absolute;
				border-radius: 50%;
				filter: blur(80px);
				pointer-events: none;

				&--1 {
					width: 400px;
					height: 400px;
					background: radial-gradient(circle, rgba(139, 92, 246, 0.2) 0%, transparent 70%);
					top: -100px;
					right: -100px;
					animation: pulse 8s ease-in-out infinite;
				}

				&--2 {
					width: 500px;
					height: 500px;
					background: radial-gradient(circle, rgba(59, 130, 246, 0.15) 0%, transparent 70%);
					bottom: -150px;
					left: -150px;
					animation: pulse 10s ease-in-out infinite 2s;
				}
			}

			&__container {
				position: relative;
				z-index: 1;
				display: flex;
				flex-direction: column;
				align-items: center;
				gap: 28px;
				width: 100%;
				max-width: 440px;
			}

			&__logo {
				display: flex;
				align-items: center;
				gap: 14px;
				animation: float 4s ease-in-out infinite;

				&-image {
					width: 56px;
					height: 56px;
					object-fit: contain;
					filter: drop-shadow(0 4px 12px rgba(139, 92, 246, 0.3));
				}

				&-text {
					font-size: 36px;
					font-weight: 700;
					color: #f8fafc;
					letter-spacing: -0.5px;
					margin: 0;
					background: linear-gradient(135deg, #f8fafc 0%, #cbd5e1 100%);
					-webkit-background-clip: text;
					-webkit-text-fill-color: transparent;
					background-clip: text;
				}
			}

			&__card {
				width: 100%;
				padding: 40px;
				border-radius: 24px;
				background: rgba(17, 24, 39, 0.85);
				backdrop-filter: blur(24px);
				border: 1px solid rgba(148, 163, 184, 0.12);
				box-shadow:
					0 4px 6px -1px rgba(0, 0, 0, 0.4),
					0 20px 40px -8px rgba(0, 0, 0, 0.5),
					inset 0 1px 0 rgba(255, 255, 255, 0.05);

				mat-card-header {
					display: block;
					text-align: center;
					margin-bottom: 32px;
				}

				mat-card-title {
					font-size: 28px;
					font-weight: 700;
					color: #f8fafc;
					margin-bottom: 8px;
				}

				mat-card-subtitle {
					font-size: 15px;
					color: #94a3b8;
				}

				mat-card-content {
					padding: 0;
				}

				mat-card-actions {
					padding: 20px 0 0;
					margin: 0;
				}
			}

			&__divider {
				display: flex;
				align-items: center;
				gap: 16px;
				margin: 28px 0;

				&-line {
					flex: 1;
					height: 1px;
					background: linear-gradient(90deg, transparent, rgba(148, 163, 184, 0.25), transparent);
				}

				span {
					font-size: 13px;
					color: #64748b;
					text-transform: uppercase;
					letter-spacing: 1.5px;
					font-weight: 500;
				}
			}

			&__social {
				display: flex;
				flex-direction: column;
				gap: 14px;

				&-item {
					transition: transform 0.2s ease;

					&:hover {
						transform: translateY(-2px);
					}
				}
			}

			&__footer {
				text-align: center;
				font-size: 14px;
				margin: 0;
				color: #94a3b8;
			}

			&__link {
				color: #a78bfa;
				text-decoration: none;
				font-weight: 600;
				transition: all 0.2s ease;
				position: relative;

				&::after {
					content: "";
					position: absolute;
					bottom: -2px;
					left: 0;
					width: 0;
					height: 2px;
					background: #a78bfa;
					transition: width 0.3s ease;
				}

				&:hover {
					color: #c4b5fd;

					&::after {
						width: 100%;
					}
				}
			}
		}
	`,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RegisterPageComponent {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);

  readonly loading = signal(false);
  readonly errorMessage = signal("");
  readonly oauthLoading = signal<OAuthProvider | undefined>(undefined);

  onRegister(data: RegisterFormData): void {
    this.loading.set(true);
    this.errorMessage.set("");

    this.authService.register(data).subscribe({
      next: () => {
        this.loading.set(false);
        // Auto-login after registration
        this.authService.login({ username: data.username, password: data.password }).subscribe({
          next: () => {
            this.router.navigate(["/"]);
          },
          error: () => {
            this.router.navigate(["/login"]);
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
    window.location.href = this.authService.getOAuthUrl(provider);
  }
}
