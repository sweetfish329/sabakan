import { ChangeDetectionStrategy, Component, input, output } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";

/**
 * A component for displaying empty state with icon, title, description, and action.
 *
 * @example
 * ```html
 * <app-empty-state
 *   icon="dns"
 *   title="No Containers Found"
 *   description="There are no containers running on this system."
 *   actionLabel="Create Container"
 *   (actionClicked)="onCreate()"
 * />
 * ```
 */
@Component({
  selector: "app-empty-state",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatButtonModule, MatIconModule],
  template: `
		<div class="empty-state" role="status">
			<div class="empty-icon-container">
				<mat-icon class="empty-icon">{{ icon() }}</mat-icon>
			</div>
			<h2 class="empty-title">{{ title() }}</h2>
			@if (description()) {
				<p class="empty-description">{{ description() }}</p>
			}
			@if (actionLabel()) {
				<button
					mat-flat-button
					color="primary"
					class="empty-action"
					(click)="actionClicked.emit()"
				>
					@if (actionIcon()) {
						<mat-icon>{{ actionIcon() }}</mat-icon>
					}
					{{ actionLabel() }}
				</button>
			}
		</div>
	`,
  styles: `
		.empty-state {
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
			padding: 64px 32px;
			text-align: center;
			animation: fadeInUp 0.4s ease-out;
		}

		.empty-icon-container {
			display: flex;
			align-items: center;
			justify-content: center;
			width: 120px;
			height: 120px;
			border-radius: 50%;
			background: linear-gradient(
				135deg,
				rgba(156, 39, 176, 0.1) 0%,
				rgba(103, 58, 183, 0.1) 100%
			);
			margin-bottom: 24px;
		}

		.empty-icon {
			font-size: 56px;
			height: 56px;
			width: 56px;
			color: rgba(0, 0, 0, 0.38);
		}

		.empty-title {
			margin: 0 0 12px;
			font-size: 1.5rem;
			font-weight: 500;
			color: rgba(0, 0, 0, 0.87);
		}

		.empty-description {
			margin: 0 0 24px;
			font-size: 1rem;
			color: rgba(0, 0, 0, 0.6);
			max-width: 400px;
			line-height: 1.5;
		}

		.empty-action {
			border-radius: 24px;
			padding: 0 24px;
			height: 44px;
		}

		@keyframes fadeInUp {
			from {
				opacity: 0;
				transform: translateY(16px);
			}
			to {
				opacity: 1;
				transform: translateY(0);
			}
		}
	`,
})
export class EmptyStateComponent {
  /** Material icon name to display */
  readonly icon = input<string>("inbox");

  /** Main title text */
  readonly title = input.required<string>();

  /** Optional description text */
  /** Optional description text */
  readonly description = input<string>();

  /** Optional action button label */
  readonly actionLabel = input<string>();

  /** Optional icon for the action button */
  readonly actionIcon = input<string>();

  /** Emitted when the action button is clicked */
  readonly actionClicked = output();
}
