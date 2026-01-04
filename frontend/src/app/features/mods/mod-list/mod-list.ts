import { Component, inject, signal, type OnInit, ChangeDetectionStrategy } from "@angular/core";
import { RouterLink } from "@angular/router";
import { MatCardModule } from "@angular/material/card";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";

import { ModService, type Mod } from "../../../services/mod.service";

/**
 * Component displaying a list of mods.
 *
 * @example
 * ```html
 * <app-mod-list />
 * ```
 */
@Component({
  selector: "app-mod-list",
  standalone: true,
  imports: [RouterLink, MatCardModule, MatButtonModule, MatIconModule, MatProgressSpinnerModule],
  template: `
		<div class="mod-list">
			<header class="mod-list__header">
				<h1 class="mod-list__title">MOD Catalog</h1>
				<button mat-flat-button routerLink="/mods/new">
					<mat-icon>add</mat-icon>
					Add Mod
				</button>
			</header>

			@if (loading()) {
				<div class="mod-list__loading">
					<mat-spinner diameter="48" />
				</div>
			} @else if (mods().length === 0) {
				<div class="mod-list__empty">
					<mat-icon style="font-size: 48px; width: 48px; height: 48px; opacity: 0.5">
						extension_off
					</mat-icon>
					<p>No mods available</p>
				</div>
			} @else {
				<div class="mod-list__grid">
					@for (mod of mods(); track mod.ID) {
						<article class="mod-card">
							<div class="mod-card__header">
								<h2 class="mod-card__name">{{ mod.name }}</h2>
								@if (mod.version) {
									<span class="mod-card__version">v{{ mod.version }}</span>
								}
							</div>
							<p class="mod-card__description">
								{{ mod.description || "No description" }}
							</p>
							<div class="mod-card__actions">
								<button mat-stroked-button [routerLink]="['/mods', mod.ID]">
									<mat-icon>visibility</mat-icon>
									View
								</button>
								@if (mod.sourceUrl) {
									<a mat-stroked-button [href]="mod.sourceUrl" target="_blank" rel="noopener">
										<mat-icon>link</mat-icon>
										Source
									</a>
								}
							</div>
						</article>
					}
				</div>
			}
		</div>
	`,
  styleUrl: "./mod-list.scss",
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ModListComponent implements OnInit {
  private readonly modService = inject(ModService);

  readonly loading = signal(true);
  readonly mods = signal<Mod[]>([]);
  readonly error = signal<string | null>(null);

  ngOnInit(): void {
    this.loadMods();
  }

  private loadMods(): void {
    this.loading.set(true);
    this.modService.list().subscribe({
      next: (mods) => {
        this.mods.set(mods);
        this.loading.set(false);
      },
      error: (err: Error) => {
        this.error.set(err.message);
        this.loading.set(false);
      },
    });
  }
}
