package chessboard

import (
	"image/color"

	"fyne.io/fyne"
)

// Renderer renders a ChessBoard.
type Renderer struct {
}

// Layout layouts the board elements.
func (renderer Renderer) Layout(size fyne.Size) {

}

// MinSize computes the minimum size.
func (renderer Renderer) MinSize() fyne.Size {
	return fyne.NewSize(300, 300)
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
	return []fyne.CanvasObject{}
}

// Destroy cleans up the renderer.
func (renderer Renderer) Destroy() {

}
