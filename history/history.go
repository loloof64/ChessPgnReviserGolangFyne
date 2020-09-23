package history

import (
	"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/loloof64/chess-pgn-reviser-fyne/commonTypes"
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
			pos = pos.Add(fyne.NewPos(0, currMaxH+l.gap.Height))
			pos.X = 0
			o.Move(pos)
			pos = pos.Add(fyne.NewPos(size.Width+l.gap.Width, 0))
			w = size.Width + l.gap.Width
			h += currMaxH + l.gap.Height
		} else {
			o.Move(pos)
			pos = pos.Add(fyne.NewPos(size.Width+l.gap.Width, 0))
			w += size.Width + l.gap.Width
		}
	}
}

func newHistoryLayout(width int) HistoryLayout {
	return HistoryLayout{width: width, gap: fyne.NewSize(5, 8)}
}

// History is a widget that shows the played moves, and is intended to
// load selected position on the board if game is not in progress.
type History struct {
	widget.BaseWidget

	preferredSize     fyne.Size
	currentMoveNumber int

	container         *fyne.Container
	onPositionRequest func(moveData commonTypes.GameMove) bool

	currentHighlightedButtonIndex int
	currentMoveDataIndex          int

	allMovesData []commonTypes.GameMove
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

// SetOnPositionRequestHandler
// It should return whether the move could be processed (good data and game not in progress).
func (history *History) SetOnPositionRequestHandler(handler func(moveData commonTypes.GameMove) bool) {
	history.onPositionRequest = handler
}

// CreateRenderer creates the Renderer for History widget.
func (history *History) CreateRenderer() fyne.WidgetRenderer {
	renderer := &historyRenderer{history: history}
	return renderer
}

// AddMove adds a move to the History widget.
func (history *History) AddMove(moveData commonTypes.GameMove) {
	var moveComponent *widget.Button
	history.currentHighlightedButtonIndex += 1
	thisButtonIndex := history.currentHighlightedButtonIndex

	moveComponent = widget.NewButton(moveData.Fan, func() {
		if history.onPositionRequest != nil {
			if history.onPositionRequest(moveData) {
				history.currentHighlightedButtonIndex = thisButtonIndex
				history.updateButtonsStyles()
			}
		}
	})
	history.container.AddObject(moveComponent)
	history.container.Resize(history.preferredSize)
	history.allMovesData = append(history.allMovesData, moveData)
	if moveData.IsBlackMove {
		history.currentMoveNumber += 1
		// Though it is not a button, we need to update this index
		// as buttons and labels are in the same array.
		history.currentHighlightedButtonIndex += 1
		numberComponent := widget.NewLabel(fmt.Sprintf("%v.", history.currentMoveNumber))
		history.container.AddObject(numberComponent)
		history.container.Resize(history.preferredSize)
	}
	history.Refresh()
}

// Clear clears all moves from the History widget.
func (history *History) Clear(startMoveNumber int) {
	history.container.Objects = nil
	history.allMovesData = nil
	// Though it is not a button, we need to update this index
	// as buttons and labels are in the same array.
	history.currentHighlightedButtonIndex = 0
	history.currentMoveDataIndex = -1
	history.currentMoveNumber = startMoveNumber
	numberComponent := widget.NewLabel(fmt.Sprintf("%v.", history.currentMoveNumber))
	history.container.AddObject(numberComponent)
	history.container.Resize(history.preferredSize)
	history.Refresh()
}

// Tries to select the last element.
func (history *History) RequestLastItemSelection() {
	lastButtonIndex := history.findLastButtonIndex()
	lastMoveDataIndex := history.findLastMoveDataIndex()

	if lastButtonIndex < 0 || lastMoveDataIndex < 0 {
		return
	}
	if history.onPositionRequest != nil {
		if history.onPositionRequest(history.allMovesData[lastMoveDataIndex]) {
			history.currentHighlightedButtonIndex = lastButtonIndex
			history.updateButtonsStyles()
		}
	}
}

func (history *History) findLastMoveDataIndex() int {
	var lastMoveDataIndex int
	for index, _ := range history.allMovesData {
		lastMoveDataIndex = index
	}
	return lastMoveDataIndex
}

func (history *History) findLastButtonIndex() int {
	var lastButtonIndex int
	for index, currentObject := range history.container.Objects {
		_, ok := currentObject.(*widget.Button)
		if ok {
			lastButtonIndex = index
		}
	}

	return lastButtonIndex
}

func (history *History) updateButtonsStyles() {
	for index, currentObject := range history.container.Objects {
		currentButton, ok := currentObject.(*widget.Button)
		if ok {
			if index == history.currentHighlightedButtonIndex {
				currentButton.Style = widget.PrimaryButton
			} else {
				currentButton.Style = widget.DefaultButton
			}
		}
	}
	history.Refresh()
}
