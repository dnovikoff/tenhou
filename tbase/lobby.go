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

func (this LobbyRules) Extract(mask LobbyRules) LobbyRules {
	return this & mask
}

func (this LobbyRules) Check(f LobbyRules) bool {
	return (f & this) == f
}

func (this LobbyRules) String() string {
	return fmt.Sprintf("%04x", int(this))
}

func (this LobbyRules) DebugString() string {
	buf := &bytes.Buffer{}
	w := func(r LobbyRules, t string, f string) {
		if this.Check(r) {
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
	dan := this.Extract(MaskDanAll)
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
