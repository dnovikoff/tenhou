package stats

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tenhou/genproto/stats"
)

func TestParseLobbyInfo(t *testing.T) {
	p := func(in string) *stats.Record {
		p := &statParser{}
		r := p.Next()
		require.NoError(t, p.LobbyConfig(in))
		return r
	}

	assert.Equal(t, &stats.Record{
		Type:   stats.GameType_GAME_TYPE_3,
		Lobby:  stats.LobbyType_LOBBY_TYPE_KU,
		Length: stats.GameLength_GAME_LENGTH_SOUTH,
		Tanyao: stats.Tanyao_TANYAO_YES,
		Akkas:  stats.Akkas_AKKAS_YES,
		IsDz:   false,
	}, p("三般南喰赤－"))

	assert.Equal(t, &stats.Record{
		Type:   stats.GameType_GAME_TYPE_4,
		Lobby:  stats.LobbyType_LOBBY_TYPE_KU,
		Length: stats.GameLength_GAME_LENGTH_SOUTH,
		Tanyao: stats.Tanyao_TANYAO_YES,
		Akkas:  stats.Akkas_AKKAS_YES,
		IsDz:   false,
	}, p("四般南喰赤－"))

	assert.Equal(t, &stats.Record{
		Type:   stats.GameType_GAME_TYPE_4,
		Lobby:  stats.LobbyType_LOBBY_TYPE_DAN,
		Length: stats.GameLength_GAME_LENGTH_SOUTH,
		Tanyao: stats.Tanyao_TANYAO_YES,
		Akkas:  stats.Akkas_AKKAS_YES,
		IsDz:   false,
	}, p("四上南喰赤－"))

	assert.Equal(t, &stats.Record{
		Type:   stats.GameType_GAME_TYPE_3,
		Lobby:  stats.LobbyType_LOBBY_TYPE_PHOENIX,
		Length: stats.GameLength_GAME_LENGTH_SOUTH,
		Tanyao: stats.Tanyao_TANYAO_YES,
		Akkas:  stats.Akkas_AKKAS_YES,
		IsDz:   false,
	}, p("三鳳南喰赤－"))

	assert.Equal(t, &stats.Record{
		Type:   stats.GameType_GAME_TYPE_4,
		Lobby:  stats.LobbyType_LOBBY_TYPE_DZ,
		Length: stats.GameLength_GAME_LENGTH_ONE,
		Tanyao: stats.Tanyao_TANYAO_YES,
		Akkas:  stats.Akkas_AKKAS_YES,
		IsDz:   false,
	}, p("四技－喰赤－"))

	assert.Equal(t, &stats.Record{
		Type:   stats.GameType_GAME_TYPE_4,
		Lobby:  stats.LobbyType_LOBBY_TYPE_KU,
		Length: stats.GameLength_GAME_LENGTH_SOUTH,
		Tanyao: stats.Tanyao_TANYAO_YES,
		Akkas:  stats.Akkas_AKKAS_YES,
		IsDz:   false,
	}, p("四般南喰赤－"))

	assert.Equal(t, &stats.Record{
		Type:       stats.GameType_GAME_TYPE_3,
		Lobby:      stats.LobbyType_LOBBY_TYPE_X1,
		Length:     stats.GameLength_GAME_LENGTH_EAST,
		Tanyao:     stats.Tanyao_TANYAO_YES,
		Akkas:      stats.Akkas_AKKAS_YES,
		IsDz:       true,
		NumberType: stats.NumberType_NUMBER_5,
	}, p("三若東喰赤祝５"))
}

func TestParseLobbyError(t *testing.T) {
	e := func(in string) error {
		p := &statParser{}
		p.Next()
		err := p.LobbyConfig(in)
		require.Error(t, err)
		return err
	}

	assert.Error(t, e("三般南喰赤－12"))
	assert.Error(t, e("123456"))
	assert.Error(t, e("三般南喰赤 "))
	assert.Error(t, e("3般南喰赤－"))
}

func TestParsePlayer(t *testing.T) {
	p := func(in string) *stats.Player {
		x := &stats.Player{}
		require.NoError(t, parsePlayer(in, x))
		return x
	}
	assert.Equal(t, &stats.Player{Name: "小狂三", Score: 510}, p("小狂三(+51.0)"))
	assert.Equal(t, &stats.Player{Name: "exia_ang", Score: -570}, p("exia_ang(-57.0)"))
	assert.Equal(t, &stats.Player{Name: "初心者です", Score: 26}, p("初心者です(+2.6)"))
	assert.Equal(t, &stats.Player{Name: "ブルース・リー", Score: -420, Coins: -1}, p("ブルース・リー(-42.0,-1枚)"))
}

