package parser

import (
	"fmt"

	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/yaku"

	"github.com/dnovikoff/tenhou/tbase"
	"github.com/facebookgo/stackerr"
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

	machi := node.GetHai("machi")
	if len(machi) != 1 {
		err = stackerr.Newf("Expected machi have length of 1, but %v found", len(machi))
		return
	}
	agari.WinTile = machi[0]
	agari.FinalScores = node.GetFinalScores()

	ints = node.IntList("yakuman")
	if len(ints) > 0 {
		agari.TenhouYakumans = ints
		yakuman := make(tbase.Yakumans, len(ints))
		for k, v := range ints {
			yakuman[k] = tbase.YakumanMap[v]
		}
		agari.Yakumans = yakuman
	}

	ints = node.IntList("yaku")
	if len(ints) > 0 {
		agari.TenhouYakus = ints
		l := len(ints) / 2
		yaku := make([]tbase.YakuResult, l)
		for k := 0; k < l; k++ {
			key := ints[k*2]
			value := ints[k*2+1]
			yaku[k] = tbase.YakuResult{tbase.YakuMap[key], value}
		}
		agari.Yakus = yaku
	}

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
