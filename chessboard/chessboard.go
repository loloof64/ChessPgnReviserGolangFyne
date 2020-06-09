package chessboard

import (
	"fyne.io/fyne"
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
	return Renderer{
		boardWidget: board,
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
