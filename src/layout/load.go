package layout

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/imroc/req/v3"
	"github.com/inhies/go-bytesize"
	"github.com/milkcoke/auto-setup-gui/src/app"
	filehandle "github.com/milkcoke/auto-setup-gui/src/file"
	"image/color"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type AppConfig struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	// TODO
	//  setDirectory 할때 모든 appWidget 에 대한 경로 업데이트 해줘야함. (eventListener 혹은 동적인 다운로드 경로)
	DownloadPath string
	SaveMenuItem *fyne.MenuItem
	AppWidgets   []*appWidget
	Container    *fyne.Container
}

// Related data needs to be handled or manipulated as a unit.
// If method should accept multiple variables,  it's good that wrapping data idea as struct .
type appWidget struct {
	ImageButton     *widget.Button
	installerConfig app.InstallerConfig
	progressBar     *widget.ProgressBar
	parentWidget    fyne.Window
}

var Config AppConfig

func (appConfig *AppConfig) initDownloadDir(pathLabel *widget.Label) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Println("Error getting user home directory ", err)
		return
	}

	downloadDir := filepath.Join(homeDir, "Downloads")
	appConfig.DownloadPath = downloadDir
	pathLabel.SetText(downloadDir)
}

// asyncRetryDownload
// This is called only when file not-exist or exist but partial.
func asyncRetryDownload(readFileFD *os.File, appWidget *appWidget, fullFileLength int64) {

	// Check file existence
	retryFileInfo, err := readFileFD.Stat()
	if err != nil {
		log.Fatalln("Failed to open file : ", appWidget.installerConfig.Name)
	}

	// Check file download complete
	// This is for protecting code for recursive function
	if fullFileLength == retryFileInfo.Size() {
		readFileFD.Close()
		// only this printed when download complete without checking file size
		dialog.ShowInformation("Success after retrying : ", appWidget.installerConfig.Name+" download complete", appWidget.parentWidget)
		return
	}

	streamFile, err := os.OpenFile(readFileFD.Name(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open readFileFD : ", err)
	}
	defer streamFile.Close()
	appWidget.progressBar.Show()

	//var startSize = retryFileInfo.Size()
	callback := func(info req.DownloadInfo) {
		if info.Response.Response != nil {
			appWidget.progressBar.SetValue(float64(info.DownloadedSize) / float64(info.Response.ContentLength))
		}
	}

	res, err := req.R().
		SetDownloadCallbackWithInterval(callback, 300*time.Millisecond).
		SetHeader("Range", fmt.Sprintf("bytes=%d-", retryFileInfo.Size())).
		SetOutput(streamFile).
		Get(appWidget.installerConfig.Url)

	defer res.Body.Close()

	if err != nil {
		log.Println("Http request fail : ", err)
		time.Sleep(2 * time.Second)
		go asyncRetryDownload(readFileFD, appWidget, fullFileLength)
		return
	}

	readFileFD.Close()
	dialog.ShowInformation("Success after retrying : ", appWidget.installerConfig.Name+" download complete", appWidget.parentWidget)

}

func (appWidget *appWidget) setEventListener(appConfig *AppConfig) {

	// Tool installer total size init
	headerRes, err := req.R().Head(appWidget.installerConfig.Url)
	if err != nil {
		log.Println("Failed to request head ", err)
	}
	contentLength := headerRes.GetHeader("Content-Length")
	fullFileSize, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		log.Println("Failed to parse int of total size of tool installer ", err)
	}

	appWidget.ImageButton.OnTapped = func() {
		appWidget.ImageButton.Disable()

		defer appWidget.ImageButton.Enable()
		defer appWidget.progressBar.Hide()
		// return type func() 로 하면 method 로도 eventListener 등록할 수 있음.

		// C: Global client
		// Which have corresponding global wrappers
		// Just treat the package name req as Client to test, set up the Client without create any Client explicitly.
		// So don't use req.C() not for global configuration.
		// It affects on global request created by req.R()
		client := req.C().SetOutputDirectory(appConfig.DownloadPath)

		var fileFullPath = filepath.Join(appConfig.DownloadPath, appWidget.installerConfig.Name+appWidget.installerConfig.Ext)

		// Already file exists
		if _, err := os.Stat(fileFullPath); err == nil {
			readFileFD, err := os.Open(fileFullPath)

			if err != nil {
				log.Fatalln(err)
			}

			fileInfo, err := readFileFD.Stat()
			if err != nil {
				log.Fatalln(err)
			}
			log.Println("전체 파일 길이: ", bytesize.New(float64(fullFileSize)))
			log.Println("현재 파일 크기: ", bytesize.New(float64(fileInfo.Size())))

			// Already installer is installed.
			if fullFileSize == fileInfo.Size() {
				filehandle.NavigateToDir(fileFullPath)
				if err != nil {
					log.Println("Invalid directory path : ", err)
				}
				return
			} else {
				go asyncRetryDownload(readFileFD, appWidget, fullFileSize)
				return
			}
		}

		// First download request
		callback := func(info req.DownloadInfo) {
			if info.Response.Response != nil {
				appWidget.progressBar.SetValue(float64(info.DownloadedSize) / float64(info.Response.ContentLength))
			}
		}

		go func() {
			appWidget.progressBar.Show()
			res, err := client.R().
				SetDownloadCallbackWithInterval(callback, 300*time.Millisecond).
				SetOutputFile(appWidget.installerConfig.Name + appWidget.installerConfig.Ext).
				Get(appWidget.installerConfig.Url)

			if err != nil {
				log.Println("Failed to download : ", err)
				readFileFD, err := os.Open(fileFullPath)
				defer readFileFD.Close()
				if err != nil {
					log.Println("Failed to open readFileFD : ", readFileFD)
					log.Println("error : ", err)
				}
				go asyncRetryDownload(readFileFD, appWidget, fullFileSize)
				return
			}

			// 요청 응답 코드가 다른 경우 그냥 중단.
			if res.GetStatusCode() != http.StatusOK {
				log.Printf("Status code : %d\n", res.GetStatusCode())
				dialog.ShowInformation("Error", appWidget.installerConfig.Name+" download failed", appWidget.parentWidget)
				return
			}

			dialog.ShowInformation("Success", appWidget.installerConfig.Name+" download complete", appWidget.parentWidget)
		}()
	}

}

