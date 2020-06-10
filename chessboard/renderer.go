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

	cells [][]*canvas.Rectangle

	objects []fyne.CanvasObject
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
		}
	}
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
	return renderer.objects
}

// Destroy cleans up the renderer.
func (renderer Renderer) Destroy() {

}
