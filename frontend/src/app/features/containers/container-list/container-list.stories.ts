import {
	applicationConfig,
	type Meta,
	type StoryObj,
} from '@storybook/angular';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideRouter } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { of } from 'rxjs';
import { ContainerListComponent } from './container-list';
import { ContainerService } from '../../../services/container.service';
import type { Container } from '../../../models/container.model';

const mockContainers: Container[] = [
	{
		id: 'abc123def456',
		name: 'minecraft-server-1',
		image: 'itzg/minecraft-server:latest',
		state: 'running',
		status: 'Up 2 hours',
		created: '2024-01-01T00:00:00Z',
		ports: [{ hostPort: 25565, containerPort: 25565, protocol: 'tcp' }],
		labels: { game: 'minecraft', type: 'paper' },
	},
	{
		id: 'xyz789ghi012',
		name: 'valheim-server',
		image: 'lloesche/valheim-server:latest',
		state: 'stopped',
		status: 'Exited (0) 3 hours ago',
		created: '2024-01-01T00:00:00Z',
		ports: [
			{ hostPort: 2456, containerPort: 2456, protocol: 'udp' },
			{ hostPort: 2457, containerPort: 2457, protocol: 'udp' },
		],
		labels: { game: 'valheim' },
	},
	{
		id: 'pause123',
		name: 'terraria-server',
		image: 'ryshe/terraria:latest',
		state: 'paused',
		status: 'Paused',
		created: '2024-01-01T00:00:00Z',
		ports: [{ hostPort: 7777, containerPort: 7777, protocol: 'tcp' }],
		labels: {},
	},
];

/**
 * Creates a mock ContainerService for Storybook.
 */
function createMockContainerService(
	containers: Container[] = mockContainers,
	shouldError = false,
): Partial<ContainerService> {
	return {
		list: () => (shouldError ? of([]) : of(containers)),
		start: () => of(void 0),
		stop: () => of(void 0),
	};
}

const meta: Meta<ContainerListComponent> = {
	title: 'Features/Containers/ContainerList',
	component: ContainerListComponent,
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
type Story = StoryObj<ContainerListComponent>;

/**
 * Default view with multiple containers in various states.
 */
export const Default: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ContainerService,
					useValue: createMockContainerService(),
				},
			],
		}),
	],
};

/**
 * Empty state when no containers are found.
 */
export const Empty: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ContainerService,
					useValue: createMockContainerService([]),
				},
			],
		}),
	],
};

/**
 * Single running container.
 */
export const SingleRunning: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ContainerService,
					useValue: createMockContainerService(mockContainers.slice(0, 1)),
				},
			],
		}),
	],
};

/**
 * All containers in stopped state.
 */
export const AllStopped: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ContainerService,
					useValue: createMockContainerService(
						mockContainers.map((c) => ({ ...c, state: 'stopped' as const })),
					),
				},
			],
		}),
	],
};

/**
 * Many containers to test grid layout.
 */
export const ManyContainers: Story = {
	decorators: [
		applicationConfig({
			providers: [
				provideAnimationsAsync(),
				provideRouter([]),
				provideHttpClient(),
				{
					provide: ContainerService,
					useValue: createMockContainerService([
						...mockContainers,
						{
							id: 'ark123',
							name: 'ark-survival',
							image: 'hermsi/ark-server:latest',
							state: 'running',
							status: 'Up 1 day',
							created: '2024-01-01T00:00:00Z',
							ports: [{ hostPort: 27015, containerPort: 27015, protocol: 'udp' }],
							labels: { game: 'ark' },
						},
						{
							id: 'rust456',
							name: 'rust-server',
							image: 'didstopia/rust-server:latest',
							state: 'running',
							status: 'Up 5 hours',
							created: '2024-01-01T00:00:00Z',
							ports: [{ hostPort: 28015, containerPort: 28015, protocol: 'udp' }],
							labels: { game: 'rust' },
						},
						{
							id: 'satisfactory789',
							name: 'satisfactory-dedicated',
							image: 'wolveix/satisfactory-server:latest',
							state: 'exited',
							status: 'Exited (1) 2 days ago',
							created: '2024-01-01T00:00:00Z',
							ports: [{ hostPort: 7777, containerPort: 7777, protocol: 'udp' }],
							labels: { game: 'satisfactory' },
						},
					]),
				},
			],
		}),
	],
};
