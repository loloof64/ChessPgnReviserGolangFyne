package chessboard

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/gookit/ini/v2"
	"github.com/notnil/chess"
)

// BlackSide defines the side of the black side on the board.
type BlackSide int

const (
	// BlackAtTop sets black side is at top of the board.
	BlackAtTop BlackSide = iota

	// BlackAtBottom sets black side is at bottom of the board.
	BlackAtBottom
)

// GameEndStatus defines the finished status of the game.
type GameEndStatus int

const (
	// NotFinished says that the game is neither won, neither draw.
	NotFinished GameEndStatus = iota

	// WhiteWon says that white side won.
	WhiteWon

	// BlackWon says that black side won.
	BlackWon

	// Draw says that game ended in a draw.
	Draw
)

type cell struct {
	file int8
	rank int8
}

type lastMove struct {
	originCell cell
	targetCell cell

	baseline       canvas.Line
	leftArrowLine  canvas.Line
	rightArrowLine canvas.Line
}

type movedPiece struct {
	location   fyne.Position
	pieceValue chess.Piece
	pieceImage fyne.CanvasObject
	startCell  cell
	endCell    cell
}

// ChessBoard is a chess board widget.
type ChessBoard struct {
	widget.BaseWidget

	parent    *fyne.Window
	game      chess.Game
	blackSide BlackSide
	length    int
	lastMove  *lastMove

	movedPiece          *movedPiece
	gameInProgress      bool
	dragndropInProgress bool
	pendingPromotion    bool
	promotionDialog     dialog.Dialog

	onWhiteWin func()
	onBlackWin func()
	onDraw     func()
	onMoveDone func(fan string)

	pieces [8][8]*canvas.Image
}

// CreateRenderer creates the board renderer.
func (board *ChessBoard) CreateRenderer() fyne.WidgetRenderer {
	board.ExtendBaseWidget(board)

	cells := [8][8]*canvas.Rectangle{}
	pieces := [8][8]*canvas.Image{}
	filesCoords := [2][8]*canvas.Text{}
	ranksCoords := [2][8]*canvas.Text{}

	board.buildCellsAndPieces(&cells, &pieces)
	board.buildFilesCoordinates(&filesCoords)
	board.buildRanksCoordinates(&ranksCoords)

	playerTurn := board.buildPlayerTurn()

	board.pieces = pieces

	return Renderer{
		boardWidget: board,
		cells:       cells,
		filesCoords: filesCoords,
		ranksCoords: ranksCoords,
		playerTurn:  playerTurn,
	}
}

// SetOnWhiteWinHandler sets the handler for white side win.
func (board *ChessBoard) SetOnWhiteWinHandler(handler func()) {
	board.onWhiteWin = handler
}

// SetOnBlackWinHandler sets the handler for black side win.
func (board *ChessBoard) SetOnBlackWinHandler(handler func()) {
	board.onBlackWin = handler
}

// SetOnDrawHandler sets the handler for draw.
func (board *ChessBoard) SetOnDrawHandler(handler func()) {
	board.onDraw = handler
}

// SetOnMoveDoneHandler sets the handler for moves done on the chess board widget.
func (board *ChessBoard) SetOnMoveDoneHandler(handler func(fan string)) {
	board.onMoveDone = handler
}

func imageResourceFromPiece(piece chess.Piece) fyne.StaticResource {
	var result fyne.StaticResource

	switch piece {
	case chess.WhitePawn:
		result = *resourceChessplt45Svg
	case chess.WhiteKnight:
		result = *resourceChessnlt45Svg
	case chess.WhiteBishop:
		result = *resourceChessblt45Svg
	case chess.WhiteRook:
		result = *resourceChessrlt45Svg
	case chess.WhiteQueen:
		result = *resourceChessqlt45Svg
	case chess.WhiteKing:
		result = *resourceChessklt45Svg
	case chess.BlackPawn:
		result = *resourceChesspdt45Svg
	case chess.BlackKnight:
		result = *resourceChessndt45Svg
	case chess.BlackBishop:
		result = *resourceChessbdt45Svg
	case chess.BlackRook:
		result = *resourceChessrdt45Svg
	case chess.BlackQueen:
		result = *resourceChessqdt45Svg
	case chess.BlackKing:
		result = *resourceChesskdt45Svg
	}

	return result
}

