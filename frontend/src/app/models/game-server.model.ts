/**
 * Status of a game server.
 */
export type GameServerStatus = "stopped" | "running" | "creating" | "error";

/**
 * Port mapping for a game server.
 */
export interface GameServerPort {
  id?: number;
  gameServerId?: number;
  hostPort: number;
  containerPort: number;
  protocol: string;
}

/**
 * Environment variable for a game server.
 */
export interface GameServerEnv {
  id?: number;
  gameServerId?: number;
  key: string;
  value?: string;
  isSecret: boolean;
}

/**
 * Game server model.
 */
export interface GameServer {
  id: number;
  slug: string;
  name: string;
  description?: string;
  image: string;
  status: GameServerStatus;
  containerId?: string;
  ownerId: number;
  ports: GameServerPort[];
  envs: GameServerEnv[];
  createdAt?: string;
  updatedAt?: string;
}

/**
 * Request payload for creating a game server.
 */
export interface CreateGameServerRequest {
  slug: string;
  name: string;
  game: string;
  description?: string;
  ports?: Omit<GameServerPort, "id" | "gameServerId">[];
  envs?: Omit<GameServerEnv, "id" | "gameServerId">[];
}

/**
 * Request payload for updating a game server.
 */
export interface UpdateGameServerRequest {
  name?: string;
  description?: string;
  ports?: Omit<GameServerPort, "id" | "gameServerId">[];
  envs?: Omit<GameServerEnv, "id" | "gameServerId">[];
}
