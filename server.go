package main

import (
	"net/http"

	// "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func handleWs(h *HubStruct, w http.ResponseWriter, r *http.Request) {
	nickname := r.Header.Get("iht-nickname")
	if nickname == "" {
		nickname = "anonim"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &ClientStruct{
		nickname: nickname,
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func RunServer(cmd *cobra.Command, args []string) {
	hub := &HubStruct{
		broadcast:  make(chan []byte),
		register:   make(chan *ClientStruct),
		unregister: make(chan *ClientStruct),
		clients:    make(map[*ClientStruct]bool),
	}
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWs(hub, w, r)
	})

	log.WithField("port", Flags.Port).Info("Starting server")
	log.Fatal(http.ListenAndServe(Flags.Host+":"+Flags.Port, nil))
}
