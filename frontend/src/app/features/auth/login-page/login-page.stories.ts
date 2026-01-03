import type { Meta, StoryObj } from "@storybook/angular";
import { LoginPageComponent } from "./login-page";

const meta: Meta<LoginPageComponent> = {
  title: "Features/Auth/LoginPage",
  component: LoginPageComponent,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
  },
};

export default meta;
type Story = StoryObj<LoginPageComponent>;

/**
 * Default login page.
 */
export const Default: Story = {};
