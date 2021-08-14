package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
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

func RunClient(cmd *cobra.Command, args []string) {
	var (
		header    http.Header = make(http.Header)
		serverUrl url.URL     = url.URL{
			Scheme: "ws",
			Host:   Flags.Host + ":" + Flags.Port,
			Path:   "/",
		}
	)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter nickname: ")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	header.Add("iht-nickname", nickname)

	log.Printf("connecting to %s", serverUrl.String())
	ws, _, err := websocket.DefaultDialer.Dial(serverUrl.String(), header)
	if err != nil {
		log.Fatal(err)
	}

	chatLayout, err := RenderLayout()
	if err != nil {
		log.Fatal(err)
	}

	if nickname == "" {
		chatLayout.room.Rows = append(chatLayout.room.Rows, "anonim")
	} else {
		chatLayout.room.Rows = append(chatLayout.room.Rows, nickname)
	}

	defer ws.Close()
	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			messageType, messageRaw, err := ws.ReadMessage()
			if err != nil {
				log.Warn(err)
				return
			}

			switch messageType {
			case websocket.TextMessage:
				// fmt.Println(string(messageRaw))
				chatLayout.output.Rows = append(chatLayout.output.Rows, string(messageRaw))
			default:
				log.Printf("%d, %s", messageType, messageRaw)
			}
		}
	}()

	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Escape>", "<C-c>":
				ui.Close()
				os.Exit(0)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				width, height := payload.Width, payload.Height
				chatLayout.grid.SetRect(0, 0, width, height)
				ui.Render(chatLayout.grid)
			}
			switch e.Type {
			case ui.KeyboardEvent:
				switch e.ID {
				case "<Up>":
					chatLayout.output.ScrollUp()
				case "<Down>":
					chatLayout.output.ScrollDown()
				case "<PageUp>":
					chatLayout.output.ScrollPageUp()
				case "<PageDown>":
					chatLayout.output.ScrollPageDown()
				case "<Home>":
					chatLayout.output.ScrollTop()
				case "<End>":
					chatLayout.output.ScrollBottom()
				case "<Backspace>":
					if len(chatLayout.input.Text) > 0 {
						chatLayout.input.Text = chatLayout.input.Text[:len(chatLayout.input.Text)-1]
					}
				case "<Enter>":
					if err := ws.WriteMessage(websocket.TextMessage, []byte(chatLayout.input.Text)); err != nil {
						log.Warn(err)
					}
					chatLayout.input.Text = ""
				case "<Space>":
					chatLayout.input.Text += " "
				default:
					chatLayout.input.Text += e.ID
				}
				ui.Render(chatLayout.grid)
			}

		case <-ticker.C:
			ui.Render(chatLayout.grid)
		}
	}
}
