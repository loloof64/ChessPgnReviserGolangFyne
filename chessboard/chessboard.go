package chessboard

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"
)

// ChessBoard is a chess board widget.
type ChessBoard struct {
	widget.BaseWidget

	size int
}

// CreateRenderer creates the board renderer.
func (board *ChessBoard) CreateRenderer() fyne.WidgetRenderer {
	board.ExtendBaseWidget(board)

	cells := make([][]*canvas.Rectangle, 0, 8)
	objects := make([]fyne.CanvasObject, 0, 161)

	for line := 0; line < 8; line++ {
		cellsLine := make([]*canvas.Rectangle, 0, 8)
		for col := 0; col < 8; col++ {
			isWhiteCell := (line+col)%2 != 0
			var cellColor color.Color
			if isWhiteCell {
				cellColor = color.RGBA{255, 206, 158, 0xff}
			} else {
				cellColor = color.RGBA{209, 139, 71, 0xff}
			}
			cellRef := canvas.NewRectangle(cellColor)
			cellsLine = append(cellsLine, cellRef)
			objects = append(objects, cellRef)
		}
		cells = append(cells, cellsLine)
	}

	return Renderer{
		boardWidget: board,
		cells:       cells,
		objects:     objects,
	}
}

// NewChessBoard creates a new chess board.
func NewChessBoard(size int) *ChessBoard {
	chessBoard := &ChessBoard{
		size: size,
	}
	chessBoard.ExtendBaseWidget(chessBoard)

	return chessBoard
}
