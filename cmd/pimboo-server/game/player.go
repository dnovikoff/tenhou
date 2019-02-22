package game

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/tbase"
)

type Player struct {
	Hand    compact.Instances
	Melds   tbase.CalledList
	Discard compact.Instances
	Score   score.Money
	tempai  *tempai.TempaiResults
	furiten bool
	first   bool
}

func NewPlayer(x score.Money) *Player {
	return &Player{Score: x}
}

func (p *Player) Waits() compact.Tiles {
	return tempai.GetWaits(p.tempai)
}

func (p *Player) update() {
	p.tempai = tempai.Calculate(p.Hand, calc.Declared(p.Melds.Core()))

	checker := p.Discard.UniqueTiles()
	p.furiten = (checker & p.Waits()) != 0
}

func (p *Player) win(ctx *yaku.Context) *yaku.Result {
	ctx.IsFirstTake = p.first
	ctx.Rules = rules
	return yaku.Win(p.tempai, ctx, nil)
}

func (p *Player) drop(t tile.Instance) {
	p.first = false
	p.Hand.Remove(t)
	p.Discard.Set(t)
	p.update()
}

func (p *Player) take(t tile.Instance) bool {
	p.Hand.Set(t)
	return p.Waits().Check(t.Tile())
}

func (p *Player) Init(hand compact.Instances) {
	p.Hand = hand
	p.Melds = nil
	p.first = true
	p.Discard = compact.NewInstances()
	p.update()
}
