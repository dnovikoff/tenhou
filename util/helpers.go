package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type Exceptions string

func (e Exceptions) Check(b byte) bool {
	for _, v := range []byte(e) {
		if v == b {
			return true
		}
	}
	return false
}

func EscapeWithExceptions(s string, except Exceptions) string {
	hexCount := len(s)
	t := make([]byte, len(s)+2*hexCount)
	j := 0

	for i := 0; i < len(s); i++ {
		c := s[i]
		if except.Check(c) {
			t[j] = c
			j += 1
			continue
		}
		t[j] = '%'
		t[j+1] = "0123456789ABCDEF"[c>>4]
		t[j+2] = "0123456789ABCDEF"[c&15]
		j += 3
	}
	return string(t[:j])
}

// Reduced version of url.escape golang function
func Escape(s string) string {
	return EscapeWithExceptions(s, "")
}

func InstanceString(in tile.Instances) string {
	tmp := make([]string, len(in))
	for k, v := range in {
		tmp[k] = strconv.Itoa(InstanceToTenhou(v))
	}
	return strings.Join(tmp, ",")
}

func MeldString(in tbase.Melds) string {
	tmp := make([]string, len(in))
	for k, v := range in {
		tmp[k] = strconv.Itoa(int(v))
	}
	return strings.Join(tmp, ",")
}

func ScoreString(in []score.Money) string {
	tmp := make([]string, len(in))
	for k, v := range in {
		tmp[k] = strconv.Itoa(int(v / 100))
	}
	return strings.Join(tmp, ",")
}

func IntsString(in []int) string {
	tmp := make([]string, len(in))
	for k, v := range in {
		tmp[k] = strconv.Itoa(v)
	}
	return strings.Join(tmp, ",")
}

func ParseJoinString(in string) (n, t int, rejoin bool, err error) {
	lst := tbase.StringList(in)
	switch len(lst) {
	case 2:
	case 3:
		if lst[2] == "r" {
			rejoin = true
		} else {
			err = stackerr.Newf("Expected letter to b 'r', but got '%v'", lst[2])
			return
		}
	default:
		err = stackerr.Newf("Expected t to have 2 or 3 values, but got '%v'", len(lst))
		return
	}
	n, err = strconv.Atoi(lst[0])
	if err != nil {
		err = stackerr.Wrap(err)
		return
	}
	t, err = strconv.Atoi(lst[1])
	if err != nil {
		err = stackerr.Wrap(err)
		return
	}
	return
}

func YakuString(y tbase.Yakus) string {
	return IntsString(y.Ints())
}

func YakumanString(y tbase.Yakumans) string {
	return IntsString(y.Ints())
}

func FinalsString(ch tbase.FinalScoreChanges) string {
	tmp := make([]string, len(ch))
	for k, v := range ch {
		if v.Diff.IsInt {
			tmp[k] = fmt.Sprintf("%d,%d", v.Score, int(v.Diff.Value))
		} else {
			tmp[k] = fmt.Sprintf("%d,%.1f", v.Score, v.Diff.Value)
		}

	}
	return strings.Join(tmp, ",")
}

func InstanceToTenhou(i tile.Instance) int {
	return int(i - tile.InstanceBegin)
}

func InstanceFromTenhou(i int) tile.Instance {
	return tile.Instance(i) + tile.InstanceBegin
}
