package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

func TestSMessageDrop(t *testing.T) {
	c := NewXMLWriter()
	test := func(o tbase.Opponent, i tile.Instance, isTg bool, s Suggest) {
		params := Drop{}
		params.Opponent = o
		params.Instance = i
		params.IsTsumogiri = isTg
		params.Suggest = s
		c.Drop(params)
	}
	test(tbase.Self, tile.Green.Instance(0), false, 0)
	test(tbase.Self, tile.Green.Instance(0), true, 0)

	test(tbase.Right, tile.Man1.Instance(0), false, 0)
	test(tbase.Right, tile.Man1.Instance(0), true, 0)

	test(tbase.Front, tile.Man2.Instance(0), false, 0)
	test(tbase.Front, tile.Man2.Instance(0), true, 0)

	test(tbase.Left, tile.Man1.Instance(1), false, 0)
	test(tbase.Left, tile.Man1.Instance(1), true, 1)

	assert.Equal(t, `<D128/> <d128/> <E0/> <e0/> <F4/> <f4/> <G1/> <g1 t="1"/>`, c.String())
}

func TestSMessageTake(t *testing.T) {
	c := NewXMLWriter()
	test := func(o tbase.Opponent, i tile.Instance, s Suggest) {
		params := Take{}
		params.Opponent = o
		params.Instance = i
		params.Suggest = s
		c.Take(params)
	}
	test(tbase.Self, tile.Pin1.Instance(0), 0)
	test(tbase.Self, tile.Pin1.Instance(0), 1)
	test(tbase.Right, tile.Pin1.Instance(0), 0)
	test(tbase.Front, tile.Pin1.Instance(0), 0)
	test(tbase.Left, tile.Pin1.Instance(0), 1)
	assert.Equal(t, `<T36/> <T36 t="1"/> <U/> <V/> <W/>`, c.String())
}

func TestSMessageReach(t *testing.T) {
	c := NewXMLWriter()
	test := func(o tbase.Opponent, step int, sc []score.Money) {
		params := Reach{}
		params.Opponent = o
		params.Step = step
		params.Score = sc
		c.Reach(params)
	}
	test(tbase.Front, 1, nil)
	test(tbase.Self, 2, []score.Money{25000, 25000, 25000, 24000})

	assert.Equal(t, `<REACH who="2" step="1"/> <REACH who="0" ten="250,250,250,240" step="2"/>`, c.String())
}

func TestSDeclare(t *testing.T) {
	c := NewXMLWriter()
	x := &tbase.Called{
		Type:     tbase.Pon,
		Opponent: tbase.Front,
		Called:   tile.Sou4.Instance(3),
		Upgraded: tile.Sou4.Instance(2),
	}
	params := Declare{}
	params.Opponent = tbase.Self
	params.Meld = tbase.EncodeCalled(x)
	c.Declare(params)
	assert.Equal(t, `<N who="0" m="33354"/>`, c.String())
}

func TestSInit(t *testing.T) {
	c := NewXMLWriter()
	tg := compact.NewTileGenerator()
	tiles, err := tg.InstancesFromString("579m577p234679s57z")
	require.NoError(t, err)
	b := tbase.Init{
		tbase.Seed{
			RoundNumber: 1,
			Honba:       2,
			Sticks:      3,
			Dice:        [2]int{5, 4},
			Indicator:   tg.Instance(tile.Man4),
		},
		[]score.Money{25000, 25000, 25000, 24000},
		tbase.Front,
		nil,
	}
	c.Init(Init{b, tiles})

	assert.Equal(t, `<INIT seed="1,2,3,5,4,12" ten="250,250,250,240" oya="2" hai="16,24,32,52,60,61,76,80,84,92,96,104,124,132"/>`, c.String())
}

func TestSLogInfo(t *testing.T) {
	c := NewXMLWriter()
	c.LogInfo(LogInfo{WithDealer{tbase.Front}, "2018011712gm-0009-0000-84f0883b"})
	assert.Equal(t, `<TAIKYOKU oya="2" log="2018011712gm-0009-0000-84f0883b"/>`, c.String())
}

func TestSGo(t *testing.T) {
	c := NewXMLWriter()
	c.Go(Go{WithLobby{0, 9}, "12D32385-4987D01D"})
	assert.Equal(t, `<GO type="9" lobby="0" gpid="12D32385-4987D01D"/>`, c.String())
}

func TestSUN(t *testing.T) {
	c := NewXMLWriter()
	c.UserList(UserList{tbase.UserList{
		tbase.User{Name: "NoName", Rate: 1500.00, Sex: tbase.SexMale, Dan: 0},
		tbase.User{Name: "@yukimi", Rate: 1499.99, Sex: tbase.SexFemale, Dan: 2},
		tbase.User{Name: "さあ、行くぞ", Rate: 1467.57, Sex: tbase.SexMale, Dan: 10},
		tbase.User{Name: "pjgwpdjw", Rate: 1500.00, Sex: tbase.SexMale, Dan: 0},
	}})

	assert.Equal(t, `<UN n0="%4E%6F%4E%61%6D%65" n1="%40%79%75%6B%69%6D%69" n2="%E3%81%95%E3%81%82%E3%80%81%E8%A1%8C%E3%81%8F%E3%81%9E" n3="%70%6A%67%77%70%64%6A%77" dan="0,2,10,0" rate="1500.00,1499.99,1467.57,1500.00" sx="M,F,M,M"/>`,
		c.String())
}

func TestSHello(t *testing.T) {
	c := NewXMLWriter()
	c.Hello(Hello{
		Name:        "NoName",
		Auth:        "20180117-e7b5e83e",
		RatingScale: DefaultRatingScale})
	assert.Equal(t, `<HELO uname="%4E%6F%4E%61%6D%65" auth="20180117-e7b5e83e" ratingscale="PF3=1.000000&PF4=1.000000&PF01C=0.582222&PF02C=0.501632&PF03C=0.414869&PF11C=0.823386&PF12C=0.709416&PF13C=0.586714&PF23C=0.378722&PF33C=0.535594&PF1C00=8.000000"/>`,
		c.String())
}

// Get: <AGARI ba="0,1" hai="14,18,19,21,23,26,36,37,57,63,66,82,84,91" machi="91" ten="20,5200,0" yaku="1,1,0,1,7,1,52,1,53,0" doraHai="79,33" doraHaiUra="29,50" who="3" fromWho="3" sc="250,-13,250,-13,250,-26,240,62" />
// <PROF lobby="0" type="9" add="-39.0,0,0,0,1,0,9,0,2,4,1"/>
// <RYUUKYOKU ba="0,1" sc="120,-10,404,-10,246,-10,220,30" hai3="19,22,26,29,30,62,66,81,84,91,98,102,105" owari="110,-39.0,404,50.0,236,-16.0,250,5.0" />
// <PROF lobby="0" type="9" add="-38.0,0,0,0,1,0,8,1,1,4,2"/> <AGARI ba="0,0" hai="23,24,31,56,62,66,83,86,89,100,101" m="15977" machi="66" ten="30,1100,0" yaku="8,1" doraHai="14" who="1" fromWho="1" sc="124,-3,372,11,141,-3,363,-5" owari="121,-38.0,383,48.0,138,-26.0,358,16.0" />
