package tbase

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/meld"
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

func (this Meld) Convert() meld.Meld {
	return NewCoreMeld(this)
}

func (this Melds) Convert() meld.Melds {
	x := make(meld.Melds, len(this))
	for k, v := range this {
		x[k] = v.Convert()
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

func (this Meld) Who() base.Opponent {
	return base.Opponent(this.Extract(14, 16))
}

func (this Meld) Instance(f, l int) tile.Instance {
	return tile.Instance(this.Extract(f, l))
}

func getBase(base int) tile.Instance {
	t := tile.Tile(base/4) + tile.TileBegin
	copyId := tile.CopyID(base % 4)
	return t.Instance(copyId)
}

func newCoreMeldKan(in Meld) meld.Meld {
	base := in.Extract(0, 8)
	t := getBase(base)
	kan := meld.NewKan(t)
	if in.Who() == 0 {
		return kan.Meld()
	}
	return meld.NewKanOpened(t, in.Who()).Meld()
}

func newCoreMeldPon(in Meld) meld.Meld {
	isUpgraded := in.Extract(11, 12) == 1
	notInPonIdx := tile.CopyID(in.Extract(9, 11))

	baseAndCalled := in.Extract(0, 7)
	t := tile.Tile(baseAndCalled/3) + tile.TileBegin
	calledIndex := tile.CopyID(baseAndCalled % 3)
	if calledIndex >= notInPonIdx {
		calledIndex++
	}

	i := t.Instance(calledIndex)
	if isUpgraded {
		return meld.NewKanUpgraded(i, notInPonIdx, in.Who()).Meld()
	}
	return meld.NewPonOpened(i, notInPonIdx, in.Who()).Meld()
}

func newCoreMeldChi(in Meld) meld.Meld {
	pattern := in.Extract(0, 6)
	calledIndex := pattern % 3
	pattern /= 3
	kind := pattern / 7
	first := pattern%7 + 1
	var kindCore tile.Type
	switch kind {
	case 0:
		kindCore = tile.TypeMan
	case 1:
		kindCore = tile.TypePin
	case 2:
		kindCore = tile.TypeSou
	}
	copyAt := func(idx int) tile.CopyID {
		start := 11 - idx*2
		copyId := in.Extract(start, start+2)
		return tile.CopyID(copyId)
	}
	var tmp [3]tile.CopyID
	for k := 0; k < 3; k++ {
		tmp[k] = copyAt(k)
	}
	tmp[calledIndex] = meld.OpenCopy(tmp[calledIndex])
	base := kindCore.Tile(first)
	return meld.NewSeq(base, tmp[0], tmp[1], tmp[2]).Meld()
}

func newTenhouSeqMeld(in meld.Seq) Meld {
	calledIndex := in.OpenedIndex() - 1
	base := in.Base() - tile.TileBegin
	tiles := in.Instances()

	m := Meld(int((base/9)*7+base%9)*3 + calledIndex)
	m = (m << 1)
	m = (m << 2) | Meld(tiles[2].CopyID())
	m = (m << 2) | Meld(tiles[1].CopyID())
	m = (m << 2) | Meld(tiles[0].CopyID())
	m = (m << 1) | 1
	m = (m << 2) | Meld(in.Opponent())
	return m
}

func newTenhouSameMeld(in meld.Same) Meld {
	cp := in.OpenedCopy()
	if cp > in.NotInPonCopy() {
		cp--
	}
	m := Meld(in.Base()-tile.TileBegin)*3 + Meld(cp)
	m = (m << 2)
	m = (m << 2) | Meld(in.NotInPonCopy())
	m = (m << 1)
	if in.IsUpgraded() {
		m |= 1
	}
	m = (m << 1)
	if !in.IsUpgraded() {
		m |= 1
	}
	m = (m << 1)
	m = (m << 2) | Meld(in.Opponent())
	return m
}

func newTenhouKanMeld(in meld.Same) Meld {
	m := Meld(in.Base()-tile.TileBegin)*4 + Meld(in.OpenedCopy())
	m = (m << 6)
	m = (m << 2) | Meld(in.Opponent())
	return m
}

func NewTenhouMeld(in meld.Meld) Meld {
	switch in.Type() {
	case meld.TypeSame:
		same := meld.Same(in)
		if same.IsKan() && !same.IsUpgraded() {
			return newTenhouKanMeld(same)
		}
		return newTenhouSameMeld(same)
	case meld.TypeSeq:
		return newTenhouSeqMeld(meld.Seq(in))
	}
	return 0
}

func NewTenhouMelds(in meld.Melds) Melds {
	x := make(Melds, len(in))
	for k, v := range in {
		x[k] = NewTenhouMeld(v)
	}
	return x
}

func NewCoreMeld(in Meld) meld.Meld {
	switch in.Type() {
	case MeldTypeKan:
		return newCoreMeldKan(in)
	case MeldTypePon:
		return newCoreMeldPon(in)
	case MeldTypeChi:
		return newCoreMeldChi(in)
	}
	return 0
}
