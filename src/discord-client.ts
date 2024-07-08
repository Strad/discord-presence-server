import { Client, SetActivity } from '@xhayper/discord-rpc';
import { EventEmitter } from 'node:events';

const TIME_TILL_ACTIVITY_CLEAR = 10000;

class DiscordClientManager extends EventEmitter {
	clientMap: Map<string, Client>;
	lastActivityMap: Map<string, number>;
	private isUpdating = false;

	constructor() {
		super();

		this.clientMap = new Map();
		this.lastActivityMap = new Map();
		process.on('exit', () =>
			this.clientMap.forEach((client) => client.destroy()),
		);
	}

	private async createClient(clientId: string) {
		let client = this.clientMap.get(clientId);
		if (client && client.isConnected) {
			await this.disconnect(client);
		}

		client = new Client({ clientId });
		this.clientMap.set(clientId, client);

		console.log('Connecting to Discord client');

		try {
			await client.connect();
		} catch (error) {
			console.error(
				'Failed to connect to Discord IPC socket, will wait for next update for reattempt: ',
				error,
			);
		}

		return client;
	}

	async setActivity(activity: SetActivity, clientId: string) {
		let client = this.clientMap.get(clientId);

		if (!client || !client.isConnected) {
			client = await this.createClient(clientId);
		}

		console.log(
			`Setting activity on clientId ${clientId}: `,
			JSON.stringify(activity, undefined, 2),
		);
		await client.user?.setActivity(activity);

		if (!this.isUpdating) {
			this.isUpdating = true;
			this.emit('connected');
		}

		this.lastActivityMap.set(clientId, Date.now());

		setTimeout(() => {
			const lastUpdate = this.lastActivityMap.get(clientId) ?? 0;

			if (Date.now() - lastUpdate > TIME_TILL_ACTIVITY_CLEAR) {
				this.clearActivity(clientId);

				this.emit('disconnected');
				this.isUpdating = false;
			}
		}, TIME_TILL_ACTIVITY_CLEAR);
	}

	clearActivity(clientId: string) {
		const client = this.clientMap.get(clientId);
		if (client && client.user) {
			console.log(`Clearing activity on clientId ${clientId}`);
			client.user
				.clearActivity()
				.catch((error) =>
					console.error('Error trying to clear activity: ', error),
				);
		}
	}

	async disconnect(client: Client) {
		if (client && client.isConnected) {
			console.log('Disconnecting from Discord');
			return client
				.destroy()
				.catch((error) =>
					console.error('Error trying to disconnect from Discord IPC: ', error),
				);
		}
	}
}

export const discordClientManager = new DiscordClientManager();
