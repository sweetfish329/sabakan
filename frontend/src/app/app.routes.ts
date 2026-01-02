import type { Routes } from '@angular/router';

/**
 * Application routes.
 */
export const routes: Routes = [
	{ path: '', redirectTo: 'containers', pathMatch: 'full' },
	{
		path: 'containers',
		loadComponent: () =>
			import('./features/containers/container-list/container-list').then(
				(m) => m.ContainerListComponent,
			),
	},
	{
		path: 'containers/:id',
		loadComponent: () =>
			import('./features/containers/container-detail/container-detail').then(
				(m) => m.ContainerDetailComponent,
			),
	},
];
