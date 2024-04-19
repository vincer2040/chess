package game

import (
	"math"
)

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

func (tm *TrackedMove) IsEnPassant() bool {
	piece := tm.Piece & PIECEMASK
	// we won't be able to tell initially
	// that it is a capture even though
	// it is because en passant moves
	// go to a None square
	if tm.IsCapture() {
		return false
	}
	if piece != Pawn {
		return false
	}
	amtMoved := int(math.Abs(float64(tm.To - tm.From)))
	return amtMoved == 7 || amtMoved == 9
}
