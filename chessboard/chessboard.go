package chessboard

import (
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/notnil/chess"
)

// ChessBoard is a chess board widget.
type ChessBoard struct {
	widget.BaseWidget

	game chess.Game
	size int
}

// CreateRenderer creates the board renderer.
func (board *ChessBoard) CreateRenderer() fyne.WidgetRenderer {
	board.ExtendBaseWidget(board)

	whiteCellColor := color.RGBA{255, 206, 158, 0xff}
	blackCellColor := color.RGBA{209, 139, 71, 0xff}

	cells := [8][8]*canvas.Rectangle{}
	pieces := [8][8]*canvas.Image{}

	generalObjects := make([]fyne.CanvasObject, 0, 100)
	piecesObjects := make([]fyne.CanvasObject, 0, 64)

	for line := 0; line < 8; line++ {
		for col := 0; col < 8; col++ {
			isWhiteCell := (line+col)%2 != 0
			var cellColor color.Color
			if isWhiteCell {
				cellColor = whiteCellColor
			} else {
				cellColor = blackCellColor
			}
			cellRef := canvas.NewRectangle(cellColor)
			cells[line][col] = cellRef
			generalObjects = append(generalObjects, cellRef)

			square := chess.Square(col + 8*line)
			pieceValue := board.game.Position().Board().Piece(square)
			if pieceValue != chess.NoPiece {
				imageResource := imageResourceFromPiece(pieceValue)
				image := canvas.NewImageFromResource(&imageResource)
				image.FillMode = canvas.ImageFillContain

				pieces[line][col] = image
				piecesObjects = append(piecesObjects, image)
			}
		}
	}

	return Renderer{
		boardWidget:    board,
		cells:          cells,
		pieces:         pieces,
		generalObjects: generalObjects,
		piecesObjects:  piecesObjects,
	}
}

func imageResourceFromPiece(piece chess.Piece) fyne.StaticResource {
	var result fyne.StaticResource

	switch piece {
	case chess.WhitePawn:
		result = *resourceChessplt45Svg
	case chess.WhiteKnight:
		result = *resourceChessnlt45Svg
	case chess.WhiteBishop:
		result = *resourceChessblt45Svg
	case chess.WhiteRook:
		result = *resourceChessrlt45Svg
	case chess.WhiteQueen:
		result = *resourceChessqlt45Svg
	case chess.WhiteKing:
		result = *resourceChessklt45Svg
	case chess.BlackPawn:
		result = *resourceChesspdt45Svg
	case chess.BlackKnight:
		result = *resourceChessndt45Svg
	case chess.BlackBishop:
		result = *resourceChessbdt45Svg
	case chess.BlackRook:
		result = *resourceChessrdt45Svg
	case chess.BlackQueen:
		result = *resourceChessqdt45Svg
	case chess.BlackKing:
		result = *resourceChesskdt45Svg
	}

	return result
}

// NewChessBoard creates a new chess board.
func NewChessBoard(size int) *ChessBoard {
	chessBoard := &ChessBoard{
		size: size,
		game: *chess.NewGame(),
	}
	chessBoard.ExtendBaseWidget(chessBoard)

	return chessBoard
}
