package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

func (node *Node) GetOpponent(key string) tbase.Opponent {
	return tbase.Opponent(node.Int(key))
}

func (node *Node) GetWho() tbase.Opponent {
	return node.GetOpponent("who")
}

func (node *Node) Check(key string) bool {
	_, ok := node.Attributes[key]
	return ok
}

func (node *Node) GetDealer() tbase.Opponent {
	return node.GetOpponent("oya")
}

func (node *Node) GetScoreChanges() tbase.ScoreChanges {
	sc := node.IntList("sc")
	c := len(sc) / 2
	changes := make(tbase.ScoreChanges, c)
	for i := 0; i < c; i++ {
		cur := &changes[i]
		cur.Score = sc[i*2]
		cur.Diff = sc[i*2+1]
	}
	return changes
}

func (node *Node) GetHands() tbase.Hands {
	hands := make(tbase.Hands, 4)
	for k := 0; k < 4; k++ {
		hands[k] = node.GetHaiNum(k)
	}
	return hands
}

func (node *Node) GetTableStatus() (status tbase.TableStatus, err error) {
	ba := node.IntList("ba")
	if len(ba) != 2 {
		err = stackerr.Newf("'ba' element should contain exactly two elements. Result: '%v'", ba)
		return
	}
	status.Honba = score.Honba(ba[0])
	status.Sticks = score.RiichiSticks(ba[1])
	return
}

type Node struct {
	Name       string
	Attributes map[string]string
}

func (node *Node) ValidateUnused() error {
	unused := node.Keys()
	if len(unused) == 0 {
		return nil
	}
	return fmt.Errorf("Unused keys for node '%v': %v", node.Name, unused)
}

func (node *Node) Clone() *Node {
	x := *node
	attrs := make(map[string]string, len(x.Attributes))
	for k, v := range node.Attributes {
		attrs[k] = v
	}
	x.Attributes = attrs
	return &x
}

func (node *Node) Unused(name string) {
	delete(node.Attributes, name)
}

func (node *Node) String(name string) string {
	x := node.Attributes[name]
	node.Unused(name)
	return x
}

func (node *Node) PString(name string) *string {
	if !node.Check(name) {
		return nil
	}
	x := node.String(name)
	return &x
}

func (node *Node) StringList(name string) []string {
	return tbase.StringList(node.String(name))
}

func (node *Node) GetMeld() tbase.Meld {
	return tbase.Meld(node.Int("m"))
}

func (node *Node) Int(name string) int {
	i, _ := strconv.Atoi(node.String(name))
	return i
}

func (node *Node) IntPointer(name string) *int {
	x := node.String(name)
	if x == "" {
		return nil
	}
	i, _ := strconv.Atoi(node.String(name))
	return &i
}

func (node *Node) IntList(name string) []int {
	return tbase.IntList(node.String(name))
}

func (node *Node) Keys() []string {
	x := make([]string, 0, len(node.Attributes))
	for k := range node.Attributes {
		x = append(x, k)
	}
	return x
}

func (node *Node) FloatList(name string) []tbase.Float {
	s := node.StringList(name)
	if s == nil {
		return nil
	}
	x := make([]tbase.Float, len(s))
	for k, v := range s {
		i, _ := strconv.ParseFloat(v, 64)
		dotIndex := strings.Index(v, ".")
		x[k].Value = i
		if dotIndex == -1 {
			x[k].IsInt = true
		}
	}
	return x
}

func (node *Node) GetHaiNum(id int) tile.Instances {
	return node.GetHai("hai" + strconv.Itoa(id))
}

func (node *Node) GetFinalScores() tbase.FinalScoreChanges {
	fl := node.FloatList("owari")
	ret := make(tbase.FinalScoreChanges, len(fl)/2)
	for k := range ret {
		ret[k] = tbase.FinalScoreChange{
			Score: int(fl[k*2].Value),
			Diff:  fl[k*2+1],
		}
	}
	return ret
}

func (node *Node) GetScores() tbase.Scores {
	return node.GetScoresByName("ten")
}

func (node *Node) GetScoresByName(name string) tbase.Scores {
	ten := node.IntList(name)
	p := make(tbase.Scores, len(ten))
	for k, v := range ten {
		p[k] = score.Money(v) * 100
	}
	return p
}

func (node *Node) GetTiles(name string) tile.Tiles {
	ints := node.IntList(name)
	if len(ints) == 0 {
		return nil
	}
	ret := make(tile.Tiles, len(ints))
	for k, v := range ints {
		ret[k] = tile.Tile(v)
	}
	return ret
}

func (node *Node) GetInstance(name string) tile.Instance {
	if !node.Check(name) {
		return tile.InstanceNull
	}
	return util.InstanceFromTenhou(node.Int(name))
}

func (node *Node) GetHai(name string) tile.Instances {
	ints := node.IntList(name)
	if len(ints) == 0 {
		return nil
	}
	ret := make(tile.Instances, len(ints))
	for k, v := range ints {
		ret[k] = util.InstanceFromTenhou(v)
	}
	return ret
}
