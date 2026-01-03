import { ChangeDetectionStrategy, Component, computed, inject, signal } from "@angular/core";
import { Router } from "@angular/router";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { MatSnackBar, MatSnackBarModule } from "@angular/material/snack-bar";
import { MatToolbarModule } from "@angular/material/toolbar";

import { ContainerCardComponent } from "../container-card/container-card";
import { ContainerService } from "../../../services/container.service";
import type { Container } from "../../../models/container.model";

/**
 * Displays a list of all containers with management controls.
 */
@Component({
  selector: "app-container-list",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatProgressBarModule,
    MatSnackBarModule,
    ContainerCardComponent,
  ],
  template: `
    <div class="container-list-page">
      <mat-toolbar color="primary" class="page-toolbar">
        <span>Containers</span>
        <span class="spacer"></span>
        <button mat-icon-button (click)="loadContainers()" [disabled]="loading()">
          <mat-icon>refresh</mat-icon>
        </button>
      </mat-toolbar>

      @if (loading()) {
        <mat-progress-bar mode="indeterminate"></mat-progress-bar>
      }

      @if (error()) {
        <div class="error-banner">
          <mat-icon>error</mat-icon>
          <span>{{ error() }}</span>
          <button mat-button (click)="loadContainers()">Retry</button>
        </div>
      }

      <div class="container-grid">
        @for (container of containers(); track container.id) {
          <app-container-card
            [container]="container"
            [loading]="loadingIds().has(container.id)"
            (startClicked)="onStart($event)"
            (stopClicked)="onStop($event)"
            (detailsClicked)="onViewDetails($event)"
          />
        } @empty {
          @if (!loading()) {
            <div class="empty-state">
              <mat-icon class="empty-icon">dns</mat-icon>
              <h2>No Containers Found</h2>
              <p>There are no containers running on this system.</p>
            </div>
          }
        }
      </div>
    </div>
  `,
  styleUrl: "./container-list.scss",
})
export class ContainerListComponent {
  private readonly containerService = inject(ContainerService);
  private readonly router = inject(Router);
  private readonly snackBar = inject(MatSnackBar);

  /** List of containers */
  readonly containers = signal<Container[]>([]);

  /** Whether the list is loading */
  readonly loading = signal(false);

  /** Error message if any */
  // eslint-disable-next-line unicorn/no-null
  readonly error = signal<string | null>(null);

  /** Set of container IDs currently being operated on */
  readonly loadingIds = signal(new Set<string>());

  /** Whether the list is empty */
  readonly isEmpty = computed(() => !this.loading() && this.containers().length === 0);

  constructor() {
    this.loadContainers();
  }

  /**
   * Loads the container list from the API.
   */
  loadContainers(): void {
    this.loading.set(true);
    // eslint-disable-next-line unicorn/no-null
    this.error.set(null);

    this.containerService.list().subscribe({
      next: (containers) => {
        this.containers.set(containers);
        this.loading.set(false);
      },
      error: (err: unknown) => {
        let message = "Failed to load containers";
        if (err instanceof Error) {
          ({ message } = err);
        }
        this.error.set(message);
        this.loading.set(false);
      },
    });
  }

  /**
   * Starts a container.
   * @param {string} id - Container ID
   */
  onStart(id: string): void {
    this.setLoading(id, true);

    this.containerService.start(id).subscribe({
      next: () => {
        this.snackBar.open("Container started", "Close", { duration: 3000 });
        this.setLoading(id, false);
        this.loadContainers();
      },
      error: (err: unknown) => {
        let message = "Failed to start container";
        if (err instanceof Error) {
          ({ message } = err);
        }
        this.snackBar.open(message, "Close", { duration: 5000 });
        this.setLoading(id, false);
      },
    });
  }

  /**
   * Stops a container.
   * @param {string} id - Container ID
   */
  onStop(id: string): void {
    this.setLoading(id, true);

    this.containerService.stop(id).subscribe({
      next: () => {
        this.snackBar.open("Container stopped", "Close", { duration: 3000 });
        this.setLoading(id, false);
        this.loadContainers();
      },
      error: (err: unknown) => {
        let message = "Failed to stop container";
        if (err instanceof Error) {
          ({ message } = err);
        }
        this.snackBar.open(message, "Close", { duration: 5000 });
        this.setLoading(id, false);
      },
    });
  }

  /**
   * Navigates to container details page.
   * @param {string} id - Container ID
   */
  onViewDetails(id: string): void {
    void this.router.navigate(["/containers", id]);
  }

  private setLoading(id: string, isLoading: boolean): void {
    const current = this.loadingIds();
    const updated = new Set(current);
    if (isLoading) {
      updated.add(id);
    } else {
      updated.delete(id);
    }
    this.loadingIds.set(updated);
  }
}
