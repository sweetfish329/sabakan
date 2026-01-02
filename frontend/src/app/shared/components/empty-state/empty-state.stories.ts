import type { Meta, StoryObj } from "@storybook/angular";
import { fn } from "@storybook/test";
import { EmptyStateComponent } from "./empty-state";

const meta: Meta<EmptyStateComponent> = {
  title: "Shared/EmptyState",
  component: EmptyStateComponent,
  tags: ["autodocs"],
  argTypes: {
    icon: {
      control: "text",
      description: "Material icon name to display",
    },
    title: {
      control: "text",
      description: "Main title text",
    },
    description: {
      control: "text",
      description: "Optional description text",
    },
    actionLabel: {
      control: "text",
      description: "Optional action button label",
    },
    actionIcon: {
      control: "text",
      description: "Optional icon for the action button",
    },
    actionClicked: { action: "actionClicked" },
  },
  args: {
    actionClicked: fn(),
  },
};

export default meta;
type Story = StoryObj<EmptyStateComponent>;

export const Default: Story = {
  args: {
    icon: "inbox",
    title: "No Items Found",
  },
};

export const WithDescription: Story = {
  args: {
    icon: "dns",
    title: "No Containers",
    description: "There are no containers running on this system.",
  },
};

export const WithAction: Story = {
  args: {
    icon: "folder_open",
    title: "No Projects Yet",
    description: "Create your first project to get started.",
    actionLabel: "Create Project",
    actionIcon: "add",
  },
};

export const NoContainers: Story = {
  args: {
    icon: "dns",
    title: "No Containers Found",
    description:
      "You have not created any containers yet. Start by pulling an image and creating a new container.",
    actionLabel: "Create Container",
    actionIcon: "add_circle",
  },
};

export const SearchNoResults: Story = {
  args: {
    icon: "search_off",
    title: "No Results",
    description: "No items match your search criteria. Try adjusting your filters.",
    actionLabel: "Clear Filters",
    actionIcon: "filter_list_off",
  },
};

export const ErrorState: Story = {
  args: {
    icon: "cloud_off",
    title: "Connection Lost",
    description: "Unable to connect to the server. Please check your internet connection.",
    actionLabel: "Retry",
    actionIcon: "refresh",
  },
};

export const NoNotifications: Story = {
  args: {
    icon: "notifications_off",
    title: "All Caught Up!",
    description: "You have no new notifications.",
  },
};

export const CustomIcon: Story = {
  args: {
    icon: "rocket_launch",
    title: "Ready for Launch",
    description: "Your application is configured and ready to deploy.",
    actionLabel: "Deploy Now",
    actionIcon: "cloud_upload",
  },
};