// NewChessBoard creates a new chess board.
func NewChessBoard(length int, parent *fyne.Window) *ChessBoard {
	customFen, _ := chess.FEN("8/8/8/8/8/8/8/8 w - - 0 1")

	chessBoard := &ChessBoard{
		length:    length,
		blackSide: BlackAtTop,
		game:      *chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}), customFen),
		parent:    parent,
	}
	chessBoard.ExtendBaseWidget(chessBoard)

	return chessBoard
}

// NewGame starts a new game
func (board *ChessBoard) NewGame() {
	standardFen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	startFen, _ := chess.FEN(standardFen)

	board.game = *chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}), startFen)
	board.gameInProgress = true
	board.lastMove = nil
	board.pendingPromotion = false

	board.updatePieces()
	board.Refresh()
}

// SetOrientation sets the orientation of the board, putting the black side at the requested side.
func (board *ChessBoard) SetOrientation(orientation BlackSide) {
	board.blackSide = orientation
	board.Refresh()
}

// Dragged handles the dragged event for the chess board.
func (board *ChessBoard) Dragged(event *fyne.DragEvent) {
	if board.pendingPromotion {
		return
	}

	if board.dragndropInProgress == false {
		board.startDragAndDrop(event)
	} else {
		board.updateDragAndDrop(event)
	}
}

// DragEnd handles the drag end event for the chess board
func (board *ChessBoard) DragEnd() {
	if !board.gameInProgress {
		return
	}

	if board.movedPiece == nil {
		return
	}

	file := board.movedPiece.endCell.file
	rank := board.movedPiece.endCell.rank

	inBounds := file >= 0 && file <= 7 && rank >= 0 && rank <= 7
	if !inBounds {
		board.resetDragAndDrop()
		board.Refresh()
		return
	}

	rank1 := int8(0)
	rank8 := int8(7)

	isPromotionMove :=
		board.movedPiece != nil &&
			(board.movedPiece.pieceValue == chess.WhitePawn && rank == rank8) ||
			(board.movedPiece.pieceValue == chess.BlackPawn && rank == rank1)

	if isPromotionMove {

		fakeMoveToBeDone := board.getMatchingMove(chess.Queen)
		if fakeMoveToBeDone == nil {
			board.resetDragAndDrop()
			board.Refresh()
			return
		}

		currentFen, _ := chess.FEN(board.game.FEN())
		gameClone := chess.NewGame(chess.UseNotation(chess.LongAlgebraicNotation{}), currentFen)

		err := gameClone.Move(fakeMoveToBeDone)

		if err != nil {
			board.resetDragAndDrop()
			board.Refresh()
			return
		}

		board.pendingPromotion = true
		board.launchPromotionDialog()
		return
	}

	moveToBeDone := board.getMatchingMove(chess.NoPieceType)
	if moveToBeDone == nil {
		board.resetDragAndDrop()
		board.Refresh()
		return
	}

	moveSan := chess.AlgebraicNotation{}.Encode(board.game.Position(), moveToBeDone)
	moveFan := convertSanToFan(moveSan, board.game.Position().Turn() == chess.White)

	err := board.game.Move(moveToBeDone)
	if err == nil {
		board.lastMove = &lastMove{
			originCell: board.movedPiece.startCell,
			targetCell: board.movedPiece.endCell,
		}
		if board.onMoveDone != nil {
			board.onMoveDone(moveFan)
		}
	}
	board.resetDragAndDrop()
	board.updatePieces()
	board.Refresh()

	board.handleGameEndedStatus()
}

