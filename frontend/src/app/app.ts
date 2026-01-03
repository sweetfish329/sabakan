import { Component } from "@angular/core";
import { RouterOutlet } from "@angular/router";

/**
 * Root component of the application.
 */
@Component({
  selector: "app-root",
  imports: [RouterOutlet],
  template: `<router-outlet />`,
  styles: [],
})
// eslint-disable-next-line @typescript-eslint/no-extraneous-class -- Angular component class
export class App {}
