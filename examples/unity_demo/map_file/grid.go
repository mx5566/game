package map_file

import (
	"github.com/xiaonanln/goworld/engine/astar"
)

type Grid struct {
	Pos Position
	W   *Map
}

// PathNeighbors returns the neighbors of the tile, excluding blockers and
// tiles off the edge of the board.
func (t *Grid) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	for _, offset := range [][]int32{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
		x := t.Pos.X + offset[0]
		y := t.Pos.Y + offset[1]
		if n := t.W.Tile(x, y); n != nil {
			if !n.W.IsHasBlockGrid(x, y) {
				neighbors = append(neighbors, n)
			}
		}
	}
	return neighbors
}

// PathNeighborCost returns the movement cost of the directly neighboring tile.
func (t *Grid) PathNeighborCost(to astar.Pather) float64 {
	return KindCosts[KindPlain]
}

// PathEstimatedCost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *Grid) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Grid)
	absX := toT.Pos.X - t.Pos.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Pos.Y - t.Pos.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}
