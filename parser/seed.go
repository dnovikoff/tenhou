package parser

import (
	"fmt"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

func ParseSeed(in string) (ret tbase.Seed, err error) {
	seed := tbase.IntList(in)
	if len(seed) != 6 {
		err = stackerr.Newf("Expected 6 elements for seed")
		return
	}
	ret.RoundNumber = seed[0]
	ret.Honba = score.Honba(seed[1])
	ret.Sticks = score.RiichiSticks(seed[2])
	ret.Dice = [2]int{seed[3], seed[4]}
	ret.Indicator = util.InstanceFromTenhou(seed[5])
	return
}

func SeedString(seed *tbase.Seed) string {
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d",
		seed.RoundNumber,
		seed.Honba,
		seed.Sticks,
		seed.Dice[0],
		seed.Dice[1],
		util.InstanceToTenhou(seed.Indicator),
	)
}
