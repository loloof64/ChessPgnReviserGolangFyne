package chessboard

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

// ChessBoard is a chess board widget.
type ChessBoard struct {
	widget.BaseWidget
}

// CreateRenderer creates the board renderer.
func (board *ChessBoard) CreateRenderer() fyne.WidgetRenderer {
	board.ExtendBaseWidget(board)
	return Renderer{}
}

// NewChessBoard creates a new chess board.
func NewChessBoard() *ChessBoard {
	chessBoard := &ChessBoard{}
	chessBoard.ExtendBaseWidget(chessBoard)

	return chessBoard
}
