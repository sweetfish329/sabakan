 
 
 
 
 
import type { Meta, StoryObj } from "@storybook/angular";
import type { InfoItem } from "./info-card";
import { InfoCardComponent } from "./info-card";

const containerInfo: InfoItem[] = [
  { label: "Container ID", value: "abc123def456789", monospace: true },
  { label: "Name", value: "minecraft-server-1" },
  { label: "Image", value: "itzg/minecraft-server:latest" },
  { label: "Status", value: "Running", icon: "play_circle" },
  { label: "Created", value: "2024-01-01 12:00:00" },
];

const networkInfo: InfoItem[] = [
  { label: "IP Address", value: "172.17.0.2", monospace: true, icon: "lan" },
  { label: "Gateway", value: "172.17.0.1", monospace: true },
  { label: "Network", value: "bridge" },
  { label: "MAC Address", value: "02:42:ac:11:00:02", monospace: true },
];

const meta: Meta<InfoCardComponent> = {
  title: "Shared/InfoCard",
  component: InfoCardComponent,
  tags: ["autodocs"],
  argTypes: {
    title: {
      control: "text",
      description: "Optional card title",
    },
    subtitle: {
      control: "text",
      description: "Optional card subtitle",
    },
    icon: {
      control: "text",
      description: "Optional card icon",
    },
    items: {
      control: "object",
      description: "Items to display",
    },
    compact: {
      control: "boolean",
      description: "Whether to use compact layout",
    },
  },
};

export default meta;
type Story = StoryObj<InfoCardComponent>;

export const Default: Story = {
  args: {
    items: containerInfo,
  },
};

export const WithTitle: Story = {
  args: {
    title: "Container Information",
    items: containerInfo,
  },
};

export const WithTitleAndIcon: Story = {
  args: {
    title: "Container Details",
    subtitle: "minecraft-server-1",
    icon: "dns",
    items: containerInfo,
  },
};

export const NetworkInfo: Story = {
  args: {
    title: "Network Configuration",
    icon: "settings_ethernet",
    items: networkInfo,
  },
};

export const Compact: Story = {
  args: {
    title: "Quick Info",
    items: [
      { label: "Status", value: "Active" },
      { label: "Uptime", value: "14 days" },
      { label: "Memory", value: "512 MB" },
    ],
    compact: true,
  },
};

export const WithIcons: Story = {
  args: {
    title: "System Status",
    icon: "monitor_heart",
    items: [
      { label: "CPU Usage", value: "45%", icon: "memory" },
      { label: "Memory Usage", value: "2.4 GB / 8 GB", icon: "storage" },
      { label: "Disk Usage", value: "120 GB / 500 GB", icon: "hard_drive" },
      { label: "Network I/O", value: "1.2 MB/s", icon: "network_check" },
    ],
  },
};

export const MixedContent: Story = {
  args: {
    title: "Server Configuration",
    subtitle: "Last updated: 2024-01-01",
    icon: "settings",
    items: [
      { label: "Server Type", value: "Paper 1.20.4" },
      { label: "Max Players", value: "20" },
      { label: "Game Mode", value: "Survival" },
      { label: "Difficulty", value: "Normal", icon: "sports_esports" },
      { label: "World Seed", value: "1234567890", monospace: true },
      {
        label: "Server Address",
        value: "play.example.com:25565",
        monospace: true,
        icon: "link",
      },
    ],
  },
};

export const MultipleCards: Story = {
  render: () => ({
    template: `
			<div style="display: grid; grid-template-columns: repeat(2, 1fr); gap: 16px; max-width: 800px;">
				<app-info-card
					title="Container"
					icon="dns"
					[items]="[
						{ label: 'Name', value: 'minecraft-server' },
						{ label: 'Status', value: 'Running', icon: 'play_circle' },
						{ label: 'Uptime', value: '14 days' }
					]"
				/>
				<app-info-card
					title="Resources"
					icon="memory"
					[items]="[
						{ label: 'CPU', value: '45%' },
						{ label: 'Memory', value: '2.4 GB' },
						{ label: 'Disk', value: '120 GB' }
					]"
				/>
			</div>
		`,
  }),
};
