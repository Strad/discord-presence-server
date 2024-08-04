//go:build !windows
// +build !windows

package ipc

import (
	"net"
	"os"
	"time"
)

type SocketConn struct {
	socket net.Conn
}

func OpenSocket() (*SocketConn, error) {
	sock, err := net.DialTimeout("unix", getIpcPath()+"/discord-ipc-0", time.Second*2)
	if err != nil {
		return nil, err
	}

	connection := SocketConn{
		socket: sock,
	}

	return &connection, nil
}

func getIpcPath() string {
	variablesnames := []string{"XDG_RUNTIME_DIR", "TMPDIR", "TMP", "TEMP"}

	if _, err := os.Stat("/run/user/1000/snap.discord"); err == nil {
		return "/run/user/1000/snap.discord"
	}

	if _, err := os.Stat("/run/user/1000/.flatpak/com.discordapp.Discord/xdg-run"); err == nil {
		return "/run/user/1000/.flatpak/com.discordapp.Discord/xdg-run"
	}

	for _, variablename := range variablesnames {
		path, exists := os.LookupEnv(variablename)

		if exists {
			return path
		}
	}

	return "/tmp"
}
