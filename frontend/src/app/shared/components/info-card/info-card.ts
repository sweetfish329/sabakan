import { ChangeDetectionStrategy, Component, input } from "@angular/core";
import { MatCardModule } from "@angular/material/card";
import { MatIconModule } from "@angular/material/icon";

/**
 * A single info item for display.
 */
export interface InfoItem {
  /** Label for the item */
  label: string;
  /** Value to display */
  value: string;
  /** Optional icon */
  icon?: string;
  /** Whether to use monospace font for the value */
  monospace?: boolean;
}

/**
 * A card component for displaying key-value information.
 *
 * @example
 * ```html
 * <app-info-card
 *   title="Container Info"
 *   [items]="[
 *     { label: 'ID', value: 'abc123', monospace: true },
 *     { label: 'Status', value: 'Running', icon: 'check_circle' },
 *   ]"
 * />
 * ```
 */
@Component({
  selector: "app-info-card",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatCardModule, MatIconModule],
  template: `
		<mat-card class="info-card" [class.compact]="compact()">
			@if (title()) {
				<mat-card-header>
					@if (icon()) {
						<mat-icon mat-card-avatar class="card-icon">{{ icon() }}</mat-icon>
					}
					<mat-card-title>{{ title() }}</mat-card-title>
					@if (subtitle()) {
						<mat-card-subtitle>{{ subtitle() }}</mat-card-subtitle>
					}
				</mat-card-header>
			}
			<mat-card-content>
				<dl class="info-list">
					@for (item of items(); track item.label) {
						<div class="info-item">
							<dt class="info-label">
								@if (item.icon) {
									<mat-icon class="item-icon">{{ item.icon }}</mat-icon>
								}
								{{ item.label }}
							</dt>
							<dd class="info-value" [class.monospace]="item.monospace">
								{{ item.value }}
							</dd>
						</div>
					}
				</dl>
			</mat-card-content>
		</mat-card>
	`,
  styles: `
		.info-card {
			background: linear-gradient(
				135deg,
				rgba(255, 255, 255, 0.95) 0%,
				rgba(255, 255, 255, 0.85) 100%
			);
			backdrop-filter: blur(10px);
			border: 1px solid rgba(0, 0, 0, 0.06);
			transition: transform 0.2s ease, box-shadow 0.2s ease;

			&:hover {
				transform: translateY(-2px);
				box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
			}

			&.compact {
				.info-item {
					padding: 8px 0;
				}
			}
		}

		.card-icon {
			background: linear-gradient(
				135deg,
				rgba(156, 39, 176, 0.1) 0%,
				rgba(103, 58, 183, 0.1) 100%
			);
			border-radius: 50%;
			padding: 8px;
			width: 40px;
			height: 40px;
			font-size: 24px;
		}

		.info-list {
			margin: 0;
			padding: 0;
		}

		.info-item {
			display: flex;
			justify-content: space-between;
			align-items: flex-start;
			padding: 12px 0;
			border-bottom: 1px solid rgba(0, 0, 0, 0.06);

			&:last-child {
				border-bottom: none;
			}
		}

		.info-label {
			display: flex;
			align-items: center;
			gap: 8px;
			font-size: 0.875rem;
			color: rgba(0, 0, 0, 0.6);
			font-weight: 500;
		}

		.item-icon {
			font-size: 18px;
			height: 18px;
			width: 18px;
			color: rgba(0, 0, 0, 0.4);
		}

		.info-value {
			margin: 0;
			font-size: 0.875rem;
			color: rgba(0, 0, 0, 0.87);
			text-align: right;
			word-break: break-word;

			&.monospace {
				font-family: 'Roboto Mono', monospace;
				background: rgba(0, 0, 0, 0.04);
				padding: 2px 8px;
				border-radius: 4px;
			}
		}
	`,
})
export class InfoCardComponent {
  /** Optional card title */
  readonly title = input<string>();

  /** Optional card subtitle */
  readonly subtitle = input<string>();

  /** Optional card icon */
  readonly icon = input<string>();

  /** Items to display */
  readonly items = input.required<InfoItem[]>();

  /** Whether to use compact layout */
  readonly compact = input<boolean>(false);
}
