package types

type DataType int

const (
	IllegalType DataType = iota
	PositionType
	MoveType
	ErrorType
	CommandType
)

type DataInterface interface {
	data()
}

type Data struct {
	Type DataType
	Data DataInterface
}

type Command string

type Position string

type Move struct {
	From int
	To   int
}

type Error string

func (c Command) data()  {}
func (p Position) data() {}
func (m Move) data()     {}
func (e Error) data()    {}
