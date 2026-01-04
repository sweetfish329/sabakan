import { applicationConfig } from "@storybook/angular";
import type { Meta, StoryObj } from "@storybook/angular";
import { provideAnimations } from "@angular/platform-browser/animations";
import { provideRouter } from "@angular/router";
import { provideHttpClient } from "@angular/common/http";
import { of } from "rxjs";
import { GameServerListComponent } from "./game-server-list";
import { GameServerService } from "../../../services/game-server.service";
import type { GameServer } from "../../../models/game-server.model";

const mockServers: GameServer[] = [
  {
    id: 1,
    slug: "my-minecraft",
    name: "My Minecraft Server",
    description: "A fun survival server for friends",
    image: "itzg/minecraft-server:latest",
    status: "running",
    containerId: "abc123",
    ownerId: 1,
    ports: [{ hostPort: 25_565, containerPort: 25_565, protocol: "tcp" }],
    envs: [],
  },
  {
    id: 2,
    slug: "palworld-adventure",
    name: "Palworld Adventure",
    description: "Catch them all!",
    image: "thijsvanloef/palworld-server-docker:latest",
    status: "stopped",
    ownerId: 1,
    ports: [{ hostPort: 8211, containerPort: 8211, protocol: "udp" }],
    envs: [],
  },
  {
    id: 3,
    slug: "rust-survival",
    name: "Rust Survival",
    image: "max-pfeiffer/rust-game-server-docker:latest",
    status: "running",
    ownerId: 1,
    ports: [{ hostPort: 28_015, containerPort: 28_015, protocol: "udp" }],
    envs: [],
  },
];

/**
 * Creates a mock GameServerService for Storybook.
 * @param {GameServer[]} servers - List of servers to return
 * @returns {Partial<GameServerService>} Mock service
 */
const createMockService = (servers: GameServer[] = mockServers): Partial<GameServerService> => ({
  list: () => of(servers),
  delete: () => of(void 0),
});

const meta: Meta<GameServerListComponent> = {
  title: "Features/GameServers/GameServerList",
  component: GameServerListComponent,
  tags: ["autodocs"],
  decorators: [
    applicationConfig({
      providers: [provideAnimations(), provideRouter([]), provideHttpClient()],
    }),
  ],
};

export default meta;
type Story = StoryObj<GameServerListComponent>;

/**
 * Default view with multiple game servers.
 */
export const Default: Story = {
  decorators: [
    applicationConfig({
      providers: [
        provideAnimations(),
        provideRouter([]),
        provideHttpClient(),
        { provide: GameServerService, useValue: createMockService() },
      ],
    }),
  ],
};

/**
 * Empty state when no servers exist.
 */
export const Empty: Story = {
  decorators: [
    applicationConfig({
      providers: [
        provideAnimations(),
        provideRouter([]),
        provideHttpClient(),
        { provide: GameServerService, useValue: createMockService([]) },
      ],
    }),
  ],
};

/**
 * Single running server.
 */
export const SingleServer: Story = {
  decorators: [
    applicationConfig({
      providers: [
        provideAnimations(),
        provideRouter([]),
        provideHttpClient(),
        { provide: GameServerService, useValue: createMockService(mockServers.slice(0, 1)) },
      ],
    }),
  ],
};

/**
 * All servers stopped.
 */
export const AllStopped: Story = {
  decorators: [
    applicationConfig({
      providers: [
        provideAnimations(),
        provideRouter([]),
        provideHttpClient(),
        {
          provide: GameServerService,
          useValue: createMockService(
            mockServers.map((server) => ({ ...server, status: "stopped" as const })),
          ),
        },
      ],
    }),
  ],
};
