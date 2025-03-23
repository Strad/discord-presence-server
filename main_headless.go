//go:build darwin

package main

import "log"

func main() {
	log.Println("Starting headless Discord Presence Server (macOS)")
	go ServeWs()
	select {} // Keeps the app alive
}
