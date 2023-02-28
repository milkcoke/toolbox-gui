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
	"github.com/milkcoke/auto-setup-gui/src/app"
	"image/color"
	"io"
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

func (appWidget *appWidget) setEventListener(appConfig *AppConfig) {
	// C: Global client
	// Which have corresponding global wrappers
	// Just treat the package name req as Client to test, set up the Client without create any Client explicitly.
	// So don't use req.C() not for global configuration.
	// It affects on global request created by req.R()
	client := req.C().SetOutputDirectory(appConfig.DownloadPath)

	var fileFullPath = filepath.Join(appConfig.DownloadPath, appWidget.installerConfig.Name+appWidget.installerConfig.Ext)

	// 파일이 이미 존재
	if _, err := os.Stat(fileFullPath); err == nil {
		file, err := os.Open(fileFullPath)

		if err != nil {
			log.Fatalln(err)
		}
		defer file.Close()

		fileInfo, err := file.Stat()
		if err != nil {
			log.Fatalln(err)
		}

		headerRes, err := req.R().Head(appWidget.installerConfig.Url)
		if err != nil {
			log.Println("헤더 응답 오류 ", err)
		}
		// Convert string to int64
		contentLength := headerRes.GetHeader("Content-Length")
		fullFileLength, err := strconv.ParseInt(contentLength, 10, 64)
		if err != nil {
			log.Println("파일 크기 변환 에러 ", err)
		}
		log.Println("전체 파일 길이: ", contentLength)
		log.Println("현재 파일 크기: ", fileInfo.Size())

		// 도중에 중단된 파일 존재.
		if fullFileLength == fileInfo.Size() {
			// 이미 완성된 파일이라면, 해당 파일을 열어주고 끝내야함.
			// 단, 이벤트 등록은 여전히!
			os.Chdir(fileFullPath)
			if err != nil {
				log.Println("잘못된 디렉토리 경로")
			}
		} else {
			// 미완성 파일 존재

			// TODO
			//  여기서 SetHeader 를 갱신하고 재시도하는 로직이 필요함.
			res, err := req.R().SetHeader("Range", fmt.Sprintf("bytes=%d-", fileInfo.Size())).Get(appWidget.installerConfig.Url)

			// err 날 때 마다 재시도 조져야되나 하...
			if err != nil {
				log.Fatalln("또 ㅈㄹ이네 : ", err)
			}
			defer res.Body.Close()

			file, err = os.OpenFile(fileFullPath, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			_, err = io.Copy(file, res.Body)
			if err != nil {
				log.Fatalln("덮어쓰다 ㅈ창남 : ", err)
			}
		}

		return
	}

	callback := func(info req.DownloadInfo) {
		if info.Response.Response != nil {
			appWidget.progressBar.SetValue(float64(info.DownloadedSize) / float64(info.Response.ContentLength))
		}
	}

	// return type func() 로 하면 method 로도 eventListener 등록할 수 있음.
	appWidget.ImageButton.OnTapped = func() {
		appWidget.ImageButton.Disable()
		appWidget.progressBar.Show()

		go func() {
			defer appWidget.ImageButton.Enable()
			defer appWidget.progressBar.Hide()
			// TODO
			//  Read https://req.cool/docs/tutorial/download/
			res, err := client.R().
				SetRetryCount(5).
				SetRetryFixedInterval(2*time.Second).
				SetOutputFile(appWidget.installerConfig.Name+appWidget.installerConfig.Ext).
				SetDownloadCallbackWithInterval(callback, 300*time.Millisecond).
				Get(appWidget.installerConfig.Url)

			if err != nil {
				log.Printf("Failed to download : %s\n", appWidget.installerConfig.Name)
				dialog.ShowInformation("Error", appWidget.installerConfig.Name+" download failed", appWidget.parentWidget)
				return
			}

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
