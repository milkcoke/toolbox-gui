package markdown

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"io"
	"log"
	"strings"
)

type config struct {
	EditWidget      *widget.Entry
	PreviewWidget   *widget.RichText
	CurrentFile     fyne.URI
	SaveMenuItem    *fyne.MenuItem
	ButtonContainer *fyne.Container
}

var Config config

func (widgetConfig *config) MakeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("기모링딩동")

	widgetConfig.EditWidget = edit
	widgetConfig.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (widgetConfig *config) LoadImageButtons() (buttonContainer *fyne.Container) {
	// Load the icon image from a file
	nodeIcon, err := fyne.LoadResourceFromPath("assets/nodejs_logo.svg")
	if err != nil {
		log.Printf("Error to retrieve %s svg\n", nodeIcon.Name())
	}

	golangIcon, err := fyne.LoadResourceFromPath("assets/go_logo_aqua.svg")
	if err != nil {
		log.Printf("Error to retrieve %s svg\n", golangIcon.Name())
	}

	notionIcon, err := fyne.LoadResourceFromPath("assets/notion.svg")
	if err != nil {
		log.Printf("Error to retrieve %s svg\n", notionIcon.Name())
	}

	dockerIcon, err := fyne.LoadResourceFromPath("assets/docker.svg")
	if err != nil {
		log.Printf("Error to retrieve %s svg\n", dockerIcon.Name())
	}

	nodeImgBtn := widget.NewButtonWithIcon("Node.js", nodeIcon, func() {
		log.Printf("%s 실행\n", nodeIcon.Name())
	})

	goImgBtn := widget.NewButtonWithIcon("Go", golangIcon, func() {
		log.Printf("%s 실행\n", golangIcon.Name())
	})

	notionImgBtn := widget.NewButtonWithIcon("Notion", notionIcon, func() {
		log.Printf("%s 실행\n", notionIcon.Name())
	})

	dockerImgBtn := widget.NewButtonWithIcon("Docker", dockerIcon, func() {
		log.Printf("%s 실행\n", dockerIcon.Name())
	})

	var iconSize = fyne.Size{Width: 640, Height: 480}

	nodeImgBtn.Resize(iconSize)
	goImgBtn.Resize(iconSize)
	notionImgBtn.Resize(iconSize)
	dockerImgBtn.Resize(iconSize)

	buttonsContainer := container.New(layout.NewGridLayout(3),
		nodeImgBtn,
		goImgBtn,
		notionImgBtn,
		dockerImgBtn,
	)
	widgetConfig.ButtonContainer = buttonsContainer

	return widgetConfig.ButtonContainer
}

var fileExtensionFilter = storage.NewExtensionFileFilter([]string{".md", ".MD"})

func (widgetConfig *config) CreateMenuItems(win fyne.Window) {
	openMenuItem := fyne.NewMenuItem("Open ...", widgetConfig.openFunc(win))

	saveMenuItem := fyne.NewMenuItem("Save", widgetConfig.saveFunc(win))
	// binding app widget to men item
	widgetConfig.SaveMenuItem = saveMenuItem
	// default disabled
	widgetConfig.SaveMenuItem.Disabled = true

	saveAsMenuItem := fyne.NewMenuItem("Save as", widgetConfig.saveAsFunc(win))

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}

func (widgetConfig *config) saveFunc(win fyne.Window) func() {
	return func() {
		// current file 열려있는 상태에만 저장 가능.
		if widgetConfig.CurrentFile != nil {
			writer, err := storage.Writer(widgetConfig.CurrentFile)

			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			writer.Write([]byte(widgetConfig.EditWidget.Text))
			dialog.ShowInformation("File", "The file saved successfully!", win)

			defer writer.Close()
		}
	}
}

func (widgetConfig *config) saveAsFunc(win fyne.Window) func() {
	return func() {
		saveDialog := dialog.NewFileSave(func(uriWriteCloser fyne.URIWriteCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// Check if the user cancel this operation
			if uriWriteCloser == nil {
				// user cancel
				return // application would not die
			}

			if !strings.HasSuffix(strings.ToLower(uriWriteCloser.URI().String()), ".md") {
				dialog.ShowInformation("[Error]", "This file extension is not allowed to save.", win)
				return
			}

			// save the file
			uriWriteCloser.Write([]byte(widgetConfig.EditWidget.Text))
			// keep track of what file is currently open
			// Supporting cross-platform,  we need using repository URI pattern.
			// That's why we use Fyne.URI pattern.
			widgetConfig.CurrentFile = uriWriteCloser.URI()

			defer uriWriteCloser.Close()

			// window title change.
			win.SetTitle(win.Title() + " - " + uriWriteCloser.URI().Name())
			widgetConfig.SaveMenuItem.Disabled = false

		}, win)

		// Only markdown file is allowed to save
		saveDialog.SetFilter(fileExtensionFilter)
		saveDialog.Show()
	}
}

func (widgetConfig *config) openFunc(win fyne.Window) func() {
	return func() {
		openDialog := dialog.NewFileOpen(func(uriReadCloser fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			if uriReadCloser == nil {
				return
			}

			defer uriReadCloser.Close()

			data, err := io.ReadAll(uriReadCloser)

			if err != nil {
				dialog.ShowError(err, win)
				return
			}

			// current file data set (text)
			widgetConfig.EditWidget.SetText(string(data))

			// Set current file name on the widget
			widgetConfig.CurrentFile = uriReadCloser.URI()

			// 앱 위젯 제목
			win.SetTitle(win.Title() + " - " + uriReadCloser.URI().Name())
			widgetConfig.SaveMenuItem.Disabled = false

		}, win)

		// Only markdown files are shown and open.
		openDialog.SetFilter(fileExtensionFilter)
		openDialog.Show()
	}
}
