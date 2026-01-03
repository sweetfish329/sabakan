import { ChangeDetectionStrategy, Component, inject, type OnInit } from "@angular/core";
import { Router, ActivatedRoute } from "@angular/router";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { AuthService } from "../../../services/auth.service";

/**
 * OAuth callback handler component.
 * Processes tokens from URL and redirects to home.
 */
@Component({
  selector: "app-oauth-callback",
  standalone: true,
  imports: [MatProgressSpinnerModule],
  template: `
		<div class="oauth-callback">
			<mat-spinner />
			<p>Completing authentication...</p>
		</div>
	`,
  styles: `
		.oauth-callback {
			min-height: 100vh;
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
			gap: 24px;
			background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);

			p {
				font-size: 16px;
				color: rgba(255, 255, 255, 0.7);
			}
		}
	`,
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class OAuthCallbackComponent implements OnInit {
  private readonly authService = inject(AuthService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);

  ngOnInit(): void {
    // Get tokens from URL query params
    this.route.queryParams.subscribe((params) => {
      const { access_token: accessToken, refresh_token: refreshToken, error } = params;

      if (error) {
        console.error("OAuth error:", error);
        this.router.navigate(["/login"], {
          queryParams: { error: "OAuth authentication failed" },
        });
        return;
      }

      if (accessToken && refreshToken) {
        this.authService.storeOAuthTokens(accessToken, refreshToken);
        this.router.navigate(["/"]);
      } else {
        this.router.navigate(["/login"]);
      }
    });
  }
}
