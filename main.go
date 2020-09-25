package main

import (
	"fmt"
	"strings"

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
	"github.com/loloof64/chess-pgn-reviser-fyne/pgnLoader"
	"github.com/notnil/chess"
)

func loadLocales() {
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
}

func buildAppInstance() fyne.App {
	app := app.New()
	currentTheme := app.Settings().Theme()
	if currentTheme == theme.LightTheme() {
		app.Settings().SetTheme(&CustomLightTheme{})
	} else {
		app.Settings().SetTheme(&CustomDarkTheme{})
	}
	app.SetIcon(resourceChessPng)

	return app
}

func buildMainWindow(app fyne.App) fyne.Window {
	title := ini.String("general.title")
	mainWindow := app.NewWindow(title)
	return mainWindow
}

func buildMainContent(mainWindow fyne.Window) fyne.CanvasObject {

	boardOrientation := chessboard.BlackAtBottom
	chessboardComponent := chessboard.NewChessBoard(400, &mainWindow)
	historyComponent := history.NewHistory(fyne.NewSize(400, 400))

	gotoPreviousHistoryButton := widget.NewButtonWithIcon("", resourcePreviousSvg, func() {
		historyComponent.RequestPreviousItemSelection()
	})

	gotoNextHistoryButton := widget.NewButtonWithIcon("", resourceNextSvg, func() {
		historyComponent.RequestNextItemSelection()
	})

	gotoStartPositionButton := widget.NewButtonWithIcon("", resourceFirstSvg, func() {
		historyComponent.RequestStartPositionSelection()
	})

	gotoLastHistoryButton := widget.NewButtonWithIcon("", resourceLastSvg, func() {
		historyComponent.RequestLastItemSelection()
	})

	historyButtonsZone := fyne.NewContainerWithLayout(
		layout.NewCenterLayout(),
		fyne.NewContainerWithLayout(
			layout.NewHBoxLayout(),
			gotoStartPositionButton,
			gotoPreviousHistoryButton,
			gotoNextHistoryButton,
			gotoLastHistoryButton,
		),
	)

	hideHistoryNavigationToolbar := func() {
		historyButtonsZone.Hide()
	}

	showHistoryNavigationToolbar := func() {
		historyButtonsZone.Show()
	}

	historyMainContent := widget.NewVScrollContainer(historyComponent)
	historyMainContent.Resize(fyne.NewSize(400, 400))
	historyZone := fyne.NewContainerWithLayout(
		layout.NewBorderLayout(historyButtonsZone, nil, nil, nil),
		historyButtonsZone,
		historyMainContent,
	)
	hideHistoryNavigationToolbar()

	errorOpeningFileTitle := ini.String("serialization.errorOpeningFileTitle")
	errorOpeningFileMessage := ini.String("serialization.errorOpeningFileMessage")

	startGameItem := widget.NewToolbarAction(resourceStartSvg, func() {
		openFileDialog := dialog.NewFileOpen(func(fileData fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println(err)
				dialog.ShowInformation(errorOpeningFileTitle, errorOpeningFileMessage, mainWindow)
				return
			}

			if fileData == nil {
				return
			}

			fileNameRune := []rune(fmt.Sprintf("%v", fileData.URI()))
			// Stripping "file://" prefix
			filePath := string(fileNameRune[7:])

			pgnLoader, err := pgnLoader.LoadPgnFile(filePath)

			if err != nil {
				fmt.Println(err)
				dialog.ShowInformation(errorOpeningFileTitle, errorOpeningFileMessage, mainWindow)
				return
			}

			selectedGamePgn := pgnLoader.Games[0]
			reader := strings.NewReader(selectedGamePgn)
			/*selectedGameParsed, err := chess.PGN(reader)

			if err != nil {
				fmt.Println(err)
				dialog.ShowInformation(errorOpeningFileTitle, errorOpeningFileMessage, mainWindow)
				return
			}*/

			hideHistoryNavigationToolbar()
			historyComponent.Clear(1)
			chessboardComponent.NewGame()
		}, mainWindow)
		openFileDialog.Show()
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

		confirmDialog := dialog.NewCustomConfirm(dialogTitle, confirmButtonText,
			cancelButtonText, dialogComponent, func(confirmed bool) {
				if confirmed {
					showHistoryNavigationToolbar()
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
		showHistoryNavigationToolbar()
		dialog.ShowInformation(gameFinished, whiteWon, mainWindow)
	})

	chessboardComponent.SetOnBlackWinHandler(func() {
		showHistoryNavigationToolbar()
		dialog.ShowInformation(gameFinished, blackWon, mainWindow)
	})

	chessboardComponent.SetOnDrawHandler(func() {
		showHistoryNavigationToolbar()
		dialog.ShowInformation(gameFinished, draw, mainWindow)
	})

	chessboardComponent.SetOnMoveDoneHandler(func(moveData commonTypes.GameMove) {
		historyComponent.AddMove(moveData)
	})

	chessboardComponent.SetOnRequestLastHistoryPositionHandler(func() {
		historyComponent.RequestLastItemSelection()
	})

	historyComponent.SetOnPositionRequestHandler(
		func(moveData commonTypes.GameMove) bool {
			return chessboardComponent.RequestHistoryPosition(moveData)
		})

	claimDrawItem := widget.NewToolbarAction(resourceAgreementSvg, func() {
		drawAcceptedMessage := ini.String("drawClaim.accepted")

		drawRefusedMessage := ini.String("drawClaim.rejected")

		accepted := chessboardComponent.ClaimDraw()
		if accepted {
			showHistoryNavigationToolbar()
			dialog.ShowInformation(gameFinished, drawAcceptedMessage, mainWindow)
		} else {
			dialog.ShowInformation(gameFinished, drawRefusedMessage, mainWindow)
		}
	})

	toolbar := widget.NewToolbar(startGameItem, reverseBoardItem, claimDrawItem,
		stopGameItem)

	gameZone := fyne.NewContainerWithLayout(layout.NewHBoxLayout(),
		chessboardComponent, historyZone)

	mainLayout := layout.NewVBoxLayout()
	mainContent := fyne.NewContainerWithLayout(
		mainLayout,
		toolbar,
		gameZone,
	)

	return mainContent
}

func main() {
	loadLocales()
	app := buildAppInstance()
	mainWindow := buildMainWindow(app)
	mainContent := buildMainContent(mainWindow)
	mainWindow.SetContent(mainContent)
	mainWindow.ShowAndRun()
}
