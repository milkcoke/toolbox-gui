package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("M-Falcon")

	w.SetContent(widget.NewLabel("Hi laydies"))
	w.ShowAndRun()
}
