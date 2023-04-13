package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SseHandler(w http.ResponseWriter, r *http.Request, clientChannelMap map[string]chan string) {
	fmt.Println("Client connected")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	clientId := mux.Vars(r)["clientId"]
	clientChannelMap[clientId] = make(chan string)

	defer func() {
		if clientChannel, ok := clientChannelMap[clientId]; ok {
			close(clientChannel)
			delete(clientChannelMap, clientId)
		}
		fmt.Println("Client closed connection")
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Println("Could not init http.Flusher")
	}

	for {
		select {
		case message := <-clientChannelMap[clientId]:
			fmt.Println("case message... sending message")
			fmt.Println(message)
			fmt.Fprintf(w, "data: %s\n\n", message)
			flusher.Flush()
		case <-r.Context().Done():
			fmt.Println("Client closed connection")
			return
		}
	}

}
