package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/milkcoke/auto-setup-gui/src/layout"
	"testing"
)

// This is integration test

func Test_DownlaodNode(t *testing.T) {
	testApp := test.NewApp()

	// create a window for the testApp
	testWindow := testApp.NewWindow("Test-Installer")

	var config = layout.Config
	var container = config.LoadImageButtons(testWindow)

	// show window and run testApp
	testWindow.Resize(fyne.Size{Width: 800, Height: 480})

	testWindow.SetContent(container)

	testApp.Run()
	config.ImageButtons[0].Tapped(&fyne.PointEvent{Position: fyne.Position{X: 46, Y: 30}})
}
