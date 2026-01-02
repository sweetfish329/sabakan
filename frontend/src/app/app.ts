import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

/**
 * Root component of the application.
 */
@Component({
	selector: 'app-root',
	imports: [RouterOutlet],
	template: `<router-outlet />`,
	styles: [],
})
export class App {}
