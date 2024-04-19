package protocol

import (
	"bytes"
	"errors"
	"strconv"

	"github.com/vincer2040/chess/internal/types"
)

const (
	POSITION_BYTE = '+'
	MOVE_BYTE     = '$'
	COMMAND_BYTE  = '#'
	SEPARATOR     = ':'
	ERROR_BYTE    = '-'

	// client should never contain these two:
	LEGAL_MOVES_BYTE = '~'
	ARRAY_BYTE       = '*'
)

type Parser struct {
	input []byte
	pos   int
	ch    byte
}

func NewParser(input []byte) Parser {
	p := Parser{
		input: input,
		pos:   0,
		ch:    0,
	}
	p.readByte()
	return p
}

func (p *Parser) Parse() types.Data {
	res := types.Data{Type: types.IllegalType, Data: nil}
	switch p.ch {
	case POSITION_BYTE:
		pos := p.parsePosition()
		if pos == "" {
			break
		}
		res.Data = pos
		res.Type = types.PositionType
		break
	case COMMAND_BYTE:
		cmd := p.parseCommand()
		if cmd == "" {
			break
		}
		res.Data = cmd
		res.Type = types.CommandType
		break
	case MOVE_BYTE:
		move, err := p.parseMove()
		if err != nil {
			break
		}
		res.Data = move
		res.Type = types.MoveType
		break
	}
	return res
}

func (p *Parser) parsePosition() types.Position {
	p.readByte()
	buf := bytes.NewBufferString("")
	for p.ch != '\r' && p.ch != 0 {
		buf.WriteByte(p.ch)
		p.readByte()
	}
	if !p.expectEnd() {
		return ""
	}
	return types.Position(buf.String())
}

func (p *Parser) parseCommand() types.Command {
	p.readByte()
	buf := bytes.NewBufferString("")
	for p.ch != '\r' && p.ch != 0 {
		buf.WriteByte(p.ch)
		p.readByte()
	}
	if !p.expectEnd() {
		return ""
	}
	return types.Command(buf.String())
}

func (p *Parser) parseMove() (types.Move, error) {
	p.readByte()
	fromBuf := bytes.NewBufferString("")
	toBuf := bytes.NewBufferString("")
	for p.ch != SEPARATOR && p.ch != 0 {
		fromBuf.WriteByte(p.ch)
		p.readByte()
	}
	p.readByte()
	for p.ch != '\r' && p.ch != 0 {
		toBuf.WriteByte(p.ch)
		p.readByte()
	}
	if !p.expectEnd() {
		return types.Move{}, errors.New("expected end")
	}
	from, err := strconv.Atoi(fromBuf.String())
	if err != nil {
		return types.Move{}, err
	}
	to, err := strconv.Atoi(toBuf.String())
	if err != nil {
		return types.Move{}, err
	}
	return types.Move{
		From: from,
		To:   to,
	}, nil
}

func (p *Parser) expectEnd() bool {
	if p.ch != '\r' {
		return false
	}
	p.readByte()
	if p.ch != '\n' {
		return false
	}
	return true
}

func (p *Parser) readByte() {
	if p.pos >= len(p.input) {
		p.ch = 0
		return
	}
	p.ch = p.input[p.pos]
	p.pos++
}
