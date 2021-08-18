package client

import (
	"os"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/gorilla/websocket"
)

type chatLayout struct {
	grid   *ui.Grid
	room   *widgets.List
	output *widgets.List
	input  *widgets.Paragraph
}

func RenderLayout() (chatLayout, error) {
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

func handleUiEvents(cl *chatLayout, ws *websocket.Conn) error {
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
					if err := ws.WriteMessage(websocket.TextMessage, []byte(cl.input.Text)); err != nil {
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
