package tbase

import (
	"bytes"
	"fmt"
)

type LobbyRules int

const (
	FlagOnline LobbyRules = 1 << iota
	FlagNoAkkas
	FlagNoKuitan
	FlagHanchan
	Flag3Man
	FlagDan1
	FlagFast
	FlagDan2
	FlagEnd
)

const (
	MaskDanAll     LobbyRules = FlagDan1 | FlagDan2
	MaskDanKu      LobbyRules = 0
	MaskDan1                  = FlagDan1
	MaskDan2                  = FlagDan2
	MaskDanPhoenix            = MaskDan1 | MaskDan2

	RulesDzjanso LobbyRules = 0x0841
)

func (r LobbyRules) Extract(mask LobbyRules) LobbyRules {
	return r & mask
}

func (r LobbyRules) Check(f LobbyRules) bool {
	return (f & r) == f
}

func (r LobbyRules) String() string {
	return fmt.Sprintf("%04x", int(r))
}

func (r LobbyRules) DebugString() string {
	buf := &bytes.Buffer{}
	w := func(x LobbyRules, t string, f string) {
		if r.Check(x) {
			fmt.Fprintf(buf, t)
		} else {
			fmt.Fprintf(buf, f)
		}
	}
	w(FlagNoAkkas, "a", "A")
	w(FlagNoKuitan, "k", "K")
	w(FlagHanchan, "H", "T")
	w(Flag3Man, "3", "4")
	w(FlagFast, "F", "f")
	dan := r.Extract(MaskDanAll)
	switch dan {
	case MaskDanKu:
		fmt.Fprintf(buf, "0")
	case MaskDan1:
		fmt.Fprintf(buf, "D")
	case MaskDan2:
		fmt.Fprintf(buf, "U")
	case MaskDanPhoenix:
		fmt.Fprintf(buf, "X")
	default:
		fmt.Fprintf(buf, "[Strange value %b]", dan)
	}
	return string(buf.Bytes())
}
