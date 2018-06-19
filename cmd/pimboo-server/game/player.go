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

func (this *Player) Waits() compact.Tiles {
	return tempai.GetWaits(this.tempai)
}

func (this *Player) update() {
	this.tempai = tempai.Calculate(this.Hand, calc.Declared(this.Melds.Core()))

	checker := this.Discard.UniqueTiles()
	this.furiten = (checker & this.Waits()) != 0
}

func (this *Player) win(ctx *yaku.Context) *yaku.Result {
	ctx.IsFirstTake = this.first
	ctx.Rules = rules
	return yaku.Win(this.tempai, ctx, nil)
}

func (this *Player) drop(t tile.Instance) {
	this.first = false
	this.Hand.Remove(t)
	this.Discard.Set(t)
	this.update()
}

func (this *Player) take(t tile.Instance) bool {
	this.Hand.Set(t)
	return this.Waits().Check(t.Tile())
}

func (this *Player) Init(hand compact.Instances) {
	this.Hand = hand
	this.Melds = nil
	this.first = true
	this.Discard = compact.NewInstances()
	this.update()
}
