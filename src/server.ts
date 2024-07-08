import WebSocket from 'ws';
import { ZodError, z } from 'zod';
import { discordClientManager } from './discord-client';

const updateMessageSchema = z.object({
	clearActivity: z.boolean().optional(),
	clientId: z.string(),
	presence: z
		.object({
			state: z.string(),
			details: z.string(),
			startTimestamp: z.union([
				z.number(),
				z
					.string()
					.datetime()
					.transform((arg) => new Date(arg)),
			]),
			endTimestamp: z.union([
				z.number(),
				z
					.string()
					.datetime()
					.transform((arg) => new Date(arg)),
			]),
			largeImageKey: z.string(),
			smallImageKey: z.string(),
			largeImageText: z.string(),
			smallImageText: z.string(),
			buttons: z.array(
				z.object({
					label: z.string(),
					url: z.string(),
				}),
			),
		})
		.partial()
		.optional(),
});

const server = new WebSocket.Server({ port: 6969 });

server.on('connection', (socket) => {
	socket.on('open', () => console.log('Client connected to RPC'));

	socket.on('message', async (msg) => {
		try {
			const message = JSON.parse(msg.toString());
			const payload = updateMessageSchema.parse(message);

			if (
				payload.clearActivity ||
				Object.keys(payload?.presence ?? {}).length === 0
			) {
				discordClientManager.clearActivity(payload.clientId);
				return;
			}

			await discordClientManager.setActivity(
				payload.presence ?? {},
				payload.clientId,
			);
		} catch (error) {
			if (error instanceof ZodError) {
				console.error('Invalid presence update: ', error.issues);
			} else {
				console.error(error);
			}
		}
	});

	socket.on('error', (error) => console.error(error));
});
