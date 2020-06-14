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

type cell struct {
	file int
	rank int
}

type lastMove struct {
	originCell cell
	targetCell cell

	baseline       canvas.Line
	leftArrowLine  canvas.Line
	rightArrowLine canvas.Line
}

type movedPiece struct {
	location  fyne.Position
	piece     fyne.CanvasObject
	startCell cell
	endCell   cell
}

// ChessBoard is a chess board widget.
type ChessBoard struct {
	widget.BaseWidget

	game        chess.Game
	blackSide   BlackSide
	length      int
	cellsLength int
	lastMove    *lastMove

	movedPiece          movedPiece
	dragndropInProgress bool

	pieces [8][8]*canvas.Image
}

// CreateRenderer creates the board renderer.
func (board *ChessBoard) CreateRenderer() fyne.WidgetRenderer {
	board.ExtendBaseWidget(board)

	cells := [8][8]*canvas.Rectangle{}
	pieces := [8][8]*canvas.Image{}
	filesCoords := [2][8]*canvas.Text{}
	ranksCoords := [2][8]*canvas.Text{}

	board.buildCellsAndPieces(&cells, &pieces)
	board.buildFilesCoordinates(&filesCoords)
	board.buildRanksCoordinates(&ranksCoords)

	playerTurn := board.buildPlayerTurn()

	board.pieces = pieces

	return Renderer{
		boardWidget: board,
		cells:       cells,
		filesCoords: filesCoords,
		ranksCoords: ranksCoords,
		playerTurn:  playerTurn,
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
		game:      *chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{})),
		movedPiece: movedPiece{
			startCell: cell{file: -1, rank: -1},
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
	if board.dragndropInProgress == false {
		board.startDragAndDrop(event)
	} else {
		board.updateDragAndDrop(event)
	}
}

// DragEnd handles the drag end event for the chess board
func (board *ChessBoard) DragEnd() {
	file := board.movedPiece.endCell.file
	rank := board.movedPiece.endCell.rank

	inBounds := file >= 0 && file <= 7 && rank >= 0 && rank <= 7
	if !inBounds {
		board.resetDragAndDrop()
		board.Refresh()
		return
	}

	moveStr := board.getMoveString()

	err := board.game.MoveStr(moveStr)
	if err == nil {
		board.lastMove = &lastMove{
			originCell: board.movedPiece.startCell,
			targetCell: board.movedPiece.endCell,
		}
	}
	board.resetDragAndDrop()
	board.updatePieces()
	board.Refresh()
}

func (board *ChessBoard) startDragAndDrop(event *fyne.DragEvent) {
	cellsLength := int(float64(board.length) / 9)
	halfCellsLength := int(float64(cellsLength) / 2)

	position := event.Position
	// This is really needed to be coded as is !
	// First file and rank, then bounds test, then adjust values with the board orientation
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
	board.movedPiece.startCell = cell{file: file, rank: rank}
	board.movedPiece.endCell = cell{file: file, rank: rank}
	board.Refresh()
}

func (board *ChessBoard) updateDragAndDrop(event *fyne.DragEvent) {
	cellsLength := int(float64(board.length) / 9)
	halfCellsLength := int(float64(cellsLength) / 2)

	position := event.Position
	var file, rank int
	if board.blackSide == BlackAtTop {
		file = int(math.Floor(float64(position.X-halfCellsLength) / float64(cellsLength)))
		rank = 7 - int(math.Floor(float64(position.Y-halfCellsLength)/float64(cellsLength)))
	} else {
		file = 7 - int(math.Floor(float64(position.X-halfCellsLength)/float64(cellsLength)))
		rank = int(math.Floor(float64(position.Y-halfCellsLength) / float64(cellsLength)))
	}

	board.movedPiece.location = fyne.Position{X: position.X - halfCellsLength, Y: position.Y - halfCellsLength}
	board.movedPiece.endCell = cell{file: file, rank: rank}
	board.Refresh()
}

func (board *ChessBoard) resetDragAndDrop() {
	board.dragndropInProgress = false
	board.movedPiece.location = fyne.Position{X: -1000, Y: -1000}
	board.movedPiece.startCell = cell{file: -1, rank: -1}
}

func (board *ChessBoard) getMoveString() string {
	asciiLowerA := 97
	asciiOne := 49

	return fmt.Sprintf("%c%c%c%c",
		asciiLowerA+board.movedPiece.startCell.file,
		asciiOne+board.movedPiece.startCell.rank,
		asciiLowerA+board.movedPiece.endCell.file,
		asciiOne+board.movedPiece.endCell.rank,
	)
}

func (board *ChessBoard) buildCellsAndPieces(cells *[8][8]*canvas.Rectangle, pieces *[8][8]*canvas.Image) {
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

			square := chess.Square(col + 8*line)
			pieceValue := board.game.Position().Board().Piece(square)
			if pieceValue != chess.NoPiece {
				imageResource := imageResourceFromPiece(pieceValue)
				image := canvas.NewImageFromResource(&imageResource)
				image.FillMode = canvas.ImageFillContain

				pieces[line][col] = image
			}
		}
	}
}

func (board *ChessBoard) buildFilesCoordinates(filesCoords *[2][8]*canvas.Text) {
	coordsColor := color.RGBA{255, 199, 0, 0xff}
	asciiLowerA := 97

	for file := 0; file < 8; file++ {
		coord := fmt.Sprintf("%c", asciiLowerA+file)
		topCoord := canvas.NewText(coord, coordsColor)
		bottomCoord := canvas.NewText(coord, coordsColor)

		filesCoords[0][file] = topCoord
		filesCoords[1][file] = bottomCoord
	}
}

func (board *ChessBoard) buildRanksCoordinates(ranksCoords *[2][8]*canvas.Text) {
	coordsColor := color.RGBA{255, 199, 0, 0xff}
	asciiOne := 49

	for rank := 0; rank < 8; rank++ {
		coord := fmt.Sprintf("%c", asciiOne+(7-rank))
		leftCoord := canvas.NewText(coord, coordsColor)
		rightCoord := canvas.NewText(coord, coordsColor)

		ranksCoords[0][rank] = leftCoord
		ranksCoords[1][rank] = rightCoord
	}

}

func (board *ChessBoard) buildPlayerTurn() *canvas.Circle {
	var playerTurnColor color.Color
	gameTurn := board.game.Position().Turn()
	if gameTurn == chess.White {
		playerTurnColor = color.White
	} else {
		playerTurnColor = color.Black
	}

	return canvas.NewCircle(playerTurnColor)
}

func (board *ChessBoard) updatePieces() {
	for line := 0; line < 8; line++ {
		for col := 0; col < 8; col++ {
			board.pieces[line][col] = nil

			square := chess.Square(col + 8*line)
			pieceValue := board.game.Position().Board().Piece(square)
			if pieceValue != chess.NoPiece {
				imageResource := imageResourceFromPiece(pieceValue)
				image := canvas.NewImageFromResource(&imageResource)
				image.FillMode = canvas.ImageFillContain

				board.pieces[line][col] = image
			}
		}
	}
}
