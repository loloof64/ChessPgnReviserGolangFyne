package chessboard

import (
	"image/color"
	"math"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
)

// Renderer renders a ChessBoard.
type Renderer struct {
	boardWidget *ChessBoard

	cells       [8][8]*canvas.Rectangle
	pieces      [8][8]*canvas.Image
	filesCoords [2][8]*canvas.Text
	ranksCoords [2][8]*canvas.Text
	playerTurn  *canvas.Circle

	cellsObjects  []fyne.CanvasObject
	piecesObjects []fyne.CanvasObject
	coordsObjects []fyne.CanvasObject
}

// Layout layouts the board elements.
func (renderer Renderer) Layout(size fyne.Size) {
	minSize := math.Min(float64(size.Width), float64(size.Height))
	cellsLength := int(minSize / 9.0)
	halfCellsLength := cellsLength / 2
	cellsSize := fyne.Size{Width: int(cellsLength), Height: int(cellsLength)}

	for lineIndex, lineValues := range renderer.cells {
		for colIndex, cellValue := range lineValues {
			x := halfCellsLength + colIndex*cellsLength
			y := halfCellsLength + (7-lineIndex)*cellsLength
			cellPosition := fyne.Position{X: x, Y: y}

			cellValue.Resize(cellsSize)
			cellValue.Move(cellPosition)

			currentPiece := renderer.pieces[lineIndex][colIndex]

			if currentPiece != nil {
				currentPiece.Resize(cellsSize)
				currentPiece.Move(cellPosition)
			}
		}
	}

	coordsFontSize := int(float64(cellsLength) * 0.25)

	fileCoordsOffset := int(float64(cellsLength) * 0.95)
	rankCoordsOffset := int(float64(cellsLength) * 0.8)

	for file := 0; file < 8; file++ {
		x := fileCoordsOffset + cellsLength*file
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

	for rank := 0; rank < 8; rank++ {
		y := rankCoordsOffset + cellsLength*rank
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

	turnCirclePlace := int(float64(cellsLength) * 8.5)
	turnCircle := renderer.playerTurn
	turnCircle.Resize(fyne.Size{Width: halfCellsLength, Height: halfCellsLength})
	turnCircle.Move(fyne.Position{X: turnCirclePlace, Y: turnCirclePlace})
}

// MinSize computes the minimum size.
func (renderer Renderer) MinSize() fyne.Size {
	size := renderer.boardWidget.size
	return fyne.NewSize(size, size)
}

// Refresh refreshes the board.
func (renderer Renderer) Refresh() {

}

// BackgroundColor sets the board background color.
func (renderer Renderer) BackgroundColor() color.Color {
	return color.RGBA{20, 110, 200, 0xff}
}

// Objects returns the objects of the canvas of the renderer.
func (renderer Renderer) Objects() []fyne.CanvasObject {
	result := make([]fyne.CanvasObject, 0, 170)

	for _, object := range renderer.cellsObjects {
		result = append(result, object)
	}

	for _, object := range renderer.piecesObjects {
		result = append(result, object)
	}

	for _, object := range renderer.coordsObjects {
		result = append(result, object)
	}

	result = append(result, renderer.playerTurn)

	return result
}

// Destroy cleans up the renderer.
func (renderer Renderer) Destroy() {

}
