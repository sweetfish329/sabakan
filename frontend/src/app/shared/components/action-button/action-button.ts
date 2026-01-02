import { ChangeDetectionStrategy, Component, input, output } from "@angular/core";
import { NgTemplateOutlet } from "@angular/common";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";

/**
 * Button variant types.
 * Button size options.
 */
type ActionButtonVariant = "primary" | "secondary" | "danger" | "subtle";
type ActionButtonSize = "small" | "medium" | "large";

/**
 * An action button component with loading state and icon support.
 *
 * @example
 * ```html
 * <app-action-button
 *   label="Start Container"
 *   icon="play_arrow"
 *   variant="primary"
 *   [loading]="isStarting"
 *   (clicked)="onStart()"
 * />
 * ```
 */
@Component({
  selector: "app-action-button",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatButtonModule, MatIconModule, MatProgressSpinnerModule, NgTemplateOutlet],
  template: `
		@switch (variant()) {
			@case ('primary') {
				<button
					mat-flat-button
					color="primary"
					class="action-button"
					[class]="sizeClass()"
					[disabled]="disabled() || loading()"
					(click)="onClick()"
				>
					<ng-container *ngTemplateOutlet="buttonContent" />
				</button>
			}
			@case ('danger') {
				<button
					mat-flat-button
					color="warn"
					class="action-button"
					[class]="sizeClass()"
					[disabled]="disabled() || loading()"
					(click)="onClick()"
				>
					<ng-container *ngTemplateOutlet="buttonContent" />
				</button>
			}
			@case ('secondary') {
				<button
					mat-stroked-button
					color="primary"
					class="action-button"
					[class]="sizeClass()"
					[disabled]="disabled() || loading()"
					(click)="onClick()"
				>
					<ng-container *ngTemplateOutlet="buttonContent" />
				</button>
			}
			@case ('subtle') {
				<button
					mat-button
					class="action-button subtle"
					[class]="sizeClass()"
					[disabled]="disabled() || loading()"
					(click)="onClick()"
				>
					<ng-container *ngTemplateOutlet="buttonContent" />
				</button>
			}
		}

		<ng-template #buttonContent>
			@if (loading()) {
				<mat-spinner diameter="20" class="button-spinner" />
			} @else if (icon()) {
				<mat-icon>{{ icon() }}</mat-icon>
			}
			@if (label()) {
				<span class="button-label">{{ label() }}</span>
			}
		</ng-template>
	`,
  styles: `
		.action-button {
			border-radius: 24px;
			transition: transform 0.2s ease, box-shadow 0.2s ease;

			&:not(:disabled):hover {
				transform: translateY(-1px);
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
			}

			&:not(:disabled):active {
				transform: translateY(0);
			}

			&.small {
				height: 32px;
				padding: 0 12px;
				font-size: 0.8125rem;
			}

			&.medium {
				height: 40px;
				padding: 0 20px;
			}

			&.large {
				height: 48px;
				padding: 0 28px;
				font-size: 1rem;
			}

			&.subtle {
				color: rgba(0, 0, 0, 0.6);
			}
		}

		.button-spinner {
			margin-right: 8px;
		}

		.button-label {
			margin-left: 4px;
		}

		mat-icon {
			margin-right: 4px;
		}
	`,
})
export class ActionButtonComponent {
  /** Button label text */
  readonly label = input<string>();

  /** Material icon name */
  readonly icon = input<string>();

  /** Button variant */
  readonly variant = input<ActionButtonVariant>("primary");

  /** Button size */
  readonly size = input<ActionButtonSize>("medium");

  /** Whether the button is in loading state */
  readonly loading = input(false);

  /** Whether the button is disabled */
  readonly disabled = input(false);

  /** Emitted when the button is clicked */
  readonly clicked = output();

  /**
   * Computed size class
   * @returns {string} The CSS class for the button size
   */
  protected sizeClass(): string {
    return this.size();
  }

  protected onClick(): void {
    if (!this.loading() && !this.disabled()) {
      this.clicked.emit();
    }
  }
}

export type { ActionButtonVariant, ActionButtonSize };
