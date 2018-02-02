package util

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/tbase"
)

type Exceptions string

func (this Exceptions) Check(b byte) bool {
	for _, v := range []byte(this) {
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
		tmp[k] = strconv.Itoa(int(v))
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

func ParseXML(input string) (ret parser.Nodes, err error) {
	// Dirty hack
	input = "<mjloggm>" + input + "</mjloggm>"
	d := xml.NewDecoder(strings.NewReader(input))
	d.Strict = false
	var root parser.Root
	err = stackerr.Wrap(d.Decode(&root))
	for k, v := range root.Nodes {
		if len(v.Attributes) == 0 {
			root.Nodes[k].Attributes = nil
		}
	}
	ret = root.Nodes
	return
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
		err = stackerr.Newf("Expected t to have 2 or 3 values, but got", len(lst))
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

func YakuString(y tbase.YakuResults) string {
	tmp := make([]string, len(y))
	for k, v := range y {
		tmp[k] = fmt.Sprintf("%d,%d", tbase.ReverseYakuMap[v.Key], v.Value)
	}
	return strings.Join(tmp, ",")
}

func YakumanString(y tbase.Yakumans) string {
	tmp := make([]string, len(y))
	for k, v := range y {
		tmp[k] = strconv.Itoa(tbase.ReverseYakumanMap[v])
	}
	return strings.Join(tmp, ",")
}

func FinalsString(ch tbase.ScoreChanges, floats bool) string {
	tmp := make([]string, len(ch))
	for k, v := range ch {
		diff := float64(v.Diff/100) / 10
		if floats {
			tmp[k] = fmt.Sprintf("%d,%.1f", v.Score/100, diff)
		} else {
			tmp[k] = fmt.Sprintf("%d,%d", v.Score/100, int(diff))
		}

	}
	return strings.Join(tmp, ",")
}
