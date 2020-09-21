package commonTypes

// GameMove defines a move of an History widget.
type GameMove struct {
	Fan                string
	Fen                string
	LastMoveOriginCell Cell
	LastMoveTargetCell Cell
}

// Cell defines a coordinate of a Chess Board widget.
type Cell struct {
	File int8
	Rank int8
}
