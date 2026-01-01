import { Component, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';

/**
 * root component of the application.
 */
@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  template: `
    <h1>Welcome to {{ title() }}!</h1>

    <router-outlet />
  `,
  styles: [],
})
export class App {
  /**
   * Title of the application.
   */
  protected readonly title = signal('frontend');
}
