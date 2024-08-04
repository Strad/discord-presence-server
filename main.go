package main

import (
	"discord-rpc-server/icon"
	"log"
	"os"

	"github.com/getlantern/systray"
)

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

	mQuitOrig := systray.AddMenuItem("Quit", "Quit presence server")
	go func() {
		<-mQuitOrig.ClickedCh
		systray.Quit()
	}()

	systray.SetIcon(icon.Data_disconnected)
}

func SetTrayIconDisconnected() {
	systray.SetIcon(icon.Data_disconnected)
}

func SetTrayIconActive() {
	systray.SetIcon(icon.Data_connected)
}

func SetTrayIconConnected() {
	systray.SetIcon(icon.Data_icon)
}
