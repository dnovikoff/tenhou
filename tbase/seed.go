package tbase

import (
	"fmt"

	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/facebookgo/stackerr"
)

type Seed struct {
	RoundNumber int
	Honba       score.Honba
	Sticks      score.RiichiSticks
	Dice        [2]int
	Indicator   tile.Instance
}

func (this Seed) String() string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d",
		this.RoundNumber,
		this.Honba,
		this.Sticks,
		this.Dice[0],
		this.Dice[1],
		this.Indicator,
	)
}

func ParseSeed(in string) (ret Seed, err error) {
	seed := IntList(in)
	if len(seed) != 6 {
		err = stackerr.Newf("Expected 6 elements for seed")
		return
	}
	ret.RoundNumber = seed[0]
	ret.Honba = score.Honba(seed[1])
	ret.Sticks = score.RiichiSticks(seed[2])
	ret.Dice = [2]int{seed[3], seed[4]}
	ret.Indicator = tile.Instance(seed[5])
	return
}
