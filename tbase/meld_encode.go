package tbase

import (
	"fmt"

	"github.com/dnovikoff/tempai-core/tile"
)

func EncodeCalled(in *Called) Meld {
	switch in.Type {
	case Kan, ClosedKan:
		return encodeKan(in)
	case Pon, UpgrdedKan:
		return encodeSame(in)
	case Chi:
		return encodeChi(in)
	}
	panic(fmt.Sprintf("Unexpected type %v", in.Type))
}

func EncodeCalledList(in CalledList) Melds {
	x := make(Melds, len(in))
	for k, v := range in {
		x[k] = EncodeCalled(v)
	}
	return x
}

func encodeSame(in *Called) Meld {
	cp := in.Called.CopyID()
	if cp > in.Upgraded.CopyID() {
		cp--
	}
	base := in.Called.Tile() - tile.TileBegin
	m := Meld(base)*3 + Meld(cp)
	m = (m << 2)
	m = (m << 2) | Meld(in.Upgraded.CopyID())
	m = (m << 1)
	if in.Type == UpgrdedKan {
		m |= 1
	}
	m = (m << 1)
	if in.Type != UpgrdedKan {
		m |= 1
	}
	m = (m << 1)
	m = (m << 2) | Meld(in.Opponent)
	return m
}

func encodeKan(in *Called) Meld {
	m := Meld(in.Called.Tile()-tile.TileBegin)*4 + Meld(in.Called.CopyID())
	m = (m << 6)
	m = (m << 2) | Meld(in.Opponent)
	return m
}

func encodeChi(in *Called) Meld {
	base := in.Core.Tile() - tile.TileBegin
	calledIndex := int(in.Called.Tile() - in.Core.Tile())
	tiles := in.Tiles

	m := Meld(int((base/9)*7+base%9)*3 + calledIndex)
	m = (m << 1)
	m = (m << 2) | Meld(tiles[2].CopyID())
	m = (m << 2) | Meld(tiles[1].CopyID())
	m = (m << 2) | Meld(tiles[0].CopyID())
	m = (m << 1) | 1
	m = (m << 2) | Meld(in.Opponent)
	return m
}
