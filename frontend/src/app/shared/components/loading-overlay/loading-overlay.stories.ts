import type { Meta, StoryObj } from "@storybook/angular";
import { LoadingOverlayComponent } from "./loading-overlay";

const meta: Meta<LoadingOverlayComponent> = {
  title: "Shared/LoadingOverlay",
  component: LoadingOverlayComponent,
  tags: ["autodocs"],
  argTypes: {
    visible: {
      control: "boolean",
      description: "Whether the overlay is visible",
    },
    message: {
      control: "text",
      description: "Optional message to display below the spinner",
    },
    fullScreen: {
      control: "boolean",
      description: "Whether to cover the entire screen",
    },
    diameter: {
      control: { type: "number", min: 20, max: 100 },
      description: "Diameter of the spinner in pixels",
    },
  },
  decorators: [
    (story) => ({
      ...story(),
      template: `
				<div style="position: relative; height: 300px; border: 1px dashed #ccc; background: #f5f5f5;">
					<p style="padding: 16px; color: #666;">Content behind the overlay</p>
					${story().template ?? ""}
				</div>
			`,
    }),
  ],
};

export default meta;
type Story = StoryObj<LoadingOverlayComponent>;

export const Default: Story = {
  args: {
    visible: true,
  },
  render: (args) => ({
    props: args,
    template: `<app-loading-overlay [visible]="visible" />`,
  }),
};

export const WithMessage: Story = {
  args: {
    visible: true,
    message: "Loading data...",
  },
  render: (args) => ({
    props: args,
    template: `<app-loading-overlay [visible]="visible" [message]="message" />`,
  }),
};

export const SmallSpinner: Story = {
  args: {
    visible: true,
    message: "Please wait",
    diameter: 32,
  },
  render: (args) => ({
    props: args,
    template: `<app-loading-overlay [visible]="visible" [message]="message" [diameter]="diameter" />`,
  }),
};

export const LargeSpinner: Story = {
  args: {
    visible: true,
    diameter: 72,
  },
  render: (args) => ({
    props: args,
    template: `<app-loading-overlay [visible]="visible" [diameter]="diameter" />`,
  }),
};

export const Hidden: Story = {
  args: {
    visible: false,
  },
  render: (args) => ({
    props: args,
    template: `<app-loading-overlay [visible]="visible" message="This should not be visible" />`,
  }),
};

export const FullScreen: Story = {
  args: {
    visible: true,
    fullScreen: true,
    message: "Loading application...",
  },
  decorators: [],
  render: (args) => ({
    props: args,
    template: `<app-loading-overlay [visible]="visible" [fullScreen]="fullScreen" [message]="message" />`,
  }),
};

export const CustomMessage: Story = {
  args: {
    visible: true,
    message: "Saving changes, please do not close this window...",
  },
  render: (args) => ({
    props: args,
    template: `<app-loading-overlay [visible]="visible" [message]="message" />`,
  }),
};
