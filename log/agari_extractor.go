package log

import (
	"fmt"
	"os"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/tbase"
)

type AgariExtractor struct {
	NullController
	Round     int
	Dealer    base.Opponent
	DoubleRon bool
	Info      *Info
	Defected  bool
	KanFlag   bool

	scoringRules score.Rules
	yakuRules    yaku.Rules

	Callback func(*Info, *tbase.Agari, *yaku.Context, base.Wind, score.Rules)
}

var _ Controller = &AgariExtractor{}

func NewAgariExtractor(x func(*Info, *tbase.Agari, *yaku.Context, base.Wind, score.Rules)) *AgariExtractor {
	return &AgariExtractor{
		Callback:     x,
		scoringRules: score.RulesTenhou(),
		yakuRules:    yaku.RulesTenhouRed(),
	}
}

func (this *AgariExtractor) Open(info Info) bool {
	this.Info = &info
	this.Defected = false
	rules := tbase.FlagDan1 | tbase.FlagDan2 | tbase.FlagHanchan | tbase.FlagOnline
	return info.Lobby == 0 && (info.Rules.Extract(tbase.FlagEnd-1) == rules)
}

func (this *AgariExtractor) Close() {
	if this.Defected {
		fmt.Fprintf(os.Stderr, "Defect log %v\n", this.Info.FullName)
	}
}

func (this *AgariExtractor) UserList(params client.UserList) {
	if len(params.Users) == 3 {
		this.Defected = true
	}
}

func (this *AgariExtractor) Init(params Init) {
	this.Dealer = params.Dealer
	this.Round = params.RoundNumber
	this.DoubleRon = false
}

func (this *AgariExtractor) Discard(WithOpponentAndInstance) {
	this.KanFlag = false
}

func (this *AgariExtractor) Declare(params Declare) {
	this.KanFlag = params.Meld.Convert().IsKan()
}

func (this *AgariExtractor) Agari(agari tbase.Agari) {
	if this.Defected {
		return
	}
	ctx := &yaku.Context{}

	ctx.DoraTiles = yaku.IndicatorsToDoraTiles(agari.DoraIndicators)
	ctx.UraTiles = yaku.IndicatorsToDoraTiles(agari.UraIndicators)
	yakus := agari.Yakus
	yakumans := agari.Yakumans

	ctx.IsDaburi = yakus.CheckCore(yaku.YakuDaburi)
	ctx.IsFirstTake = yakus.CheckCore(yaku.YakuRenhou) || yakumans.CheckCore(yaku.YakumanRenhou)
	ctx.IsRiichi = ctx.IsDaburi || yakus.CheckCore(yaku.YakuRiichi)
	ctx.IsIpatsu = yakus.CheckCore(yaku.YakuIppatsu)
	ctx.IsLastTile = yakus.CheckCore(yaku.YakuHaitei) || yakus.CheckCore(yaku.YakuHoutei)
	ctx.RoundWind = base.WindEast + base.Wind(this.Round/4)
	ctx.Tile = agari.WinTile
	ctx.SelfWind = base.WindEast.Advance(int(agari.Who - this.Dealer))
	otherWind := base.WindEast.Advance(int(agari.From - this.Dealer))
	ctx.IsTsumo = ctx.SelfWind == otherWind
	ctx.Rules = this.yakuRules
	ctx.IsChankan = (this.KanFlag && ctx.IsRon()) || yakus.CheckCore(yaku.YakuChankan)
	ctx.IsRinshan = (this.KanFlag && ctx.IsTsumo) || yakus.CheckCore(yaku.YakuRinshan)
	if this.DoubleRon {
		agari.Status.Sticks = 0
		agari.Status.Honba = 0
	}
	this.Callback(this.Info, &agari, ctx, otherWind, this.scoringRules)
	this.DoubleRon = true
}
