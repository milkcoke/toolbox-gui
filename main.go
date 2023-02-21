package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/milkcoke/fyne-app/src/markdown"
)

func main() {
	// create a fyne app
	app := app.New()

	// create a window for the app
	win := app.NewWindow("Markdown")

	markdown.Config.CreateMenuItems(win)
	var container = markdown.Config.LoadImageButtons(win)

	win.SetContent(container)

	// show window and run app
	win.Resize(fyne.Size{Width: 800, Height: 480})
	//win.CenterOnScreen()
	win.ShowAndRun()
}
