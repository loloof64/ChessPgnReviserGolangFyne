package chessboard

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"github.com/notnil/chess"
)

// Renderer renders a ChessBoard.
type Renderer struct {
	boardWidget *ChessBoard

	cells       [8][8]*canvas.Rectangle
	filesCoords [2][8]*canvas.Text
	ranksCoords [2][8]*canvas.Text
	playerTurn  *canvas.Circle
}

// Layout layouts the board elements.
func (renderer Renderer) Layout(size fyne.Size) {
	renderer.layoutCells(size)
	renderer.boardWidget.LayoutLastMoveArrowIfNeeded(size)
	renderer.layoutPieces(size)
	renderer.layoutMovedPieceIfAny(size)
	renderer.layoutFilesCoordinates(size)
	renderer.layoutRanksCoordinates(size)
	renderer.layoutPlayerTurn(size)

	renderer.updatePlayerTurn()
	renderer.updateCellsForDragAndDrop()
}

// MinSize computes the minimum size.
func (renderer Renderer) MinSize() fyne.Size {
	size := renderer.boardWidget.length
	return fyne.NewSize(size, size)
}

// Refresh refreshes the board.
func (renderer Renderer) Refresh() {
	renderer.Layout(renderer.boardWidget.Size())
	canvas.Refresh(renderer.boardWidget)
}

// BackgroundColor sets the board background color.
func (renderer Renderer) BackgroundColor() color.Color {
	return color.RGBA{20, 110, 200, 0xff}
}

// Objects returns the objects of the canvas of the renderer.
func (renderer Renderer) Objects() []fyne.CanvasObject {
	result := make([]fyne.CanvasObject, 0, 170)

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			cell := renderer.cells[rank][file]
			if cell != nil {
				result = append(result, cell)
			}
		}
	}

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			piece := renderer.boardWidget.pieces[rank][file]
			if piece == nil {
				continue
			}

			isDraggedPiece := renderer.boardWidget.dragndropInProgress &&
				renderer.boardWidget.movedPiece.startCell.File == int8(file) &&
				renderer.boardWidget.movedPiece.startCell.Rank == int8(rank)

			if isDraggedPiece {
				continue
			}
			result = append(result, piece)
		}
	}

	for col := 0; col < 8; col++ {
		result = append(result, renderer.filesCoords[0][col])
		result = append(result, renderer.filesCoords[1][col])
		result = append(result, renderer.ranksCoords[0][col])
		result = append(result, renderer.ranksCoords[1][col])
	}

	result = append(result, renderer.playerTurn)

	if renderer.boardWidget.dragndropInProgress {
		result = append(result, renderer.boardWidget.movedPiece.pieceImage)
	}

	if renderer.boardWidget.lastMove != nil {
		result = append(result, &renderer.boardWidget.lastMove.baseline)
		result = append(result, &renderer.boardWidget.lastMove.leftArrowLine)
		result = append(result, &renderer.boardWidget.lastMove.rightArrowLine)
	}

	return result
}

// Destroy cleans up the renderer.
func (renderer Renderer) Destroy() {

}

func (renderer Renderer) layoutCells(size fyne.Size) {
	minSize := math.Min(float64(size.Width), float64(size.Height))
	cellsLength := int(minSize / 9.0)
	halfCellsLength := cellsLength / 2
	cellsSize := fyne.Size{Width: int(cellsLength), Height: int(cellsLength)}

	for lineIndex, lineValues := range renderer.cells {
		for colIndex, cellValue := range lineValues {
			var x, y int
			if renderer.boardWidget.blackSide == BlackAtTop {
				x = halfCellsLength + colIndex*cellsLength
				y = halfCellsLength + (7-lineIndex)*cellsLength
			} else {
				x = halfCellsLength + (7-colIndex)*cellsLength
				y = halfCellsLength + lineIndex*cellsLength
			}
			cellPosition := fyne.Position{X: x, Y: y}

			cellValue.Resize(cellsSize)
			cellValue.Move(cellPosition)
		}
	}
}

func (renderer Renderer) layoutPieces(size fyne.Size) {
	minSize := math.Min(float64(size.Width), float64(size.Height))
	cellsLength := int(minSize / 9.0)
	halfCellsLength := cellsLength / 2
	cellsSize := fyne.Size{Width: int(cellsLength), Height: int(cellsLength)}

	for lineIndex, lineValues := range renderer.cells {
		for colIndex := range lineValues {
			var x, y int
			if renderer.boardWidget.blackSide == BlackAtTop {
				x = halfCellsLength + colIndex*cellsLength
				y = halfCellsLength + (7-lineIndex)*cellsLength
			} else {
				x = halfCellsLength + (7-colIndex)*cellsLength
				y = halfCellsLength + lineIndex*cellsLength
			}
			cellPosition := fyne.Position{X: x, Y: y}

			currentPiece := renderer.boardWidget.pieces[lineIndex][colIndex]

			if currentPiece != nil {
				currentPiece.Resize(cellsSize)
				currentPiece.Move(cellPosition)
			}
		}
	}
}

