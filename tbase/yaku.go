package tbase

import (
	"sort"
	"strconv"
	"strings"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/yaku"
)

var YakuMap = map[int]yaku.Yaku{
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

var ReverseYakuMap = func() map[yaku.Yaku]int {
	ret := make(map[yaku.Yaku]int, len(YakuMap))
	for k, v := range YakuMap {
		ret[v] = k
	}
	return ret
}()

var YakumanMap = map[int]yaku.Yakuman{
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

var ReverseYakumanMap = func() map[yaku.Yakuman]int {
	ret := make(map[yaku.Yakuman]int, len(YakumanMap))
	for k, v := range YakumanMap {
		ret[v] = k
	}
	return ret
}()

func ParseYakuList(in string) (ret yaku.YakuSet, err error) {
	s := strings.Split(in, ",")
	ret = make(yaku.YakuSet, len(s))
	var i int
	var key yaku.Yaku
	var ok bool
	for k, v := range s {
		i, err = strconv.Atoi(v)
		if err != nil {
			ret = nil
			err = stackerr.Wrap(err)
			return
		}
		if k%2 == 0 {
			key, ok = YakuMap[i]
			if !ok {
				err = stackerr.Newf("Yaku with value %v not found in map", i)
				ret = nil
				return
			}
		} else {
			ret[key] = yaku.HanPoints(i)
		}
	}
	return
}

func ParseYakumanList(in string) (ret yaku.YakumanSet, err error) {
	s := strings.Split(in, ",")
	ret = make(yaku.YakumanSet, len(s))
	var i int
	var key yaku.Yakuman
	var ok bool
	for k, v := range s {
		i, err = strconv.Atoi(v)
		if err != nil {
			ret = nil
			err = stackerr.Wrap(err)
			return
		}
		if k%2 == 0 {
			key, ok = YakumanMap[i]
			if !ok {
				err = stackerr.Newf("Yakuman with value %v not found in map", i)
				ret = nil
				return
			}
		} else {
			ret[key] = i
		}
	}
	return
}

type YakuResult struct {
	Key   yaku.Yaku
	Value int
}

type YakuResults []YakuResult
type Yakumans []yaku.Yakuman

type YakuResultSet map[yaku.Yaku]int
type YakumanSet map[yaku.Yakuman]bool

var _ sort.Interface = YakuResults{}
var _ sort.Interface = Yakumans{}

func (a YakuResults) Len() int           { return len(a) }
func (a YakuResults) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a YakuResults) Less(i, j int) bool { return a[i].Key < a[j].Key }

func (a Yakumans) Len() int           { return len(a) }
func (a Yakumans) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Yakumans) Less(i, j int) bool { return a[i] < a[j] }

func YakusFromCore(in yaku.YakuSet) (ret YakuResults) {
	for k, v := range in {
		ret = append(ret, YakuResult{k, int(v)})
	}
	sort.Sort(ret)
	return
}

func YakumansFromCore(in yaku.YakumanSet) (ret Yakumans) {
	for k := range in {
		ret = append(ret, k)
	}
	sort.Sort(ret)
	return
}

func (this YakuResults) ToSet() YakuResultSet {
	result := make(YakuResultSet, len(this))
	for _, v := range this {
		result[v.Key] = v.Value
	}
	return result
}

func (this Yakumans) ToSet() YakumanSet {
	result := make(YakumanSet, len(this))
	for _, v := range this {
		result[v] = true
	}
	return result
}

func (this YakuResultSet) Check(v yaku.Yaku) bool {
	return this[v] != 0
}
