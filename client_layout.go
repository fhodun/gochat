package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type chatLayoutStruct struct {
	grid   *ui.Grid
	room   *widgets.List
	output *widgets.List
	input  *widgets.Paragraph
}

func RenderLayout() (chatLayoutStruct, error) {
	if err := ui.Init(); err != nil {
		return chatLayoutStruct{}, err
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

	chatLayout := chatLayoutStruct{
		grid:   grid,
		room:   chatRoom,
		output: chatOutput,
		input:  chatInput,
	}

	ui.Render(chatLayout.grid)

	return chatLayout, nil
}
