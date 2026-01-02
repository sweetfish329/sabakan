import {
	ChangeDetectionStrategy,
	Component,
	computed,
	inject,
	signal,
} from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatChipsModule } from '@angular/material/chips';
import { MatDividerModule } from '@angular/material/divider';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { ContainerService } from '../../../services/container.service';
import type { Container, ContainerLogEntry } from '../../../models/container.model';

/**
 * Displays detailed information about a specific container.
 */
@Component({
	selector: 'app-container-detail',
	changeDetection: ChangeDetectionStrategy.OnPush,
	imports: [
		RouterLink,
		MatToolbarModule,
		MatButtonModule,
		MatIconModule,
		MatCardModule,
		MatChipsModule,
		MatDividerModule,
		MatTabsModule,
		MatProgressBarModule,
		MatSnackBarModule,
	],
	template: `
		<div class="container-detail-page">
			<mat-toolbar color="primary" class="page-toolbar">
				<button mat-icon-button routerLink="/containers">
					<mat-icon>arrow_back</mat-icon>
				</button>
				<span>{{ container()?.name || 'Container Details' }}</span>
				<span class="spacer"></span>
				@if (container(); as c) {
					@if (canStart()) {
						<button mat-raised-button color="accent" (click)="onStart()" [disabled]="actionLoading()">
							<mat-icon>play_arrow</mat-icon>
							Start
						</button>
					}
					@if (canStop()) {
						<button mat-raised-button color="warn" (click)="onStop()" [disabled]="actionLoading()">
							<mat-icon>stop</mat-icon>
							Stop
						</button>
					}
				}
			</mat-toolbar>

			@if (loading()) {
				<mat-progress-bar mode="indeterminate"></mat-progress-bar>
			}

			@if (error()) {
				<div class="error-banner">
					<mat-icon>error</mat-icon>
					<span>{{ error() }}</span>
				</div>
			}

			@if (container(); as c) {
				<div class="content">
					<mat-card class="info-card">
						<mat-card-header>
							<mat-card-title>Container Information</mat-card-title>
						</mat-card-header>
						<mat-card-content>
							<div class="info-grid">
								<div class="info-item">
									<span class="label">ID</span>
									<span class="value monospace">{{ c.id }}</span>
								</div>
								<div class="info-item">
									<span class="label">Name</span>
									<span class="value">{{ c.name }}</span>
								</div>
								<div class="info-item">
									<span class="label">Image</span>
									<span class="value">{{ c.image }}</span>
								</div>
								<div class="info-item">
									<span class="label">State</span>
									<mat-chip [class]="'state-chip ' + c.state">{{ c.state }}</mat-chip>
								</div>
								<div class="info-item">
									<span class="label">Status</span>
									<span class="value">{{ c.status }}</span>
								</div>
								<div class="info-item">
									<span class="label">Created</span>
									<span class="value">{{ c.created }}</span>
								</div>
							</div>

							@if (c.ports.length > 0) {
								<mat-divider></mat-divider>
								<h3>Ports</h3>
								<div class="ports-list">
									@for (port of c.ports; track port.containerPort) {
										<span class="port-badge">
											{{ port.hostIp || '0.0.0.0' }}:{{ port.hostPort }} → {{ port.containerPort }}/{{ port.protocol }}
										</span>
									}
								</div>
							}

							@if (Object.keys(c.labels).length > 0) {
								<mat-divider></mat-divider>
								<h3>Labels</h3>
								<div class="labels-list">
									@for (label of labelEntries(); track label[0]) {
										<span class="label-badge">{{ label[0] }}: {{ label[1] }}</span>
									}
								</div>
							}
						</mat-card-content>
					</mat-card>

					<mat-card class="logs-card">
						<mat-card-header>
							<mat-card-title>Logs</mat-card-title>
							<button mat-icon-button (click)="loadLogs()">
								<mat-icon>refresh</mat-icon>
							</button>
						</mat-card-header>
						<mat-card-content>
							@if (logsLoading()) {
								<mat-progress-bar mode="indeterminate"></mat-progress-bar>
							}
							<pre class="logs-content">@for (entry of logs(); track $index) {<span [class]="'log-' + entry.stream">{{ entry.message }}</span>
}@empty {<span class="no-logs">No logs available</span>}</pre>
						</mat-card-content>
					</mat-card>
				</div>
			}
		</div>
	`,
	styles: `
		.container-detail-page {
			min-height: 100vh;
			background-color: #f5f5f5;
		}

		.page-toolbar {
			position: sticky;
			top: 0;
			z-index: 100;
			gap: 16px;
		}

		.spacer {
			flex: 1 1 auto;
		}

		.content {
			padding: 16px;
			display: flex;
			flex-direction: column;
			gap: 16px;
		}

		.info-card,
		.logs-card {
			max-width: 1200px;
		}

		.info-grid {
			display: grid;
			grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
			gap: 16px;
			margin-bottom: 16px;
		}

		.info-item {
			display: flex;
			flex-direction: column;
			gap: 4px;
		}

		.label {
			font-size: 0.75rem;
			color: rgba(0, 0, 0, 0.54);
			text-transform: uppercase;
		}

		.value {
			font-size: 1rem;
		}

		.monospace {
			font-family: monospace;
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
		}

		.ports-list,
		.labels-list {
			display: flex;
			flex-wrap: wrap;
			gap: 8px;
			margin-top: 8px;
		}

		.port-badge {
			background-color: #e8eaf6;
			color: #3f51b5;
			padding: 4px 12px;
			border-radius: 16px;
			font-family: monospace;
			font-size: 0.875rem;
		}

		.label-badge {
			background-color: #fff3e0;
			color: #e65100;
			padding: 4px 12px;
			border-radius: 16px;
			font-size: 0.875rem;
		}

		.logs-card {
			mat-card-header {
				display: flex;
				align-items: center;

				mat-card-title {
					flex: 1;
				}
			}
		}

		.logs-content {
			background-color: #1e1e1e;
			color: #d4d4d4;
			padding: 16px;
			border-radius: 8px;
			overflow-x: auto;
			max-height: 400px;
			overflow-y: auto;
			font-family: monospace;
			font-size: 0.875rem;
			line-height: 1.5;
			margin: 0;
		}

		.log-stdout {
			color: #d4d4d4;
		}

		.log-stderr {
			color: #f48771;
		}

		.no-logs {
			color: #808080;
			font-style: italic;
		}

		.error-banner {
			display: flex;
			align-items: center;
			gap: 12px;
			padding: 12px 16px;
			background-color: #ffebee;
			color: #c62828;
		}

		h3 {
			margin-top: 16px;
			margin-bottom: 8px;
			font-weight: 500;
		}
	`,
})
export class ContainerDetailComponent {
	protected readonly Object = Object;

