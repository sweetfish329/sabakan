/**
 * State of a container.
 */
export type ContainerState =
	| 'running'
	| 'stopped'
	| 'created'
	| 'paused'
	| 'restarting'
	| 'exited'
	| 'unknown';

/**
 * Represents a port mapping for a container.
 */
export interface PortMapping {
	/** Host IP address */
	hostIp?: string;
	/** Port on the host */
	hostPort: number;
	/** Port inside the container */
	containerPort: number;
	/** Protocol (tcp, udp) */
	protocol: string;
}

/**
 * Represents a game server container.
 */
export interface Container {
	/** Unique identifier of the container */
	id: string;
	/** Human-readable name of the container */
	name: string;
	/** Container image used */
	image: string;
	/** Current state of the container */
	state: ContainerState;
	/** Human-readable status string */
	status: string;
	/** Creation timestamp */
	created: string;
	/** List of port mappings */
	ports: PortMapping[];
	/** Container labels */
	labels: Record<string, string>;
}

/**
 * Represents a single log entry from a container.
 */
export interface ContainerLogEntry {
	/** When the log was generated */
	timestamp?: string;
	/** Output stream (stdout or stderr) */
	stream: string;
	/** Log message content */
	message: string;
}
