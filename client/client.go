package client

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os"

	// "github.com/fhodun/gochat/models"
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
	}
	return nil
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

	// log.Printf("connecting to %s", serverUrl.String())
	var resp *http.Response
	c.conn, resp, err = websocket.DefaultDialer.Dial(serverUrl.String(), header)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	newStr := buf.String()
	chatLayout.room.Rows = append(chatLayout.room.Rows, newStr)

	defer c.conn.Close()
	done := make(chan struct{}, 1)

	go func() {
		defer close(done)
		for {
			messageType, messageRaw, err := c.conn.ReadMessage()
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

	log.Fatal(handleUiEvents(&chatLayout, c.conn))
}
