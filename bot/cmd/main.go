package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"mattermost-vote-bot/internal/handlers"
	"mattermost-vote-bot/internal/storage"

	"github.com/mattermost/mattermost-server/v6/model"
)

func main() {
	mmClient := model.NewAPIv4Client(os.Getenv("MATTERMOST_URL"))
	mmClient.SetToken(os.Getenv("MATTERMOST_TOKEN"))

	tntClient, err := storage.InitTarantool(os.Getenv("TARANTOOL_URI"))
	if err != nil {
		log.Fatal("Tarantool connection failed:", err)
	}

	wsClient, err := model.NewWebSocketClient4("ws://mattermost:8065", mmClient.AuthToken)
	if err != nil {
		log.Fatal("WebSocket error:", err)
	}
	wsClient.Listen()

	for event := range wsClient.EventChannel {
		if event.EventType() == model.WebsocketEventPosted {
			var post model.Post
			if err := json.Unmarshal([]byte(event.GetData()["post"].(string)), &post); err != nil {
				log.Printf("Error parsing post: %v", err)
				continue
			}

			if strings.HasPrefix(post.Message, "/vote") {
				handlers.HandleCommand(&post, mmClient, tntClient)
			}
		}
	}
}
