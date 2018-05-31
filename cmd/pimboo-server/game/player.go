package game

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
)

func NewPlayer(x score.Money) *Player {
	return &Player{Score: x}
}

func (this *Player) update() {
	this.tempai = tempai.Calculate(this.Hand, calc.Melds(this.Melds)).Index()

	checker := this.Discard.UniqueTiles()
	this.furiten = (checker & this.tempai.Waits()) != 0
}

func (this *Player) win(ctx *yaku.Context) *yaku.YakuResult {
	ctx.IsFirstTake = this.first
	ctx.Rules = rules
	return yaku.Win(this.tempai, ctx)
}

func (this *Player) drop(t tile.Instance) {
	this.first = false
	this.Hand.Remove(t)
	this.Discard.Set(t)
	this.update()
}

func (this *Player) take(t tile.Instance) bool {
	this.Hand.Set(t)
	return this.tempai.Waits().Check(t.Tile())
}

func (this *Player) Init(hand compact.Instances) {
	this.Hand = hand
	this.Melds = nil
	this.first = true
	this.Discard = compact.NewInstances()
	this.update()
}
