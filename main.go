package main

import (
	"fyne.io/fyne/app"
	"github.com/loloof64/chess-pgn-reviser-fyne/chessboard"
)

func main() {
	app := app.New()

	mainWindow := app.NewWindow("Chess Pgn Reviser")
	mainWindow.SetContent(
		chessboard.NewChessBoard(400),
	)

	mainWindow.ShowAndRun()
}
