// build.js
const exe = require('@angablue/exe');

const build = exe({
	entry: './out/entry.js',
	out: './out/Discord Presence Server (Debug).exe',
	version: '1.0.0',
	target: 'latest-win-x64',
	icon: './build/icon.ico', // Application icons must be in .ico format
	executionLevel: 'asInvoker',
	properties: {
		FileDescription:
			'Update Discord Rich Presence via a local WebSocket server',
		ProductName: 'Discord Presence Server',
		LegalCopyright: 'github.com/NotForMyCV',
	},
});

build.then(() => {
	require('create-nodew-exe')({
		src: './out/Discord Presence Server (Debug).exe',
		dst: './out/Discord Presence Server.exe',
	});
});
