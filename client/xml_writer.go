package client

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

type XMLWriter struct {
	util.XMLWriter
}

var _ Controller = XMLWriter{}

func NewXMLWriter() XMLWriter {
	return XMLWriter{util.NewXMLWriter()}
}

func (this XMLWriter) Drop(o base.Opponent, t tile.Instance, isTsumogiri bool, suggest Suggest) {
	letter := 'D'
	if isTsumogiri {
		letter = 'd'
	}
	letter += rune(o)
	if suggest > 0 {
		this.WriteBody(`%s%d t="%d"`, string(letter), t, suggest)
		return
	}
	this.WriteBody("%s%d", string(letter), t)
}

func (this XMLWriter) WriteTake(o base.Opponent, t tile.Instance, suggest Suggest, hideOthers bool) {
	letter := 'T'
	letter += rune(o)
	if hideOthers && o != base.Self {
		this.WriteBody(string(letter))
		return
	}
	if suggest > 0 {
		this.WriteBody(`%s%d t="%d"`, string(letter), t, suggest)
		return
	}
	this.WriteBody("%s%d", string(letter), t)
}

func (this XMLWriter) Take(o base.Opponent, t tile.Instance, s Suggest) {
	this.WriteTake(o, t, s, true)
}

func (this XMLWriter) Reach(o base.Opponent, step int, score []score.Money) {
	ten := ""
	if len(score) > 0 {
		ten = ` ten="` + util.ScoreString(score) + `"`
	}
	this.WriteBody(`REACH who="%d"%s step="%d"`, o, ten, step)
}

func (this XMLWriter) Declare(o base.Opponent, m tbase.Meld, t Suggest) {
	w := this.Begin("N").
		WriteOpponent("who", o).
		WriteIntArg("m", int(m))

	if t > 0 {
		w.WriteIntArg("t", int(t))
	}
	w.End()
}

func (w XMLWriter) writeInit(in Init) {
	w.WriteArg("seed", in.Seed.String()).
		WriteArg("ten", util.ScoreString(in.Scores))
	if in.Chip != nil {
		w.WriteArg("chip", util.IntsString(in.Chip))
	}
	w.WriteDealer(in.Dealer).
		WriteArg("hai", util.InstanceString(in.Hand))
}

func (this XMLWriter) Init(in Init) {
	this.Begin("INIT")
	this.writeInit(in)
	this.End()
}

func (this XMLWriter) Reinit(in Reinit) {
	this.Begin("REINIT")
	this.writeInit(in.Init)
	for k, v := range in.Melds {
		if len(v) == 0 {
			continue
		}
		this.WriteArg("m"+strconv.Itoa(k), util.MeldString(v))
	}
	for k, v := range in.Discard {
		tiles := v
		if len(tiles) == 0 {
			continue
		}
		r := in.Riichi[k]
		if r > 0 {
			tiles = make(tile.Instances, 0, len(v)+1)
			tiles = append(tiles, v[:r]...)
			tiles = append(tiles, 255)
			tiles = append(tiles, v[r:]...)
		}
		this.WriteArg("kawa"+strconv.Itoa(k), util.InstanceString(tiles))
	}
	this.End()
}

func (w XMLWriter) LogInfo(dealer base.Opponent, hash string) {
	w.Begin("TAIKYOKU")
	w.WriteDealer(dealer)
	// No hash for game with bots
	if hash != "" {
		w.WriteArg("log", hash)
	}
	w.End()
}

func (w XMLWriter) Go(gpid string, gameType int, lobby int) {
	w.Begin("GO")
	w.WriteIntArg("type", gameType)
	if lobby != -1 {
		w.WriteIntArg("lobby", lobby)
	}
	if gpid != "" {
		w.WriteArg("gpid", gpid)
	}
	w.End()
}

func (w XMLWriter) Hello(name string, auth string, stats HelloStats) {
	w.Begin("HELO")
	w.WriteArg("uname", util.Escape(name))
	w.WriteArg("auth", auth)
	if stats.Message != "" {
		em := util.EscapeWithExceptions(stats.Message, " :0123456789")
		w.WriteArg("nintei", em)
	}
	if stats.PF4 != "" {
		w.WriteArg("PF4", stats.PF4)
	}
	if stats.Expire != 0 {
		w.WriteIntArg("expire", stats.Expire)
	}
	if stats.ExpireDays != 0 {
		w.WriteIntArg("expiredays", stats.ExpireDays)
	}
	if stats.RatingScale != "" {
		w.WriteArg("ratingscale", stats.RatingScale)
	}
	if stats.RR != "" {
		w.WriteArg("rr", stats.RR)
	}
	w.End()
}

func (w XMLWriter) UserList(list tbase.UserList) {
	w.WriteUserList(list, true)
}

