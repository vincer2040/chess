package game

import "math"

type TrackedMove struct {
	Piece    Piece
	Captured Piece
	From     int
	To       int
}

func newTrackedMove(piece, captured Piece, from, to int) TrackedMove {
	return TrackedMove{
		Piece:    piece,
		Captured: captured,
		From:     from,
		To:       to,
	}
}

func (tm *TrackedMove) IsCastle() bool {
	piece := tm.Piece & PIECEMASK
	if piece != King {
		return false
	}
	amtMoved := int(math.Abs(float64(tm.To - tm.From)))
	return amtMoved == 2
}

func (tm *TrackedMove) IsDoublePawnPush() bool {
	piece := tm.Piece & PIECEMASK
	if piece != Pawn {
		return false
	}
	amtMoved := int(math.Abs(float64(tm.To - tm.From)))
	return amtMoved == 16
}

func (tm *TrackedMove) IsCapture() bool {
    return tm.Captured != None
}
