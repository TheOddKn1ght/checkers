package main

type Piece int

const (
	Empty Piece = iota
	BlackMan
	BlackKing
	WhiteMan
	WhiteKing
)

func (p Piece) IsEmpty() bool { return p == Empty }
func (p Piece) IsBlack() bool { return p == BlackMan || p == BlackKing }
func (p Piece) IsWhite() bool { return p == WhiteMan || p == WhiteKing }
func (p Piece) IsKing() bool  { return p == BlackKing || p == WhiteKing }

type Position struct {
	Row int
	Col int
}

func (p Position) InBounds() bool {
	return p.Row >= 0 && p.Row < BoardSize && p.Col >= 0 && p.Col < BoardSize
}

type Move struct {
	From    Position
	To      Position
	Capture *Position
}
