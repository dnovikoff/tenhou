package tbase

import (
	"fmt"
	"sort"

	"github.com/dnovikoff/tempai-core/yaku"
)

type Yaku int

type YakuRecord struct {
	Yaku  Yaku
	Value yaku.HanPoints
}

type Yakus []YakuRecord
type Yakuman int
type Yakumans []Yakuman

func (this Yakus) Ints() []int {
	ret := make([]int, len(this)*2)
	for k, v := range this {
		ret[k*2] = int(v.Yaku)
		ret[k*2+1] = int(v.Value)
	}
	return ret
}

func (this Yakus) CheckCore(x yaku.Yaku) bool {
	y, ok := ReverseYakuMap[x]
	if !ok {
		return false
	}
	return this.Check(y)
}

func (this Yakumans) CheckCore(x yaku.Yakuman) bool {
	y, ok := ReverseYakumanMap[x]
	if !ok {
		return false
	}
	return this.Check(y)
}

func (this Yakus) Check(x Yaku) bool {
	for _, v := range this {
		if x == v.Yaku {
			return true
		}
	}
	return false
}

func (this Yakumans) Check(x Yakuman) bool {
	for _, v := range this {
		if x == v {
			return true
		}
	}
	return false
}

func YakusFromInts(in []int) Yakus {
	if len(in) == 0 {
		return nil
	}
	x := make(Yakus, len(in)/2)
	for k := range x {
		x[k] = YakuRecord{Yaku(in[k*2]), yaku.HanPoints(in[k*2+1])}
	}
	return x
}

func YakumansFromInts(in []int) Yakumans {
	if len(in) == 0 {
		return nil
	}
	x := make(Yakumans, len(in))
	for k, v := range in {
		x[k] = Yakuman(v)
	}
	return x
}

func (this Yakumans) Ints() []int {
	ret := make([]int, len(this))
	for k, v := range this {
		ret[k] = int(v)
	}
	return ret
}

var YakuMap = map[Yaku]yaku.Yaku{
	0: yaku.YakuTsumo,
	1: yaku.YakuRiichi,
	2: yaku.YakuIppatsu,
	3: yaku.YakuChankan,
	4: yaku.YakuRinshan,
	5: yaku.YakuHaitei,
	6: yaku.YakuHoutei,
	7: yaku.YakuPinfu,
	8: yaku.YakuTanyao,
	9: yaku.YakuIppeiko,

	10: yaku.YakuTonSelf,
	11: yaku.YakuNanSelf,
	12: yaku.YakuSjaSelf,
	13: yaku.YakuPeiSelf,
	14: yaku.YakuTonRound,
	15: yaku.YakuNanRound,
	16: yaku.YakuSjaRound,
	17: yaku.YakuPeiRound,

	18: yaku.YakuHaku,
	19: yaku.YakuHatsu,
	20: yaku.YakuChun,

	21: yaku.YakuDaburi,
	22: yaku.YakuChiitoi,
	23: yaku.YakuChanta,
	24: yaku.YakuItsuu,
	25: yaku.YakuSanshoku,
	26: yaku.YakuSanshokuDoukou,
	27: yaku.YakuSankantsu,
	28: yaku.YakuToitoi,
	29: yaku.YakuSanankou,
	30: yaku.YakuShousangen,
	31: yaku.YakuHonrouto,
	32: yaku.YakuRyanpeikou,
	33: yaku.YakuJunchan,
	34: yaku.YakuHonitsu,
	35: yaku.YakuChinitsu,
	52: yaku.YakuDora,
	53: yaku.YakuUraDora,
	54: yaku.YakuAkaDora,
}

var ReverseYakuMap = func() map[yaku.Yaku]Yaku {
	ret := make(map[yaku.Yaku]Yaku, len(YakuMap))
	for k, v := range YakuMap {
		ret[v] = k
	}
	return ret
}()

var YakumanMap = map[Yakuman]yaku.Yakuman{
	36: yaku.YakumanRenhou,
	37: yaku.YakumanTenhou,
	38: yaku.YakumanChihou,
	39: yaku.YakumanDaisangen,
	40: yaku.YakumanSuuankou,
	41: yaku.YakumanSuuankouTanki,
	42: yaku.YakumanTsuiisou,
	43: yaku.YakumanRyuuiisou,
	44: yaku.YakumanChinrouto,
	45: yaku.YakumanChuurenpooto,
	46: yaku.YakumanChuurenpooto9,
	47: yaku.YakumanKokushi,
	48: yaku.YakumanKokushi13,
	49: yaku.YakumanDaisuushi,
	50: yaku.YakumanShousuushi,
	51: yaku.YakumanSuukantsu,
}

var ReverseYakumanMap = func() map[yaku.Yakuman]Yakuman {
	ret := make(map[yaku.Yakuman]Yakuman, len(YakumanMap))
	for k, v := range YakumanMap {
		ret[v] = k
	}
	return ret
}()

var _ sort.Interface = Yakus{}
var _ sort.Interface = Yakumans{}

func (a Yakus) Len() int           { return len(a) }
func (a Yakus) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Yakus) Less(i, j int) bool { return a[i].Yaku < a[j].Yaku }

func (a Yakumans) Len() int           { return len(a) }
func (a Yakumans) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Yakumans) Less(i, j int) bool { return a[i] < a[j] }

func YakusFromCore(in yaku.YakuSet) (ret Yakus, err error) {
	for k, v := range in {
		y, ok := ReverseYakuMap[k]
		if !ok {
			err = fmt.Errorf("No yaku '%v' in reverse map", k)
			ret = nil
			return
		}
		ret = append(ret, YakuRecord{y, v})
	}
	sort.Sort(ret)
	return
}

func YakumansFromCore(in yaku.YakumanSet) (ret Yakumans, err error) {
	for k := range in {
		y, ok := ReverseYakumanMap[k]
		if !ok {
			err = fmt.Errorf("No yakuman '%v' in reverse map", k)
			ret = nil
			return
		}
		ret = append(ret, y)
	}
	sort.Sort(ret)
	return
}

func (this Yakus) ToCore() (ret yaku.YakuSet, err error) {
	result := make(yaku.YakuSet, len(this))
	for _, v := range this {
		y, ok := YakuMap[v.Yaku]
		if !ok {
			err = fmt.Errorf("No yaku '%v' in map", v.Yaku)
			return
		}
		result[y] = v.Value
	}
	ret = result
	return
}

func (this Yakumans) ToCore() (ret yaku.YakumanSet, err error) {
	result := make(yaku.YakumanSet, len(this))
	for _, v := range this {
		y, ok := YakumanMap[v]
		if !ok {
			err = fmt.Errorf("No yakuman '%v' in map", v)
			return
		}
		result[y] = 1
	}
	ret = result
	return
}
