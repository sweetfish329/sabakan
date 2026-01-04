import { applicationConfig } from "@storybook/angular";
import type { Meta, StoryObj } from "@storybook/angular";
import { provideAnimations } from "@angular/platform-browser/animations";
import { GameServerCardComponent } from "./game-server-card";
import type { GameServer } from "../../../models/game-server.model";

const mockMinecraftServer: GameServer = {
  id: 1,
  slug: "my-minecraft",
  name: "My Minecraft Server",
  description: "A fun survival server for friends",
  image: "itzg/minecraft-server:latest",
  status: "running",
  containerId: "abc123",
  ownerId: 1,
  ports: [
    { hostPort: 25_565, containerPort: 25_565, protocol: "tcp" },
    { hostPort: 25_575, containerPort: 25_575, protocol: "tcp" },
  ],
  envs: [
    { key: "EULA", value: "TRUE", isSecret: false },
    { key: "MEMORY", value: "4G", isSecret: false },
  ],
};

const mockStoppedServer: GameServer = {
  id: 2,
  slug: "palworld-server",
  name: "Palworld Adventure",
  description: "Catch them all!",
  image: "thijsvanloef/palworld-server-docker:latest",
  status: "stopped",
  ownerId: 1,
  ports: [{ hostPort: 8211, containerPort: 8211, protocol: "udp" }],
  envs: [],
};

const mockErrorServer: GameServer = {
  id: 3,
  slug: "broken-server",
  name: "Broken Server",
  description: "Something went wrong",
  image: "unknown:latest",
  status: "error",
  ownerId: 1,
  ports: [],
  envs: [],
};

const mockCreatingServer: GameServer = {
  id: 4,
  slug: "new-ark-server",
  name: "ARK Survival",
  image: "hermsi/ark-server:latest",
  status: "creating",
  ownerId: 1,
  ports: [{ hostPort: 27_015, containerPort: 27_015, protocol: "udp" }],
  envs: [],
};

const meta: Meta<GameServerCardComponent> = {
  title: "Features/GameServers/GameServerCard",
  component: GameServerCardComponent,
  tags: ["autodocs"],
  decorators: [
    applicationConfig({
      providers: [provideAnimations()],
    }),
  ],
  argTypes: {
    startClicked: { action: "startClicked" },
    stopClicked: { action: "stopClicked" },
    detailsClicked: { action: "detailsClicked" },
    deleteClicked: { action: "deleteClicked" },
  },
};

export default meta;
type Story = StoryObj<GameServerCardComponent>;

/**
 * Running Minecraft server with ports and description.
 */
export const Running: Story = {
  args: {
    server: mockMinecraftServer,
    loading: false,
  },
};

/**
 * Stopped Palworld server.
 */
export const Stopped: Story = {
  args: {
    server: mockStoppedServer,
    loading: false,
  },
};

/**
 * Server in error state.
 */
export const ErrorState: Story = {
  args: {
    server: mockErrorServer,
    loading: false,
  },
};

/**
 * Server being created.
 */
export const Creating: Story = {
  args: {
    server: mockCreatingServer,
    loading: false,
  },
};

/**
 * Card in loading state.
 */
export const Loading: Story = {
  args: {
    server: mockMinecraftServer,
    loading: true,
  },
};
