/* eslint-disable import/exports-last */
import { ChangeDetectionStrategy, Component, input, output } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";

/**
 * Severity level for the error banner.
 */
export type ErrorSeverity = "error" | "warning" | "info";

/**
 * Configuration for severity appearance.
 */
interface SeverityConfig {
  icon: string;
  bgColor: string;
  textColor: string;
  borderColor: string;
}

const SEVERITY_CONFIG: Record<ErrorSeverity, SeverityConfig> = {
  error: {
    icon: "error",
    bgColor: "#ffebee",
    textColor: "#c62828",
    borderColor: "#ef9a9a",
  },
  warning: {
    icon: "warning",
    bgColor: "#fff3e0",
    textColor: "#e65100",
    borderColor: "#ffcc80",
  },
  info: {
    icon: "info",
    bgColor: "#e3f2fd",
    textColor: "#1565c0",
    borderColor: "#90caf9",
  },
};

/**
 * A banner component for displaying error, warning, or info messages.
 *
 * @example
 * ```html
 * <app-error-banner
 *   message="Failed to load containers"
 *   severity="error"
 *   [showRetry]="true"
 *   (retryClicked)="onRetry()"
 *   (dismissClicked)="onDismiss()"
 * />
 * ```
 */
@Component({
  selector: "app-error-banner",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatButtonModule, MatIconModule],
  template: `
		<div
			class="error-banner"
			[style.background-color]="config.bgColor"
			[style.color]="config.textColor"
			[style.border-color]="config.borderColor"
			role="alert"
		>
			<mat-icon class="banner-icon" [style.color]="config.textColor">
				{{ config.icon }}
			</mat-icon>
			<span class="banner-message">{{ message() }}</span>
			<div class="banner-actions">
				@if (showRetry()) {
					<button
						mat-button
						class="retry-button"
						[style.color]="config.textColor"
						(click)="retryClicked.emit()"
					>
						<mat-icon>refresh</mat-icon>
						{{ retryLabel() }}
					</button>
				}
				@if (showDismiss()) {
					<button
						mat-icon-button
						class="dismiss-button"
						[style.color]="config.textColor"
						(click)="dismissClicked.emit()"
						aria-label="Dismiss"
					>
						<mat-icon>close</mat-icon>
					</button>
				}
			</div>
		</div>
	`,
  styles: `
		.error-banner {
			display: flex;
			align-items: center;
			gap: 12px;
			padding: 12px 16px;
			border-left: 4px solid;
			border-radius: 4px;
			animation: slideIn 0.3s ease-out;
		}

		.banner-icon {
			flex-shrink: 0;
		}

		.banner-message {
			flex: 1;
			font-size: 0.875rem;
			font-weight: 500;
		}

		.banner-actions {
			display: flex;
			align-items: center;
			gap: 4px;
		}

		.retry-button {
			font-weight: 500;
		}

		.dismiss-button {
			opacity: 0.7;
			transition: opacity 0.2s ease;

			&:hover {
				opacity: 1;
			}
		}

		@keyframes slideIn {
			from {
				opacity: 0;
				transform: translateY(-8px);
			}
			to {
				opacity: 1;
				transform: translateY(0);
			}
		}
	`,
})
export class ErrorBannerComponent {
  /** The message to display */
  readonly message = input.required<string>();

  /** Severity level of the banner */
  readonly severity = input<ErrorSeverity>("error");

  /** Whether to show the retry button */
  readonly showRetry = input(false);

  /** Label for the retry button */
  readonly retryLabel = input<string>("Retry");

  /** Whether to show the dismiss button */
  readonly showDismiss = input(true);

  /** Emitted when the retry button is clicked */
  /** Emitted when the retry button is clicked */
  readonly retryClicked = output();

  /** Emitted when the dismiss button is clicked */
  readonly dismissClicked = output();

  /**
   * Get configuration based on severity
   * @returns The configuration object
   */
  protected get config(): SeverityConfig {
    return SEVERITY_CONFIG[this.severity()];
  }
}
