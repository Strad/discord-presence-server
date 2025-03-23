//go:build !darwin

package main

import "log"

func UpdateTrayStatus() {
	log.Println("Ready to update tray")
	for {
		update := <-rpcClient.Updates
		switch update {
		case "connected":
			go SetTrayIconConnected()
		case "cleared":
			go SetTrayIconConnected()
		case "disconnected":
			go SetTrayIconDisconnected()
		case "updated":
			go SetTrayIconActive()
		}
	}
}
