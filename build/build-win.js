// build.js
const exe = require('@angablue/exe');
const pkgJson = require('../package.json');

const build = exe({
	entry: './dist-win/index.js',
	out: './out/Discord Presence Server (Debug).exe',
	version: pkgJson.version,
	pkg: ['-c', './build/pkg-win.json'],
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
