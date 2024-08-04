package main

import (
	"discord-rpc-server/discord-rpc/client"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"slices"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
)

type Presence struct {
	Type           *int     `json:"type,omitempty"`
	State          *string  `json:"state,omitempty"`
	Details        *string  `json:"details,omitempty"`
	StartTimestamp *float64 `json:"startTimestamp,omitempty"`
	EndTimestamp   *float64 `json:"endTimestamp,omitempty"`
	LargeImageKey  *string  `json:"largeImageKey,omitempty"`
	SmallImageKey  *string  `json:"smallImageKey,omitempty"`
	LargeImageText *string  `json:"largeImageText,omitempty"`
	SmallImageText *string  `json:"smallImageText,omitempty"`
	Buttons        *[]struct {
		Label *string `json:"label,omitempty"`
		Url   *string `json:"url,omitempty"`
	} `json:"buttons,omitempty"`
}

type UpdateMessageSchema struct {
	ClientID      *string   `json:"clientId,omitempty" validate:"required"`
	ClearActivity *bool     `json:"clearActivity,omitempty" validate:"required_without=Presence"`
	Presence      *Presence `json:"presence,omitempty" validate:"required_without=ClearActivity"`
}

var (
	port      = flag.String("port", "6969", "ws server port")
	upgrader  = websocket.Upgrader{} // use default options
	validate  = validator.New()
	rpcClient = client.Create()
	// https://discord.com/developers/docs/topics/rpc#setactivity-set-activity-argument-structure
	activityTypes = []int{0, 2, 3, 5}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		updateMessage := UpdateMessageSchema{}
		json.Unmarshal(message, &updateMessage)

		if err := validate.Struct(updateMessage); err != nil {
			log.Printf("wsHandler(): Error validating JSON, recv: %s", err)
			c.WriteMessage(mt, []byte(fmt.Sprintf("Error validating JSON: %s", err)))
			continue
		}

		if updateMessage.ClearActivity != nil && *updateMessage.ClearActivity {
			resp, err := rpcClient.ClearActivity(*updateMessage.ClientID)
			if err != nil {
				log.Printf("wsHandler(): Error clearing activity: %s", err)
			}

			c.WriteMessage(mt, []byte(resp))
			continue
		}

		activity := client.Activity{}

		if updateMessage.Presence.Type != nil && slices.Contains(activityTypes, *updateMessage.Presence.Type) {
			activity.Type = *updateMessage.Presence.Type
		} else {
			activity.Type = 0
		}

		if updateMessage.Presence.Details != nil {
			activity.Details = *updateMessage.Presence.Details
		}

		if updateMessage.Presence.State != nil {
			activity.State = *updateMessage.Presence.State
		}

		if updateMessage.Presence.LargeImageKey != nil {
			activity.LargeImage = *updateMessage.Presence.LargeImageKey
		}

		if updateMessage.Presence.LargeImageText != nil {
			activity.LargeText = *updateMessage.Presence.LargeImageText
		}

		if updateMessage.Presence.SmallImageKey != nil {
			activity.SmallImage = *updateMessage.Presence.SmallImageKey
		}

		if updateMessage.Presence.SmallImageText != nil {
			activity.SmallText = *updateMessage.Presence.SmallImageText
		}

		if updateMessage.Presence.StartTimestamp != nil || updateMessage.Presence.EndTimestamp != nil {
			activity.Timestamps = &client.Timestamps{}
		}

		if updateMessage.Presence.StartTimestamp != nil {
			start := time.UnixMilli(int64(math.Round(*updateMessage.Presence.StartTimestamp)))
			activity.Timestamps.Start = &start
		}

		if updateMessage.Presence.EndTimestamp != nil {
			end := time.UnixMilli(int64(math.Round(*updateMessage.Presence.EndTimestamp)))
			activity.Timestamps.End = &end
		}

		if updateMessage.Presence.Buttons != nil {
			buttonArr := make([]*client.Button, len(*updateMessage.Presence.Buttons))

			for i := 0; i < len(buttonArr); i++ {
				buttonArr[i] = &client.Button{
					Label: *(*updateMessage.Presence.Buttons)[i].Label,
					Url:   *(*updateMessage.Presence.Buttons)[i].Url,
				}
			}

			activity.Buttons = buttonArr
		}

		resp, err := rpcClient.SetActivity(*updateMessage.ClientID, activity)

		if err != nil {
			log.Printf("wsHandler(): Error setting activity: %s", err)
		}

		c.WriteMessage(mt, []byte(resp))
	}
}

func ServeWs() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", wsHandler)

	log.Printf("Starting websocket server")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", *port), nil))
}

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
