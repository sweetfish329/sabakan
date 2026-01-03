import type { Meta, StoryObj } from "@storybook/angular";
import { fn } from "@storybook/test";
import { ContainerCardComponent } from "./container-card";
import type { Container } from "../../../models/container.model";

const mockRunningContainer: Container = {
  id: "abc123def456",
  name: "minecraft-server-1",
  image: "itzg/minecraft-server:latest",
  state: "running",
  status: "Up 2 hours",
  created: "2024-01-01T00:00:00Z",
  ports: [{ hostPort: 25_565, containerPort: 25_565, protocol: "tcp" }],
  labels: { game: "minecraft", type: "paper" },
};

const mockStoppedContainer: Container = {
  id: "xyz789ghi012",
  name: "valheim-server",
  image: "lloesche/valheim-server:latest",
  state: "stopped",
  status: "Exited (0) 3 hours ago",
  created: "2024-01-01T00:00:00Z",
  ports: [
    { hostPort: 2456, containerPort: 2456, protocol: "udp" },
    { hostPort: 2457, containerPort: 2457, protocol: "udp" },
  ],
  labels: { game: "valheim" },
};

const mockPausedContainer: Container = {
  id: "pause123",
  name: "terraria-server",
  image: "ryshe/terraria:latest",
  state: "paused",
  status: "Paused",
  created: "2024-01-01T00:00:00Z",
  ports: [{ hostPort: 7777, containerPort: 7777, protocol: "tcp" }],
  labels: {},
};

const meta: Meta<ContainerCardComponent> = {
  title: "Features/Containers/ContainerCard",
  component: ContainerCardComponent,
  tags: ["autodocs"],
  argTypes: {
    startClicked: { action: "startClicked" },
    stopClicked: { action: "stopClicked" },
    detailsClicked: { action: "detailsClicked" },
  },
  args: {
    startClicked: fn(),
    stopClicked: fn(),
    detailsClicked: fn(),
  },
};

export default meta;
type Story = StoryObj<ContainerCardComponent>;

export const Running: Story = {
  args: {
    container: mockRunningContainer,
    loading: false,
  },
};

export const Stopped: Story = {
  args: {
    container: mockStoppedContainer,
    loading: false,
  },
};

export const Paused: Story = {
  args: {
    container: mockPausedContainer,
    loading: false,
  },
};

export const Loading: Story = {
  args: {
    container: mockRunningContainer,
    loading: true,
  },
};

export const NoName: Story = {
  args: {
    container: {
      ...mockRunningContainer,
      name: "",
    },
    loading: false,
  },
};

export const NoPorts: Story = {
  args: {
    container: {
      ...mockRunningContainer,
      ports: [],
    },
    loading: false,
  },
};
