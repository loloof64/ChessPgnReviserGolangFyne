package history

import (
	"fmt"
	"image/color"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/loloof64/chess-pgn-reviser-fyne/commonTypes"
)

// HistoryLayout defines the layout of the History widget.
type HistoryLayout struct {
	width float32
	gap   fyne.Size
}

func (l HistoryLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	currMaxH, w, h := float32(0), float32(0), float32(0)

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
	w, h, currMaxH := float32(0), float32(0), float32(0)

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

func newHistoryLayout(width float32) HistoryLayout {
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

	allMovesData  []commonTypes.GameMove
	startPosition string
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
	var moveComponent *fyne.Container
	history.currentHighlightedButtonIndex += 1
	thisButtonIndex := history.currentHighlightedButtonIndex
	// Indeed, for now, we haven't added the current move yet.
	thisMoveDataIndex := len(history.allMovesData)

	moveButton := widget.NewButton(moveData.Fan, func() {
		if history.onPositionRequest != nil {
			if history.onPositionRequest(moveData) {
				history.currentHighlightedButtonIndex = thisButtonIndex
				history.currentMoveDataIndex = thisMoveDataIndex
				history.updateButtonsStyles()
			}
		}
	})
	moveComponent = container.New(
		layout.NewMaxLayout(),
		canvas.NewRectangle(color.Transparent),
		moveButton,
	)

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
		history.Refresh()
	}
}

// Clear clears all moves from the History widget.
func (history *History) Clear(startPositionFen string) {
	history.container.Objects = nil
	history.allMovesData = nil
	history.startPosition = startPositionFen
	positionParts := strings.Split(startPositionFen, " ")
	startMoveNumber, err := strconv.Atoi(positionParts[len(positionParts)-1])
	if err != nil {
		startMoveNumber = 1
	}
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

// Tries to select the start position.
func (history *History) RequestStartPositionSelection() {
	history.currentHighlightedButtonIndex = -1
	history.currentMoveDataIndex = -1
	history.updateButtonsStyles()

	positionToRequest := commonTypes.GameMove{}
	positionToRequest.Fen = history.startPosition
	if history.onPositionRequest != nil {
		_ = history.onPositionRequest(positionToRequest)
	}
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
			history.currentMoveDataIndex = lastMoveDataIndex
			history.updateButtonsStyles()
		}
	}
}

// Tries to select the previous element, or to load start position
func (history *History) RequestPreviousItemSelection() {
	previousButtonIndex := history.findPreviousButtonIndex()
	previousMoveDataIndex := history.findPreviousLastMoveDataIndex()

	if history.onPositionRequest != nil {
		weAreOutsideMoves := previousButtonIndex < 0 || previousMoveDataIndex < 0
		if weAreOutsideMoves {
			positionToRequest := commonTypes.GameMove{}
			positionToRequest.Fen = history.startPosition
			if history.onPositionRequest(positionToRequest) {
				history.currentHighlightedButtonIndex = -1
				history.currentMoveDataIndex = -1
				history.updateButtonsStyles()
			}
		} else {
			if history.onPositionRequest(history.allMovesData[previousMoveDataIndex]) {
				history.currentHighlightedButtonIndex = previousButtonIndex
				history.currentMoveDataIndex = previousMoveDataIndex
				history.updateButtonsStyles()
			}
		}
	}
}

// Tries to select the next element
func (history *History) RequestNextItemSelection() {
	nextButtonIndex := history.findNextButtonIndex()
	nextMoveDataIndex := history.findNextLastMoveDataIndex()

	if history.onPositionRequest != nil {
		weAreOutsideMoves := nextButtonIndex < 0 || nextMoveDataIndex < 0
		if weAreOutsideMoves {
			return
		} else {
			if history.onPositionRequest(history.allMovesData[nextMoveDataIndex]) {
				history.currentHighlightedButtonIndex = nextButtonIndex
				history.currentMoveDataIndex = nextMoveDataIndex
				history.updateButtonsStyles()
			}
		}
	}
}

func (history *History) findLastMoveDataIndex() int {
	var lastMoveDataIndex int
	for index := range history.allMovesData {
		lastMoveDataIndex = index
	}
	return lastMoveDataIndex
}

func (history *History) findLastButtonIndex() int {
	var lastButtonIndex int
	for index, currentObject := range history.container.Objects {
		_, ok := currentObject.(*fyne.Container)
		if ok {
			lastButtonIndex = index
		}
	}

	return lastButtonIndex
}

func (history *History) findPreviousLastMoveDataIndex() int {
	if history.currentMoveDataIndex < 0 {
		return -1
	} else {
		return history.currentMoveDataIndex - 1
	}
}

func (history *History) findPreviousButtonIndex() int {
	if history.currentHighlightedButtonIndex < 0 {
		return -1
	}
	for index := history.currentHighlightedButtonIndex - 1; index >= 0; index-- {
		_, ok := history.container.Objects[index].(*fyne.Container)
		if ok {
			return index
		}
	}
	return -1
}

func (history *History) findNextLastMoveDataIndex() int {
	if history.currentMoveDataIndex >= len(history.allMovesData)-1 {
		return len(history.allMovesData) - 1
	} else if len(history.allMovesData) == 0 {
		return -1
	} else {
		return history.currentMoveDataIndex + 1
	}
}

func (history *History) findNextButtonIndex() int {
	for index := history.currentHighlightedButtonIndex + 1; index < len(history.container.Objects); index++ {
		_, ok := history.container.Objects[index].(*fyne.Container)
		if ok {
			return index
		}
	}

	return -1
}

func (history *History) updateButtonsStyles() {
	for index, currentObject := range history.container.Objects {
		currentButtonZone, buttonZoneOk := currentObject.(*fyne.Container)
		if buttonZoneOk {
			_, buttonOk := currentButtonZone.Objects[1].(*widget.Button)
			if buttonOk {
				var fillColor color.Color
				if index == history.currentHighlightedButtonIndex {
					fillColor = color.NRGBA{R: 100, G: 30, B: 255, A: 255}
				} else {
					fillColor = color.Transparent
				}
				var currentButtonBackground = currentButtonZone.Objects[0].(*canvas.Rectangle)
				currentButtonBackground.FillColor = fillColor
			}
		}
	}
	history.Refresh()
}
