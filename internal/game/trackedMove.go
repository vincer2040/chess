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

func (tm *TrackedMove) isCastle() bool {
	piece := tm.Piece & PIECEMASK
	if piece != King {
		return false
	}
	amtMoved := int(math.Abs(float64(tm.To - tm.From)))
	return amtMoved == 2
}

func (tm *TrackedMove) isDoublePawnPush() bool {
	piece := tm.Piece & PIECEMASK
	if piece != Pawn {
		return false
	}
	amtMoved := int(math.Abs(float64(tm.To - tm.From)))
	return amtMoved == 16
}

func (tm *TrackedMove) isCapture() bool {
	return tm.Captured != None
}

func (tm *TrackedMove) isEnPassant() bool {
	piece := tm.Piece & PIECEMASK
	// we won't be able to tell initially
	// that it is a capture even though
	// it is because en passant moves
	// go to a None square
	if tm.isCapture() {
		return false
	}
	if piece != Pawn {
		return false
	}
	amtMoved := int(math.Abs(float64(tm.To - tm.From)))
	return amtMoved == 7 || amtMoved == 9
}

func (tm *TrackedMove) disablesCastle(castleRights *CastleRights) (bool, []DisabledCastleDirection) {
	color := tm.Piece & COLORMASK
	piece := tm.Piece & PIECEMASK
	disabled := make([]DisabledCastleDirection, 0)

	if color == White {
		if !castleRights.WhiteKing && !castleRights.WhiteQueen {
			return false, disabled
		}
	} else {
		if !castleRights.BlackKing && !castleRights.BlackQueen {
			return false, disabled
		}
	}

	if piece == King {
		if color == White {
			disabled = append(disabled, WhiteCastleKing, WhiteCastleQueen)
		} else {
			disabled = append(disabled, BlackCastleKing, BlackCastleQueen)
		}
	}

	if piece == Rook {
		switch tm.From {
		case 56:
			disabled = append(disabled, WhiteCastleQueen)
			break
		case 63:
			disabled = append(disabled, WhiteCastleKing)
			break
		case 0:
			disabled = append(disabled, BlackCastleQueen)
			break
		case 7:
			disabled = append(disabled, BlackCastleKing)
			break
		}
	}

	if !tm.isCapture() {
		if len(disabled) == 0 {
			return false, disabled
		}
		return true, disabled
	}

	captured := tm.Captured & PIECEMASK
	if captured == Rook {
		switch tm.To {
		case 56:
			disabled = append(disabled, WhiteCastleQueen)
			break
		case 63:
			disabled = append(disabled, WhiteCastleKing)
			break
		case 0:
			disabled = append(disabled, BlackCastleQueen)
			break
		case 7:
			disabled = append(disabled, BlackCastleKing)
			break
		}
	}

	if len(disabled) == 0 {
		return false, disabled
	}

	return true, disabled
}
