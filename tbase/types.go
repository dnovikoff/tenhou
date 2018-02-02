package tbase

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
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

type Agari struct {
	Who            base.Opponent
	From           base.Opponent
	Pao            *base.Opponent
	Status         TableStatus
	Scores         Scores
	FinalScores    ScoreChanges
	Changes        ScoreChanges
	Hand           tile.Instances
	DoraIndicators tile.Instances
	UraIndicators  tile.Instances
	WinTile        tile.Instance
	TenhouYakus    []int
	TenhouYakumans []int
	// TODO: switch to tenhou (do not convert)
	Yakumans Yakumans
	Yakus    YakuResults
	Melds    Melds
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
	Dealer base.Opponent
	Chip   []int
}
