package tbase

import (
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
)

type Seed struct {
	RoundNumber int
	Honba       score.Honba
	Sticks      score.RiichiSticks
	Dice        [2]int
	Indicator   tile.Instance
}
