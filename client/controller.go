package client

import (
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

func (s Suggest) Check(x Suggest) bool {
	return (s & x) == x
}

type UNController interface {
	UserList(UserList)
	Reconnect(Reconnect)
}

// This is how client looks from server's point of view
type Controller interface {
	UNController

	Drop(Drop)
	Take(Take)
	Reach(Reach)
	Declare(Declare)
	Init(Init)
	Reinit(Reinit)
	LogInfo(LogInfo)
	Go(Go)
	Hello(Hello)

	// TODO: research
	LobbyStats(LobbyStats)
	Agari(tbase.Agari)
	Indicator(WithInstance)
	// TODO: research add
	EndButton(EndButton)
	Furiten(Furiten)
	Ryuukyoku(tbase.Ryuukyoku)
	Rejoin(Rejoin)
	Disconnect(WithOpponent)
	Chat(Chat)
	// TODO: research
	Ranking(Ranking)
	Recover(Recover)
}

const DefaultRatingScale = "PF3=1.000000&PF4=1.000000&PF01C=0.582222&PF02C=0.501632&PF03C=0.414869&PF11C=0.823386&PF12C=0.709416&PF13C=0.586714&PF23C=0.378722&PF33C=0.535594&PF1C00=8.000000"
