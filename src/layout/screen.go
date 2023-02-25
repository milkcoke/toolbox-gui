package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"log"
	"os"
	"path/filepath"
)

type AppConfig struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	DownloadPath  string
	SaveMenuItem  *fyne.MenuItem
	ImageButtons  []*widget.Button
	Container     *fyne.Container
}

var Config AppConfig

func (widgetConfig *AppConfig) initDownloadDir(pathLabel *widget.Label) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error getting user home directory ", err)
		return
	}

	downloadDir := filepath.Join(homeDir, "Downloads")
	widgetConfig.DownloadPath = downloadDir
	pathLabel.SetText(downloadDir)
}

func (widgetConfig *AppConfig) MakeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("기모링딩동")

	widgetConfig.EditWidget = edit
	widgetConfig.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (widgetConfig *AppConfig) LoadImageButtons(win fyne.Window) (buttonContainer *fyne.Container) {
	pathLabel := widget.NewLabel("Download Path")
	widgetConfig.initDownloadDir(pathLabel)

	downloadDirPathBtn := widget.NewButton("Select path", widgetConfig.setDirectory(win, pathLabel))

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

	nodeImgBtn := widget.NewButtonWithIcon("Node.js", nodeIcon, func() {})
	goImgBtn := widget.NewButtonWithIcon("Go", golangIcon, func() {})
	notionImgBtn := widget.NewButtonWithIcon("Notion", notionIcon, func() {})
	dockerImgBtn := widget.NewButtonWithIcon("Docker", dockerIcon, func() {})

	/**
	 * button, space is not resized in layout and container
	 * since it's inherited from container or layout
	 */
	var imgSize = fyne.Size{Width: 640, Height: 240}
	nodeImgBtn.Resize(imgSize)

	nodeProgress := widget.NewProgressBar()
	goProgress := widget.NewProgressBar()
	notionProgress := widget.NewProgressBar()
	dockerProgress := widget.NewProgressBar()

	//nodeProgress.Hide()
	goProgress.Hide()
	notionProgress.Hide()
	dockerProgress.Hide()

	//var space = layout.NewSpacer()

	// 여기서 앱을 불러와서 New Window 를 띄워야함.
	//buttonsContainer := container.New(layout.NewGridLayout(3),
	buttonsContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(391, 240)),
		container.NewMax(nodeImgBtn, container.NewCenter(nodeProgress)),
		container.NewMax(goImgBtn, container.NewCenter(goProgress)),
		container.NewMax(notionImgBtn, container.NewCenter(notionProgress)),
		container.NewMax(dockerImgBtn, container.NewCenter(dockerProgress)),
		// space,
	)

	vboxContainer := container.NewVBox(
		container.NewHBox(downloadDirPathBtn, pathLabel),
		container.NewBorder(canvas.NewLine(color.White), nil, nil, nil),
		buttonsContainer,
	)

	widgetConfig.ImageButtons = []*widget.Button{
		nodeImgBtn, goImgBtn, notionImgBtn, dockerImgBtn,
	}

	widgetConfig.initEventListener()

	widgetConfig.Container = vboxContainer

	return widgetConfig.Container
}

func (widgetConfig *AppConfig) initEventListener() {
	for _, imgbtn := range widgetConfig.ImageButtons {
		addEventListener(imgbtn)
	}
}

func addEventListener(button *widget.Button) {
	button.OnTapped = func() {
		log.Printf("%s 실행\n", button.Icon.Name())
		if button.Disabled() {
			button.Enable()
		} else {
			button.Disable()
		}
	}
}

func (widgetConfig *AppConfig) setDirectory(win fyne.Window, pathLabel *widget.Label) func() {
	return func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if uri != nil {
				pathLabel.SetText(uri.Path())
				widgetConfig.DownloadPath = uri.Path()
			}
		}, win)
	}
}
