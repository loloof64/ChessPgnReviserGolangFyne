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

	mainWindow := app.NewWindow("Chess Pgn Reviser")

	boardOrientation := chessboard.BlackAtBottom
	chessboardComponent := chessboard.NewChessBoard(400, &mainWindow)

	startGameItem := widget.NewToolbarAction(resourceStartSvg, func() {
		chessboardComponent.NewGame()
	})

	reverseBoardItem := widget.NewToolbarAction(resourceReverseSvg, func() {
		if boardOrientation == chessboard.BlackAtBottom {
			boardOrientation = chessboard.BlackAtTop
		} else {
			boardOrientation = chessboard.BlackAtBottom
		}
		chessboardComponent.SetOrientation(boardOrientation)
	})

	toolbar := widget.NewToolbar(startGameItem, reverseBoardItem)

	mainLayout := layout.NewVBoxLayout()
	mainContent := fyne.NewContainerWithLayout(
		mainLayout,
		toolbar,
		chessboardComponent,
	)

	mainWindow.SetContent(mainContent)

	mainWindow.ShowAndRun()
}
