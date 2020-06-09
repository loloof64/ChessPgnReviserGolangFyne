package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

func main() {
	app := app.New()

	w := app.NewWindow("Chess Pgn Reviser")
	w.SetContent(widget.NewLabel("Hello World !"))

	w.ShowAndRun()
}
