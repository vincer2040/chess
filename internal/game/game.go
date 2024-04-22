package game

import (
	"strings"

	"github.com/vincer2040/chess/internal/types"
)

type DisabledCastleDirection int

const (
	WhiteCastleKing DisabledCastleDirection = iota
	WhiteCastleQueen
	BlackCastleKing
	BlackCastleQueen
)

type Game struct {
	board          Board
	trackedMoves   []TrackedMove
	toMove         byte
	castleRights   CastleRights
	enPassant      int
	legalMoves     LegalMoves
	attackingMoves AttackingMoves
}

func New(fen string) Game {
	split := strings.Split(fen, " ")
	p := split[0]
	board := newBoard(p)
	toMove := byte(split[1][0])
	castleRights := split[2]
	g := Game{
		board:          board,
		trackedMoves:   make([]TrackedMove, 0),
		toMove:         toMove,
		castleRights:   newCastleRights(castleRights),
		enPassant:      -1,
		legalMoves:     nil,
		attackingMoves: nil,
	}
	g.attackingMoves = getAttackingMoves(g.board, g.toMove)
	g.legalMoves = getLegalMoves(g.board, g.toMove, &g.castleRights, g.enPassant, g.attackingMoves)
	return g
}

func (g *Game) MakeMove(move *types.Move) {
	movedPiece := g.board[move.From]
	captured := g.board[move.To]
	trackedMove := newTrackedMove(movedPiece, captured, move.From, move.To)
	g.board[move.To] = movedPiece
	g.board[move.From] = None

	if trackedMove.isCastle() {
		g.castle(move)
	}

	disablesCastle, disabledCastleDirectoins := trackedMove.disablesCastle(&g.castleRights)
	if disablesCastle {
		g.disableCastle(disabledCastleDirectoins)
	}

	if trackedMove.isEnPassant() {
		g.board[g.enPassant] = None
	}

	if trackedMove.isDoublePawnPush() {
		g.enPassant = move.To
	} else {
		g.enPassant = -1
	}

	if g.toMove == 'w' {
		g.toMove = 'b'
	} else {
		g.toMove = 'w'
	}

	g.trackedMoves = append(g.trackedMoves, trackedMove)
	g.attackingMoves = getAttackingMoves(g.board, g.toMove)
	g.legalMoves = getLegalMoves(g.board, g.toMove, &g.castleRights, g.enPassant, g.attackingMoves)
}

func (g *Game) GetLegalMoves() LegalMoves {
	return g.legalMoves
}

func (g *Game) GetAttackingMoves() AttackingMoves {
	return g.attackingMoves
}

func (g *Game) PrintBoard() {
	g.board.print()
}

func (g *Game) disableCastle(directions []DisabledCastleDirection) {
	for _, dir := range directions {
		switch dir {
		case WhiteCastleKing:
			g.castleRights.WhiteKing = false
			break
		case WhiteCastleQueen:
			g.castleRights.WhiteQueen = false
			break
		case BlackCastleKing:
			g.castleRights.BlackKing = false
			break
		case BlackCastleQueen:
			g.castleRights.BlackQueen = false
			break
		}
	}
}

func (g *Game) castle(move *types.Move) {
	if g.toMove == 'w' {
		if move.To == 62 {
			g.board[63] = None
			g.board[61] = Rook | White
		} else {
			g.board[56] = None
			g.board[59] = Rook | White
		}
		g.castleRights.WhiteKing = false
		g.castleRights.WhiteQueen = false
	} else {
		if move.To == 6 {
			g.board[7] = None
			g.board[5] = Rook | Black
		} else {
			g.board[0] = None
			g.board[3] = Rook | Black
		}
		g.castleRights.BlackKing = false
		g.castleRights.BlackQueen = false
	}
}

type CastleRights struct {
	WhiteKing  bool
	WhiteQueen bool
	BlackKing  bool
	BlackQueen bool
}

func newCastleRights(castleRights string) CastleRights {
	whiteKing := strings.Contains(castleRights, "K")
	whiteQueen := strings.Contains(castleRights, "Q")
	blackKing := strings.Contains(castleRights, "k")
	blackQueen := strings.Contains(castleRights, "q")
	return CastleRights{
		WhiteKing:  whiteKing,
		WhiteQueen: whiteQueen,
		BlackKing:  blackKing,
		BlackQueen: blackQueen,
	}
}
