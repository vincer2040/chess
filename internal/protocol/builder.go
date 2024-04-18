package protocol

import (
	"strconv"

	"github.com/vincer2040/chess/internal/game"
	"github.com/vincer2040/chess/internal/types"
)

type Builder []byte

func NewBuilder() Builder {
	return []byte{}
}

func (b Builder) AddPosition(position string) Builder {
	b = append(b, POSITION_BYTE)
	for _, ch := range position {
		b = append(b, byte(ch))
	}
	return b.addEnd()
}

func (b Builder) AddLegalMoves(legalMoves game.LegalMoves) Builder {
	b = append(b, LEGAL_MOVES_BYTE)
	amt := strconv.Itoa(len(legalMoves))
	for _, ch := range amt {
		b = append(b, byte(ch))
	}
	b = b.addEnd()
	for k, v := range legalMoves {
        // add the key
		key := strconv.Itoa(k)
		for _, ch := range key {
			b = append(b, byte(ch))
		}
		b = b.addEnd()
        // add the number of moves for this piece
        b = append(b, ARRAY_BYTE)
        l := strconv.Itoa(len(v))
        for _, ch := range l {
            b = append(b, byte(ch))
        }
		b = b.addEnd()
        // add the moves
		for i, idx := range v {
			s := strconv.Itoa(idx)
			for _, ch := range s {
				b = append(b, byte(ch))
			}
			if i != len(v)-1 {
				b = append(b, SEPARATOR)
			}
		}
		b = b.addEnd()
	}
	return b
}

func (b Builder) AddMove(move *types.Move) Builder {
	from := strconv.Itoa(move.From)
	to := strconv.Itoa(move.To)
	b = append(b, MOVE_BYTE)
	for _, ch := range from {
		b = append(b, byte(ch))
	}
	b = append(b, SEPARATOR)
	for _, ch := range to {
		b = append(b, byte(ch))
	}
	return b.addEnd()
}

func (b Builder) AddCommand(command string) Builder {
	b = append(b, COMMAND_BYTE)
	for _, ch := range command {
		b = append(b, byte(ch))
	}
	return b.addEnd()
}

func (b Builder) AddError(error string) Builder {
	b = append(b, ERROR_BYTE)
	for _, ch := range error {
		b = append(b, byte(ch))
	}
	return b.addEnd()
}

func (b Builder) Reset() Builder {
    b = []byte{}
    return b
}

func (b Builder) addEnd() Builder {
	b = append(b, '\r')
	b = append(b, '\n')
	return b
}