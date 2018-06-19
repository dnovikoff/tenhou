package tbase

import (
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
)

type User struct {
	Num  int
	Name string
	Dan  int
	Rate float64
	Sex  Sex
	Rc   *int
}

type UserList []User

type ScoreChange struct {
	Score score.Money
	Diff  score.Money
}

type ScoreChanges []ScoreChange
type Scores []score.Money
type Hands []tile.Instances

type TableStatus struct {
	Honba  score.Honba
	Sticks score.RiichiSticks
}

type Score struct {
	Fu     yaku.FuPoints
	Total  score.Money
	Riichi score.RiichiSticks
}

type Agari struct {
	Who            Opponent
	From           Opponent
	Pao            *Opponent
	Status         TableStatus
	Score          Score
	FinalScores    ScoreChanges
	Changes        ScoreChanges
	Hand           tile.Instances
	DoraIndicators tile.Instances
	UraIndicators  tile.Instances
	WinTile        tile.Instance
	Yakus          Yakus
	Yakumans       Yakumans
	Melds          Melds
	// TODO: research
	Ratio string
}

type Ryuukyoku struct {
	DrawType     DrawType
	TableStatus  TableStatus
	ScoreChanges ScoreChanges
	Finals       ScoreChanges
	Hands        Hands
}

type Init struct {
	Seed
	Scores Scores
	Dealer Opponent
	Chip   []int
}
