package log

import (
	"encoding/json"
	"errors"
	"fmt"

	multierror "github.com/hashicorp/go-multierror"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/tbase"
)

func NewValidator(err *error) *AgariExtractor {
	return NewAgariExtractor(func(info *Info, agari *tbase.Agari, ctx *yaku.Context, w base.Wind, scoring score.Rules) {
		ValidateAgari(err, info, agari, ctx, w, scoring)
	})
}

type AgariReport struct {
	Log            *string        `json:"id,omitempty"`
	Round          int            `json:"round"`
	Wind           int            `json:"wind"`
	Score          score.Money    `json:"score"`
	Tile           tile.Instance  `json:"tile"`
	Hand           tile.Instances `json:"hand"`
	Melds          tbase.Melds    `json:"melds,omitempty"`
	DoraIndicators tile.Instances `json:"dora-indicators"`
	UraIndicators  tile.Instances `json:"ura-indicators,omitempty"`
	Yaku           tbase.Yakus    `json:"yaku,omitempty"`
	Yakuman        tbase.Yakumans `json:"yakuman,omitempty"`
}

func ValidateAgari(outError *error, info *Info, agari *tbase.Agari, ctx *yaku.Context, w base.Wind, scoring score.Rules) {
	comp := compact.NewInstances().Add(agari.Hand)
	comp.Remove(agari.WinTile)
	t := tempai.Calculate(comp, agari.Melds.Convert()).Index()
	yaku := yaku.Win(t, ctx)
	addError := func(format string, a ...interface{}) {
		id := fmt.Sprintf("%v", agari)
		(*outError) = multierror.Append((*outError), errors.New("Error at ["+id+"]: "+fmt.Sprintf(format, a...)))
	}

	totalExpected := agari.Changes[agari.Who].Diff

	if yaku == nil {
		addError("Expected win for hand %v + [%v] + %v with Score: %v. Round %v",
			agari.Hand,
			agari.Melds,
			agari.WinTile.Tile().String(),
			totalExpected,
			ctx.RoundWind,
		)
		return
	}
	var scoreFinal score.Score
	var baseScore score.Score
	if len(agari.Yakumans) == 0 {
		scoreFinal = score.GetScore(scoring, yaku.Sum(), yaku.Fus.Sum(), agari.Status.Honba)
		baseScore = score.GetScore(scoring, yaku.Sum(), yaku.Fus.Sum(), 0)
	} else {
		scoreFinal = score.GetYakumanScore(scoring, len(agari.Yakumans), agari.Status.Honba)
		baseScore = score.GetYakumanScore(scoring, len(agari.Yakumans), 0)
	}
	changes := scoreFinal.GetChanges(ctx.SelfWind, w, agari.Status.Sticks)
	total := changes.TotalWin()

	if total != totalExpected {
		addError("Money mismatch. Expected: %v, Calculated: %v (%v.%v). Debug %v + %v", totalExpected, total, yaku.Sum(), yaku.Fus.Sum(), yaku.Yaku, yaku.Bonuses)
	}

	report := &AgariReport{
		Round:          int(ctx.RoundWind),
		Wind:           int(ctx.SelfWind),
		Hand:           comp.Instances(),
		Melds:          agari.Melds,
		Tile:           ctx.Tile,
		Score:          baseScore.GetChanges(ctx.SelfWind, w, 0).TotalWin(),
		DoraIndicators: agari.DoraIndicators,
		UraIndicators:  agari.UraIndicators,
		Yaku:           agari.Yakus,
		Yakuman:        agari.Yakumans,
		Log:            &info.FullName,
	}

	bytes, _ := json.Marshal(report)
	fmt.Printf("%s\n", bytes)
}
