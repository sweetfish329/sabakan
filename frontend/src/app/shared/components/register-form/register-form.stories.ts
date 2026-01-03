import type { Meta, StoryObj } from "@storybook/angular";
import { RegisterFormComponent } from "./register-form";

const meta: Meta<RegisterFormComponent> = {
  title: "Shared/RegisterForm",
  component: RegisterFormComponent,
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
type Story = StoryObj<RegisterFormComponent>;

/**
 * Default register form state.
 */
export const Default: Story = {
  args: {
    loading: false,
    errorMessage: "",
  },
};

/**
 * Register form in loading state.
 */
export const Loading: Story = {
  args: {
    loading: true,
    errorMessage: "",
  },
};

/**
 * Register form with error message.
 */
export const WithError: Story = {
  args: {
    loading: false,
    errorMessage: "Username already exists",
  },
};
