package chessboard

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/widget"

	"github.com/notnil/chess"
)

// BlackSide defines the side of the black side on the board.
type BlackSide int

const (
	// BlackAtTop sets black side is at top of the board.
	BlackAtTop BlackSide = iota

	// BlackAtBottom sets black side is at bottom of the board.
	BlackAtBottom
)

// Cell handles a cell inside the board component.
type Cell struct {
	file int
	rank int
}

// MovedPiece handles the moved piece during a DragAndDrop, for the board component.
type MovedPiece struct {
	location  fyne.Position
	piece     fyne.CanvasObject
	startCell Cell
}

// ChessBoard is a chess board widget.
type ChessBoard struct {
	widget.BaseWidget

	game        chess.Game
	blackSide   BlackSide
	length      int
	cellsLength int

	movedPiece          MovedPiece
	dragndropInProgress bool
}

// CreateRenderer creates the board renderer.
func (board *ChessBoard) CreateRenderer() fyne.WidgetRenderer {
	board.ExtendBaseWidget(board)

	cells := [8][8]*canvas.Rectangle{}
	pieces := [8][8]*canvas.Image{}
	filesCoords := [2][8]*canvas.Text{}
	ranksCoords := [2][8]*canvas.Text{}

	cellsObjects := make([]fyne.CanvasObject, 0, 64)
	piecesObjects := make([]fyne.CanvasObject, 0, 64)
	coordsObjects := make([]fyne.CanvasObject, 0, 32)

	cellsObjects, piecesObjects = board.buildCellsAndPieces(&cells, &pieces, cellsObjects, piecesObjects)

	coordsColor := color.RGBA{255, 199, 0, 0xff}
	asciiLowerA := 97
	asciiOne := 49

	for file := 0; file < 8; file++ {
		coord := fmt.Sprintf("%c", asciiLowerA+file)
		topCoord := canvas.NewText(coord, coordsColor)
		bottomCoord := canvas.NewText(coord, coordsColor)

		filesCoords[0][file] = topCoord
		filesCoords[1][file] = bottomCoord
		coordsObjects = append(coordsObjects, topCoord)
		coordsObjects = append(coordsObjects, bottomCoord)
	}

	for rank := 0; rank < 8; rank++ {
		coord := fmt.Sprintf("%c", asciiOne+(7-rank))
		leftCoord := canvas.NewText(coord, coordsColor)
		rightCoord := canvas.NewText(coord, coordsColor)

		ranksCoords[0][rank] = leftCoord
		ranksCoords[1][rank] = rightCoord
		coordsObjects = append(coordsObjects, leftCoord)
		coordsObjects = append(coordsObjects, rightCoord)
	}

	var playerTurnColor color.Color
	gameTurn := board.game.Position().Turn()
	if gameTurn == chess.White {
		playerTurnColor = color.White
	} else {
		playerTurnColor = color.Black
	}
	playerTurn := canvas.NewCircle(playerTurnColor)

	return Renderer{
		boardWidget:   board,
		cells:         cells,
		pieces:        pieces,
		filesCoords:   filesCoords,
		ranksCoords:   ranksCoords,
		playerTurn:    playerTurn,
		cellsObjects:  cellsObjects,
		piecesObjects: piecesObjects,
		coordsObjects: coordsObjects,
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
func NewChessBoard(length int) *ChessBoard {
	chessBoard := &ChessBoard{
		length:    length,
		blackSide: BlackAtTop,
		game:      *chess.NewGame(),
		movedPiece: MovedPiece{
			startCell: Cell{file: -1, rank: -1},
			location:  fyne.Position{X: -1000, Y: -1000},
		},
	}
	chessBoard.ExtendBaseWidget(chessBoard)

	return chessBoard
}

// SetOrientation sets the orientation of the board, putting the black side at the requested side.
func (board *ChessBoard) SetOrientation(orientation BlackSide) {
	board.blackSide = orientation
	board.Refresh()
}

// Dragged handles the dragged event for the chess board.
func (board *ChessBoard) Dragged(event *fyne.DragEvent) {
	cellsLength := int(float64(board.length) / 9)
	halfCellsLength := int(float64(cellsLength) / 2)

	if board.dragndropInProgress == false {
		position := event.Position
		file := int(math.Floor(float64(position.X-halfCellsLength) / float64(cellsLength)))
		rank := int(math.Floor(float64(position.Y-halfCellsLength) / float64(cellsLength)))

		inBounds := file >= 0 && file <= 7 && rank >= 0 && rank <= 7
		if !inBounds {
			return
		}

		if board.blackSide == BlackAtTop {
			rank = 7 - rank
		} else {
			file = 7 - file
		}

		square := chess.Square(file + 8*rank)
		pieceValue := board.game.Position().Board().Piece(square)

		if pieceValue == chess.NoPiece {
			return
		}

		board.dragndropInProgress = true

		imageResource := imageResourceFromPiece(pieceValue)
		image := canvas.NewImageFromResource(&imageResource)
		image.FillMode = canvas.ImageFillContain
		board.movedPiece.piece = image
		board.movedPiece.location = fyne.Position{X: position.X - halfCellsLength, Y: position.Y - halfCellsLength}
		board.Refresh()
	} else {
		position := event.Position
		board.movedPiece.location = fyne.Position{X: position.X - halfCellsLength, Y: position.Y - halfCellsLength}
		board.Refresh()
	}
}

// DragEnd handles the drag end event for the chess board
func (board *ChessBoard) DragEnd() {
	board.dragndropInProgress = false
	board.movedPiece.location = fyne.Position{X: -1000, Y: -1000}
	board.movedPiece.startCell = Cell{file: -1, rank: -1}

	board.Refresh()
}

func (board *ChessBoard) buildCellsAndPieces(cells *[8][8]*canvas.Rectangle, pieces *[8][8]*canvas.Image,
	cellsObjects []fyne.CanvasObject, piecesObjects []fyne.CanvasObject) (
	[]fyne.CanvasObject, []fyne.CanvasObject) {
	whiteCellColor := color.RGBA{255, 206, 158, 0xff}
	blackCellColor := color.RGBA{209, 139, 71, 0xff}

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
			cellsObjects = append(cellsObjects, cellRef)

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

	return cellsObjects, piecesObjects
}
