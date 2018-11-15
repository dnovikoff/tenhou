package log

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/tbase"
)

type AgariExtractor struct {
	NullController
	Round     int
	Dealer    tbase.Opponent
	DoubleRon bool
	Info      *Info
	KanFlag   bool
	IsSanma   bool

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

func (e *AgariExtractor) Open(info Info) bool {
	e.IsSanma = true
	e.Info = &info
	rules := tbase.FlagDan1 | tbase.FlagDan2 | tbase.FlagHanchan | tbase.FlagOnline
	return info.Lobby == 0 && (info.Rules.Extract(tbase.FlagEnd-1) == rules)
}

func (e *AgariExtractor) Close() {
}

func (e *AgariExtractor) UserList(params client.UserList) {
	e.IsSanma = len(params.Users.Names) == 3
}

func (e *AgariExtractor) Init(params Init) {
	e.Dealer = params.Dealer
	e.Round = params.RoundNumber
	e.DoubleRon = false
}

func (e *AgariExtractor) Discard(WithOpponentAndInstance) {
	e.KanFlag = false
}

func (e *AgariExtractor) Declare(params Declare) {
	e.KanFlag = params.Meld.Decode().IsKan()
}

func (e *AgariExtractor) Agari(agari tbase.Agari) {
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
	ctx.RoundWind = base.WindEast + base.Wind(e.Round/4)
	ctx.Tile = agari.WinTile
	ctx.SelfWind = base.WindEast.Advance(int(agari.Who - e.Dealer))
	otherWind := base.WindEast.Advance(int(agari.From - e.Dealer))
	ctx.IsTsumo = ctx.SelfWind == otherWind
	ctx.Rules = e.yakuRules
	ctx.IsChankan = (e.KanFlag && !ctx.IsTsumo) || yakus.CheckCore(yaku.YakuChankan)
	ctx.IsRinshan = (e.KanFlag && ctx.IsTsumo) || yakus.CheckCore(yaku.YakuRinshan)
	if e.DoubleRon {
		agari.Status.Sticks = 0
		agari.Status.Honba = 0
	}
	e.Callback(e.Info, &agari, ctx, otherWind, e.scoringRules)
	e.DoubleRon = true
}