func convertSanToFan(san string, whiteMove bool) string {
	fan := san
	var kingChange, queenChange, rookChange, bishopChange, knightChange string
	if whiteMove {
		kingChange = "\u2654"
		queenChange = "\u2655"
		rookChange = "\u2656"
		bishopChange = "\u2657"
		knightChange = "\u2658"
	} else {
		kingChange = "\u265A"
		queenChange = "\u265B"
		rookChange = "\u265C"
		bishopChange = "\u265D"
		knightChange = "\u265E"
	}

	fan = strings.ReplaceAll(fan, "K", kingChange)
	fan = strings.ReplaceAll(fan, "Q", queenChange)
	fan = strings.ReplaceAll(fan, "R", rookChange)
	fan = strings.ReplaceAll(fan, "B", bishopChange)
	fan = strings.ReplaceAll(fan, "N", knightChange)

	return fan
}

// ClaimDraw emits a draw claim (for 3-folds repetitions, or for 50-moves rule).
// Returns true if the draw has been accepted, otherwise false.
func (board *ChessBoard) ClaimDraw() bool {
	possibleDraws := board.game.EligibleDraws()
	drawByThreefoldRepetition, drawByFiftyMovesRule := false, false

	for _, item := range possibleDraws {
		if item == chess.ThreefoldRepetition {
			drawByThreefoldRepetition = true
			break
		} else if item == chess.FiftyMoveRule {
			drawByFiftyMovesRule = true
			break
		}
	}

	if drawByThreefoldRepetition {
		err := board.game.Draw(chess.ThreefoldRepetition)
		if err != nil {
			return true
		}
	}

	if drawByFiftyMovesRule {
		err := board.game.Draw(chess.FiftyMoveRule)
		if err != nil {
			return true
		}
	}

	return false
}

func (board *ChessBoard) startDragAndDrop(event *fyne.DragEvent) {

	if !board.gameInProgress {
		return
	}

	cellsLength := int(float64(board.length) / 9)
	halfCellsLength := int(float64(cellsLength) / 2)

	position := event.Position
	// This is really needed to be coded as is !
	// First file and rank, then bounds test, then adjust values with the board orientation
	file := int8(math.Floor(float64(position.X-halfCellsLength) / float64(cellsLength)))
	rank := int8(math.Floor(float64(position.Y-halfCellsLength) / float64(cellsLength)))

	inBounds := file >= 0 && file <= 7 && rank >= 0 && rank <= 7
	if !inBounds {
		return
	}

	if board.blackSide == BlackAtTop {
		rank = 7 - rank
	} else {
		file = 7 - file
	}

	square := chess.Square(file + 8*rank)
	pieceValue := board.game.Position().Board().Piece(square)

	pieceSide := pieceValue.Color()
	pieceBelongsToSideInTurn := pieceSide == board.game.Position().Turn()
	if !pieceBelongsToSideInTurn {
		return
	}

	if pieceValue == chess.NoPiece {
		return
	}

	board.dragndropInProgress = true

	imageResource := imageResourceFromPiece(pieceValue)
	image := canvas.NewImageFromResource(&imageResource)
	image.FillMode = canvas.ImageFillContain
	movedPiece := movedPiece{}
	movedPiece.pieceImage = image
	movedPiece.location = fyne.Position{X: position.X - halfCellsLength, Y: position.Y - halfCellsLength}
	movedPiece.startCell = cell{file: file, rank: rank}
	movedPiece.endCell = cell{file: file, rank: rank}
	movedPiece.pieceValue = pieceValue
	board.movedPiece = &movedPiece
	board.Refresh()
}

func (board *ChessBoard) updateDragAndDrop(event *fyne.DragEvent) {
	if !board.gameInProgress {
		return
	}

	cellsLength := int(float64(board.length) / 9)
	halfCellsLength := int(float64(cellsLength) / 2)

	position := event.Position
	var file, rank int8
	if board.blackSide == BlackAtTop {
		file = int8(math.Floor(float64(position.X-halfCellsLength) / float64(cellsLength)))
		rank = int8(7 - int8(math.Floor(float64(position.Y-halfCellsLength)/float64(cellsLength))))
	} else {
		file = int8(7 - int8(math.Floor(float64(position.X-halfCellsLength)/float64(cellsLength))))
		rank = int8(math.Floor(float64(position.Y-halfCellsLength) / float64(cellsLength)))
	}

	board.movedPiece.location = fyne.Position{X: position.X - halfCellsLength, Y: position.Y - halfCellsLength}
	board.movedPiece.endCell = cell{file: file, rank: rank}
	board.Refresh()
}

