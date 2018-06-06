package base

import (
	"github.com/dnovikoff/tempai-core/tile"
)

//go:generate stringer -type=Wind
type Wind int

const (
	WindEast Wind = iota
	WindSouth
	WindWest
	WindNorth
	WindEnd
)

func (this Wind) tile() tile.Tile {
	return tile.Tile(this) + tile.East
}

func (this Wind) Opponent(other Wind) Opponent {
	diff := other - this
	if diff < 0 {
		diff += 4
	}
	return Self + Opponent(diff)
}

func (this Wind) CheckTile(t tile.Tile) bool {
	return this.tile() == t
}

func (this Wind) fix() Wind {
	if this < 0 {
		return (this + (4 * (this / -4)) + 4) % 4
	}
	return this % 4
}

func (this Wind) Advance(num int) Wind {
	return (this + Wind(num)).fix()
}
