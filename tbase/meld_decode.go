package tbase

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/tile"
)

func decodeKan(in Meld) *Called {
	base := in.Extract(0, 8)
	t := getBase(base)
	x := &Called{
		Type:     ClosedKan,
		Opponent: in.Who(),
		Called:   t,
		Core:     calc.Kan(t.Tile()),
	}
	x.Tiles = compact.NewMask(0, t.Tile()).SetCount(4).Instances()
	if x.Opponent != Self {
		x.Type = Kan
		x.Core = calc.Open(x.Core)
	}
	return x
}

func decodePon(in Meld) *Called {
	isUpgraded := in.Extract(11, 12) == 1
	notInPonIdx := tile.CopyID(in.Extract(9, 11))

	baseAndCalled := in.Extract(0, 7)
	t := tile.Tile(baseAndCalled/3) + tile.TileBegin
	calledIndex := tile.CopyID(baseAndCalled % 3)
	if calledIndex >= notInPonIdx {
		calledIndex++
	}

	i := t.Instance(calledIndex)
	tiles := compact.NewMask(0, t).SetCount(4).UnsetCopyBit(notInPonIdx).Instances()
	x := &Called{
		Type:     Pon,
		Opponent: in.Who(),
		Tiles:    tiles,
		Called:   i,
		Upgraded: t.Instance(notInPonIdx),
	}

	if isUpgraded {
		x.Type = UpgrdedKan
		x.Core = calc.Open(calc.Kan(t))
	} else {
		x.Core = calc.Open(calc.Pon(t))
	}
	return x
}

func decodeChi(in Meld) *Called {
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
	base := kindCore.Tile(first)
	x := &Called{
		Type:     Chi,
		Opponent: in.Who(),
		Tiles: tile.Instances{
			(base + 0).Instance(tmp[0]),
			(base + 1).Instance(tmp[1]),
			(base + 2).Instance(tmp[2]),
		},
		Called: (base + tile.Tile(calledIndex)).Instance(tmp[calledIndex]),
		Core:   calc.Chi(base),
	}
	return x
}

func DecodeCalled(in Meld) *Called {
	switch in.Type() {
	case MeldTypeKan:
		return decodeKan(in)
	case MeldTypePon:
		return decodePon(in)
	case MeldTypeChi:
		return decodeChi(in)
	}
	return nil
}
