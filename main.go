package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/milkcoke/toolbox-gui/src/layout"
)

func main() {
	// create a fyne app
	app := app.New()

	// create a window for the app
	win := app.NewWindow("Falcon Toolbox")

	var config = layout.Config
	var container = config.LoadImageButtons(win)

	// show window and run app
	win.Resize(fyne.Size{Width: 800, Height: 480})

	win.SetContent(container)

	//win.CenterOnScreen()
	win.ShowAndRun()
}
