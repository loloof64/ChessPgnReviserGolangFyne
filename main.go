package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/cloudfoundry-attic/jibber_jabber"
	"github.com/gookit/ini/v2"
	"github.com/loloof64/chess-pgn-reviser-fyne/chessboard"
	"github.com/loloof64/chess-pgn-reviser-fyne/history"

	"fmt"
)

func main() {
	lang, err := jibber_jabber.DetectLanguage()
	if err != nil {
		lang = "en"
	}

	langFiles := map[string]string{
		"en": "config/locales/en.ini",
		"fr": "config/locales/fr.ini",
		"es": "config/locales/es.ini",
	}

	langFileToLoad, found := langFiles[lang]
	if !found {
		langFileToLoad = langFiles["en"]
	}

	err = ini.LoadExists(langFileToLoad)
	if err != nil {
		panic(err)
	}

	app := app.New()
	app.SetIcon(resourceChessPng)

	title := ini.String("general.title")
	mainWindow := app.NewWindow(title)

	boardOrientation := chessboard.BlackAtBottom
	chessboardComponent := chessboard.NewChessBoard(400, &mainWindow)
	historyComponent := history.NewHistory(fyne.NewSize(400, 400))
	historyZone := widget.NewVScrollContainer(historyComponent)

	//
	for i := 0; i < 500; i++ {
		historyComponent.AddMove(history.GameMove{San: fmt.Sprintf("#%v", i+1)})
	}
	historyComponent.Clear()
	for i := 0; i < 300; i++ {
		historyComponent.AddMove(history.GameMove{San: fmt.Sprintf("#%v", i+1)})
	}
	//

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

	gameZone := fyne.NewContainerWithLayout(layout.NewHBoxLayout(), chessboardComponent, historyZone)

	mainLayout := layout.NewVBoxLayout()
	mainContent := fyne.NewContainerWithLayout(
		mainLayout,
		toolbar,
		gameZone,
	)

	mainWindow.SetContent(mainContent)

	mainWindow.ShowAndRun()
}