func (w XMLWriter) WriteUserList(list tbase.UserList, isFloatFormat bool) {
	w.Begin("UN")
	dan := make([]string, len(list))
	rate := make([]string, len(list))
	sex := make([]string, len(list))
	rc := make([]string, len(list))

	haveRc := false
	for k, v := range list {
		w.WriteArg("n"+strconv.Itoa(k), util.Escape(v.Name))
		dan[k] = strconv.Itoa(v.Dan)
		if isFloatFormat {
			rate[k] = fmt.Sprintf("%.2f", v.Rate)
		} else {
			rate[k] = strconv.Itoa(int(v.Rate))
		}
		if v.Rc != nil {
			rc[k] = strconv.Itoa(*v.Rc)
			haveRc = true
		}
		sex[k] = v.Sex.Letter()
	}
	w.WriteArg("dan", strings.Join(dan, ","))
	if haveRc {
		w.WriteArg("rc", strings.Join(rc, ","))
	}
	w.WriteArg("rate", strings.Join(rate, ","))
	w.WriteArg("sx", strings.Join(sex, ","))
	w.End()
}

func (this XMLWriter) LobbyStats(n, j, g string) {
	this.WriteBody(`LN n="%s" j="%s" g="%s"`, n, j, g)
}

func (w XMLWriter) Agari(a *tbase.Agari) {
	w.WriteAgari(a, true)
}

func (w XMLWriter) WriteAgari(a *tbase.Agari, floatFormat bool) {
	w.Begin("AGARI")
	w.WriteTableStatus(a.Status)
	w.WriteArg("hai", util.InstanceString(a.Hand))
	if len(a.Melds) > 0 {
		w.WriteArg("m", util.MeldString(a.Melds))
	}
	w.WriteIntArg("machi", int(a.WinTile))
	w.WriteArg("ten", util.ScoreString(a.Scores))
	if len(a.Yakus) > 0 {
		w.WriteArg("yaku", util.YakuString(a.Yakus))
	}
	if len(a.Yakumans) > 0 {
		w.WriteArg("yakuman", util.YakumanString(a.Yakumans))
	}
	w.WriteArg("doraHai", util.InstanceString(a.DoraIndicators))
	if len(a.UraIndicators) > 0 {
		w.WriteArg("doraHaiUra", util.InstanceString(a.UraIndicators))
	}
	w.WriteOpponent("who", a.Who)
	w.WriteOpponent("fromWho", a.From)
	if a.Pao != nil {
		w.WriteOpponent("paoWho", *a.Pao)
	}
	w.WriteScoreChanges(a.Changes)
	if len(a.FinalScores) > 0 {
		w.WriteArg("owari", util.FinalsString(a.FinalScores, floatFormat))
	}
	if a.Ratio != "" {
		w.WriteArg("ratio", a.Ratio)
	}
	w.AddTrailingSpace().End()
}

func (this XMLWriter) Indicator(i tile.Instance) {
	this.WriteBody(`DORA hai="%d" `, i)
}

func (this XMLWriter) EndButton(lobby int, tp int, add string) {
	this.WriteBody(`PROF lobby="%d" type="%d" add="%s"`, lobby, tp, add)
}

func (this XMLWriter) Furiten(show bool) {
	val := 0
	if show {
		val = 1
	}
	this.WriteBody(`FURITEN show="%d" `, val)
}

func (w XMLWriter) Ryuukyoku(a *tbase.Ryuukyoku) {
	w.WriteRyuukyoku(a, true)
}

func (w XMLWriter) WriteRyuukyoku(a *tbase.Ryuukyoku, floats bool) {
	w.Begin("RYUUKYOKU")
	if a.DrawType != tbase.DrawEnd {
		w.WriteArg("type", tbase.ReverseDrawMap[a.DrawType])
	}
	w.WriteTableStatus(a.TableStatus)
	w.WriteScoreChanges(a.ScoreChanges)
	for k, v := range a.Hands {
		if v == nil {
			continue
		}
		w.WriteArg("hai"+strconv.Itoa(k), util.InstanceString(v))
	}
	if len(a.Finals) > 0 {
		w.WriteArg("owari", util.FinalsString(a.Finals, floats))
	}
	w.AddTrailingSpace().End()
}

func (this XMLWriter) Rejoin(lobbyNumber int, lobbyType int, rejoin bool) {
	r := ""
	if rejoin {
		r = ",r"
	}
	this.WriteBody(`REJOIN t="%d,%d%s"`, lobbyNumber, lobbyType, r)
}

func (this XMLWriter) Disconnect(o base.Opponent) {
	this.WriteBody(`BYE who="%d" `, o)
}

func (w XMLWriter) Reconnect(o base.Opponent, name string) {
	w.Begin("UN")
	w.WriteArg("n"+strconv.Itoa(int(o)), util.Escape(name))
	w.AddTrailingSpace().End()
}

func (w XMLWriter) Chat(name string, text string) {
	w.Begin("CHAT")
	w.WriteArg("uname", util.Escape(name))
	w.WriteArg("text", util.Escape(text))
	w.AddTrailingSpace().End()
}

func (this XMLWriter) Ranking(v2 string) {
	this.WriteBody(`RANKING v2="%s"`, v2)
}

func (w XMLWriter) Recover(status tbase.TableStatus, dealer base.Opponent, sc tbase.ScoreChanges) {
	w.Begin("SAIKAI")
	w.WriteTableStatus(status)
	w.WriteDealer(dealer)
	w.WriteScoreChanges(sc)
	w.End()
}
