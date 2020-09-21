package commonTypes

// GameMove defines a move of an History widget.
type GameMove struct {
	// Fan is the move notation with figurines.
	Fan string

	// Fen is the position after the move, in Forsyth-Edwards Notation.
	Fen string

	// LastMoveOriginCell is the origin cell of the move.
	LastMoveOriginCell Cell

	// LastMoveTargetCell is the target cell of the move.
	LastMoveTargetCell Cell

	// MoveNumberMarker, if defined, makes this not a real move, but a marker for
	// adding a move number.
	MoveNumberMarker int

	// IsBlackMove says whether it is a black move.
	IsBlackMove bool
}

// Cell defines a coordinate of a Chess Board widget.
type Cell struct {
	File int8
	Rank int8
}
