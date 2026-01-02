import { ChangeDetectionStrategy, Component, input, output } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";

/**
 * A page header component with title, optional subtitle, back button, and action area.
 *
 * @example
 * ```html
 * <app-page-header
 *   title="Container Details"
 *   subtitle="minecraft-server-1"
 *   [showBack]="true"
 *   (backClicked)="goBack()"
 * >
 *   <button mat-button color="primary">Edit</button>
 * </app-page-header>
 * ```
 */
@Component({
  selector: "app-page-header",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatButtonModule, MatIconModule],
  template: `
		<header class="page-header">
			<div class="header-left">
				@if (showBack()) {
					<button
						mat-icon-button
						class="back-button"
						(click)="backClicked.emit()"
						aria-label="Go back"
					>
						<mat-icon>arrow_back</mat-icon>
					</button>
				}
				<div class="header-titles">
					<h1 class="header-title">{{ title() }}</h1>
					@if (subtitle()) {
						<span class="header-subtitle">{{ subtitle() }}</span>
					}
				</div>
			</div>
			<div class="header-actions">
				<ng-content />
			</div>
		</header>
	`,
  styles: `
		.page-header {
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 16px;
			padding: 16px 24px;
			background: linear-gradient(
				135deg,
				rgba(255, 255, 255, 0.9) 0%,
				rgba(255, 255, 255, 0.7) 100%
			);
			backdrop-filter: blur(10px);
			border-bottom: 1px solid rgba(0, 0, 0, 0.08);
		}

		.header-left {
			display: flex;
			align-items: center;
			gap: 12px;
		}

		.back-button {
			transition: transform 0.2s ease;

			&:hover {
				transform: translateX(-2px);
			}
		}

		.header-titles {
			display: flex;
			flex-direction: column;
			gap: 2px;
		}

		.header-title {
			margin: 0;
			font-size: 1.5rem;
			font-weight: 600;
			color: rgba(0, 0, 0, 0.87);
			line-height: 1.2;
		}

		.header-subtitle {
			font-size: 0.875rem;
			color: rgba(0, 0, 0, 0.6);
		}

		.header-actions {
			display: flex;
			align-items: center;
			gap: 8px;
		}
	`,
})
export class PageHeaderComponent {
  /** Main title text */
  readonly title = input.required<string>();

  /** Optional subtitle text */
  readonly subtitle = input<string>();

  /** Whether to show the back button */
  readonly showBack = input(false);

  /** Emitted when the back button is clicked */
  readonly backClicked = output();
}
