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

	app.Settings().SetTheme(&markdown.CustomTheme{})

	//get the user interface
	// 공유 Config 는 무조건 하나의 값일까?
	//edit, preview := markdown.Config.MakeUI()
	//imgBtn := markdown.Config.LoadImageButtons()

	markdown.Config.CreateMenuItems(win)
	var container = markdown.Config.LoadImageButtons(win)

	// set the content of the window
	// Split two widgets
	//win.SetContent(container.NewHSplit(edit, preview))

	// set container
	win.SetContent(container)

	// show window and run app
	win.Resize(fyne.Size{Width: 800, Height: 500})
	//win.CenterOnScreen()
	win.ShowAndRun()
}
