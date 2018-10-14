package log

import (
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/tbase"
)

type Controller interface {
	// Should return false if not interested
	Open(Info) bool
	Close()

	Shuffle(Shuffle)
	Go(client.WithLobby)
	Start(client.WithDealer)
	Init(Init)
	Draw(WithOpponentAndInstance)
	Discard(WithOpponentAndInstance)
	Declare(Declare)
	Ryuukyoku(tbase.Ryuukyoku)
	Reach(client.Reach)
	Agari(tbase.Agari)
	Indicator(client.WithInstance)
	Disconnect(client.WithOpponent)

	client.UNController
}

type Shuffle struct {
	Seed string
	Ref  string
}

type Declare struct {
	client.WithOpponent
	Meld tbase.Meld
}

type WithOpponentAndInstance struct {
	client.WithOpponent
	client.WithInstance
}

type Init struct {
	tbase.Init
	Hands tbase.Hands
	// TODO: research
	Shuffle string
}
