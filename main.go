package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/BurntSushi/toml"
	"github.com/cloudfoundry-attic/jibber_jabber"
	"github.com/loloof64/chess-pgn-reviser-fyne/chessboard"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.MustLoadMessageFile("active.en.toml")
	bundle.MustLoadMessageFile("active.fr.toml")
	bundle.MustLoadMessageFile("active.es.toml")

	lang, err := jibber_jabber.DetectLanguage()
	if err != nil {
		lang = "en"
	}

	localizer := i18n.NewLocalizer(
		bundle, lang,
	)

	app := app.New()

	title := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "AppTitle",
	})
	mainWindow := app.NewWindow(title)

	boardOrientation := chessboard.BlackAtBottom
	chessboardComponent := chessboard.NewChessBoard(400, &mainWindow, localizer)

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

	gameFinished := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "GameFinished",
	})

	whiteWon := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "WhiteWon",
	})

	blackWon := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "BlackWon",
	})

	draw := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "Draw",
	})

	chessboardComponent.SetOnWhiteWinHandler(func() {
		dialog.ShowInformation(gameFinished, whiteWon, mainWindow)
	})

	chessboardComponent.SetOnBlackWinHandler(func() {
		dialog.ShowInformation(gameFinished, blackWon, mainWindow)
	})

	chessboardComponent.SetOnDrawHandler(func() {
		dialog.ShowInformation(gameFinished, draw, mainWindow)
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
