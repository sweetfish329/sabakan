import type { Meta, StoryObj } from "@storybook/angular";
import { fn } from "@storybook/test";
import { MatButtonModule } from "@angular/material/button";
import { MatIconModule } from "@angular/material/icon";
import { PageHeaderComponent } from "./page-header";

const meta: Meta<PageHeaderComponent> = {
  title: "Shared/PageHeader",
  component: PageHeaderComponent,
  tags: ["autodocs"],
  argTypes: {
    title: {
      control: "text",
      description: "Main title text",
    },
    subtitle: {
      control: "text",
      description: "Optional subtitle text",
    },
    showBack: {
      control: "boolean",
      description: "Whether to show the back button",
    },
    backClicked: { action: "backClicked" },
  },
  args: {
    backClicked: fn(),
  },
};

export default meta;
type Story = StoryObj<PageHeaderComponent>;

export const Default: Story = {
  args: {
    title: "Containers",
  },
};

export const WithSubtitle: Story = {
  args: {
    title: "Container Details",
    subtitle: "minecraft-server-1",
  },
};

export const WithBackButton: Story = {
  args: {
    title: "Edit Container",
    subtitle: "minecraft-server-1",
    showBack: true,
  },
};

export const WithActions: Story = {
  args: {
    title: "Container Details",
    subtitle: "abc123def456",
    showBack: true,
  },
  render: (args) => ({
    props: args,
    moduleMetadata: {
      imports: [PageHeaderComponent, MatButtonModule, MatIconModule],
    },
    template: `
			<app-page-header
				[title]="title"
				[subtitle]="subtitle"
				[showBack]="showBack"
				(backClicked)="backClicked($event)"
			>
				<button mat-button color="primary">
					<mat-icon>edit</mat-icon>
					Edit
				</button>
				<button mat-button color="warn">
					<mat-icon>delete</mat-icon>
					Delete
				</button>
			</app-page-header>
		`,
  }),
};

export const DashboardHeader: Story = {
  args: {
    title: "Dashboard",
  },
  render: (args) => ({
    props: args,
    moduleMetadata: {
      imports: [PageHeaderComponent, MatButtonModule, MatIconModule],
    },
    template: `
			<app-page-header [title]="title">
				<button mat-flat-button color="primary">
					<mat-icon>add</mat-icon>
					New Container
				</button>
			</app-page-header>
		`,
  }),
};

export const SettingsPage: Story = {
  args: {
    title: "Settings",
    subtitle: "Configure your application preferences",
    showBack: true,
  },
  render: (args) => ({
    props: args,
    moduleMetadata: {
      imports: [PageHeaderComponent, MatButtonModule, MatIconModule],
    },
    template: `
			<app-page-header
				[title]="title"
				[subtitle]="subtitle"
				[showBack]="showBack"
				(backClicked)="backClicked($event)"
			>
				<button mat-button>
					<mat-icon>refresh</mat-icon>
					Reset to Defaults
				</button>
				<button mat-flat-button color="primary">
					<mat-icon>save</mat-icon>
					Save Changes
				</button>
			</app-page-header>
		`,
  }),
};

export const LongTitle: Story = {
  args: {
    title: "Very Long Page Title That Might Need to Be Truncated",
    subtitle: "This is a subtitle with additional context information",
    showBack: true,
  },
};
