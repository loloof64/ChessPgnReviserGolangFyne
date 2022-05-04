package chessboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

// IconButton defines a clickable image
type IconButton struct {
	widget.Icon
	action func()
	size   fyne.Size
	icon   fyne.Resource
}

// NewIconButton created a new IconButton
func NewIconButton(icon fyne.Resource, size fyne.Size, action func()) *IconButton {
	button := &IconButton{
		action: action,
		size:   size,
		icon:   icon,
	}

	button.ExtendBaseWidget(button)
	button.SetResource(icon)

	return button
}

// Tapped defines the primary action
func (button *IconButton) Tapped(_ *fyne.PointEvent) {
	button.action()
}

// TappedSecondary defines the secondary action
func (button *IconButton) TappedSecondary(_ *fyne.PointEvent) {

}

// MinSize computes the minimum size of the button
func (button *IconButton) MinSize() (size fyne.Size) {
	size = button.size
	return
}

// Layout layouts the content of the button
func (button *IconButton) Layout(size fyne.Size) {
	button.Resize(button.size)
}
