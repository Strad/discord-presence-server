//go:build !darwin

package main

import (
	"discord-rpc-server/icon"
	"log"
	"os"

	"github.com/getlantern/systray"
)

var status *systray.MenuItem

func main() {
	log.Printf("Discord Presence Server")
	go ServeWs()
	go UpdateTrayStatus()
	systray.Run(trayOnReady, func() {
		os.Exit(0)
	})
}

func trayOnReady() {
	systray.SetTemplateIcon(icon.Data_icon, icon.Data_icon)
	systray.SetTitle("Discord Rich Presence Server")
	systray.SetTooltip("Discord Rich Presence Server")

	status = systray.AddMenuItem("Waiting to receive presence updates", "")
	status.Disable()

	systray.AddSeparator()
	mQuitOrig := systray.AddMenuItem("Quit", "Quit presence server")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

	systray.SetIcon(icon.Data_disconnected)
}

func SetTrayIconDisconnected() {
	status.SetTitle("Presence server disconnected from Discord")
	status.SetTooltip("An error occurred and the server has lost connection with the Discord client")
	systray.SetIcon(icon.Data_disconnected)
}

func SetTrayIconActive() {
	status.SetTitle("Relaying presence updates to Discord")
	status.SetTooltip("")
	systray.SetIcon(icon.Data_connected)
}

func SetTrayIconConnected() {
	status.SetTitle("Ready to relay presence updates to Discord")
	status.SetTooltip("")
	systray.SetIcon(icon.Data_icon)
}
