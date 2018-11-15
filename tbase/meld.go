package tbase

import (
	"github.com/dnovikoff/tempai-core/tile"
)

type Meld int
type Melds []Meld
type MeldKan Meld
type MeldPon Meld
type MeldChi Meld

type MeldType int

const (
	MeldTypeBad MeldType = iota
	MeldTypeChi
	MeldTypePon
	MeldTypeKan
)

func (m Meld) Extract(first, last int) int {
	shift := uint(16 - last)
	mask := 0
	for k := first; k < last; k++ {
		mask = mask*2 + 1
	}
	val := int(m>>shift) & mask
	return val
}

func (m Meld) Decode() *Called {
	return DecodeCalled(m)
}

func (m Melds) Decode() CalledList {
	x := make(CalledList, len(m))
	for k, v := range m {
		x[k] = v.Decode()
	}
	return x
}

func (m Meld) Type() MeldType {
	if m.Extract(13, 14) == 1 {
		return MeldTypeChi
	}
	if m.Extract(8, 13) == 0 {
		return MeldTypeKan
	}
	return MeldTypePon
}

func (m Meld) Who() Opponent {
	return Opponent(m.Extract(14, 16))
}

func (m Meld) Instance(f, l int) tile.Instance {
	return tile.Instance(m.Extract(f, l))
}

func getBase(base int) tile.Instance {
	t := tile.Tile(base/4) + tile.TileBegin
	copyId := tile.CopyID(base % 4)
	return t.Instance(copyId)
}
