import type { Meta, StoryObj } from "@storybook/angular";
import { fn } from "@storybook/test";
import { ActionButtonComponent } from "./action-button";

const meta: Meta<ActionButtonComponent> = {
  title: "Shared/ActionButton",
  component: ActionButtonComponent,
  tags: ["autodocs"],
  argTypes: {
    label: {
      control: "text",
      description: "Button label text",
    },
    icon: {
      control: "text",
      description: "Material icon name",
    },
    variant: {
      control: "select",
      options: ["primary", "secondary", "danger", "subtle"],
      description: "Button variant",
    },
    size: {
      control: "select",
      options: ["small", "medium", "large"],
      description: "Button size",
    },
    loading: {
      control: "boolean",
      description: "Whether the button is in loading state",
    },
    disabled: {
      control: "boolean",
      description: "Whether the button is disabled",
    },
    clicked: { action: "clicked" },
  },
  args: {
    clicked: fn(),
  },
};

export default meta;
type Story = StoryObj<ActionButtonComponent>;

export const Primary: Story = {
  args: {
    label: "Primary Action",
    icon: "check",
    variant: "primary",
    size: "medium",
    loading: false,
    disabled: false,
  },
};

export const Secondary: Story = {
  args: {
    label: "Secondary Action",
    icon: "edit",
    variant: "secondary",
    size: "medium",
  },
};

export const Danger: Story = {
  args: {
    label: "Delete",
    icon: "delete",
    variant: "danger",
    size: "medium",
  },
};

export const Subtle: Story = {
  args: {
    label: "Cancel",
    variant: "subtle",
    size: "medium",
  },
};

export const Loading: Story = {
  args: {
    label: "Saving...",
    variant: "primary",
    size: "medium",
    loading: true,
  },
};

export const Disabled: Story = {
  args: {
    label: "Disabled",
    icon: "block",
    variant: "primary",
    size: "medium",
    disabled: true,
  },
};

export const IconOnly: Story = {
  args: {
    icon: "play_arrow",
    variant: "primary",
    size: "medium",
  },
};

export const Small: Story = {
  args: {
    label: "Small",
    icon: "add",
    variant: "primary",
    size: "small",
  },
};

export const Large: Story = {
  args: {
    label: "Large Action",
    icon: "rocket_launch",
    variant: "primary",
    size: "large",
  },
};

export const StartContainer: Story = {
  args: {
    label: "Start",
    icon: "play_arrow",
    variant: "primary",
    size: "medium",
  },
};

export const StopContainer: Story = {
  args: {
    label: "Stop",
    icon: "stop",
    variant: "danger",
    size: "medium",
  },
};

export const AllVariants: Story = {
  render: () => ({
    template: `
			<div style="display: flex; flex-wrap: wrap; gap: 16px; align-items: center;">
				<app-action-button label="Primary" icon="check" variant="primary" />
				<app-action-button label="Secondary" icon="edit" variant="secondary" />
				<app-action-button label="Danger" icon="delete" variant="danger" />
				<app-action-button label="Subtle" variant="subtle" />
			</div>
		`,
  }),
};

export const AllSizes: Story = {
  render: () => ({
    template: `
			<div style="display: flex; flex-wrap: wrap; gap: 16px; align-items: center;">
				<app-action-button label="Small" icon="add" size="small" />
				<app-action-button label="Medium" icon="add" size="medium" />
				<app-action-button label="Large" icon="add" size="large" />
			</div>
		`,
  }),
};

export const LoadingStates: Story = {
  render: () => ({
    template: `
			<div style="display: flex; flex-wrap: wrap; gap: 16px; align-items: center;">
				<app-action-button label="Loading Primary" variant="primary" [loading]="true" />
				<app-action-button label="Loading Secondary" variant="secondary" [loading]="true" />
				<app-action-button label="Loading Danger" variant="danger" [loading]="true" />
			</div>
		`,
  }),
};
