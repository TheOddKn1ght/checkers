package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 28, G: 29, B: 34, A: 255})
	g.drawBoard(screen)
	g.drawPieces(screen)
	g.drawHighlights(screen)
	g.drawInfo(screen)
}

func (g *Game) drawBoard(screen *ebiten.Image) {
	light := color.RGBA{R: 236, G: 220, B: 191, A: 255}
	dark := color.RGBA{R: 127, G: 91, B: 58, A: 255}
	for row := range BoardSize {
		for col := range BoardSize {
			c := light
			if (row+col)%2 == 1 {
				c = dark
			}
			x := float64(padding + col*cellSize)
			y := float64(padding + row*cellSize)
			ebitenutil.DrawRect(screen, x, y, float64(cellSize), float64(cellSize), c)
		}
	}
}

func (g *Game) drawPieces(screen *ebiten.Image) {
	whiteFill := color.RGBA{R: 248, G: 248, B: 252, A: 255}
	whiteEdge := color.RGBA{R: 180, G: 182, B: 190, A: 255}

	for row := range BoardSize {
		for col := range BoardSize {
			piece := g.board[row][col]
			if piece == Empty {
				continue
			}
			cx := float64(padding + col*cellSize + cellSize/2)
			cy := float64(padding + row*cellSize + cellSize/2)
			radius := float64(cellSize) * 0.38

			if piece.IsBlack() {
				drawFilledCircle(screen, cx, cy, radius, color.RGBA{R: 34, G: 34, B: 38, A: 255})
			} else {
				drawFilledCircle(screen, cx, cy, radius+1, whiteEdge)
				drawFilledCircle(screen, cx, cy, radius, whiteFill)
			}

			if piece.IsKing() {
				drawFilledCircle(screen, cx, cy, radius*0.45, color.RGBA{R: 246, G: 206, B: 88, A: 255})
			}
		}
	}
}

func (g *Game) drawHighlights(screen *ebiten.Image) {
	if g.selected != nil {
		x := float64(padding + g.selected.Col*cellSize)
		y := float64(padding + g.selected.Row*cellSize)
		ebitenutil.DrawRect(screen, x, y, float64(cellSize), float64(cellSize), color.RGBA{R: 68, G: 133, B: 255, A: 110})
	}
	for _, mv := range g.validMoves {
		x := float64(padding + mv.To.Col*cellSize + cellSize/2)
		y := float64(padding + mv.To.Row*cellSize + cellSize/2)
		drawFilledCircle(screen, x, y, float64(cellSize)*0.12, color.RGBA{R: 85, G: 180, B: 90, A: 220})
	}
	if g.lastMovedTo != nil {
		x := float64(padding + g.lastMovedTo.Col*cellSize)
		y := float64(padding + g.lastMovedTo.Row*cellSize)
		ebitenutil.DrawRect(screen, x, y, float64(cellSize), float64(cellSize), color.RGBA{R: 17, G: 73, B: 105, A: 70})
	}
}

func (g *Game) drawInfo(screen *ebiten.Image) {
	top := padding + BoardSize*cellSize + 16
	turn := "Black"
	if g.whiteTurn {
		turn = "White"
	}
	if g.winnerMessage != "" {
		turn = g.winnerMessage
	}

	drawText(screen, fmt.Sprintf("Turn: %s", turn), padding, top, g.fontFace, color.White)
	drawText(screen, g.statusMessage, padding, top+34, g.fontFace, color.RGBA{R: 211, G: 215, B: 219, A: 255})
	drawText(screen, "Controls: Click to move, R to reset", padding, top+68, g.fontFace, color.RGBA{R: 168, G: 174, B: 180, A: 255})
}

func drawFilledCircle(screen *ebiten.Image, cx, cy, r float64, c color.Color) {
	for y := int(cy - r); y <= int(cy+r); y++ {
		for x := int(cx - r); x <= int(cx+r); x++ {
			dx := float64(x) - cx
			dy := float64(y) - cy
			if dx*dx+dy*dy <= r*r {
				screen.Set(x, y, c)
			}
		}
	}
}

func drawText(screen *ebiten.Image, msg string, x, y int, face font.Face, clr color.Color) {
	text.Draw(screen, msg, face, x, y, clr)
}
