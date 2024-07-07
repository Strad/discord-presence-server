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

		const reconnection = async () => {
			if (client.isConnected) {
				await this.disconnect(client);
			}

			client.once('disconnected', reconnection);
			client.once('ERROR', reconnection);

			console.log('Connecting to Discord client');
			await client.connect();
		};

		await reconnection();
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
				this.emit('disconnected');
				this.isUpdating = false;
				this.clearActivity(clientId);
			}
		}, TIME_TILL_ACTIVITY_CLEAR);
	}

	async clearActivity(clientId: string) {
		const client = this.clientMap.get(clientId);
		if (client) {
			console.log(`Clearing activity on clientId ${clientId}`);
			await client.user?.clearActivity();
		}
	}

	async disconnect(client: Client) {
		if (client && client.isConnected) {
			console.log('Disconnecting from Discord');

			await client.destroy();
		}
	}
}

export const discordClientManager = new DiscordClientManager();
