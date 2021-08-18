package server

import (
	"net/http"

	"github.com/fhodun/gochat/packets"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func handleWs(h *hub, w http.ResponseWriter, r *http.Request) {
	var (
		upgrader websocket.Upgrader = websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}
		nickname string = r.Header.Get("nickname")
	)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(err)
		return
	}

	c := &client{
		Nickname: nickname,
		Hub:      h,
		Conn:     conn,
		Send:     make(chan packets.Message, 256),
	}

	c.Hub.Register <- c

	go c.writePump()
	go c.readPump()
}

func RunServer(cmd *cobra.Command, args []string) {
	var (
		port string = args[0]
		h    *hub   = &hub{
			Broadcast:  make(chan packets.Message),
			Register:   make(chan *client),
			Unregister: make(chan *client),
			Clients:    make(map[*client]bool),
		}
	)

	go h.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWs(h, w, r)
	})

	log.WithField("port", port).Info("Starting server")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
