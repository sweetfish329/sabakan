import type { Meta, StoryObj } from "@storybook/angular";
import { RegisterPageComponent } from "./register-page";

const meta: Meta<RegisterPageComponent> = {
  title: "Features/Auth/RegisterPage",
  component: RegisterPageComponent,
  tags: ["autodocs"],
  parameters: {
    layout: "fullscreen",
  },
};

export default meta;
type Story = StoryObj<RegisterPageComponent>;

/**
 * Default register page.
 */
export const Default: Story = {};
