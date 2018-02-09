package client

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type Suggest int

const (
	SuggestNone Suggest = 0
)

const (
	SuggestKan      Suggest = 1 << iota
	SuggestPon              // 2
	SuggestChi              // 4
	SuggestRon              // 8
	SuggestTsumo            // 16
	SuggestRiichi           // 32
	SuggestDraw             // 64
	SuggetSanmaDora         // 128
)

func (this Suggest) Check(x Suggest) bool {
	return (this & x) == x
}

type HelloStats struct {
	RatingScale string
	// TODO: research
	PF4 string
	RR  string

	// TODO: change to date
	Expire     int
	ExpireDays int
	// nintei
	// TODO: check what is possible to display for client
	Message string
}

type Init struct {
	tbase.Init
	Hand tile.Instances
}

type Reinit struct {
	Init
	Melds   []tbase.Melds
	Discard []tile.Instances
	Riichi  []int
}

type UNController interface {
	UserList(list tbase.UserList)
	Reconnect(o base.Opponent, name string)
}

// This is how client looks from server's point of view
type Controller interface {
	Drop(o base.Opponent, t tile.Instance, isTsumogiri bool, suggest Suggest)
	Take(o base.Opponent, t tile.Instance, suggest Suggest)
	Reach(o base.Opponent, step int, score []score.Money)
	Declare(o base.Opponent, m tbase.Meld, s Suggest)
	Init(in Init)
	Reinit(in Reinit)
	LogInfo(dealer base.Opponent, hash string)
	Go(gpid string, gameType int, lobby int)
	Hello(name string, auth string, stats HelloStats)

	// TODO: research
	LobbyStats(n, j, g string)
	Agari(data *tbase.Agari)
	Indicator(i tile.Instance)
	// TODO: research add
	EndButton(lobby int, tp int, add string)
	Furiten(show bool)
	Ryuukyoku(*tbase.Ryuukyoku)
	Rejoin(lobbyNumber int, lobbyType int, rejoin bool)
	Disconnect(o base.Opponent)
	Chat(name string, text string)
	// TODO: research
	Ranking(v2 string)
	Recover(status tbase.TableStatus, dealer base.Opponent, ch tbase.ScoreChanges)

	UNController
}

const DefaultRatingScale = "PF3=1.000000&PF4=1.000000&PF01C=0.582222&PF02C=0.501632&PF03C=0.414869&PF11C=0.823386&PF12C=0.709416&PF13C=0.586714&PF23C=0.378722&PF33C=0.535594&PF1C00=8.000000"