func (board *ChessBoard) resetDragAndDrop() {
	board.dragndropInProgress = false
	board.movedPiece.location = fyne.Position{X: -1000, Y: -1000}
	board.movedPiece.startCell = cell{file: -1, rank: -1}
}

func (board *ChessBoard) getMatchingMove(promotionPiece chess.PieceType) *chess.Move {
	possibleMoves := board.game.ValidMoves()

	for _, currentMove := range possibleMoves {
		if int8(currentMove.S1().File()) == board.movedPiece.startCell.file &&
			int8(currentMove.S1().Rank()) == board.movedPiece.startCell.rank &&
			int8(currentMove.S2().File()) == board.movedPiece.endCell.file &&
			int8(currentMove.S2().Rank()) == board.movedPiece.endCell.rank &&
			currentMove.Promo() == promotionPiece {
			return currentMove
		}
	}

	return nil
}

func (board ChessBoard) buildCellsAndPieces(cells *[8][8]*canvas.Rectangle, pieces *[8][8]*canvas.Image) {
	whiteCellColor := color.RGBA{255, 206, 158, 0xff}
	//blackCellColor := color.RGBA{209, 139, 71, 0xff}

	for line := 0; line < 8; line++ {
		for col := 0; col < 8; col++ {
			isWhiteCell := (line+col)%2 == 0
			var cellColor color.Color
			if isWhiteCell {
				cellColor = whiteCellColor
			} else {
				cellColor = whiteCellColor
			}

			cellRef := canvas.NewRectangle(cellColor)
			cells[line][col] = cellRef

			square := chess.Square(col + 8*line)
			pieceValue := board.game.Position().Board().Piece(square)
			if pieceValue != chess.NoPiece {
				imageResource := imageResourceFromPiece(pieceValue)
				image := canvas.NewImageFromResource(&imageResource)
				image.FillMode = canvas.ImageFillContain

				pieces[line][col] = image
			}
		}
	}
}

func (board ChessBoard) buildFilesCoordinates(filesCoords *[2][8]*canvas.Text) {
	coordsColor := color.RGBA{255, 199, 0, 0xff}
	asciiLowerA := 97

	for file := 0; file < 8; file++ {
		coord := fmt.Sprintf("%c", asciiLowerA+file)
		topCoord := canvas.NewText(coord, coordsColor)
		bottomCoord := canvas.NewText(coord, coordsColor)

		filesCoords[0][file] = topCoord
		filesCoords[1][file] = bottomCoord
	}
}

func (board ChessBoard) buildRanksCoordinates(ranksCoords *[2][8]*canvas.Text) {
	coordsColor := color.RGBA{255, 199, 0, 0xff}
	asciiOne := 49

	for rank := 0; rank < 8; rank++ {
		coord := fmt.Sprintf("%c", asciiOne+(7-rank))
		leftCoord := canvas.NewText(coord, coordsColor)
		rightCoord := canvas.NewText(coord, coordsColor)

		ranksCoords[0][rank] = leftCoord
		ranksCoords[1][rank] = rightCoord
	}

}

func (board ChessBoard) buildPlayerTurn() *canvas.Circle {
	var playerTurnColor color.Color
	gameTurn := board.game.Position().Turn()
	if gameTurn == chess.White {
		playerTurnColor = color.White
	} else {
		playerTurnColor = color.Black
	}

	return canvas.NewCircle(playerTurnColor)
}

