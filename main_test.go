package main

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	markdown "github.com/milkcoke/fyne-app/src/layout"
	"testing"
)

// This is integration test

func Test_makeUI(t *testing.T) {
	var testConfig = markdown.Config

	edit, preview := testConfig.MakeUI()

	var txtMsg = "Falcon never can't die!"
	test.Type(edit, txtMsg)

	if preview.String() != txtMsg {
		t.Error("Failed -- didn't find expected value in preview")
	}
}

func Test_RunApp(t *testing.T) {
	var testConfig = markdown.Config

	testApp := test.NewApp()
	testWindow := testApp.NewWindow("Test layout")

	edit, preview := testConfig.MakeUI()
	testConfig.CreateMenuItems(testWindow)

	testWindow.SetContent(container.NewHSplit(edit, preview))
	testApp.Run()

	var txtMsg = "Some test"
	test.Type(edit, txtMsg)
	if preview.String() != txtMsg {
		t.Error("Failed to run app with correct text message.")
	}
}