func (appConfig *AppConfig) MakeUI() (*widget.Entry, *widget.RichText) {
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("기모링딩동")

	appConfig.EditWidget = edit
	appConfig.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
}

func (appConfig *AppConfig) LoadImageButtons(win fyne.Window) (buttonContainer *fyne.Container) {
	pathLabel := widget.NewLabel("Download Path")
	appConfig.initDownloadDir(pathLabel)

	downloadDirPathBtn := widget.NewButton("Select path", appConfig.setDirectory(win, pathLabel))

	// Load the icon image from a file
	nodeIcon, err := fyne.LoadResourceFromPath("assets/nodejs_logo.svg")
	if err != nil {
		log.Printf("Failed to load %s icon\n", nodeIcon.Name())
	}

	golangIcon, err := fyne.LoadResourceFromPath("assets/go_logo_aqua.svg")
	if err != nil {
		log.Printf("Failed to load %s icon\n", golangIcon.Name())
	}

	notionIcon, err := fyne.LoadResourceFromPath("assets/notion.svg")
	if err != nil {
		log.Printf("Failed to load %s icon\n", notionIcon.Name())
	}

	dockerIcon, err := fyne.LoadResourceFromPath("assets/docker.svg")
	if err != nil {
		log.Printf("Failed to load %s icon\n", dockerIcon.Name())
	}

	/**
	 * button, space is not resized in layout and container
	 * since it's inherited from container or layout
	 */

	nodeProgress := widget.NewProgressBar()
	goProgress := widget.NewProgressBar()
	notionProgress := widget.NewProgressBar()
	dockerProgress := widget.NewProgressBar()

	nodeProgress.Hide()
	goProgress.Hide()
	notionProgress.Hide()
	dockerProgress.Hide()

	nodeImgBtn := widget.NewButtonWithIcon("Node.js", nodeIcon, func() {})
	goImgBtn := widget.NewButtonWithIcon("Go", golangIcon, func() {})
	notionImgBtn := widget.NewButtonWithIcon("Notion", notionIcon, func() {})
	dockerImgBtn := widget.NewButtonWithIcon("Docker", dockerIcon, func() {})

	nodeAppWidget := &appWidget{
		nodeImgBtn, app.NodeInstaller, nodeProgress, win,
	}
	goAppWidget := &appWidget{
		goImgBtn, app.GoInstaller, goProgress, win,
	}
	notionAppWidget := &appWidget{
		notionImgBtn, app.NotionInstaller, notionProgress, win,
	}

	dockerAppWidget := &appWidget{
		dockerImgBtn, app.DockerInstaller, dockerProgress, win,
	}

	nodeAppWidget.setEventListener(appConfig)
	goAppWidget.setEventListener(appConfig)
	notionAppWidget.setEventListener(appConfig)
	dockerAppWidget.setEventListener(appConfig)

	buttonsContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(391, 240)),
		container.NewMax(nodeImgBtn, container.NewCenter(nodeProgress)),
		container.NewMax(goImgBtn, container.NewCenter(goProgress)),
		container.NewMax(notionImgBtn, container.NewCenter(notionProgress)),
		container.NewMax(dockerImgBtn, container.NewCenter(dockerProgress)),
	)

	vboxContainer := container.NewVBox(
		container.NewHBox(downloadDirPathBtn, pathLabel),
		container.NewBorder(canvas.NewLine(color.White), nil, nil, nil),
		buttonsContainer,
	)

	appConfig.AppWidgets = []*appWidget{
		nodeAppWidget, goAppWidget, notionAppWidget, dockerAppWidget,
	}

	appConfig.Container = vboxContainer

	return appConfig.Container
}

func (appConfig *AppConfig) setDirectory(win fyne.Window, pathLabel *widget.Label) func() {
	return func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if uri == nil {
				return
			}

			pathLabel.SetText(uri.Path())
			appConfig.DownloadPath = uri.Path()

			// refresh new download path
			for _, appWidget := range appConfig.AppWidgets {
				appWidget.setEventListener(appConfig)
			}
		}, win)
	}
}
