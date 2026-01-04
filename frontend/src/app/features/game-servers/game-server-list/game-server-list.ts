import { ChangeDetectionStrategy, Component, computed, inject, signal } from "@angular/core";
import { Router } from "@angular/router";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { MatSnackBar, MatSnackBarModule } from "@angular/material/snack-bar";
import { MatToolbarModule } from "@angular/material/toolbar";
import { MatDialog, MatDialogModule } from "@angular/material/dialog";

import { GameServerCardComponent } from "../game-server-card/game-server-card";
import { GameServerService } from "../../../services/game-server.service";
import type { GameServer } from "../../../models/game-server.model";

/**
 * Displays a list of all game servers with management controls.
 */
@Component({
  selector: "app-game-server-list",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    MatToolbarModule,
    MatButtonModule,
    MatIconModule,
    MatProgressBarModule,
    MatSnackBarModule,
    MatDialogModule,
    GameServerCardComponent,
  ],
  templateUrl: "./game-server-list.html",
  styleUrl: "./game-server-list.scss",
})
export class GameServerListComponent {
  private readonly gameServerService = inject(GameServerService);
  private readonly router = inject(Router);
  private readonly snackBar = inject(MatSnackBar);
  private readonly dialog = inject(MatDialog);

  /** List of game servers */
  readonly servers = signal<GameServer[]>([]);

  /** Whether the list is loading */
  readonly loading = signal(false);

  /** Error message if any */
  // eslint-disable-next-line unicorn/no-null
  readonly error = signal<string | null>(null);

  /** Set of server slugs currently being operated on */
  readonly loadingIds = signal(new Set<string>());

  /** Whether the list is empty */
  readonly isEmpty = computed(() => !this.loading() && this.servers().length === 0);

  constructor() {
    this.loadServers();
  }

  /**
   * Loads the game server list from the API.
   */
  loadServers(): void {
    this.loading.set(true);
    // eslint-disable-next-line unicorn/no-null
    this.error.set(null);

    this.gameServerService.list().subscribe({
      next: (servers) => {
        this.servers.set(servers);
        this.loading.set(false);
      },
      error: (err: unknown) => {
        let message = "Failed to load game servers";
        if (err instanceof Error) {
          ({ message } = err);
        }
        this.error.set(message);
        this.loading.set(false);
      },
    });
  }

  /**
   * Opens the create server dialog.
   */
  async onCreate(): Promise<void> {
    const { GameServerCreateComponent } = await import("../game-server-create/game-server-create");
    const dialogRef = this.dialog.open(GameServerCreateComponent, {
      width: "500px",
    });

    dialogRef.afterClosed().subscribe((result: GameServer | undefined) => {
      if (result) {
        this.snackBar.open("Server created successfully", "Close", { duration: 3000 });
        this.loadServers();
      }
    });
  }

  /**
   * Handles start button click.
   * @param {string} slug - Server slug
   */
  onStart(_slug: string): void {
    this.snackBar.open("Start functionality coming soon", "Close", { duration: 3000 });
  }

  /**
   * Handles stop button click.
   * @param {string} slug - Server slug
   */
  onStop(_slug: string): void {
    this.snackBar.open("Stop functionality coming soon", "Close", { duration: 3000 });
  }

  /**
   * Navigates to server details page.
   * @param {string} slug - Server slug
   */
  onViewDetails(slug: string): void {
    void this.router.navigate(["/game-servers", slug]);
  }

  /**
   * Deletes a game server.
   * @param {string} slug - Server slug
   */
  onDelete(slug: string): void {
    // eslint-disable-next-line no-alert
    if (!confirm(`Are you sure you want to delete "${slug}"?`)) {
      return;
    }

    this.setLoading(slug, true);

    this.gameServerService.delete(slug).subscribe({
      next: () => {
        this.snackBar.open("Server deleted", "Close", { duration: 3000 });
        this.setLoading(slug, false);
        this.loadServers();
      },
      error: (err: unknown) => {
        let message = "Failed to delete server";
        if (err instanceof Error) {
          ({ message } = err);
        }
        this.snackBar.open(message, "Close", { duration: 5000 });
        this.setLoading(slug, false);
      },
    });
  }

  private setLoading(slug: string, isLoading: boolean): void {
    const current = this.loadingIds();
    const updated = new Set(current);
    if (isLoading) {
      updated.add(slug);
    } else {
      updated.delete(slug);
    }
    this.loadingIds.set(updated);
  }
}