	private readonly route = inject(ActivatedRoute);
	private readonly containerService = inject(ContainerService);
	private readonly snackBar = inject(MatSnackBar);

	readonly container = signal<Container | null>(null);
	readonly logs = signal<ContainerLogEntry[]>([]);
	readonly loading = signal(false);
	readonly logsLoading = signal(false);
	readonly actionLoading = signal(false);
	readonly error = signal<string | null>(null);

	readonly labelEntries = computed(() => {
		const c = this.container();
		return c ? Object.entries(c.labels) : [];
	});

	readonly canStart = computed(() => {
		const state = this.container()?.state;
		return state === 'stopped' || state === 'exited' || state === 'created';
	});

	readonly canStop = computed(() => {
		const state = this.container()?.state;
		return state === 'running' || state === 'paused';
	});

	constructor() {
		const id = this.route.snapshot.paramMap.get('id');
		if (id) {
			this.loadContainer(id);
			this.loadLogs(id);
		}
	}

	loadContainer(id?: string): void {
		const containerId = id ?? this.route.snapshot.paramMap.get('id');
		if (!containerId) return;

		this.loading.set(true);
		this.error.set(null);

		this.containerService.get(containerId).subscribe({
			next: (container) => {
				this.container.set(container);
				this.loading.set(false);
			},
			error: (err: unknown) => {
				const message =
					err instanceof Error ? err.message : 'Failed to load container';
				this.error.set(message);
				this.loading.set(false);
			},
		});
	}

	loadLogs(id?: string): void {
		const containerId = id ?? this.container()?.id ?? this.route.snapshot.paramMap.get('id');
		if (!containerId) return;

		this.logsLoading.set(true);

		this.containerService.logs(containerId).subscribe({
			next: (logs) => {
				this.logs.set(logs);
				this.logsLoading.set(false);
			},
			error: () => {
				this.logs.set([]);
				this.logsLoading.set(false);
			},
		});
	}

	onStart(): void {
		const id = this.container()?.id;
		if (!id) return;

		this.actionLoading.set(true);

		this.containerService.start(id).subscribe({
			next: () => {
				this.snackBar.open('Container started', 'Close', { duration: 3000 });
				this.actionLoading.set(false);
				this.loadContainer();
			},
			error: (err: unknown) => {
				const message =
					err instanceof Error ? err.message : 'Failed to start container';
				this.snackBar.open(message, 'Close', { duration: 5000 });
				this.actionLoading.set(false);
			},
		});
	}

	onStop(): void {
		const id = this.container()?.id;
		if (!id) return;

		this.actionLoading.set(true);

		this.containerService.stop(id).subscribe({
			next: () => {
				this.snackBar.open('Container stopped', 'Close', { duration: 3000 });
				this.actionLoading.set(false);
				this.loadContainer();
			},
			error: (err: unknown) => {
				const message =
					err instanceof Error ? err.message : 'Failed to stop container';
				this.snackBar.open(message, 'Close', { duration: 5000 });
				this.actionLoading.set(false);
			},
		});
	}
}
