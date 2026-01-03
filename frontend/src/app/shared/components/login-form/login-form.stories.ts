import type { Meta, StoryObj } from "@storybook/angular";
import { LoginFormComponent } from "./login-form";

const meta: Meta<LoginFormComponent> = {
  title: "Shared/LoginForm",
  component: LoginFormComponent,
  tags: ["autodocs"],
  argTypes: {
    loading: {
      control: "boolean",
      description: "Whether the form is in loading state",
    },
    errorMessage: {
      control: "text",
      description: "Error message to display",
    },
  },
};

export default meta;
type Story = StoryObj<LoginFormComponent>;

/**
 * Default login form state.
 */
export const Default: Story = {
  args: {
    loading: false,
    errorMessage: "",
  },
};

/**
 * Login form in loading state.
 */
export const Loading: Story = {
  args: {
    loading: true,
    errorMessage: "",
  },
};

/**
 * Login form with error message.
 */
export const WithError: Story = {
  args: {
    loading: false,
    errorMessage: "Invalid username or password",
  },
};
