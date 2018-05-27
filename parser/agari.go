package parser

import (
	"fmt"

	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/tbase"
)

func ParseAgari(node *Node) (result *tbase.Agari, err error) {
	agari := &tbase.Agari{}
	agari.Who = node.GetWho()
	agari.Status, err = node.GetTableStatus()
	if err != nil {
		return
	}
	agari.From = node.GetOpponent("fromWho")
	if node.Check("paoWho") {
		x := node.GetOpponent("paoWho")
		agari.Pao = &x
	}
	ints := node.IntList("ten")
	if len(ints) != 3 {
		err = fmt.Errorf("ten length for agari should be 3 != %v", len(ints))
		return
	}
	agari.Score = tbase.Score{
		yaku.FuPoints(ints[0]),
		score.Money(ints[1]),
		score.RiichiSticks(ints[2])}
	agari.Changes = node.GetScoreChanges()
	agari.Hand = node.GetHai("hai")
	agari.DoraIndicators = node.GetHai("doraHai")
	agari.UraIndicators = node.GetHai("doraHaiUra")

	agari.WinTile = node.GetInstance("machi")
	agari.FinalScores = node.GetFinalScores()

	agari.Yakumans = tbase.YakumansFromInts(node.IntList("yakuman"))
	agari.Yakus = tbase.YakusFromInts(node.IntList("yaku"))

	ints = node.IntList("m")
	if len(ints) > 0 {
		melds := make(tbase.Melds, len(ints))
		for k, v := range ints {
			melds[k] = tbase.Meld(v)
		}
		agari.Melds = melds
	}
	agari.Ratio = node.String("ratio")
	result = agari
	return
}
