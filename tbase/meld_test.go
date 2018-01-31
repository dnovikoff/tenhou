package tbase

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/tile"
)

func TestMeldChi(t *testing.T) {
	instance := tile.Instance(57)
	require.Equal(t, 1, Meld(37199).Extract(13, 14))
	assert.Equal(t, "6p", instance.StringOrNull())
	assert.Equal(t, "1001000101001111", fmt.Sprintf("%b", 37199))
	fmt.Printf(NewCoreMeld(37199).Instances().String())
	assert.EqualValues(t, meld.NewSeq(tile.Pin6, meld.OpenCopy(1), 2, 2), NewCoreMeld(37199))
}

func TestMeldsBase(t *testing.T) {
	base := func(x int) tile.Tile {
		return NewCoreMeld(Meld(x)).Base()
	}

	assert.Equal(t, tile.North, base(46633))
	assert.Equal(t, tile.White, base(47690))
	assert.Equal(t, tile.Green, base(49705))
}

func TestMelds(t *testing.T) {
	m := func(x int) meld.Meld {
		return NewCoreMeld(Meld(x))
	}
	assert.EqualValues(t, meld.NewPonOpened(tile.North, 2, 1, base.Right), m(46633))
	assert.EqualValues(t, meld.NewPonOpened(tile.White, 0, 2, base.Front), m(47690))
	assert.EqualValues(t, meld.NewPonOpened(tile.Green, 2, 1, base.Right), m(49705))
	assert.EqualValues(t, meld.NewKanUpgraded(tile.White, 0, 1, base.Front), m(47666))
	assert.EqualValues(t, meld.NewKanOpened(tile.Man2, 2, base.Left), m(1539))
	assert.EqualValues(t, meld.NewKanOpened(tile.East, 1, base.Front), m(27906))
	assert.EqualValues(t, meld.NewKanOpened(tile.Pin1, 0, base.Right), m(9217))

	assert.EqualValues(t, meld.NewPonOpened(tile.Sou4, 3, 2, base.Front), m(33354))
}

func TestMeldsReconvert(t *testing.T) {
	t2c := func(x Meld) meld.Meld {
		return NewCoreMeld(Meld(x))
	}
	c2t := func(x meld.Meld) Meld {
		return NewTenhouMeld(x)
	}
	reconvert := func(x Meld) Meld {
		return c2t(t2c(x))
	}

	// chi [6]78p
	assert.EqualValues(t, 37199, reconvert(37199))
	assert.EqualValues(t, 33354, reconvert(33354))

	// Closed green kan with copyId = 1
	assert.EqualValues(t, 33024, reconvert(33024))
	// Closed green kan with copyId = 0
	assert.EqualValues(t, 32768, reconvert(32768))
}

func TestMeldsReconvert2(t *testing.T) {
	t2c := func(x Meld) meld.Meld {
		return NewCoreMeld(Meld(x))
	}
	c2t := func(x meld.Meld) Meld {
		return NewTenhouMeld(x)
	}
	tst := func(i meld.Interface) bool {
		x := i.Meld()
		x2 := t2c(c2t(x))
		assert.Equal(t, x, x2)
		return x == x2
	}

	assert.True(t, tst(meld.NewPonOpened(tile.North, 2, 1, base.Right)))
	assert.True(t, tst(meld.NewPonOpened(tile.White, 0, 2, base.Front)))
	assert.True(t, tst(meld.NewPonOpened(tile.Green, 2, 1, base.Right)))
	assert.True(t, tst(meld.NewKanUpgraded(tile.White, 0, 1, base.Front)))
	assert.True(t, tst(meld.NewKanOpened(tile.Man2, 2, base.Left)))
	assert.True(t, tst(meld.NewKanOpened(tile.East, 1, base.Front)))
	assert.True(t, tst(meld.NewKanOpened(tile.Pin1, 0, base.Right)))
	assert.True(t, tst(meld.NewPonOpened(tile.Sou4, 3, 2, base.Front)))
}
