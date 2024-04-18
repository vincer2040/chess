package game

import (
	"strings"

	"github.com/vincer2040/chess/internal/types"
)

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

func (g *Game) MakeMove(move *types.Move) {
    movedPiece := g.board[move.From]
    g.board[move.To] = movedPiece
    g.board[move.From] = None
    if g.toMove == 'w' {
        g.toMove = 'b'
    } else {
        g.toMove = 'w'
    }
    g.legalMoves = getLegalMoves(g.board, g.toMove, &g.castleRights, g.enPessant)
}

func (g *Game) GetLegalMoves() LegalMoves {
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
