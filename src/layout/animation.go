package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"image/color"
	"time"
)

type AnimationApp struct {
	App fyne.App
}

// Only left-up and right-down allowed not shaking.
func (appWidget *appWidget) startShakeOnImageButton() {
	animation := canvas.NewPositionAnimation(fyne.NewPos(-5, -5), fyne.NewPos(5, 5), time.Millisecond*100, appWidget.ImageButton.Move)
	animation.AutoReverse = true
	animation.RepeatCount = 20
	animation.Start()
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
