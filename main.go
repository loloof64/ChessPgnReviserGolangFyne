package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/gookit/ini/v2"
	"github.com/loloof64/chess-pgn-reviser-fyne/chessboard"
)

func main() {
	err := ini.LoadExists("config/locales/en.ini", "config/locales/es.ini", "config/locales/fr.ini")
	if err != nil {
		panic(err)
	}

	app := app.New()

	title := ini.String("general.title")
	mainWindow := app.NewWindow(title)

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

	gameFinished := ini.String("general.gameFinished")

	whiteWon := ini.String("gameResult.whiteWon")

	blackWon := ini.String("gameResult.blackWon")

	draw := ini.String("gameResult.draw")

	chessboardComponent.SetOnWhiteWinHandler(func() {
		dialog.ShowInformation(gameFinished, whiteWon, mainWindow)
	})

	chessboardComponent.SetOnBlackWinHandler(func() {
		dialog.ShowInformation(gameFinished, blackWon, mainWindow)
	})

	chessboardComponent.SetOnDrawHandler(func() {
		dialog.ShowInformation(gameFinished, draw, mainWindow)
	})

	claimDrawItem := widget.NewToolbarAction(resourceAgreementSvg, func() {
		drawAcceptedMessage := ini.String("drawClaim.accepted")

		drawRefusedMessage := ini.String("drawClaim.rejected")

		accepted := chessboardComponent.ClaimDraw()
		if accepted {
			dialog.ShowInformation(gameFinished, drawAcceptedMessage, mainWindow)
		} else {
			dialog.ShowInformation(gameFinished, drawRefusedMessage, mainWindow)
		}
	})

	toolbar := widget.NewToolbar(startGameItem, reverseBoardItem, claimDrawItem)

	mainLayout := layout.NewVBoxLayout()
	mainContent := fyne.NewContainerWithLayout(
		mainLayout,
		toolbar,
		chessboardComponent,
	)

	mainWindow.SetContent(mainContent)

	mainWindow.ShowAndRun()
}
