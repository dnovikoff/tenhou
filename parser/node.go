package parser

import (
	"fmt"
	"strconv"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

func (this *Node) GetOpponent(key string) base.Opponent {
	return base.Opponent(this.Int(key))
}

func (this *Node) GetWho() base.Opponent {
	return this.GetOpponent("who")
}

func (this *Node) Check(key string) bool {
	_, ok := this.Attributes[key]
	return ok
}

func (this *Node) GetDealer() base.Opponent {
	return this.GetOpponent("oya")
}

func (this *Node) GetScoreChanges() tbase.ScoreChanges {
	sc := this.IntList("sc")
	c := len(sc) / 2
	changes := make(tbase.ScoreChanges, c)
	for i := 0; i < c; i++ {
		cur := &changes[i]
		cur.Score = score.Money(sc[i*2]) * 100
		cur.Diff = score.Money(sc[i*2+1]) * 100
	}
	return changes
}

func (this *Node) GetHands() tbase.Hands {
	hands := make(tbase.Hands, 4)
	for k := 0; k < 4; k++ {
		hands[k] = this.GetHaiNum(k)
	}
	return hands
}

func (this *Node) GetTableStatus() (status tbase.TableStatus, err error) {
	ba := this.IntList("ba")
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

func (this *Node) ValidateUnused() error {
	unused := this.Keys()
	if len(unused) == 0 {
		return nil
	}
	return fmt.Errorf("Unused keys for node '%v': %v", this.Name, unused)
}

func (this *Node) Clone() *Node {
	x := *this
	attrs := make(map[string]string, len(x.Attributes))
	for k, v := range this.Attributes {
		attrs[k] = v
	}
	x.Attributes = attrs
	return &x
}

func (this *Node) Unused(name string) {
	delete(this.Attributes, name)
}

func (this *Node) String(name string) string {
	x := this.Attributes[name]
	this.Unused(name)
	return x
}

func (this *Node) StringList(name string) []string {
	return tbase.StringList(this.String(name))
}

func (this *Node) GetMeld() tbase.Meld {
	return tbase.Meld(this.Int("m"))
}

func (this *Node) Int(name string) int {
	i, _ := strconv.Atoi(this.String(name))
	return i
}

func (this *Node) IntPointer(name string) *int {
	x := this.String(name)
	if x == "" {
		return nil
	}
	i, _ := strconv.Atoi(this.String(name))
	return &i
}

func (this *Node) IntList(name string) []int {
	return tbase.IntList(this.String(name))
}

func (this *Node) Keys() []string {
	x := make([]string, 0, len(this.Attributes))
	for k := range this.Attributes {
		x = append(x, k)
	}
	return x
}

func (this *Node) FloatList(name string) []float64 {
	s := this.StringList(name)
	if s == nil {
		return nil
	}
	x := make([]float64, len(s))
	for k, v := range s {
		i, _ := strconv.ParseFloat(v, 64)
		x[k] = i
	}
	return x
}

func (this *Node) GetHaiNum(id int) tile.Instances {
	return this.GetHai("hai" + strconv.Itoa(id))
}

//TODO: change
func (this *Node) GetFinalScores() tbase.ScoreChanges {
	fl := this.FloatList("owari")
	ret := make(tbase.ScoreChanges, len(fl)/2)
	for k := range ret {
		ret[k] = tbase.ScoreChange{
			Score: score.Money(fl[k*2] * 100),
			Diff:  score.Money(fl[k*2+1] * 1000),
		}
	}
	return ret
}

func (this *Node) GetScores() tbase.Scores {
	return this.GetScoresByName("ten")
}

func (this *Node) GetScoresByName(name string) tbase.Scores {
	ten := this.IntList(name)
	p := make(tbase.Scores, len(ten))
	for k, v := range ten {
		p[k] = score.Money(v) * 100
	}
	return p
}

func (this *Node) GetTiles(name string) tile.Tiles {
	ints := this.IntList(name)
	if len(ints) == 0 {
		return nil
	}
	ret := make(tile.Tiles, len(ints))
	for k, v := range ints {
		ret[k] = tile.Tile(v)
	}
	return ret
}

func (this *Node) GetInstance(name string) tile.Instance {
	if !this.Check(name) {
		return tile.InstanceNull
	}
	return tile.Instance(this.Int(name))
}

func (this *Node) GetHai(name string) tile.Instances {
	ints := this.IntList(name)
	if len(ints) == 0 {
		return nil
	}
	ret := make(tile.Instances, len(ints))
	for k, v := range ints {
		ret[k] = tile.Instance(v)
	}
	return ret
}
