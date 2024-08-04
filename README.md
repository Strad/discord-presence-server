# Discord Rich Presence Server

Opens a local WebSocket server to update your Discord Rich Presence.

A re-write of the server component of this project https://github.com/lolamtisch/Discord-RPC-Extension/tree/master/Server

## Usage

By default, the WebSocket server runs on port `6969`.

#### Expected payload structure

```jsonc
{
  "clientId": string,
  // One of the below are required, if clearActivity is true, it will take priority
  "clearActivity": boolean,
  "presence": Object <Presence payload structure>
}
```

#### Presence payload structure

| field            | type   | description                                                                  | example                                                        |
| ---------------- | ------ | ---------------------------------------------------------------------------- | -------------------------------------------------------------- |
| `state`          | string | the user's current party status                                              | `"Looking to Play"`, `"Playing Solo"`, `"In a Group"`          |
| `details`        | string | what the player is currently doing                                           | `"Competitive - Captain's Mode"`, `"In Queue", "Unranked PvP"` |
| `startTimestamp` | number | epoch milliseconds of activity start - including will show time as "elapsed" | `1707138793000` (Mon Feb 05 2024 13:13:13.000)                 |
| `endTimestamp`   | number | epoch milliseconds of activity end - including will show time as "remaining" | `1707142393000` (Mon Feb 05 2024 14:13:13.000)                 |
| `largeImageKey`  | string | name of the uploaded image for the large profile artwork                     | `"default"`                                                    |
| `largeImageText` | string | tooltip for the largeImageKey                                                | `"Blade's Edge Arena"`, `"Numbani"`, `"Danger Zone"`           |
| `smallImageKey`  | string | name of the uploaded image for the small profile artwork                     | `"rogue"`                                                      |
| `smallImageText` | string | tooltip for the smallImageKey                                                | `"Rogue - Level 100"`                                          |

### Example

Here is an example of connecting using a JS client, sending an update every 2 seconds, and clearing the activity after 20 seconds. The server will automatically clear the status of a given application/client ID if it has received no updates for 15 seconds.

```javascript
const ws = new WebSocket("ws://localhost:6969");

ws.addEventListener("open", () => {
	console.log("Connected to Discord Presence Server");

	const start = Date.now();

	const doUpdate = setInterval(() => {
		ws.send(
			JSON.stringify({
				// Client ID of your application (from Discord developer portal)
				clientId: "00000000000000000000",
				presence: {
					details: "Some text",
					state: "Some more text",
					largeImageKey: "Asset key, from your Discord developer portal",
					largeImageText: "Text which shows when hovering over the image",
					startTimestamp: start,
					// endTimestamp: "",
					buttons: [
						{
							label: "A button label",
							url: "https://google.co.uk",
						},
					],
				},
			})
		);
	}, 2000);

	setTimeout(() => {
		clearInterval(doUpdate);
		ws.send(
			JSON.stringify({
				clientId: "0000000000000000000",
				clearActivity: true,
			})
		);
	}, 20000);
});
```
