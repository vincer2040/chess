package game

import (
	"unicode"
)

const (
    None = 0
    Pawn = 1
    Knight = 2
    Bishop = 3
    Rook = 4
    Queen = 5
    King = 6

    White = 8
    Black = 16

    PIECEMASK = 7
    COLORMASK = 24
)

type Piece byte

func newPiece(p byte) Piece {
    switch p {
    case ' ':
        return None
    case 'P':
        return Pawn | White
    case 'N':
        return Knight | White
    case 'B':
        return Bishop | White
    case 'R':
        return Rook | White
    case 'Q':
        return Queen | White
    case 'K':
        return King | White
    case 'p':
        return Pawn | Black
    case 'n':
        return Knight | Black
    case 'b':
        return Bishop | Black
    case 'r':
        return Rook | Black
    case 'q':
        return Queen | Black
    case 'k':
        return King | Black
    }

    panic("unknown piece")
}

func (p Piece) GetPieceByte() byte {
    piece := p & PIECEMASK
    color := p & COLORMASK
    var b byte
    switch piece {
    case None:
        b = ' '
    case Pawn:
        b = 'p'
    case Knight:
        b = 'n'
    case Bishop:
        b = 'b'
    case Rook:
        b = 'r'
    case Queen:
        b = 'q'
    case King:
        b = 'k'
    }
    if color == White {
        b = byte(unicode.ToUpper(rune(b)))
    }
    return b
}