func (board *ChessBoard) updatePieces() {
	for line := 0; line < 8; line++ {
		for col := 0; col < 8; col++ {
			board.pieces[line][col] = nil

			square := chess.Square(col + 8*line)
			pieceValue := board.game.Position().Board().Piece(square)
			if pieceValue != chess.NoPiece {
				imageResource := imageResourceFromPiece(pieceValue)
				image := canvas.NewImageFromResource(&imageResource)
				image.FillMode = canvas.ImageFillContain

				board.pieces[line][col] = image
			}
		}
	}
}

func (board *ChessBoard) commitPromotion(pieceType chess.PieceType) {
	if pieceType == chess.Pawn || pieceType == chess.King {
		return
	}

	moveToBeDone := board.getMatchingMove(pieceType)

	if moveToBeDone == nil {
		board.pendingPromotion = false
		board.resetDragAndDrop()
		board.Refresh()
	}

	moveSan := chess.AlgebraicNotation{}.Encode(board.game.Position(), moveToBeDone)
	moveFan := convertSanToFan(moveSan, board.game.Position().Turn() == chess.White)

	err := board.game.Move(moveToBeDone)
	if err == nil {
		board.lastMove = &lastMove{
			originCell: board.movedPiece.startCell,
			targetCell: board.movedPiece.endCell,
		}
		if board.onMoveDone != nil {
			board.onMoveDone(moveFan)
		}
	}

	board.pendingPromotion = false
	board.resetDragAndDrop()
	board.updatePieces()
	board.Refresh()

	board.handleGameEndedStatus()
}

func (board *ChessBoard) handleGameEndedStatus() {
	gameOutcome := board.game.Outcome()
	switch gameOutcome {
	case chess.WhiteWon:
		board.gameInProgress = false
		if board.onWhiteWin != nil {
			board.onWhiteWin()
		}
	case chess.BlackWon:
		board.gameInProgress = false
		if board.onBlackWin != nil {
			board.onBlackWin()
		}
	case chess.Draw:
		board.gameInProgress = false
		if board.onDraw != nil {
			board.onDraw()
		}
	}

}

func (board *ChessBoard) launchPromotionDialog() {
	title := ini.String("promotionDialog.title")
	dismiss := ini.String("promotionDialog.dismissButton")

	var queenRes, rookRes, bishopRes, knightRes *fyne.StaticResource
	if board.game.Position().Turn() == chess.White {
		queenRes = resourceChessqlt45Svg
		rookRes = resourceChessrlt45Svg
		bishopRes = resourceChessblt45Svg
		knightRes = resourceChessnlt45Svg
	} else {
		queenRes = resourceChessqdt45Svg
		rookRes = resourceChessrdt45Svg
		bishopRes = resourceChessbdt45Svg
		knightRes = resourceChessndt45Svg
	}

	cellsLength := int(float64(board.length) / 9)
	commonButtonsSize := fyne.NewSize(cellsLength, cellsLength)

	queenButton := NewIconButton(queenRes, commonButtonsSize, func() {
		board.commitPromotion(chess.Queen)
		board.promotionDialog.Hide()
	})

	rookButton := NewIconButton(rookRes, commonButtonsSize, func() {
		board.commitPromotion(chess.Rook)
		board.promotionDialog.Hide()
	})

	bishopButton := NewIconButton(bishopRes, commonButtonsSize, func() {
		board.commitPromotion(chess.Bishop)
		board.promotionDialog.Hide()
	})

	knightButton := NewIconButton(knightRes, commonButtonsSize, func() {
		board.commitPromotion(chess.Knight)
		board.promotionDialog.Hide()
	})

	content := fyne.NewContainerWithLayout(layout.NewGridLayout(4))
	content.AddObject(queenButton)
	content.AddObject(rookButton)
	content.AddObject(bishopButton)
	content.AddObject(knightButton)

	board.promotionDialog = dialog.NewCustom(title, dismiss, content, *board.parent)
	board.promotionDialog.SetOnClosed(func() {
		board.pendingPromotion = false
		board.resetDragAndDrop()
		board.Refresh()
	})
	board.promotionDialog.Show()
}
