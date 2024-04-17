package game

import "strings"

type Game struct {
	board        Board
	toMove       byte
	castleRights CastleRights
	enPessant    int
	legalMoves   LegalMoves
}

func New(fen string) Game {
	split := strings.Split(fen, " ")
	p := split[0]
	board := newBoard(p)
	toMove := byte(split[1][0])
	castleRights := split[2]
    g :=  Game{
		board:        board,
		toMove:       toMove,
		castleRights: newCastleRights(castleRights),
		enPessant:    -1,
        legalMoves: nil,
	}
    g.legalMoves = getLegalMoves(g.board, g.toMove, &g.castleRights, g.enPessant)
    return g
}

func (g *Game) getLegalMoves() LegalMoves {
	return g.legalMoves
}

func (g *Game) PrintBoard() {
	g.board.Print()
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
