package server

import (
	"fmt"
	"net/http"

	"github.com/pubnub/go/messaging"
	"github.com/unrolled/render"
)

func broadcastHandler(messagingConfig *pubsubConfig) http.HandlerFunc {

	formatter := render.New(render.Options{
		IndentJSON: true,
	})

	pubnub := messaging.NewPubnub(messagingConfig.PublishKey, messagingConfig.SubscribeKey, "", "", false, "")
	channel := "hello_world"

	return func(w http.ResponseWriter, r *http.Request) {

		successChannel := make(chan []byte)
		errorChannel := make(chan []byte)
		go pubnub.Publish(channel, "[SYSTEM] Broadcast from the server!", successChannel, errorChannel)

		select {
		case response := <-successChannel:
			fmt.Printf("pubnub publish response: %s\n", string(response))
		case err := <-errorChannel:
			fmt.Printf("pubnub publish error: %s\n", string(err))
		case <-messaging.Timeout():
			fmt.Println("pubnub publish() timeout")
		}
		formatter.JSON(w, http.StatusOK, nil)
	}
}
