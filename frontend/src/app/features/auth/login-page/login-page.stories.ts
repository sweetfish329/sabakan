import type { Meta, StoryObj } from "@storybook/angular";
import { applicationConfig } from "@storybook/angular";
import { provideRouter } from "@angular/router";
import { provideHttpClient } from "@angular/common/http";
import { provideAnimationsAsync } from "@angular/platform-browser/animations/async";

import { LoginPageComponent } from "./login-page";

const meta: Meta<LoginPageComponent> = {
  title: "Features/Auth/LoginPage",
  component: LoginPageComponent,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
  },
  decorators: [
    applicationConfig({
      providers: [provideRouter([]), provideHttpClient(), provideAnimationsAsync()],
    }),
  ],
};

export default meta;
type Story = StoryObj<LoginPageComponent>;

/**
 * Default login page showing the initial state.
 * Includes social login buttons (Google, Discord) and email/password form.
 */
export const Default: Story = {};
