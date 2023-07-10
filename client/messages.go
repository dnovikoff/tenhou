package client

import (
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type WithLobby struct {
	LobbyNumber int
	LobbyType   int

	// TODO: parse deeper
	Title   *string
	Rule    *string
	Ranking *string
	CSRule  *string
}

type WithOpponent struct {
	Opponent tbase.Opponent
}

type WithDealer struct {
	Dealer tbase.Opponent
}

type WithInstance struct {
	Instance tile.Instance
}

type WithSuggest struct {
	Suggest Suggest
}

type Drop struct {
	WithOpponent
	WithInstance
	WithSuggest
	IsTsumogiri bool
}

type Take struct {
	WithOpponent
	WithInstance
	WithSuggest
}

type Reach struct {
	WithOpponent
	Score []score.Money
	Step  int
}

type Declare struct {
	WithOpponent
	WithSuggest
	Meld tbase.Meld
}

type LogInfo struct {
	WithDealer
	Hash string
}

type Go struct {
	WithLobby
	GpID string
}

type Hello struct {
	Name string
	Auth string

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

type LobbyStats struct {
	N string
	J string
	G string
}

type EndButton struct {
	WithLobby
	Add string
}

type Furiten struct {
	Furiten bool
}

type Rejoin struct {
	WithLobby
	Rejoin bool
}

type Chat struct {
	Name    string
	Message string
}

type Ranking struct {
	V2 string
}

type Recover struct {
	WithDealer
	Status  tbase.TableStatus
	Changes tbase.ScoreChanges
}

type Reconnect struct {
	WithOpponent
	Name string
}

type UserList struct {
	Users tbase.UserList
}
