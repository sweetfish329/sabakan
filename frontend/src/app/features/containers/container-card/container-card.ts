import { ChangeDetectionStrategy, Component, computed, input, output } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatChipsModule } from "@angular/material/chips";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatTooltipModule } from "@angular/material/tooltip";
import type { Container, ContainerState } from "../../../models/container.model";

/**
 * Displays a single container as a card with status and controls.
 */
@Component({
  selector: "app-container-card",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
  ],
  template: `
		<mat-card class="container-card" [class]="stateClass()">
			<mat-card-header>
				<mat-card-title>{{ container().name || container().id.slice(0, 12) }}</mat-card-title>
				<mat-card-subtitle>{{ container().image }}</mat-card-subtitle>
			</mat-card-header>

			<mat-card-content>
				<div class="status-row">
					<mat-chip [class]="'state-chip ' + container().state">
						<mat-icon class="status-icon">{{ stateIcon() }}</mat-icon>
						{{ container().state }}
					</mat-chip>
					<span class="status-text">{{ container().status }}</span>
				</div>

				@if (container().ports.length > 0) {
					<div class="ports-row">
						<mat-icon>lan</mat-icon>
						<span class="ports-text">
							@for (port of container().ports; track port.containerPort) {
								<span class="port-badge">{{ port.hostPort }}:{{ port.containerPort }}/{{ port.protocol }}</span>
							}
						</span>
					</div>
				}
			</mat-card-content>

			<mat-card-actions align="end">
				@if (loading()) {
					<mat-spinner diameter="24"></mat-spinner>
				} @else {
					@if (canStart()) {
						<button mat-icon-button color="primary" matTooltip="Start" (click)="onStart()">
							<mat-icon>play_arrow</mat-icon>
						</button>
					}
					@if (canStop()) {
						<button mat-icon-button color="warn" matTooltip="Stop" (click)="onStop()">
							<mat-icon>stop</mat-icon>
						</button>
					}
					<button mat-icon-button matTooltip="View Details" (click)="onViewDetails()">
						<mat-icon>info</mat-icon>
					</button>
				}
			</mat-card-actions>
		</mat-card>
	`,
  styles: `
		.container-card {
			margin: 8px;
			min-width: 300px;
			background: var(--glass-bg);
			backdrop-filter: blur(var(--glass-blur));
			-webkit-backdrop-filter: blur(var(--glass-blur));
			border: var(--glass-border);
			border-radius: 16px;
			transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
			box-shadow: var(--glass-shadow);

			&:hover {
				transform: translateY(-4px);
				box-shadow: var(--glass-shadow-hover);
			}

			&.running {
				border-left: 4px solid #4caf50;
			}

			&.stopped,
			&.exited {
				border-left: 4px solid #9e9e9e;
			}

			&.paused {
				border-left: 4px solid #ff9800;
			}

			&.restarting {
				border-left: 4px solid #2196f3;
			}
		}

		.status-row {
			display: flex;
			align-items: center;
			gap: 12px;
			margin-bottom: 12px;
		}

		.state-chip {
			text-transform: capitalize;

			&.running {
				background-color: #e8f5e9;
				color: #2e7d32;
			}

			&.stopped,
			&.exited {
				background-color: #f5f5f5;
				color: #616161;
			}

			&.paused {
				background-color: #fff3e0;
				color: #e65100;
			}

			&.restarting {
				background-color: #e3f2fd;
				color: #1565c0;
			}
		}

		.status-icon {
			font-size: 16px;
			height: 16px;
			width: 16px;
			margin-right: 4px;
		}

		.status-text {
			color: rgba(0, 0, 0, 0.6);
			font-size: 0.875rem;
		}

		.ports-row {
			display: flex;
			align-items: center;
			gap: 8px;
			color: rgba(0, 0, 0, 0.6);
		}

		.port-badge {
			background-color: #e8eaf6;
			color: #3f51b5;
			padding: 2px 8px;
			border-radius: 12px;
			font-size: 0.75rem;
			font-family: monospace;
		}
	`,
})
export class ContainerCardComponent {
  /** The container to display */
  readonly container = input.required<Container>();

  /** Whether the card is in a loading state */
  readonly loading = input<boolean>(false);

  /** Emitted when the start button is clicked */
  readonly startClicked = output<string>();

  /** Emitted when the stop button is clicked */
  readonly stopClicked = output<string>();

  /** Emitted when the details button is clicked */
  readonly detailsClicked = output<string>();

  /** CSS class based on container state */
  protected readonly stateClass = computed(() => this.container().state);

  /** Icon for the current state */
  protected readonly stateIcon = computed(() => {
    const iconMap: Record<ContainerState, string> = {
      running: "play_circle",
      stopped: "stop_circle",
      created: "add_circle",
      paused: "pause_circle",
      restarting: "refresh",
      exited: "cancel",
      unknown: "help",
    };
    return iconMap[this.container().state] ?? "help";
  });

  /** Whether the container can be started */
  protected readonly canStart = computed(() => {
    const state = this.container().state;
    return state === "stopped" || state === "exited" || state === "created";
  });

  /** Whether the container can be stopped */
  protected readonly canStop = computed(() => {
    const state = this.container().state;
    return state === "running" || state === "paused";
  });

  protected onStart(): void {
    this.startClicked.emit(this.container().id);
  }

  protected onStop(): void {
    this.stopClicked.emit(this.container().id);
  }

  protected onViewDetails(): void {
    this.detailsClicked.emit(this.container().id);
  }
}
