/* eslint-disable max-classes-per-file */
/* eslint-disable import/group-exports */
/* eslint-disable import/exports-last */
/* eslint-disable import/order */
/* eslint-disable sort-imports */
import { ChangeDetectionStrategy, Component, inject, input } from "@angular/core";
import { MatButtonModule } from "@angular/material/button";
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatIconModule } from "@angular/material/icon";

/**
 * Data passed to the confirm dialog.
 */
export interface ConfirmDialogData {
  /** Dialog title */
  title: string;
  /** Dialog message */
  message: string;
  /** Confirm button label */
  confirmLabel?: string;
  /** Cancel button label */
  cancelLabel?: string;
  /** Whether this is a destructive action */
  destructive?: boolean;
  /** Optional icon to display */
  icon?: string;
}

/**
 * A reusable confirmation dialog component.
 *
 * @example
 * ```typescript
 * const dialogRef = this.dialog.open(ConfirmDialogComponent, {
 *   data: {
 *     title: 'Delete Container',
 *     message: 'Are you sure you want to delete this container?',
 *     confirmLabel: 'Delete',
 *     destructive: true,
 *   },
 * });
 *
 * dialogRef.afterClosed().subscribe((confirmed) => {
 *   if (confirmed) {
 *     // User confirmed
 *   }
 * });
 * ```
 */
@Component({
  selector: "app-confirm-dialog",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatDialogModule, MatButtonModule, MatIconModule],
  template: `
		<div class="confirm-dialog" [class.destructive]="data.destructive">
			@if (data.icon) {
				<div class="dialog-icon-container" [class.destructive]="data.destructive">
					<mat-icon class="dialog-icon">{{ data.icon }}</mat-icon>
				</div>
			}
			<h2 mat-dialog-title>{{ data.title }}</h2>
			<mat-dialog-content>
				<p class="dialog-message">{{ data.message }}</p>
			</mat-dialog-content>
			<mat-dialog-actions align="end">
				<button mat-button [mat-dialog-close]="false">
					{{ data.cancelLabel || 'Cancel' }}
				</button>
				<button
					mat-flat-button
					[color]="data.destructive ? 'warn' : 'primary'"
					[mat-dialog-close]="true"
				>
					{{ data.confirmLabel || 'Confirm' }}
				</button>
			</mat-dialog-actions>
		</div>
	`,
  styles: `
		.confirm-dialog {
			padding: 8px;
			min-width: 320px;
		}

		.dialog-icon-container {
			display: flex;
			align-items: center;
			justify-content: center;
			width: 56px;
			height: 56px;
			border-radius: 50%;
			background: rgba(33, 150, 243, 0.1);
			margin: 0 auto 16px;

			&.destructive {
				background: rgba(244, 67, 54, 0.1);
			}
		}

		.dialog-icon {
			font-size: 32px;
			height: 32px;
			width: 32px;
			color: #1976d2;

			.destructive & {
				color: #d32f2f;
			}
		}

		.dialog-message {
			margin: 0;
			color: rgba(0, 0, 0, 0.7);
			font-size: 1rem;
			line-height: 1.5;
		}

		mat-dialog-actions {
			margin-top: 16px;
			gap: 8px;
		}
	`,
})
export class ConfirmDialogComponent {
  protected readonly dialogRef = inject(MatDialogRef<ConfirmDialogComponent>);
  protected readonly data = inject<ConfirmDialogData>(MAT_DIALOG_DATA);
}

/**
 * Standalone version for Storybook that doesn't require injection.
 */
@Component({
  selector: "app-confirm-dialog-preview",
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MatButtonModule, MatIconModule],
  template: `
		<div class="confirm-dialog" [class.destructive]="destructive()">
			@if (icon()) {
				<div class="dialog-icon-container" [class.destructive]="destructive()">
					<mat-icon class="dialog-icon">{{ icon() }}</mat-icon>
				</div>
			}
			<h2 class="dialog-title">{{ title() }}</h2>
			<div class="dialog-content">
				<p class="dialog-message">{{ message() }}</p>
			</div>
			<div class="dialog-actions">
				<button mat-button>{{ cancelLabel() }}</button>
				<button mat-flat-button [color]="destructive() ? 'warn' : 'primary'">
					{{ confirmLabel() }}
				</button>
			</div>
		</div>
	`,
  styles: `
		.confirm-dialog {
			padding: 24px;
			min-width: 320px;
			max-width: 400px;
			background: white;
			border-radius: 16px;
			box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
		}

		.dialog-icon-container {
			display: flex;
			align-items: center;
			justify-content: center;
			width: 56px;
			height: 56px;
			border-radius: 50%;
			background: rgba(33, 150, 243, 0.1);
			margin: 0 auto 16px;

			&.destructive {
				background: rgba(244, 67, 54, 0.1);
			}
		}

		.dialog-icon {
			font-size: 32px;
			height: 32px;
			width: 32px;
			color: #1976d2;

			.destructive & {
				color: #d32f2f;
			}
		}

		.dialog-title {
			margin: 0 0 8px;
			font-size: 1.25rem;
			font-weight: 500;
			text-align: center;
		}

		.dialog-content {
			margin-bottom: 24px;
		}

		.dialog-message {
			margin: 0;
			color: rgba(0, 0, 0, 0.7);
			font-size: 1rem;
			line-height: 1.5;
			text-align: center;
		}

		.dialog-actions {
			display: flex;
			justify-content: flex-end;
			gap: 8px;
		}
	`,
})
export class ConfirmDialogPreviewComponent {
  readonly title = input<string>("Confirm Action");
  readonly message = input<string>("Are you sure?");
  readonly confirmLabel = input<string>("Confirm");
  readonly cancelLabel = input<string>("Cancel");
  readonly destructive = input<boolean>(false);
  readonly icon = input<string>();
}