func TestParsePlaterError(t *testing.T) {
	p := func(in string) error {
		return parsePlayer(in, &stats.Player{})
	}
	assert.Error(t, p("小狂三(+51.0"))
	assert.Error(t, p("exia_ang-57.0)"))
	assert.Error(t, p("初心者です)+2.6("))
	assert.Error(t, p("ブルース・リー()"))
	assert.Error(t, p("ブルース・リー(1vasja)"))
}

func TestParsePlayers(t *testing.T) {
	p := func(in string) []*stats.Player {
		p := &statParser{}
		r := p.Next()
		err := p.Players(in)
		require.NoError(t, err)
		return r.GetPlayers()
	}
	assert.Equal(t, []*stats.Player{
		{Name: "初心者です", Score: 26},
		{Name: "〓アナスタシア〓", Score: 0},
		{Name: "放銃王４", Score: 0},
		{Name: "DERESUKE", Score: -26}},
		p("初心者です(+2.6) 〓アナスタシア〓(0.0) 放銃王４(0.0) DERESUKE(-2.6)<br>"))

	assert.Equal(t, []*stats.Player{
		{Name: "qzc", Score: 460, Coins: -1},
		{Name: "HUNTER", Score: -50, Coins: 1},
		{Name: "asd", Score: -410}},
		p("qzc(+46.0,-1枚) HUNTER(-5.0,+1枚) asd(-41.0,0枚)"))

	assert.Equal(t, []*stats.Player{
		{Name: "xxx", Score: 420},
		{Name: "(>\317\211<*)", Score: 40},
		{Name: "yyy", Score: -170},
		{Name: "zzz", Score: -290},
	},
		p("xxx(+42.0) (&gt;ω&lt;*)(+4.0) yyy(-17.0) zzz(-29.0)<br>"))

}

func TestParseID(t *testing.T) {
	p := func(in string) string {
		p := &statParser{}
		r := p.Next()
		err := p.Info(in)
		require.NoError(t, err)
		return r.Id
	}
	assert.Equal(t, "2015121400gm-00b1-0000-7372c860", p(`<a href="http://tenhou.net/0/?log=2015121400gm-00b1-0000-7372c860">牌譜</a>`))
}

func TestParseFiles(t *testing.T) {
	p := func(filename string) error {
		t.Log("Parsing ", filename)
		file, err := os.Open("test_data/" + filename + ".log")
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file.Close())
		}()
		return ParseStats(file, &StatFileInfo{filename, time.Now()}, func(*stats.Record) error { return nil })
	}
	assert.NoError(t, p("sca"))
	assert.NoError(t, p("scb"))
	assert.NoError(t, p("scc"))
	assert.NoError(t, p("scd"))
	assert.NoError(t, p("sce"))
	assert.NoError(t, p("scf"))
}

func TestParseFilename(t *testing.T) {
	p := func(x string) string {
		info, err := ParseStatFileName(x)
		require.NoError(t, err)
		return info.Source + " " + info.Date.String()
	}
	assert.Equal(t, "sce 2018-01-01 00:00:00 +0000 UTC", p("tenhou/stats/dat/2018/sce20180101.html.gz"))
	assert.Equal(t, "scc 2018-10-11 00:00:00 +0000 UTC", p("tenhou/stats/dat/scc2018101100.html.gz"))
}

func TestParseFilesByName(t *testing.T) {
	var last *stats.Record
	p := func(filename string) error {
		file, err := os.Open("test_data/" + filename)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, file.Close())
		}()
		return ParseStatsForFile(file, filename, func(x *stats.Record) error {
			last = x
			return nil
		})
	}
	assert.NoError(t, p("sce20180101.html"))
	assert.NoError(t, p("scc20120404.html"))
	assert.Equal(t, "Lask&Yu", last.GetPlayers()[2].Name)
}

// func TestBadLine(t *testing.T) {
// 	// L4714 | 01:35 | 四般東喰赤 | 牛乳一気飲<A3><B2>R34(+41) 田中っち(+6) へたくそ君。(-17) どつぼ(-30)
// 	nameBytes := []byte{231, 137, 155, 228, 185, 179, 228, 184, 128, 230, 176, 151, 233, 163, 178, 163, 178, 82, 51, 52}
// 	in := fmt.Sprintf("L4714 | 01:35 | 四般東喰赤 | %s(+41) 田中っち(+6) へたくそ君。(-17) どつぼ(-30)", nameBytes)
// 	var name string
// 	err := ParseFromSource(db.SourceA, bytes.NewReader([]byte(in)), time.Now(), func(r db.Record) error {
// 		name = r.Results[0].Player.Name
// 		return nil
// 	})
// 	require.NoError(t, err)
// 	// sanitized
// 	assert.Equal(t, "_牛乳一気飲R34", name)
// }

// // 07:23 | 11 | 四銀東喰赤祝５ | 和牛(+60.0,+6枚) メカゼットン(+10.0,-1枚) 前田(-20.0,-2枚) 雀達(-50.0,-3枚)
