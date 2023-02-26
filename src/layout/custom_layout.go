package layout

import (
	"fyne.io/fyne/v2"
)

var _ fyne.Layout = (*customMaxLayout)(nil)

type customMaxLayout struct {
}

func (m *customMaxLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	topLeft := fyne.NewPos(0, 0)
	for _, child := range objects {
		//log.Println("여기가 layout Size:", size)
		//log.Println("Child size Size:", child.Size())
		child.Resize(size)
		child.Move(topLeft)
	}
}

func (m *customMaxLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}

		minSize = child.Size()
		//log.Println("minSize: ", minSize)
		//log.Println("child MinSize : ", child.MinSize())
		child.Resize(minSize)
	}

	return minSize
}
