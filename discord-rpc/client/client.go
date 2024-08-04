package client

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"discord-rpc-server/discord-rpc/ipc"
)

type ClientInfo struct {
	LastMessage time.Time
	sock        *ipc.SocketConn
}

type DiscordRPC struct {
	Updates   chan string
	ClientIds map[string]*ClientInfo
}

func Create() DiscordRPC {
	return DiscordRPC{
		Updates:   make(chan string),
		ClientIds: make(map[string]*ClientInfo, 10),
	}
}

func (client *DiscordRPC) clientInfoFromId(clientId string) *ClientInfo {
	if clientInfo := client.ClientIds[clientId]; clientInfo != nil {
		return clientInfo
	}
	return nil
}

func (client *DiscordRPC) login(clientId string) error {
	if clientInfo := client.clientInfoFromId(clientId); clientInfo == nil {
		sock, err := ipc.OpenSocket()

		if err != nil {
			return err
		}

		payload, err := json.Marshal(Handshake{"1", clientId})
		if err != nil {
			return err
		}

		_, err = sock.Send(0, string(payload))

		if err != nil {
			return err
		}

		clientInfo := &ClientInfo{
			LastMessage: time.Now(),
			sock:        sock,
		}

		client.ClientIds[clientId] = clientInfo

		client.Updates <- "connected"
	}

	return nil
}

func (client *DiscordRPC) logout(clientId string) {
	if clientInfo := client.clientInfoFromId(clientId); clientInfo != nil && clientInfo.sock != nil {
		clientInfo.sock.CloseSocket()
		client.Updates <- "disconnected"
		delete(client.ClientIds, clientId)
	}
}

func (client *DiscordRPC) SetActivity(clientId string, activity Activity) (string, error) {
	if clientInfo := client.clientInfoFromId(clientId); clientInfo == nil {
		if err := client.login(clientId); err != nil {
			log.Printf("SetActivity(): Error logging into Discord: %s", err)
			client.Updates <- "disconnected"
			return "", err
		}
	}

	clientInfo := client.ClientIds[clientId]

	fakePid, err := fakePid(clientId)

	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(Frame{
		"SET_ACTIVITY",
		Args{
			fakePid,
			mapActivity(&activity),
		},
		getNonce(),
	})

	if err != nil {
		return "", err
	}

	resp, err := clientInfo.sock.Send(1, string(payload))

	if err != nil {
		client.logout(clientId)
		return "", err
	}

	clientInfo.LastMessage = time.Now()

	client.Updates <- "updated"

	go func() {
		time.Sleep(15 * time.Second)

		if clientInfo == nil {
			return
		}

		if time.Since(clientInfo.LastMessage) > time.Second*15 {
			log.Printf("SetActivity(): Clearing activity for %s (last update on %s)", clientId, clientInfo.LastMessage.Format(time.RFC3339))

			_, err = client.ClearActivity(clientId)

			if err != nil {
				log.Printf("SetActivity(): Error on scheduled activity clear: %s", err)
				return
			}

			client.Updates <- "cleared"
		}
	}()

	return resp, nil
}

func (client *DiscordRPC) ClearActivity(clientId string) (string, error) {
	if clientInfo := client.clientInfoFromId(clientId); clientInfo == nil {
		if err := client.login(clientId); err != nil {
			log.Printf("ClearActivity(): Error logging into Discord: %s", err)
			client.Updates <- "disconnected"
			return "", err
		}
	}

	clientInfo := client.ClientIds[clientId]

	fakePid, err := fakePid(clientId)

	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(Frame{
		"SET_ACTIVITY",
		Args{
			fakePid,
			nil,
		},
		getNonce(),
	})

	if err != nil {
		return "", err
	}

	resp, err := clientInfo.sock.Send(1, string(payload))

	if err != nil {
		client.logout(clientId)
		return "", err
	}

	client.Updates <- "cleared"

	return resp, nil
}

func getNonce() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println(err)
	}

	buf[6] = (buf[6] & 0x0f) | 0x40

	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:])
}

func fakePid(clientId string) (int, error) {
	return strconv.Atoi(clientId[len(clientId)-5:])
}
