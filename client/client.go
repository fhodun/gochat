package client

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fhodun/gochat/packets"
	ui "github.com/gizak/termui/v3"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type client struct {
	nickname string
	conn     *websocket.Conn
}

func retreiveVariableFromInput(n *string, message string) error {
	for *n == "" {
		var err error
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(message)
		*n, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		*n = strings.TrimSuffix(*n, "\n")
	}

	return nil
}

func findAndDelete(array []string, item string) []string {
	index := 0
	for _, i := range array {
		if i != item {
			array[index] = i
			index++
		}
	}
	return array[:index]
}

func RunClient(cmd *cobra.Command, args []string) {
	var (
		host      string      = args[0]
		header    http.Header = make(http.Header)
		serverUrl url.URL     = url.URL{
			Scheme: "ws",
			Host:   host,
			Path:   "/",
		}
		c *client = &client{
			nickname: "",
		}
	)

	if err := retreiveVariableFromInput(&c.nickname, "Enter nickname: "); err != nil {
		log.Fatal(err)
	}
	header.Add("nickname", c.nickname)

	chatLayout, err := RenderLayout()
	if err != nil {
		log.Fatal(err)
	}

	c.conn, _, err = websocket.DefaultDialer.Dial(serverUrl.String(), header)
	if err != nil {
		log.Fatal(err)
	}

	defer ui.Close()
	defer c.conn.Close()
	done := make(chan struct{}, 1)

	go func() {
		defer close(done)
		for {
			var message packets.Message
			if err := c.conn.ReadJSON(&message); err != nil {
				log.Warn(err)
				return
			}

			switch message.Type {
			case packets.NewMessage:
				chatLayout.output.Rows = append(chatLayout.output.Rows, "["+message.Author+"]: "+message.Text)
			case packets.ChattersList:
				chatLayout.room.Rows = append(chatLayout.room.Rows, message.Content...)
			case packets.RegisterChatter:
				chatLayout.room.Rows = append(chatLayout.room.Rows, message.Text)
			case packets.UnregisterChatter:
				chatLayout.room.Rows = findAndDelete(chatLayout.room.Rows, message.Text)
			default:
				log.WithFields(log.Fields{
					"message.Type":   message.Type,
					"message.Author": message.Author,
					"message.Text":   message.Text,
				}).Info("Unknown message type")
			}
		}
	}()

	log.Fatal(handleUiEvents(&chatLayout, c))
}
