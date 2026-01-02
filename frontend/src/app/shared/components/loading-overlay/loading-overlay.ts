import { ChangeDetectionStrategy, Component, input } from "@angular/core";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";

const DEFAULT_DIAMETER = 48;

/**
 * A loading overlay component that displays a spinner with optional message.
 *
 * @example
 * ```html
 * <app-loading-overlay [visible]="isLoading" message="Loading data..." />
 * <app-loading-overlay [visible]="isLoading" [fullScreen]="true" />
 * ```
 */
@Component({
  selector: "app-loading-overlay",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatProgressSpinnerModule],
  template: `
		@if (visible()) {
			<div
				class="loading-overlay"
				[class.full-screen]="fullScreen()"
				role="alert"
				aria-busy="true"
				[attr.aria-label]="message() || 'Loading'"
			>
				<div class="loading-content">
					<mat-spinner [diameter]="diameter()" />
					@if (message()) {
						<p class="loading-message">{{ message() }}</p>
					}
				</div>
			</div>
		}
	`,
  styles: `
		.loading-overlay {
			position: absolute;
			inset: 0;
			display: flex;
			align-items: center;
			justify-content: center;
			background: rgba(255, 255, 255, 0.85);
			backdrop-filter: blur(4px);
			z-index: 100;
			animation: fadeIn 0.2s ease-out;

			&.full-screen {
				position: fixed;
				z-index: 9999;
			}
		}

		.loading-content {
			display: flex;
			flex-direction: column;
			align-items: center;
			gap: 16px;
			padding: 32px;
			background: rgba(255, 255, 255, 0.9);
			border-radius: 16px;
			box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
		}

		.loading-message {
			margin: 0;
			color: rgba(0, 0, 0, 0.7);
			font-size: 0.875rem;
			font-weight: 500;
		}

		@keyframes fadeIn {
			from {
				opacity: 0;
			}
			to {
				opacity: 1;
			}
		}
	`,
})
export class LoadingOverlayComponent {
  /** Whether the overlay is visible */
  readonly visible = input<boolean>(false);

  /** Optional message to display below the spinner */
  readonly message = input<string>();

  /** Whether to cover the entire screen */
  readonly fullScreen = input<boolean>(false);

  /** Diameter of the spinner in pixels */
  readonly diameter = input<number>(DEFAULT_DIAMETER);
}
