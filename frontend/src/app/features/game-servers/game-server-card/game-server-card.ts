import { ChangeDetectionStrategy, Component, computed, input, output } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MatCardModule } from "@angular/material/card";
import { MatChipsModule } from "@angular/material/chips";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatTooltipModule } from "@angular/material/tooltip";
import type { GameServer, GameServerStatus } from "../../../models/game-server.model";

/**
 * Displays a single game server as a card with status and controls.
 */
@Component({
  selector: "app-game-server-card",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatChipsModule,
    MatTooltipModule,
    MatProgressSpinnerModule,
  ],
  templateUrl: "./game-server-card.html",
  styleUrl: "./game-server-card.scss",
})
export class GameServerCardComponent {
  /** The game server to display */
  readonly server = input.required<GameServer>();

  /** Whether the card is in a loading state */
  readonly loading = input<boolean>(false);

  /** Emitted when the start button is clicked */
  readonly startClicked = output<string>();

  /** Emitted when the stop button is clicked */
  readonly stopClicked = output<string>();

  /** Emitted when the details button is clicked */
  readonly detailsClicked = output<string>();

  /** Emitted when the delete button is clicked */
  readonly deleteClicked = output<string>();

  /** CSS class based on server status */
  protected readonly statusClass = computed(() => this.server().status);

  /** Icon for the current status */
  protected readonly statusIcon = computed(() => {
    const iconMap: Record<GameServerStatus, string> = {
      running: "play_circle",
      stopped: "stop_circle",
      creating: "hourglass_top",
      error: "error",
    };
    return iconMap[this.server().status] ?? "help";
  });

  /** Whether the server can be started */
  protected readonly canStart = computed(() => {
    const { status } = this.server();
    return status === "stopped" || status === "error";
  });

  /** Whether the server can be stopped */
  protected readonly canStop = computed(() => {
    const { status } = this.server();
    return status === "running";
  });

  /** Game type derived from image name */
  protected readonly gameType = computed(() => {
    const { image } = this.server();
    if (image.includes("minecraft")) {
      return "minecraft";
    }
    if (image.includes("palworld")) {
      return "palworld";
    }
    if (image.includes("ark")) {
      return "ark";
    }
    if (image.includes("rust")) {
      return "rust";
    }
    if (image.includes("factorio")) {
      return "factorio";
    }
    if (image.includes("satisfactory")) {
      return "satisfactory";
    }
    if (image.includes("7daystodie") || image.includes("7dtd")) {
      return "7daystodie";
    }
    return "unknown";
  });

  protected onStart(): void {
    this.startClicked.emit(this.server().slug);
  }

  protected onStop(): void {
    this.stopClicked.emit(this.server().slug);
  }

  protected onViewDetails(): void {
    this.detailsClicked.emit(this.server().slug);
  }

  protected onDelete(): void {
    this.deleteClicked.emit(this.server().slug);
  }
}
