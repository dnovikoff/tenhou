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

func (this Meld) Extract(first, last int) int {
	shift := uint(16 - last)
	mask := 0
	for k := first; k < last; k++ {
		mask = mask*2 + 1
	}
	val := int(this>>shift) & mask
	return val
}

func (this Meld) Decode() *Called {
	return DecodeCalled(this)
}

func (this Melds) Decode() CalledList {
	x := make(CalledList, len(this))
	for k, v := range this {
		x[k] = v.Decode()
	}
	return x
}

func (this Meld) Type() MeldType {
	if this.Extract(13, 14) == 1 {
		return MeldTypeChi
	}
	if this.Extract(8, 13) == 0 {
		return MeldTypeKan
	}
	return MeldTypePon
}

func (this Meld) Who() Opponent {
	return Opponent(this.Extract(14, 16))
}

func (this Meld) Instance(f, l int) tile.Instance {
	return tile.Instance(this.Extract(f, l))
}

func getBase(base int) tile.Instance {
	t := tile.Tile(base/4) + tile.TileBegin
	copyId := tile.CopyID(base % 4)
	return t.Instance(copyId)
}