func (renderer Renderer) layoutFilesCoordinates(size fyne.Size) {
	minSize := math.Min(float64(size.Width), float64(size.Height))
	cellsLength := int(minSize / 9.0)

	coordsFontSize := int(float64(cellsLength) * 0.25)

	fileCoordsOffset := int(float64(cellsLength) * 0.95)

	for file := 0; file < 8; file++ {
		var x int
		if renderer.boardWidget.blackSide == BlackAtTop {
			x = fileCoordsOffset + cellsLength*file
		} else {
			x = fileCoordsOffset + cellsLength*(7-file)
		}
		yTop := int(float64(cellsLength) * 0.015)
		yBottom := int(float64(cellsLength) * 8.515)

		topCoord := renderer.filesCoords[0][file]
		topCoord.TextStyle = fyne.TextStyle{Bold: true}
		topCoord.TextSize = coordsFontSize
		topCoord.Move(fyne.Position{X: x, Y: yTop})

		bottomCoord := renderer.filesCoords[1][file]
		bottomCoord.TextStyle = fyne.TextStyle{Bold: true}
		bottomCoord.TextSize = coordsFontSize
		bottomCoord.Move(fyne.Position{X: x, Y: yBottom})
	}
}

func (renderer Renderer) layoutRanksCoordinates(size fyne.Size) {
	minSize := math.Min(float64(size.Width), float64(size.Height))
	cellsLength := int(minSize / 9.0)

	coordsFontSize := int(float64(cellsLength) * 0.25)

	rankCoordsOffset := int(float64(cellsLength) * 0.8)

	for rank := 0; rank < 8; rank++ {
		var y int
		if renderer.boardWidget.blackSide == BlackAtTop {
			y = rankCoordsOffset + cellsLength*rank
		} else {
			y = rankCoordsOffset + cellsLength*(7-rank)
		}
		xLeft := int(float64(cellsLength) * 0.2)
		xRight := int(float64(cellsLength) * 8.7)

		leftCoord := renderer.ranksCoords[0][rank]
		leftCoord.TextStyle = fyne.TextStyle{Bold: true}
		leftCoord.TextSize = coordsFontSize
		leftCoord.Move(fyne.Position{X: xLeft, Y: y})

		rightCoord := renderer.ranksCoords[1][rank]
		rightCoord.TextStyle = fyne.TextStyle{Bold: true}
		rightCoord.TextSize = coordsFontSize
		rightCoord.Move(fyne.Position{X: xRight, Y: y})
	}
}

func (renderer Renderer) layoutPlayerTurn(size fyne.Size) {
	minSize := math.Min(float64(size.Width), float64(size.Height))
	cellsLength := int(minSize / 9.0)
	halfCellsLength := cellsLength / 2

	turnCirclePlace := int(float64(cellsLength) * 8.5)
	turnCircle := renderer.playerTurn
	turnCircle.Resize(fyne.Size{Width: halfCellsLength, Height: halfCellsLength})
	turnCircle.Move(fyne.Position{X: turnCirclePlace, Y: turnCirclePlace})
}

func (renderer Renderer) updatePlayerTurn() {
	turnCircle := renderer.playerTurn
	if renderer.boardWidget.game.Position().Turn() == chess.White {
		turnCircle.FillColor = color.White
	} else {
		turnCircle.FillColor = color.Black
	}
}

func (renderer Renderer) layoutMovedPieceIfAny(size fyne.Size) {
	minSize := math.Min(float64(size.Width), float64(size.Height))
	cellsLength := int(minSize / 9.0)
	cellsSize := fyne.Size{Width: int(cellsLength), Height: int(cellsLength)}

	if renderer.boardWidget.dragndropInProgress {
		renderer.boardWidget.movedPiece.pieceImage.Resize(cellsSize)
		renderer.boardWidget.movedPiece.pieceImage.Move(renderer.boardWidget.movedPiece.location)
	}
}

func (renderer Renderer) updateCellsForDragAndDrop() {
	whiteCellColor := color.RGBA{255, 206, 158, 0xff}
	blackCellColor := color.RGBA{209, 139, 71, 0xff}

	dndCrossCellColor := color.RGBA{255, 20, 200, 0xff}
	dndOriginCellColor := color.RGBA{255, 20, 30, 0xff}
	dndTargetCellColor := color.RGBA{20, 255, 30, 0xff}

	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			isWhiteCell := (file+rank)%2 > 0
			if isWhiteCell {
				renderer.cells[rank][file].FillColor = whiteCellColor
			} else {
				renderer.cells[rank][file].FillColor = blackCellColor
			}

			if renderer.boardWidget.dragndropInProgress == false {
				continue
			}

			isADragAndDropCrossCell := int8(file) == renderer.boardWidget.movedPiece.endCell.File ||
				int8(rank) == renderer.boardWidget.movedPiece.endCell.Rank
			if isADragAndDropCrossCell {
				renderer.cells[rank][file].FillColor = dndCrossCellColor
			}

			isOriginCell := int8(file) == renderer.boardWidget.movedPiece.startCell.File &&
				int8(rank) == renderer.boardWidget.movedPiece.startCell.Rank
			if isOriginCell {
				renderer.cells[rank][file].FillColor = dndOriginCellColor
			}

			isTargetCell := int8(file) == renderer.boardWidget.movedPiece.endCell.File &&
				int8(rank) == renderer.boardWidget.movedPiece.endCell.Rank
			if isTargetCell {
				renderer.cells[rank][file].FillColor = dndTargetCellColor
			}

		}
	}
}
