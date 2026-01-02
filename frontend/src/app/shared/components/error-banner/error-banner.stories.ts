import type { Meta, StoryObj } from "@storybook/angular";
import { fn } from "@storybook/test";
import { ErrorBannerComponent } from "./error-banner";

const meta: Meta<ErrorBannerComponent> = {
  title: "Shared/ErrorBanner",
  component: ErrorBannerComponent,
  tags: ["autodocs"],
  argTypes: {
    message: {
      control: "text",
      description: "The message to display",
    },
    severity: {
      control: "select",
      options: ["error", "warning", "info"],
      description: "Severity level of the banner",
    },
    showRetry: {
      control: "boolean",
      description: "Whether to show the retry button",
    },
    retryLabel: {
      control: "text",
      description: "Label for the retry button",
    },
    showDismiss: {
      control: "boolean",
      description: "Whether to show the dismiss button",
    },
    retryClicked: { action: "retryClicked" },
    dismissClicked: { action: "dismissClicked" },
  },
  args: {
    retryClicked: fn(),
    dismissClicked: fn(),
  },
};

export default meta;
type Story = StoryObj<ErrorBannerComponent>;

export const Error: Story = {
  args: {
    message: "Failed to load containers. Please try again.",
    severity: "error",
    showRetry: true,
    showDismiss: true,
  },
};

export const Warning: Story = {
  args: {
    message: "Your session will expire in 5 minutes.",
    severity: "warning",
    showRetry: false,
    showDismiss: true,
  },
};

export const Info: Story = {
  args: {
    message: "A new version of the application is available.",
    severity: "info",
    showRetry: false,
    showDismiss: true,
  },
};

export const WithRetry: Story = {
  args: {
    message: "Connection lost. Click retry to reconnect.",
    severity: "error",
    showRetry: true,
    retryLabel: "Reconnect",
    showDismiss: false,
  },
};

export const NoDismiss: Story = {
  args: {
    message: "Critical error: Unable to save data.",
    severity: "error",
    showRetry: true,
    showDismiss: false,
  },
};

export const MessageOnly: Story = {
  args: {
    message: "This is an informational message.",
    severity: "info",
    showRetry: false,
    showDismiss: false,
  },
};

export const LongMessage: Story = {
  args: {
    message:
      "An unexpected error occurred while processing your request. The server returned an invalid response. Please try again later or contact support if the problem persists.",
    severity: "error",
    showRetry: true,
    showDismiss: true,
  },
};

export const AllSeverities: Story = {
  render: () => ({
    template: `
			<div style="display: flex; flex-direction: column; gap: 16px;">
				<app-error-banner
					message="Error: Something went wrong"
					severity="error"
					[showRetry]="true"
				/>
				<app-error-banner
					message="Warning: Please review your settings"
					severity="warning"
				/>
				<app-error-banner
					message="Info: Updates are available"
					severity="info"
				/>
			</div>
		`,
  }),
};
