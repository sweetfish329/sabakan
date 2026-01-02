import type { Meta, StoryObj } from "@storybook/angular";
import { StatusChipComponent } from "./status-chip";

const meta: Meta<StatusChipComponent> = {
  title: "Shared/StatusChip",
  component: StatusChipComponent,
  tags: ["autodocs"],
  argTypes: {
    status: {
      control: "select",
      options: ["success", "warning", "error", "info", "neutral", "running", "stopped", "paused"],
      description: "The status type that determines the chip appearance",
    },
    label: {
      control: "text",
      description: "The label text to display",
    },
    showIcon: {
      control: "boolean",
      description: "Whether to show the status icon",
    },
    customIcon: {
      control: "text",
      description: "Custom Material icon name to override the default",
    },
  },
};

export default meta;
type Story = StoryObj<StatusChipComponent>;

export const Success: Story = {
  args: {
    status: "success",
    label: "Success",
    showIcon: true,
  },
};

export const Warning: Story = {
  args: {
    status: "warning",
    label: "Warning",
    showIcon: true,
  },
};

export const Error: Story = {
  args: {
    status: "error",
    label: "Error",
    showIcon: true,
  },
};

export const Info: Story = {
  args: {
    status: "info",
    label: "Information",
    showIcon: true,
  },
};

export const Neutral: Story = {
  args: {
    status: "neutral",
    label: "Pending",
    showIcon: true,
  },
};

export const Running: Story = {
  args: {
    status: "running",
    label: "Running",
    showIcon: true,
  },
};

export const Stopped: Story = {
  args: {
    status: "stopped",
    label: "Stopped",
    showIcon: true,
  },
};

export const Paused: Story = {
  args: {
    status: "paused",
    label: "Paused",
    showIcon: true,
  },
};

export const WithoutIcon: Story = {
  args: {
    status: "success",
    label: "No Icon",
    showIcon: false,
  },
};

export const CustomIcon: Story = {
  args: {
    status: "info",
    label: "Custom Icon",
    showIcon: true,
    customIcon: "rocket_launch",
  },
};

export const AllStatuses: Story = {
  render: () => ({
    template: `
			<div style="display: flex; flex-wrap: wrap; gap: 8px;">
				<app-status-chip status="success" label="Success" />
				<app-status-chip status="warning" label="Warning" />
				<app-status-chip status="error" label="Error" />
				<app-status-chip status="info" label="Info" />
				<app-status-chip status="neutral" label="Neutral" />
				<app-status-chip status="running" label="Running" />
				<app-status-chip status="stopped" label="Stopped" />
				<app-status-chip status="paused" label="Paused" />
			</div>
		`,
  }),
};
