package types

type DataType int

const (
	IllegalType DataType = iota
	PositionType
	MoveType
	ErrorType
	CommandType
	PromotionType
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

type PromotedTo int

const (
	KnightPromotion PromotedTo = iota
	BishopPromotion
	RookPromotion
	QueenPromotion
)

type Promotion struct {
	Move
	PromoteTo PromotedTo
}

func (c Command) data()   {}
func (p Position) data()  {}
func (m Move) data()      {}
func (e Error) data()     {}
func (p Promotion) data() {}
