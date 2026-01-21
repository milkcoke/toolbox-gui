package test

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/imroc/req/v3"
	"github.com/inhies/go-bytesize"
	"github.com/milkcoke/toolbox-gui/src/app"
	assets "github.com/milkcoke/toolbox-gui/src/assets"
	filehandle "github.com/milkcoke/toolbox-gui/src/file"
)

type AppConfig struct {
	EditWidget    *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile   fyne.URI
	DownloadPath  string
	SaveMenuItem  *fyne.MenuItem
	AppWidgets    []*appWidget
	Container     *fyne.Container
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
	appWidget.ImageButton.Disable()
	defer appWidget.ImageButton.Enable()
	// Check file existence
	retryFileInfo, err := readFileFD.Stat()
	if err != nil {
		log.Println("Failed to open file : ", appWidget.installerConfig.Name)
		time.Sleep(2 * time.Second)
		go asyncRetryDownload(readFileFD, appWidget, fullFileLength)
		return
	}

	// Check file download complete
	// This is for protecting code for recursive function
	if fullFileLength == retryFileInfo.Size() {
		readFileFD.Close()
		// only this printed when download complete without checking file size
		dialog.ShowInformation("Success after retrying : ", appWidget.installerConfig.Name+" download complete", appWidget.parentWidget)
		appWidget.updateButtonStatus(app.CompleteDownload)
		return
	}

	streamFile, err := os.OpenFile(readFileFD.Name(), os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln("Failed to open readFileFD : ", err)
	}
	defer streamFile.Close()
	appWidget.progressBar.Show()

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
	appWidget.updateButtonStatus(app.CompleteDownload)
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

		var fileFullPath = filepath.Join(appConfig.DownloadPath, appWidget.installerConfig.Name+appWidget.installerConfig.Ext)

		// C: Global client
		// Which have corresponding global wrappers
		// Just treat the package name req as Client to test, set up the Client without create any Client explicitly.
		// So don't use req.C() not for global configuration.
		// It affects on global request created by req.R()
		client := req.C().SetOutputDirectory(appConfig.DownloadPath)

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
				err := filehandle.NavigateToDir(fileFullPath)
				if err != nil {
					log.Println("Invalid directory path : ", err)
				}
				appWidget.updateButtonStatus(app.OpenInstaller)
				readFileFD.Close()
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
			appWidget.ImageButton.Disable()
			defer appWidget.ImageButton.Enable()

			appWidget.progressBar.Show()
			res, err := client.R().
				SetDownloadCallbackWithInterval(callback, 300*time.Millisecond).
				SetOutputFile(appWidget.installerConfig.Name + appWidget.installerConfig.Ext).
				Get(appWidget.installerConfig.Url)

			if err != nil {
				log.Println("Failed to download : ", err)
				readFileFD, err := os.Open(fileFullPath)
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

			//TODO: Register MouseIn() hover event
			//mouseEvent := &desktop.MouseEvent{}
			//appWidget.ImageButton.MouseIn(mouseEvent)

			dialog.ShowInformation("Success", appWidget.installerConfig.Name+" download complete", appWidget.parentWidget)
			appWidget.updateButtonStatus(app.CompleteDownload)
		}()

	}

}

func (appWidget *appWidget) startShakeOnImageButton() {
	animation := canvas.NewPositionAnimation(fyne.NewPos(-5, -5), fyne.NewPos(5, 5), time.Millisecond*100, appWidget.ImageButton.Move)
	animation.AutoReverse = true
	animation.RepeatCount = 20
	animation.Start()
}

func (appConfig *AppConfig) LoadImageButtons(win fyne.Window) (buttonContainer *fyne.Container) {
	pathLabel := widget.NewLabel("Download Path")
	appConfig.initDownloadDir(pathLabel)

	downloadDirPathBtn := widget.NewButton("Select path", appConfig.setDirectory(win, pathLabel))
	pythonIcon := fyne.NewStaticResource("Python", assets.PythonBytes)
	nodeIcon := fyne.NewStaticResource("Node.js", assets.NodeBytes)
	golangIcon := fyne.NewStaticResource("Go", assets.GoBytes)
	dockerIcon := fyne.NewStaticResource("Docker", assets.DockerBytes)

	pythonProgress := widget.NewProgressBar()
	nodeProgress := widget.NewProgressBar()
	goProgress := widget.NewProgressBar()
	dockerProgress := widget.NewProgressBar()

	goProgress.Hide()
	dockerProgress.Hide()
	pythonProgress.Hide()
	nodeProgress.Hide()

	pythonImgBtn := widget.NewButtonWithIcon("Python", pythonIcon, func() {})
	nodeImgBtn := widget.NewButtonWithIcon("Node.js", nodeIcon, func() {})
	goImgBtn := widget.NewButtonWithIcon("Go", golangIcon, func() {})
	dockerImgBtn := widget.NewButtonWithIcon("Docker", dockerIcon, func() {})

	pythonAppWidget := &appWidget{
		pythonImgBtn, app.PythonInstaller, pythonProgress, win,
	}
	nodeAppWidget := &appWidget{
		nodeImgBtn, app.NodeInstaller, nodeProgress, win,
	}
	goAppWidget := &appWidget{
		goImgBtn, app.GoInstaller, goProgress, win,
	}
	dockerAppWidget := &appWidget{
		dockerImgBtn, app.DockerInstaller, dockerProgress, win,
	}

	pythonAppWidget.setEventListener(appConfig)
	nodeAppWidget.setEventListener(appConfig)
	goAppWidget.setEventListener(appConfig)
	dockerAppWidget.setEventListener(appConfig)

	buttonsContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(391, 240)),
		container.NewMax(pythonImgBtn, container.NewCenter(pythonProgress)),
		container.NewMax(nodeImgBtn, container.NewCenter(nodeProgress)),
		container.NewMax(goImgBtn, container.NewCenter(goProgress)),
		container.NewMax(dockerImgBtn, container.NewCenter(dockerProgress)),
	)

	vboxContainer := container.NewVBox(
		container.NewHBox(downloadDirPathBtn, pathLabel),
		container.NewBorder(canvas.NewLine(color.White), nil, nil, nil),
		buttonsContainer,
	)

	appConfig.AppWidgets = []*appWidget{
		pythonAppWidget, nodeAppWidget, goAppWidget, dockerAppWidget,
	}

	appConfig.Container = vboxContainer

	return appConfig.Container
}

func (appWidget *appWidget) updateButtonStatus(appStatus app.AppStatus) {
	switch appStatus {
	case app.None, app.PartialDownloaded:
		appWidget.ImageButton.Importance = widget.MediumImportance
	case app.Downloading:
		appWidget.ImageButton.Importance = widget.LowImportance
	case app.CompleteDownload:
		appWidget.progressBar.Hide()
		appWidget.ImageButton.Importance = widget.HighImportance
	case app.OpenInstaller:
		appWidget.ImageButton.Importance = widget.LowImportance
	}
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
				appWidget.updateButtonStatus(app.None)
			}
		}, win)
	}
}
