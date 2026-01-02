import type { Meta, StoryObj } from "@storybook/angular";
import { ConfirmDialogPreviewComponent } from "./confirm-dialog";

const meta: Meta<ConfirmDialogPreviewComponent> = {
  title: "Shared/ConfirmDialog",
  component: ConfirmDialogPreviewComponent,
  tags: ["autodocs"],
  argTypes: {
    title: {
      control: "text",
      description: "Dialog title",
    },
    message: {
      control: "text",
      description: "Dialog message",
    },
    confirmLabel: {
      control: "text",
      description: "Confirm button label",
    },
    cancelLabel: {
      control: "text",
      description: "Cancel button label",
    },
    destructive: {
      control: "boolean",
      description: "Whether this is a destructive action",
    },
    icon: {
      control: "text",
      description: "Optional Material icon name",
    },
  },
  decorators: [
    (story) => ({
      ...story(),
      template: `
				<div style="display: flex; justify-content: center; padding: 32px; background: #f5f5f5;">
					${story().template ?? ""}
				</div>
			`,
    }),
  ],
};

export default meta;
type Story = StoryObj<ConfirmDialogPreviewComponent>;

export const Default: Story = {
  args: {
    title: "Confirm Action",
    message: "Are you sure you want to perform this action?",
    confirmLabel: "Confirm",
    cancelLabel: "Cancel",
    destructive: false,
  },
};

export const DeleteConfirmation: Story = {
  args: {
    title: "Delete Container",
    message: "Are you sure you want to delete this container? This action cannot be undone.",
    confirmLabel: "Delete",
    cancelLabel: "Cancel",
    destructive: true,
    icon: "delete",
  },
};

export const StopContainer: Story = {
  args: {
    title: "Stop Container",
    message: "Stopping the container will disconnect all active users. Continue?",
    confirmLabel: "Stop",
    cancelLabel: "Keep Running",
    destructive: true,
    icon: "stop_circle",
  },
};

export const SaveChanges: Story = {
  args: {
    title: "Save Changes",
    message: "Do you want to save your changes before leaving?",
    confirmLabel: "Save",
    cancelLabel: "Discard",
    destructive: false,
    icon: "save",
  },
};

export const Logout: Story = {
  args: {
    title: "Sign Out",
    message: "You will be signed out of your account.",
    confirmLabel: "Sign Out",
    cancelLabel: "Stay Signed In",
    destructive: false,
    icon: "logout",
  },
};

export const RestartServer: Story = {
  args: {
    title: "Restart Server",
    message:
      "Restarting the server will cause a brief downtime. All active connections will be terminated.",
    confirmLabel: "Restart Now",
    cancelLabel: "Cancel",
    destructive: true,
    icon: "restart_alt",
  },
};

export const WithoutIcon: Story = {
  args: {
    title: "Confirm",
    message: "Please confirm your selection.",
    confirmLabel: "OK",
    cancelLabel: "Cancel",
    destructive: false,
  },
};

export const LongMessage: Story = {
  args: {
    title: "Important Notice",
    message:
      "This operation will affect all containers in your environment. All running services will be temporarily unavailable during the migration process. Please ensure you have backed up any important data before proceeding.",
    confirmLabel: "I Understand, Proceed",
    cancelLabel: "Go Back",
    destructive: true,
    icon: "warning",
  },
};
