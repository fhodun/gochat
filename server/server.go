package server

import (
	"net/http"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func handleWs(h *hub, w http.ResponseWriter, r *http.Request) {
	nickname := r.Header.Get("iht-nickname")
	if nickname == "" {
		nickname = "anonim"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	c := &client{
		Nickname: nickname,
		Hub:      h,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}
	c.Hub.Register <- c

	go c.writePump()
	go c.readPump()
}

func RunServer(cmd *cobra.Command, args []string) {
	port := args[0]

	hub := &hub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *client),
		Unregister: make(chan *client),
		Clients:    make(map[*client]bool),
	}
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWs(hub, w, r)
	})

	log.WithField("port", port).Info("Starting server")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
