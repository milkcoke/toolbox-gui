package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type App struct {
	output *widget.Label
}

func main() {
	app := app.New()
	window := app.NewWindow("M-Falcon")

	var myApp App
	output, entry, btn := myApp.makeUI()

	window.SetContent(container.NewVBox(output, entry, btn))
	window.Resize(fyne.Size{Width: 500, Height: 500})
	window.ShowAndRun() // stop and listen event loop internally
}

func (app *App) makeUI() (*widget.Label, *widget.Entry, *widget.Button) {
	output := widget.NewLabel("Hi ladies")
	entry := widget.NewEntry()
	btn := widget.NewButton("Enter", func() {
		app.output.SetText(entry.Text)
	})

	btn.Importance = widget.HighImportance

	app.output = output

	return output, entry, btn
}
