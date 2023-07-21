package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/utkarsh-1905/conways-game/graphics"
)

var square = []float32{
	-0.5, 0.5, 0,
	-0.5, -0.5, 0,
	0.5, -0.5, 0,
	-0.5, 0.5, 0,
	0.5, 0.5, 0,
	0.5, -0.5, 0,
}

type Game struct {
	Matrix     [][]*Cell
	Iterations int
}

type Cell struct {
	Status   int
	Drawable uint32
	X        int
	Y        int
}

func (c *Cell) NewCell(x, y, dim, status int) *Cell {
	points := make([]float32, len(square))
	copy(points, square)
	for i := 0; i < len(points); i++ {
		var position float32
		var size float32
		switch i % 3 {
		case 0:
			size = 1.0 / float32(dim)
			position = float32(x) * size
		case 1:
			size = 1.0 / float32(dim)
			position = float32(y) * size
		default:
			continue
		}

		if points[i] < 0 {
			points[i] = (position * 2) - 1
		} else {
			points[i] = ((position + size) * 2) - 1
		}
	}
	return &Cell{
		Drawable: graphics.MakeVAO(points),
		X:        x,
		Y:        y,
		Status:   status,
	}
}

func (c *Cell) MakeAlive() {
	c.Status = 1
}

func (c *Cell) MakeDead() {
	c.Status = 0
}

func (c *Cell) Draw() {
	if c.Status == 0 {
		return
	}

	gl.BindVertexArray(c.Drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(square)/3))
}

func (g *Game) GetAliveNbr(x int, y int) int {
	count := 0
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			if x-i < 0 || y-j < 0 || x-i > (len(g.Matrix)-1) || y-j > (len(g.Matrix)-1) {
				continue
			}
			if g.Matrix[x-i][y-j].Status == 1 {
				count++
			}
		}
	}
	return count
}

func (g *Game) InitMatrix() {
	g.Matrix[2][3].MakeAlive()
	g.Matrix[2][4].MakeAlive()
	g.Matrix[2][5].MakeAlive()
}

func (g *Game) PrintGame() {
	for i := 0; i < len(g.Matrix); i++ {
		for j := 0; j < len(g.Matrix); j++ {
			fmt.Print(g.Matrix[i][j].Status, "\t")
		}
		fmt.Print("\n")
	}
}

func (g *Game) UpdateGame() {
	UpdatedMatrix := g.Matrix
	for i := 0; i < len(g.Matrix); i++ {
		for j := 0; j < len(g.Matrix); j++ {
			nbrs := g.GetAliveNbr(i, j)
			if nbrs < 2 || nbrs > 3 {
				UpdatedMatrix[i][j].MakeDead()
			} else if nbrs == 3 {
				UpdatedMatrix[i][j].MakeAlive()
			} else if (nbrs == 2 || nbrs == 3) && (g.Matrix[i][j].Status == 1) {
				UpdatedMatrix[i][j].MakeAlive()
			}
		}
	}
	g.Matrix = UpdatedMatrix
}

func InitGame(dim int, initalProbability float32) *Game {
	nm := make([][]*Cell, dim)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for row := 0; row < dim; row++ {
		nm[row] = make([]*Cell, dim)

		for col := 0; col < dim; col++ {
			stat := 0
			if r.Float32() < initalProbability {
				stat = 1
			}
			nm[row][col] = nm[row][col].NewCell(row, col, dim, stat)
		}
	}
	newgame := Game{
		Matrix:     nm,
		Iterations: 0,
	}

	return &newgame
}
