package tbase

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

// func TestMeldChi(t *testing.T) {
// 	instance := tile.Instance(57)
// 	require.Equal(t, 1, Meld(37199).Extract(13, 14))
// 	assert.Equal(t, "6p", instance.Tile().String())
// 	assert.Equal(t, "1001000101001111", fmt.Sprintf("%b", 37199))
// 	t.Log(NewCoreMeld(37199).Instances().String())
// 	assert.EqualValues(t, meld.NewSeq(tile.Pin6, meld.OpenCopy(1), 2, 2), NewCoreMeld(37199))
// }

func TestMeldsBase(t *testing.T) {
	base := func(x int) tile.Tile {
		return DecodeCalled(Meld(x)).Called.Tile()
	}

	assert.Equal(t, tile.North, base(46633))
	assert.Equal(t, tile.White, base(47690))
	assert.Equal(t, tile.Green, base(49705))
}

func TestMelds(t *testing.T) {
	m := func(x int) string {
		return DecodeCalled(Meld(x)).Core.Tiles().String()
	}
	assert.EqualValues(t, "444z", m(46633))
	// assert.EqualValues(t, meld.NewPonOpened(tile.White.Instance(0), 2, base.Front), m(47690))
	// assert.EqualValues(t, meld.NewPonOpened(tile.Green.Instance(2), 1, base.Right), m(49705))
	// assert.EqualValues(t, meld.NewKanUpgraded(tile.White.Instance(0), 1, base.Front), m(47666))
	// assert.EqualValues(t, meld.NewKanOpened(tile.Man2.Instance(2), base.Left), m(1539))
	// assert.EqualValues(t, meld.NewKanOpened(tile.East.Instance(1), base.Front), m(27906))
	// assert.EqualValues(t, meld.NewKanOpened(tile.Pin1.Instance(0), base.Right), m(9217))

	// assert.EqualValues(t, meld.NewPonOpened(tile.Sou4.Instance(3), 2, base.Front), m(33354))
}

func TestMeldsReconvert(t *testing.T) {
	t2c := func(x Meld) *Called {
		return DecodeCalled(Meld(x))
	}
	c2t := func(x *Called) Meld {
		return EncodeCalled(x)
	}
	reconvert := func(x Meld) Meld {
		t.Log(t2c(x).Tiles)
		t.Log(t2c(c2t(t2c(x))).Tiles)
		return c2t(t2c(x))
	}
	for _, v := range []Meld{
		// chi [6]78p
		37199,
		33354,
		// Closed green kan with copyId = 1
		33024,
		// Closed green kan with copyId = 0
		32768,
	} {
		t.Run(strconv.Itoa(int(v)), func(t *testing.T) {
			assert.Equal(t, v, reconvert(v))
		})
	}
}

// func TestMeldsReconvert2(t *testing.T) {
// 	t2c := func(x Meld) *Called {
// 		return NewCoreMeld(Meld(x))
// 	}
// 	c2t := func(x *Called) Meld {
// 		return NewTenhouMeld(x)
// 	}
// 	tst := func(i meld.Interface) bool {
// 		x := i.Meld()
// 		x2 := t2c(c2t(x))
// 		assert.Equal(t, x, x2)
// 		return x == x2
// 	}

// 	assert.True(t, tst(meld.NewPonOpened(tile.North.Instance(2), 1, base.Right)))
// 	assert.True(t, tst(meld.NewPonOpened(tile.White.Instance(0), 2, base.Front)))
// 	assert.True(t, tst(meld.NewPonOpened(tile.Green.Instance(2), 1, base.Right)))
// 	assert.True(t, tst(meld.NewKanUpgraded(tile.White.Instance(0), 1, base.Front)))
// 	assert.True(t, tst(meld.NewKanOpened(tile.Man2.Instance(2), base.Left)))
// 	assert.True(t, tst(meld.NewKanOpened(tile.East.Instance(1), base.Front)))
// 	assert.True(t, tst(meld.NewKanOpened(tile.Pin1.Instance(0), base.Right)))
// 	assert.True(t, tst(meld.NewPonOpened(tile.Sou4.Instance(3), 2, base.Front)))
// }
