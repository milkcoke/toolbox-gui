package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/imroc/req/v3"
	"github.com/milkcoke/auto-setup-gui/src/app"
	"image/color"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
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

	/**
	 * button, space is not resized in layout and container
	 * since it's inherited from container or layout
	 */
	var imgSize = fyne.Size{Width: 190, Height: 120}

	nodeImg := canvas.NewImageFromResource(nodeIcon)
	nodeImg.SetMinSize(imgSize)

	nodeProgress := widget.NewProgressBar()
	goProgress := widget.NewProgressBar()
	notionProgress := widget.NewProgressBar()
	dockerProgress := widget.NewProgressBar()

	client := req.C().SetOutputDirectory(widgetConfig.DownloadPath)

	callback := func(info req.DownloadInfo) {
		log.Printf("contentLength: %n\n", info.Response.ContentLength)
		if info.Response.Response != nil {
			log.Printf("downloadSize : %n\n", info.DownloadedSize)
			nodeProgress.SetValue(float64(info.DownloadedSize) / float64(info.Response.ContentLength))
		}
	}

	nodeImgBtn := widget.NewButtonWithIcon("Node.js", nodeIcon, func() {})
	nodeImgBtn.OnTapped = func() {
		nodeInst := app.NodeInstaller
		nodeImgBtn.Disable()
		nodeProgress.Show()
		go func() {
			defer nodeImgBtn.Enable()
			defer nodeProgress.Hide()
			// TODO
			//  Read https://req.cool/docs/tutorial/download/
			res, err := client.R().
				SetRetryCount(5).
				SetRetryFixedInterval(2*time.Second).
				SetOutputFile(nodeInst.Name+nodeInst.Ext).
				SetDownloadCallbackWithInterval(callback, 300*time.Millisecond).
				Get(nodeInst.Url)

			if err != nil {
				log.Printf("Failed to download : %s\n", nodeInst.Name)
				dialog.ShowInformation("Error", "Node.js download failed", win)
				return
			}

			if res.GetStatusCode() != http.StatusOK {
				log.Printf("Status code : %n\n", res.GetStatusCode())
				dialog.ShowInformation("Error", "Node.js download failed", win)
				return
			}

			dialog.ShowInformation("Success", "Node.js download complete", win)
		}()
	}

	goImgBtn := widget.NewButtonWithIcon("Go", golangIcon, func() {})
	notionImgBtn := widget.NewButtonWithIcon("Notion", notionIcon, func() {})
	dockerImgBtn := widget.NewButtonWithIcon("Docker", dockerIcon, func() {})

	nodeProgress.Hide()
	goProgress.Hide()
	notionProgress.Hide()
	dockerProgress.Hide()

	//var space = layout.NewSpacer()
	/**
	 * 모든 오브젝트의 크기가 CellSize (패딩포함) 으로 설정됨
	 * 그리고 MinSize 에서는 셀크기 넓이와 높이를 그대로 리턴함.
	 * 따라서 imgBtn 이 셀사이즈와 동일한 크기를 갖는지만 우선 확인해보자.
	 * MinSize 는 오브젝트가 가질 크기인데, 여기서 보면 셀사이즈 그대로 넓이를 반환함..
	 */
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

	//widgetConfig.initEventListener()

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
