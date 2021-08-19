package client

import (
	"os"
	"strings"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/gorilla/websocket"
)

type chatLayout struct {
	// Grid
	grid *ui.Grid
	// Room widget storing chat members
	room *widgets.List
	// Chat output storing chat messages
	output *widgets.List
	// Message input
	input *widgets.Paragraph
}

// Render UI layout
func renderLayout() (chatLayout, error) {
	if err := ui.Init(); err != nil {
		return chatLayout{}, err
	}

	termWidth, termHeight := ui.TerminalDimensions()
	chatRoom := widgets.NewList()
	chatOutput := widgets.NewList()
	chatInput := widgets.NewParagraph()

	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	grid.Set(
		ui.NewRow(1.0-0.2,
			ui.NewCol(0.2, chatRoom),
			ui.NewCol(1.0-0.2, chatOutput),
		),
		ui.NewRow(0.2, chatInput),
	)

	cl := chatLayout{
		grid:   grid,
		room:   chatRoom,
		output: chatOutput,
		input:  chatInput,
	}

	ui.Render(cl.grid)

	return cl, nil
}

// Handle UI keyboard events
func handleUiEvents(cl *chatLayout, c *client) error {
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "<Escape>", "<C-c>":
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				ui.Close()
				os.Exit(0)
			case "<Resize>":
				payload := e.Payload.(ui.Resize)
				width, height := payload.Width, payload.Height
				cl.grid.SetRect(0, 0, width, height)
				ui.Render(cl.grid)
			}
			switch e.Type {
			case ui.KeyboardEvent:
				switch e.ID {
				case "<Up>":
					cl.output.ScrollUp()
				case "<Down>":
					cl.output.ScrollDown()
				case "<PageUp>":
					cl.output.ScrollPageUp()
				case "<PageDown>":
					cl.output.ScrollPageDown()
				case "<Home>":
					cl.output.ScrollTop()
				case "<End>":
					cl.output.ScrollBottom()
				case "<Backspace>":
					if len(cl.input.Text) > 0 {
						cl.input.Text = cl.input.Text[:len(cl.input.Text)-1]
					}
				case "<Enter>":
					if err := c.conn.WriteMessage(websocket.TextMessage, []byte(strings.TrimSuffix(cl.input.Text, "\n"))); err != nil {
						return err
					}
					cl.input.Text = ""
				case "<Space>":
					cl.input.Text += " "
				default:
					cl.input.Text += e.ID
				}
				ui.Render(cl.grid)
			}

		case <-ticker.C:
			ui.Render(cl.grid)
		}
	}
}
