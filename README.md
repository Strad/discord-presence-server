# Discord Rich Presence Server

A re-write of the server component of this project https://github.com/lolamtisch/Discord-RPC-Extension/tree/master/Server

This is an application which creates a server which consumes presence updates and sends them to your Discord client. A use case is sending updates from a browser extension or web script.

The WebSocket server runs on port `6969`. Here is an example of connecting using a JS client, sending an update every 2 seconds, and then clearing the activity after 20 seconds. The server will automatically clear the status of a given application/client ID if it has received no updates for 10 seconds.

```javascript
// See for payload information https://discord.com/developers/docs/rich-presence/how-to#updating-presence-update-presence-payload-fields

const ws = new WebSocket('ws://localhost:6969');

ws.addEventListener('open', () => {
  console.log('Connected to Discord Presence Server');

  const start = Date.now();

  const doUpdate = setInterval(() => {
    ws.send(
      JSON.stringify({
        // Client ID of your application (from Discord developer portal)
        clientId: '00000000000000000000',
        presence: {
          details: 'Some text',
          state: 'Some more text',
          largeImageKey: 'Asset key, from your Discord developer portal',
          largeImageText: 'Text which shows when hovering over the image',
          startTimestamp: start,
          // endTimestamp: "",
          buttons: [
            {
              label: 'A button label',
              url: 'https://google.co.uk',
            },
          ],
        },
      }),
    );
  }, 2000);

  setTimeout(() => {
    clearInterval(doUpdate);
    ws.send(
      JSON.stringify({
        clientId: '0000000000000000000',
        clearActivity: true,
      }),
    );
  }, 20000);
});
```
