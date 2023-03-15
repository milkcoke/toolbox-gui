package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/milkcoke/toolbox-gui/src/layout"
)

func main() {
	app := app.New()

	win := app.NewWindow("Dev Toolbox")

	var config = layout.Config
	var container = config.LoadImageButtons(win)

	win.Resize(fyne.Size{Width: 800, Height: 480})

	win.SetContent(container)

	win.ShowAndRun()
}
