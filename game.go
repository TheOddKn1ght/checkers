package main

import (
	"fmt"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

var kingDiagonals = []Position{
	{Row: -1, Col: -1}, {Row: -1, Col: 1},
	{Row: 1, Col: -1}, {Row: 1, Col: 1},
}

var manCaptureDiagonals = []Position{
	{Row: -1, Col: -1}, {Row: -1, Col: 1},
	{Row: 1, Col: -1}, {Row: 1, Col: 1},
}

type Game struct {
	board         [BoardSize][BoardSize]Piece
	whiteTurn     bool
	selected      *Position
	validMoves    []Move
	lastMovedTo   *Position
	mustContinue  bool
	winnerMessage string
	fontFace      font.Face
	statusMessage string
}

func NewGame() *Game {
	g := &Game{
		whiteTurn: true,
		fontFace:  basicfont.Face7x13,
	}
	g.resetBoard()
	return g
}

func (g *Game) resetBoard() {
	for row := range BoardSize {
		for col := range BoardSize {
			g.board[row][col] = Empty
		}
	}

	for row := range 3 {
		for col := range BoardSize {
			if (row+col)%2 == 1 {
				g.board[row][col] = WhiteMan
			}
		}
	}
	for row := range 3 {
		targetRow := BoardSize - 1 - row
		for col := range BoardSize {
			if (targetRow+col)%2 == 1 {
				g.board[targetRow][col] = BlackMan
			}
		}
	}

	g.whiteTurn = true
	g.selected = nil
	g.validMoves = nil
	g.lastMovedTo = nil
	g.mustContinue = false
	g.winnerMessage = ""
	g.statusMessage = "White starts"
}

func (g *Game) applyMove(m Move) {
	piece := g.board[m.From.Row][m.From.Col]
	g.board[m.From.Row][m.From.Col] = Empty
	g.board[m.To.Row][m.To.Col] = piece
	g.lastMovedTo = &Position{Row: m.To.Row, Col: m.To.Col}

	wasCapture := m.Capture != nil
	if m.Capture != nil {
		g.board[m.Capture.Row][m.Capture.Col] = Empty
	}

	promoted := false
	if piece == BlackMan && m.To.Row == 0 {
		g.board[m.To.Row][m.To.Col] = BlackKing
		promoted = true
	}
	if piece == WhiteMan && m.To.Row == BoardSize-1 {
		g.board[m.To.Row][m.To.Col] = WhiteKing
		promoted = true
	}
	piece = g.board[m.To.Row][m.To.Col]

	if wasCapture {
		nextPos := Position{Row: m.To.Row, Col: m.To.Col}
		nextCaptures := g.legalCapturesFrom(nextPos)
		if len(nextCaptures) > 0 {
			g.selected = &nextPos
			g.validMoves = nextCaptures
			g.mustContinue = true
			if promoted && piece.IsKing() {
				g.statusMessage = "Promoted to king! Capture again with the same piece"
			} else {
				g.statusMessage = "Capture again with the same piece"
			}
			return
		}
	}

	g.mustContinue = false
	g.selected = nil
	g.validMoves = nil
	g.whiteTurn = !g.whiteTurn
	g.updateGameStateAfterTurn()
}

func (g *Game) updateGameStateAfterTurn() {
	blackPieces, whitePieces := g.countPieces()
	if blackPieces == 0 {
		g.winnerMessage = "White wins! Press R to restart"
		g.statusMessage = "Black has no pieces left"
		return
	}
	if whitePieces == 0 {
		g.winnerMessage = "Black wins! Press R to restart"
		g.statusMessage = "White has no pieces left"
		return
	}

	playerMoves := g.allMovesForCurrentPlayer()
	if len(playerMoves) == 0 {
		if g.whiteTurn {
			g.winnerMessage = "Black wins! Press R to restart"
			g.statusMessage = "White has no legal moves"
		} else {
			g.winnerMessage = "White wins! Press R to restart"
			g.statusMessage = "Black has no legal moves"
		}
		return
	}

	turn := "Black"
	if g.whiteTurn {
		turn = "White"
	}
	g.statusMessage = fmt.Sprintf("%s to move", turn)
}

func (g *Game) countPieces() (black, white int) {
	for row := range BoardSize {
		for col := range BoardSize {
			p := g.board[row][col]
			if p.IsBlack() {
				black++
			} else if p.IsWhite() {
				white++
			}
		}
	}
	return black, white
}

func (g *Game) allMovesForCurrentPlayer() []Move {
	var allMoves []Move
	for row := range BoardSize {
		for col := range BoardSize {
			p := g.board[row][col]
			if p == Empty {
				continue
			}
			if g.whiteTurn && !p.IsWhite() {
				continue
			}
			if !g.whiteTurn && !p.IsBlack() {
				continue
			}
			moves := g.movesForPiece(Position{Row: row, Col: col}, true)
			allMoves = append(allMoves, moves...)
		}
	}
	return allMoves
}

// takes a sequence that captures the maximum possible number of pieces
func (g *Game) legalCapturesFrom(pos Position) []Move {
	raw := rawCaptureMovesForPiece(g.board, pos)
	if len(raw) == 0 {
		return nil
	}
	return g.filterCapturesByMaxChain(g.board, pos, raw)
}

func (g *Game) movesForPiece(pos Position, respectForcedCapture bool) []Move {
	captures := g.legalCapturesFrom(pos)
	if len(captures) > 0 {
		return captures
	}
	if respectForcedCapture && g.playerHasAnyCapture() {
		return nil
	}
	return normalMovesForPiece(g.board, pos)
}

func (g *Game) playerHasAnyCapture() bool {
	return g.maxCaptureChainForPlayer(g.board) > 0
}

func (g *Game) maxCaptureChainForPlayer(b [BoardSize][BoardSize]Piece) int {
	best := 0
	for row := range BoardSize {
		for col := range BoardSize {
			p := b[row][col]
			if p == Empty {
				continue
			}
			if g.whiteTurn && !p.IsWhite() {
				continue
			}
			if !g.whiteTurn && !p.IsBlack() {
				continue
			}
			n := maxCaptureChainFrom(b, Position{Row: row, Col: col})
			if n > best {
				best = n
			}
		}
	}
	return best
}

func (g *Game) filterCapturesByMaxChain(b [BoardSize][BoardSize]Piece, pos Position, raw []Move) []Move {
	var target int
	if g.mustContinue && g.selected != nil && *g.selected == pos {
		target = maxCaptureChainFrom(b, pos)
	} else {
		target = g.maxCaptureChainForPlayer(b)
	}
	var out []Move
	for _, m := range raw {
		nb := applyMoveBoard(b, m)
		chain := 1 + maxCaptureChainFrom(nb, m.To)
		if chain == target {
			out = append(out, m)
		}
	}
	return out
}

func maxCaptureChainFrom(b [BoardSize][BoardSize]Piece, pos Position) int {
	raw := rawCaptureMovesForPiece(b, pos)
	if len(raw) == 0 {
		return 0
	}
	best := 0
	for _, m := range raw {
		nb := applyMoveBoard(b, m)
		n := 1 + maxCaptureChainFrom(nb, m.To)
		if n > best {
			best = n
		}
	}
	return best
}

func applyMoveBoard(b [BoardSize][BoardSize]Piece, m Move) [BoardSize][BoardSize]Piece {
	nb := b
	piece := nb[m.From.Row][m.From.Col]
	nb[m.From.Row][m.From.Col] = Empty
	nb[m.To.Row][m.To.Col] = piece
	if m.Capture != nil {
		nb[m.Capture.Row][m.Capture.Col] = Empty
	}
	if piece == BlackMan && m.To.Row == 0 {
		nb[m.To.Row][m.To.Col] = BlackKing
	}
	if piece == WhiteMan && m.To.Row == BoardSize-1 {
		nb[m.To.Row][m.To.Col] = WhiteKing
	}
	return nb
}

func normalMovesForPiece(b [BoardSize][BoardSize]Piece, pos Position) []Move {
	piece := b[pos.Row][pos.Col]
	if piece == Empty {
		return nil
	}
	if piece.IsKing() {
		var moves []Move
		for _, d := range kingDiagonals {
			for step := 1; ; step++ {
				to := Position{Row: pos.Row + d.Row*step, Col: pos.Col + d.Col*step}
				if !to.InBounds() {
					break
				}
				if b[to.Row][to.Col] != Empty {
					break
				}
				moves = append(moves, Move{From: pos, To: to})
			}
		}
		return moves
	}
	dirs := manForwardDirections(piece)
	moves := make([]Move, 0, len(dirs))
	for _, d := range dirs {
		to := Position{Row: pos.Row + d.Row, Col: pos.Col + d.Col}
		if !to.InBounds() {
			continue
		}
		if b[to.Row][to.Col] != Empty {
			continue
		}
		moves = append(moves, Move{From: pos, To: to})
	}
	return moves
}

func manForwardDirections(piece Piece) []Position {
	if piece.IsBlack() {
		return []Position{{Row: -1, Col: -1}, {Row: -1, Col: 1}}
	}
	return []Position{{Row: 1, Col: -1}, {Row: 1, Col: 1}}
}

func rawCaptureMovesForPiece(b [BoardSize][BoardSize]Piece, pos Position) []Move {
	piece := b[pos.Row][pos.Col]
	if piece == Empty {
		return nil
	}
	if piece.IsKing() {
		return rawKingCaptures(b, pos, piece)
	}
	return rawManCaptures(b, pos, piece)
}

func rawManCaptures(b [BoardSize][BoardSize]Piece, pos Position, piece Piece) []Move {
	moves := make([]Move, 0, 4)
	for _, d := range manCaptureDiagonals {
		mid := Position{Row: pos.Row + d.Row, Col: pos.Col + d.Col}
		to := Position{Row: pos.Row + 2*d.Row, Col: pos.Col + 2*d.Col}
		if !mid.InBounds() || !to.InBounds() {
			continue
		}
		target := b[mid.Row][mid.Col]
		if target == Empty {
			continue
		}
		if piece.IsBlack() && !target.IsWhite() {
			continue
		}
		if piece.IsWhite() && !target.IsBlack() {
			continue
		}
		if b[to.Row][to.Col] != Empty {
			continue
		}
		captured := Position{Row: mid.Row, Col: mid.Col}
		cp := captured
		moves = append(moves, Move{From: pos, To: to, Capture: &cp})
	}
	return moves
}

func rawKingCaptures(b [BoardSize][BoardSize]Piece, pos Position, piece Piece) []Move {
	var moves []Move
	for _, d := range kingDiagonals {
		step := 1
		for {
			scan := Position{Row: pos.Row + d.Row*step, Col: pos.Col + d.Col*step}
			if !scan.InBounds() {
				break
			}
			cell := b[scan.Row][scan.Col]
			if cell == Empty {
				step++
				continue
			}
			if piece.IsBlack() && !cell.IsWhite() {
				break
			}
			if piece.IsWhite() && !cell.IsBlack() {
				break
			}
			cap := Position{Row: scan.Row, Col: scan.Col}
			for land := step + 1; ; land++ {
				to := Position{Row: pos.Row + d.Row*land, Col: pos.Col + d.Col*land}
				if !to.InBounds() {
					break
				}
				if b[to.Row][to.Col] != Empty {
					break
				}
				c := cap
				moves = append(moves, Move{From: pos, To: to, Capture: &c})
			}
			break
		}
	}
	return moves
}
