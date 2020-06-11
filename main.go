package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/loloof64/chess-pgn-reviser-fyne/chessboard"
)

func main() {
	app := app.New()

	chessboardComponent := chessboard.NewChessBoard(400)
	reverseBoardButton := widget.NewButtonWithIcon("", resourceReverseSvg, func() {
		//chessboardComponent.Reverse()
		println("Clik")
	})

	mainLayout := layout.NewVBoxLayout()
	mainContent := fyne.NewContainerWithLayout(
		mainLayout,
		chessboardComponent,
		reverseBoardButton,
	)

	mainWindow := app.NewWindow("Chess Pgn Reviser")
	mainWindow.SetContent(mainContent)

	mainWindow.ShowAndRun()
}
