/* eslint-disable import/exports-last */
import { ChangeDetectionStrategy, Component, computed, input } from "@angular/core";
import { MatChipsModule } from "@angular/material/chips";
import { MatIconModule } from "@angular/material/icon";

/**
 * Available status types for the chip.
 */
export type StatusType =
  | "success"
  | "warning"
  | "error"
  | "info"
  | "neutral"
  | "running"
  | "stopped"
  | "paused";

/**
 * Configuration for status appearance.
 */
interface StatusConfig {
  /** Icon to display */
  icon: string;
  /** Background color */
  bgColor: string;
  /** Text color */
  textColor: string;
}

/**
 * Mapping of status types to their visual configuration.
 */
const STATUS_CONFIG: Record<StatusType, StatusConfig> = {
  success: {
    icon: "check_circle",
    bgColor: "rgba(76, 175, 80, 0.15)",
    textColor: "#2e7d32",
  },
  warning: {
    icon: "warning",
    bgColor: "rgba(255, 152, 0, 0.15)",
    textColor: "#e65100",
  },
  error: {
    icon: "error",
    bgColor: "rgba(244, 67, 54, 0.15)",
    textColor: "#c62828",
  },
  info: {
    icon: "info",
    bgColor: "rgba(33, 150, 243, 0.15)",
    textColor: "#1565c0",
  },
  neutral: {
    icon: "radio_button_unchecked",
    bgColor: "rgba(158, 158, 158, 0.15)",
    textColor: "#616161",
  },
  running: {
    icon: "play_circle",
    bgColor: "rgba(76, 175, 80, 0.15)",
    textColor: "#2e7d32",
  },
  stopped: {
    icon: "stop_circle",
    bgColor: "rgba(158, 158, 158, 0.15)",
    textColor: "#616161",
  },
  paused: {
    icon: "pause_circle",
    bgColor: "rgba(255, 152, 0, 0.15)",
    textColor: "#e65100",
  },
};

/**
 * A chip component for displaying status with consistent styling.
 *
 * @example
 * ```html
 * <app-status-chip status="running" label="Running" />
 * <app-status-chip status="error" label="Failed" [showIcon]="true" />
 * ```
 */
@Component({
  selector: "app-status-chip",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatChipsModule, MatIconModule],
  template: `
		<mat-chip
			class="status-chip"
			[style.background-color]="config().bgColor"
			[style.color]="config().textColor"
		>
			@if (showIcon()) {
				<mat-icon class="status-icon" [style.color]="config().textColor">
					{{ customIcon() || config().icon }}
				</mat-icon>
			}
			<span class="status-label">{{ label() }}</span>
		</mat-chip>
	`,
  styles: `
		.status-chip {
			font-weight: 500;
			text-transform: capitalize;
			transition: transform 0.2s ease, box-shadow 0.2s ease;

			&:hover {
				transform: scale(1.02);
				box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
			}
		}

		.status-icon {
			font-size: 16px;
			height: 16px;
			width: 16px;
			margin-right: 4px;
		}

		.status-label {
			font-size: 0.875rem;
		}
	`,
})
export class StatusChipComponent {
  /** The status type that determines the chip's appearance */
  readonly status = input.required<StatusType>();

  /** The label text to display */
  readonly label = input.required<string>();

  /** Whether to show the status icon */
  readonly showIcon = input<boolean>(true);

  /** Custom icon to override the default status icon */
  readonly customIcon = input<string>();

  /** Computed configuration based on status */
  protected readonly config = computed(() => STATUS_CONFIG[this.status()] ?? STATUS_CONFIG.neutral);
}
