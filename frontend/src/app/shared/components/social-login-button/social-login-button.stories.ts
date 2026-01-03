import type { Meta, StoryObj } from "@storybook/angular";
import { SocialLoginButtonComponent } from "./social-login-button";

const meta: Meta<SocialLoginButtonComponent> = {
  title: "Shared/SocialLoginButton",
  component: SocialLoginButtonComponent,
  tags: ["autodocs"],
  argTypes: {
    provider: {
      control: "select",
      options: ["google", "discord"],
      description: "OAuth provider",
    },
    disabled: {
      control: "boolean",
      description: "Whether the button is disabled",
    },
    loading: {
      control: "boolean",
      description: "Whether the button shows loading state",
    },
  },
};

export default meta;
type Story = StoryObj<SocialLoginButtonComponent>;

/**
 * Google login button in default state.
 */
export const Google: Story = {
  args: {
    provider: "google",
    disabled: false,
    loading: false,
  },
};

/**
 * Discord login button in default state.
 */
export const Discord: Story = {
  args: {
    provider: "discord",
    disabled: false,
    loading: false,
  },
};

/**
 * Google login button in loading state.
 */
export const GoogleLoading: Story = {
  args: {
    provider: "google",
    disabled: false,
    loading: true,
  },
};

/**
 * Discord login button in loading state.
 */
export const DiscordLoading: Story = {
  args: {
    provider: "discord",
    disabled: false,
    loading: true,
  },
};

/**
 * Disabled button state.
 */
export const Disabled: Story = {
  args: {
    provider: "google",
    disabled: true,
    loading: false,
  },
};
