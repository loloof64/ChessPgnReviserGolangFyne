package history

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"fmt"
)

// HistoryLayout defines the layout of the History widget.
type HistoryLayout struct {
	width int
	gap   fyne.Size
}

func (l HistoryLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	currMaxH, w, h := 0, 0, 0

	for _, obj := range objects {
		childSize := obj.MinSize()

		wontlOverflowCurrentLine := w+childSize.Width <= l.width

		if wontlOverflowCurrentLine {
			w += childSize.Width + l.gap.Width
		} else {
			h += currMaxH + l.gap.Height
			w = childSize.Width
			currMaxH = 0
		}

		mustUpdateLineMaxHeight := childSize.Height > currMaxH

		if mustUpdateLineMaxHeight {
			currMaxH = childSize.Height
		}
	}

	// We must not forget last line !
	return fyne.NewSize(l.width, h+currMaxH)
}

func (l HistoryLayout) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	pos := fyne.NewPos(0, 0)
	w, h, currMaxH := 0, 0, 0

	for _, o := range objects {
		size := o.MinSize()
		o.Resize(size)

		weMustUpdateCurrentLineMaxHeight := size.Height > currMaxH
		if weMustUpdateCurrentLineMaxHeight {
			currMaxH = size.Height
		}

		weMustGoToNextLine := w+size.Width > containerSize.Width
		if weMustGoToNextLine {
			fmt.Println("new line", "pos", pos)
			pos = pos.Add(fyne.NewPos(0, currMaxH+l.gap.Height))
			pos.X = 0
			o.Move(pos)
			w = 0
			h += currMaxH + l.gap.Height
		} else {
			fmt.Println("going on", "pos", pos)
			o.Move(pos)
			pos = pos.Add(fyne.NewPos(size.Width+l.gap.Width, 0))
			w += size.Width + l.gap.Width
		}
	}
}

func newHistoryLayout(width int) HistoryLayout {
	return HistoryLayout{width: width, gap: fyne.NewSize(5, 8)}
}

// GameMove defines a move of the History widget.
type GameMove struct {
	Fan string
}

// History is a widget that shows the played moves, and is intended to
// load selected position on the board if game is not in progress.
type History struct {
	widget.BaseWidget

	preferredSize fyne.Size
	moves         []GameMove

	container *fyne.Container
}

type historyRenderer struct {
	history *History
}

func (renderer *historyRenderer) MinSize() fyne.Size {
	return renderer.history.container.MinSize()
}

func (renderer *historyRenderer) Layout(size fyne.Size) {
	renderer.history.container.Layout.Layout(renderer.history.container.Objects, size)
}

func (renderer *historyRenderer) ApplyTheme() {

}

func (renderer *historyRenderer) BackgroundColor() color.Color {
	return theme.BackgroundColor()
}

func (renderer *historyRenderer) Refresh() {
	canvas.Refresh(renderer.history.container)
}

func (renderer *historyRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{renderer.history.container}
}

func (renderer *historyRenderer) Destroy() {

}

func NewHistory(preferredSize fyne.Size) *History {
	history := &History{preferredSize: preferredSize}
	history.ExtendBaseWidget(history)

	history.container = fyne.NewContainerWithLayout(newHistoryLayout(preferredSize.Width))

	return history
}

// CreateRenderer creates the Renderer for History widget.
func (history *History) CreateRenderer() fyne.WidgetRenderer {
	renderer := &historyRenderer{history: history}
	return renderer
}

// AddMove adds a move to the History widget.
func (history *History) AddMove(move GameMove) {
	history.moves = append(history.moves, move)
	moveComponent := widget.NewButton(move.Fan, func() {})
	history.container.AddObject(moveComponent)
	history.container.Resize(history.preferredSize)
	history.Refresh()
}

// Clear clears all moves from the History widget.
func (history *History) Clear() {
	history.moves = nil
	history.container.Objects = nil
	history.Refresh()
}
