package log

import (
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/tbase"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
)

type Controller interface {
	// Should return false if not interested
	Open(*Info) bool
	Close()
	SetFloatFormat()

	Shuffle(seed, ref string)
	Go(t int, lobby int)
	Start(base.Opponent)
	Init(*Init)
	Draw(base.Opponent, tile.Instance)
	Discard(base.Opponent, tile.Instance)
	Declare(base.Opponent, tbase.Meld)
	Ryuukyoku(*tbase.Ryuukyoku)
	Reach(base.Opponent, int, []score.Money)
	Agari(*tbase.Agari)
	Indicator(tile.Instance)
	Disconnect(base.Opponent)

	client.UNController
}

type Init struct {
	tbase.Init
	Hands tbase.Hands
	// TODO: research
	Shuffle string
}
