package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) trySelect(pos Position) {
	piece := g.board[pos.Row][pos.Col]
	if piece == Empty {
		return
	}
	if g.whiteTurn && !piece.IsWhite() {
		return
	}
	if !g.whiteTurn && !piece.IsBlack() {
		return
	}

	if g.mustContinue && (g.selected == nil || pos != *g.selected) {
		g.statusMessage = "You must continue capturing with the same piece"
		return
	}

	moves := g.movesForPiece(pos, true)
	if len(moves) == 0 {
		g.statusMessage = "This piece has no legal moves"
		return
	}

	g.selected = &Position{Row: pos.Row, Col: pos.Col}
	g.validMoves = moves
	g.statusMessage = "Piece selected"
}

func (g *Game) handleClick(x, y int) {
	if g.winnerMessage != "" {
		return
	}
	row := (y - padding) / cellSize
	col := (x - padding) / cellSize
	pos := Position{Row: row, Col: col}
	if !pos.InBounds() {
		return
	}
	if (row+col)%2 == 0 {
		return
	}

	if g.selected == nil {
		g.trySelect(pos)
		return
	}

	if pos == *g.selected {
		g.selected = nil
		g.validMoves = nil
		g.statusMessage = "Selection cleared"
		return
	}

	for _, mv := range g.validMoves {
		if mv.To == pos {
			g.applyMove(mv)
			return
		}
	}

	g.trySelect(pos)
}

type Runner struct {
	*Game
	mouseWasDown bool
}

func NewRunner(g *Game) *Runner {
	return &Runner{Game: g}
}

func (r *Runner) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		r.Game.resetBoard()
	}

	mouseDown := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
	if mouseDown && !r.mouseWasDown {
		x, y := ebiten.CursorPosition()
		r.Game.handleClick(x, y)
	}
	r.mouseWasDown = mouseDown
	return nil
}

func (r *Runner) Layout(_, _ int) (int, int) { return windowWidth, windowHeight }

func (r *Runner) Draw(screen *ebiten.Image) {
	r.Game.Draw(screen)
}
