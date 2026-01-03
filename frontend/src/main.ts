import { bootstrapApplication } from "@angular/platform-browser";
import { appConfig } from "./app/app.config";
import { App } from "./app/app";

/**
 * Bootstraps the Angular application.
 * catches any errors that occur during the bootstrap process.
 */
bootstrapApplication(App, appConfig).catch((error: unknown) => {
  console.error(error);
});
