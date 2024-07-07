import SysTray, { Conf } from 'systray';
// Start ws server
import {
	CONNECTED_ICO,
	CONNECTED_PNG,
	DISCONNECTED_ICO,
	DISCONNECTED_PNG,
} from './sys-tray-icons';
import { discordClientManager } from './discord-client';

require('./server');

const isWin32 = process.platform === 'win32';

const connectedIcon = isWin32 ? CONNECTED_ICO : CONNECTED_PNG;
const disconnectedIcon = isWin32 ? DISCONNECTED_ICO : DISCONNECTED_PNG;

const enum Actions {
	QUIT = 'Exit',
	TITLE = 'Discord Rich Presence Server',
}

const sysTrayConf: Conf = {
	menu: {
		icon: disconnectedIcon,
		title: 'Discord RPC Server',
		tooltip: 'No RPC clients connected',
		items: [
			{
				title: Actions.TITLE,
				tooltip: Actions.TITLE,
				enabled: false,
				checked: false,
			},
			{
				title: Actions.QUIT,
				tooltip: Actions.QUIT,
				checked: false,
				enabled: true,
			},
		],
	},
	copyDir: true,
};

const systray = new SysTray({ ...sysTrayConf });

systray.onClick((action) => {
	if (action.item.title === Actions.QUIT) {
		systray.kill();
	}
});

discordClientManager.on('disconnected', () => {
	systray.sendAction({
		type: 'update-menu',
		menu: {
			...sysTrayConf.menu,
			icon: disconnectedIcon,
			tooltip: 'No RPC clients connected',
		},
		seq_id: -1,
	});
});

discordClientManager.on('connected', () => {
	systray.sendAction({
		type: 'update-menu',
		menu: {
			...sysTrayConf.menu,
			icon: connectedIcon,
			tooltip: 'Receiving presence updates',
		},
		seq_id: -1,
	});
});
