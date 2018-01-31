package log

import (
	"fmt"
	"os"

	"github.com/dnovikoff/tenhou/tbase"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
)

type AgariExtractor struct {
	NullController
	Round     int
	Dealer    base.Opponent
	DoubleRon bool
	Info      *Info
	Defected  bool
	KanFlag   bool

	Callback func(*Info, *tbase.Agari, *yaku.Context, base.Wind, *score.Rules)
}

var _ Controller = &AgariExtractor{}

func NewAgariExtractor(x func(*Info, *tbase.Agari, *yaku.Context, base.Wind, *score.Rules)) *AgariExtractor {
	return &AgariExtractor{Callback: x}
}

func (this *AgariExtractor) Open(info *Info) bool {
	this.Info = info
	this.Defected = false
	rules := tbase.FlagDan1 | tbase.FlagDan2 | tbase.FlagHanchan | tbase.FlagOnline
	return info.Lobby == 0 && (info.Rules.Extract(tbase.FlagEnd-1) == rules)
}

func (this *AgariExtractor) Close() {
	if this.Defected {
		fmt.Fprintf(os.Stderr, "Defect log %v\n", this.Info.FullName)
	}
}

func (this *AgariExtractor) UserList(users tbase.UserList) {
	if len(users) == 3 {
		this.Defected = true
	}
}

func (this *AgariExtractor) Init(info *Init) {
	this.Dealer = info.Dealer
	this.Round = info.RoundNumber
	this.DoubleRon = false
}

func (this *AgariExtractor) Discard(base.Opponent, tile.Instance) {
	this.KanFlag = false
}

func (this *AgariExtractor) Declare(_ base.Opponent, x tbase.Meld) {
	this.KanFlag = x.Convert().IsKan()
}

func (this *AgariExtractor) Agari(agari *tbase.Agari) {
	if this.Defected {
		return
	}
	ctx := &yaku.Context{}

	ctx.DoraTiles = yaku.IndicatorsToDoraTiles(agari.DoraIndicators)
	ctx.UraTiles = yaku.IndicatorsToDoraTiles(agari.UraIndicators)
	yakus := agari.Yakus.ToSet()
	yakumans := agari.Yakumans.ToSet()
	ctx.IsDaburi = yakus.Check(yaku.YakuDaburi)
	ctx.IsFirstTake = yakus.Check(yaku.YakuRenhou) || yakumans[yaku.YakumanRenhou]
	ctx.IsRiichi = ctx.IsDaburi || yakus.Check(yaku.YakuRiichi)
	ctx.IsIpatsu = yakus.Check(yaku.YakuIppatsu)
	ctx.IsLastTile = yakus.Check(yaku.YakuHaitei) || yakus.Check(yaku.YakuHoutei)
	ctx.RoundWind = base.WindEast + base.Wind(this.Round/4)
	ctx.Tile = agari.WinTile
	ctx.SelfWind = base.WindEast.Advance(int(agari.Who - this.Dealer))
	otherWind := base.WindEast.Advance(int(agari.From - this.Dealer))
	ctx.IsTsumo = ctx.SelfWind == otherWind
	ctx.Rules = &yaku.RulesTenhouRed
	ctx.IsChankan = (this.KanFlag && ctx.IsRon()) || yakus.Check(yaku.YakuChankan)
	ctx.IsRinshan = (this.KanFlag && ctx.IsTsumo) || yakus.Check(yaku.YakuRinshan)
	if this.DoubleRon {
		agari.Status.Sticks = 0
		agari.Status.Honba = 0
	}
	this.Callback(this.Info, agari, ctx, otherWind, &score.RulesTenhou)
	this.DoubleRon = true
}
