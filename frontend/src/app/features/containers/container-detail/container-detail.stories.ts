import {
	applicationConfig,
	type Meta,
	type StoryObj,
} from '@storybook/angular';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideRouter, ActivatedRoute } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { of } from 'rxjs';
import { ContainerDetailComponent } from './container-detail';
import { ContainerService } from '../../../services/container.service';
import type { Container, ContainerLogEntry } from '../../../models/container.model';

const mockRunningContainer: Container = {
	id: 'abc123def456789',
	name: 'minecraft-server-1',
	image: 'itzg/minecraft-server:latest',
	state: 'running',
	status: 'Up 2 hours',
	created: '2024-01-01T00:00:00Z',
	ports: [
		{ hostPort: 25565, containerPort: 25565, protocol: 'tcp' },
		{ hostPort: 25575, containerPort: 25575, protocol: 'tcp' },
	],
	labels: { game: 'minecraft', type: 'paper', version: '1.20.4' },
};

const mockStoppedContainer: Container = {
	id: 'xyz789ghi012345',
	name: 'valheim-server',
	image: 'lloesche/valheim-server:latest',
	state: 'stopped',
	status: 'Exited (0) 3 hours ago',
	created: '2024-01-01T00:00:00Z',
	ports: [
		{ hostPort: 2456, containerPort: 2456, protocol: 'udp' },
		{ hostPort: 2457, containerPort: 2457, protocol: 'udp' },
	],
	labels: { game: 'valheim', world: 'myworld' },
};

const mockLogs: ContainerLogEntry[] = [
	{ timestamp: '2024-01-01T12:00:00Z', stream: 'stdout', message: '[Server] Starting server...' },
	{ timestamp: '2024-01-01T12:00:01Z', stream: 'stdout', message: '[Server] Loading world data...' },
	{ timestamp: '2024-01-01T12:00:05Z', stream: 'stdout', message: '[Server] World loaded successfully' },
	{ timestamp: '2024-01-01T12:00:06Z', stream: 'stdout', message: '[Server] Listening on port 25565' },
	{ timestamp: '2024-01-01T12:00:10Z', stream: 'stderr', message: '[Warning] High memory usage detected' },
	{ timestamp: '2024-01-01T12:00:15Z', stream: 'stdout', message: '[Server] Player "Steve" joined the game' },
	{ timestamp: '2024-01-01T12:00:20Z', stream: 'stdout', message: '[Server] Player "Alex" joined the game' },
];

/**
 * Creates a mock ActivatedRoute for testing.
 */
function createMockActivatedRoute(containerId: string): Partial<ActivatedRoute> {
	return {
		snapshot: {
			paramMap: {
				get: (key: string) => (key === 'id' ? containerId : null),
				has: (key: string) => key === 'id',
				getAll: () => [],
				keys: ['id'],
			},
		} as unknown as ActivatedRoute['snapshot'],
	};
}

/**
 * Creates a mock ContainerService for Storybook.
 */
function createMockContainerService(
	container: Container,
	logs: ContainerLogEntry[] = mockLogs,
): Partial<ContainerService> {
	return {
		get: () => of(container),
		logs: () => of(logs),
		start: () => of(void 0),
		stop: () => of(void 0),
	};
}

const meta: Meta<ContainerDetailComponent> = {
	title: 'Features/Containers/ContainerDetail',
	component: ContainerDetailComponent,
	tags: ['autodocs'],
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
			],
		}),
	],
};

export default meta;
type Story = StoryObj<ContainerDetailComponent>;

/**
 * Running container with logs and full information.
 */
export const Running: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ActivatedRoute,
					useValue: createMockActivatedRoute(mockRunningContainer.id),
				},
				{
					provide: ContainerService,
					useValue: createMockContainerService(mockRunningContainer),
				},
			],
		}),
	],
};

/**
 * Stopped container that can be started.
 */
export const Stopped: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ActivatedRoute,
					useValue: createMockActivatedRoute(mockStoppedContainer.id),
				},
				{
					provide: ContainerService,
					useValue: createMockContainerService(mockStoppedContainer, []),
				},
			],
		}),
	],
};

/**
 * Container with no labels.
 */
export const NoLabels: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ActivatedRoute,
					useValue: createMockActivatedRoute('nolabels123'),
				},
				{
					provide: ContainerService,
					useValue: createMockContainerService({
						...mockRunningContainer,
						id: 'nolabels123',
						labels: {},
					}),
				},
			],
		}),
	],
};

/**
 * Container with no ports.
 */
export const NoPorts: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ActivatedRoute,
					useValue: createMockActivatedRoute('noports123'),
				},
				{
					provide: ContainerService,
					useValue: createMockContainerService({
						...mockRunningContainer,
						id: 'noports123',
						ports: [],
					}),
				},
			],
		}),
	],
};

/**
 * Container with empty logs.
 */
export const EmptyLogs: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ActivatedRoute,
					useValue: createMockActivatedRoute(mockRunningContainer.id),
				},
				{
					provide: ContainerService,
					useValue: createMockContainerService(mockRunningContainer, []),
				},
			],
		}),
	],
};

/**
 * Container with stderr logs.
 */
export const WithErrors: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ActivatedRoute,
					useValue: createMockActivatedRoute(mockRunningContainer.id),
				},
				{
					provide: ContainerService,
					useValue: createMockContainerService(mockRunningContainer, [
						{ stream: 'stdout', message: '[Server] Starting...' },
						{ stream: 'stderr', message: '[Error] Failed to bind port 25565' },
						{ stream: 'stderr', message: '[Error] Address already in use' },
						{ stream: 'stdout', message: '[Server] Retrying with port 25566...' },
						{ stream: 'stderr', message: '[Warning] Running on non-standard port' },
					]),
				},
			],
		}),
	],
};
