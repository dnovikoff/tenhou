package tbase

import (
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
)

type Float struct {
	Value float64
	IsInt bool
}

type UserList struct {
	Names []string
	Dan   []int
	Rate  []Float
	Sex   []Sex
	RC    []int
	Gold  []int
}

type ScoreChange struct {
	Score int
	Diff  int
}
type ScoreChanges []ScoreChange

func (sc *ScoreChange) ScoreMoney() score.Money {
	return score.Money(sc.Score * 100)
}

func (sc *ScoreChange) DiffMoney() score.Money {
	return score.Money(sc.Diff * 100)
}

type FinalScoreChange struct {
	Score int
	Diff  Float
}
type FinalScoreChanges []FinalScoreChange

func (sc *FinalScoreChange) ScoreMoney() score.Money {
	return score.Money(sc.Score * 100)
}

func (sc *FinalScoreChange) DiffMoney() score.Money {
	return score.Money(sc.Diff.Value * 1000)
}

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
	FinalScores    FinalScoreChanges
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
	Chips []int
}

type Ryuukyoku struct {
	DrawType     DrawType
	TableStatus  TableStatus
	ScoreChanges ScoreChanges
	Finals       FinalScoreChanges
	Hands        Hands
	Ratio        []int
}

type Init struct {
	Seed
	Scores Scores
	Dealer Opponent
	Chip   []int
}
