package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"time"
)

type AnimationApp struct {
	App fyne.App
}

func (appConfig *AnimationApp) GenerateProgressBar() {
	var app = appConfig.App

	// 여기서부터
	myWindow := app.NewWindow("ProgressBar Widget")

	progress := widget.NewProgressBar()
	//infinite := widget.NewProgressBarInfinite()
	go func() {
		for i := 0.0; i <= 1.0; i += 0.1 {
			time.Sleep(time.Millisecond * 250)
			progress.SetValue(i)
		}
	}()
	myWindow.SetContent(container.NewCenter(progress))
	myWindow.Show()

	app.Settings().SetTheme(&CustomTheme{})
}

func (appConfig *AppConfig) GenerateAnimation() (animationContainer *fyne.Container) {
	var iconSize = fyne.Size{Width: 640, Height: 480}
	// Create the animation effect
	animation := canvas.NewRectangle(color.White)

	animation.Resize(iconSize)

	appConfig.Container.Add(animation)

	go func() {
		for {
			for i := 0; i < 4; i++ {
				animation.FillColor = theme.SuccessColor()
				canvas.Refresh(animation)
				time.Sleep(250 * time.Millisecond)

				animation.FillColor = color.Black
				canvas.Refresh(animation)
				time.Sleep(250 * time.Millisecond)
			}
		}
	}()

	return appConfig.Container
}
