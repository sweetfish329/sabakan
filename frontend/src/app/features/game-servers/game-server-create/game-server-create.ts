import { ChangeDetectionStrategy, Component, inject, signal } from "@angular/core";
import { FormBuilder, ReactiveFormsModule, Validators } from "@angular/forms";
import { MatButtonModule } from "@angular/material/button";
import { MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatFormFieldModule } from "@angular/material/form-field";
import { MatIconModule } from "@angular/material/icon";
import { MatInputModule } from "@angular/material/input";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { MatSelectModule } from "@angular/material/select";
import { GameServerService } from "../../../services/game-server.service";
import type { GameServer } from "../../../models/game-server.model";

/**
 * Supported game types.
 */
const GAME_OPTIONS = [
  { value: "minecraft", label: "Minecraft", icon: "sports_esports" },
  { value: "palworld", label: "Palworld", icon: "pets" },
  { value: "ark", label: "ARK: Survival Evolved", icon: "nature" },
  { value: "rust", label: "Rust", icon: "construction" },
  { value: "factorio", label: "Factorio", icon: "precision_manufacturing" },
  { value: "satisfactory", label: "Satisfactory", icon: "factory" },
  { value: "7daystodie", label: "7 Days to Die", icon: "pest_control" },
] as const;

/**
 * Dialog for creating a new game server.
 */
@Component({
  selector: "app-game-server-create",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
  ],
  templateUrl: "./game-server-create.html",
  styleUrl: "./game-server-create.scss",
})
export class GameServerCreateComponent {
  private readonly fb = inject(FormBuilder);
  private readonly gameServerService = inject(GameServerService);
  private readonly dialogRef = inject(MatDialogRef<GameServerCreateComponent>);

  readonly gameOptions = GAME_OPTIONS;
  readonly loading = signal(false);
  // eslint-disable-next-line unicorn/no-null
  readonly error = signal<string | null>(null);

  readonly form = this.fb.group({
    slug: ["", [Validators.required, Validators.pattern(/^[a-z0-9]+(?:-[a-z0-9]+)*$/)]],
    name: ["", [Validators.required, Validators.minLength(3)]],
    game: ["minecraft", Validators.required],
    description: [""],
  });

  /**
   * Generates a slug from the name input.
   */
  generateSlug(): void {
    const name = this.form.get("name")?.value ?? "";
    const slug = name
      .toLowerCase()
      .replaceAll(/[^a-z0-9]+/g, "-")
      .replaceAll(/^-+|-+$/g, "");
    this.form.get("slug")?.setValue(slug);
  }

  /**
   * Submits the form to create a new server.
   */
  onSubmit(): void {
    if (this.form.invalid) {
      return;
    }

    this.loading.set(true);
    // eslint-disable-next-line unicorn/no-null
    this.error.set(null);

    const formDescription = this.form.value.description;
    const data = {
      slug: this.form.value.slug ?? "",
      name: this.form.value.name ?? "",
      game: this.form.value.game ?? "minecraft",
      ...(formDescription ? { description: formDescription } : {}),
    };

    this.gameServerService.create(data).subscribe({
      next: (server: GameServer) => {
        this.loading.set(false);
        this.dialogRef.close(server);
      },
      error: (err: unknown) => {
        this.loading.set(false);
        let message = "Failed to create server";
        if (err instanceof Error) {
          ({ message } = err);
        }
        this.error.set(message);
      },
    });
  }

  onCancel(): void {
    this.dialogRef.close();
  }
}
