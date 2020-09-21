package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"

	"github.com/cloudfoundry-attic/jibber_jabber"
	"github.com/gookit/ini/v2"
	"github.com/loloof64/chess-pgn-reviser-fyne/chessboard"
	"github.com/loloof64/chess-pgn-reviser-fyne/commonTypes"
	"github.com/loloof64/chess-pgn-reviser-fyne/history"
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
	currentTheme := app.Settings().Theme()
	if currentTheme == theme.LightTheme() {
		app.Settings().SetTheme(&CustomLightTheme{})
	} else {
		app.Settings().SetTheme(&CustomDarkTheme{})
	}
	app.SetIcon(resourceChessPng)

	title := ini.String("general.title")
	mainWindow := app.NewWindow(title)

	boardOrientation := chessboard.BlackAtBottom
	chessboardComponent := chessboard.NewChessBoard(400, &mainWindow)
	historyComponent := history.NewHistory(fyne.NewSize(400, 400))
	historyZone := widget.NewVScrollContainer(historyComponent)

	startGameItem := widget.NewToolbarAction(resourceStartSvg, func() {
		historyComponent.Clear(1)
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

	stopGameItem := widget.NewToolbarAction(resourceStopSvg, func() {
		if !chessboardComponent.GameInProgress() {
			return
		}

		dialogTitle := ini.String("stopGameRequest.dialogTitle")
		dialogMessage := ini.String("stopGameRequest.dialogMessage")

		confirmButtonText := ini.String("general.okButton")
		cancelButtonText := ini.String("general.cancelButton")

		dialogComponent := widget.NewLabel(dialogMessage)

		confirmDialog := dialog.NewCustomConfirm(dialogTitle, confirmButtonText, cancelButtonText, dialogComponent, func(confirmed bool) {
			if confirmed {
				chessboardComponent.StopGame()
			}
		}, mainWindow)
		confirmDialog.Show()
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

	chessboardComponent.SetOnMoveDoneHandler(func(moveData commonTypes.GameMove) {
		historyComponent.AddMove(moveData)
	})

	historyComponent.SetOnPositionRequestHandler(func(moveData commonTypes.GameMove) {
		chessboardComponent.RequestHistoryPosition(moveData)
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

	toolbar := widget.NewToolbar(startGameItem, reverseBoardItem, claimDrawItem,
		stopGameItem)

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
