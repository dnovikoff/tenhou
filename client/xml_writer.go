package client

import (
	"strconv"

	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/parser"
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

func (w XMLWriter) Drop(params Drop) {
	letter := 'D'
	i := util.InstanceToTenhou(params.Instance)
	if params.IsTsumogiri {
		letter = 'd'
	}
	letter += rune(params.Opponent)
	if params.Suggest > 0 {
		w.WriteBody(`%s%d t="%d"`, string(letter), i, params.Suggest)
		return
	}
	w.WriteBody("%s%d", string(letter), i)
}

func (w XMLWriter) WriteTake(o tbase.Opponent, t tile.Instance, suggest Suggest, hideOthers bool) {
	letter := 'T'
	letter += rune(o)
	i := util.InstanceToTenhou(t)
	if hideOthers && o != tbase.Self {
		w.WriteBody(string(letter))
		return
	}
	if suggest > 0 {
		w.WriteBody(`%s%d t="%d"`, string(letter), i, suggest)
		return
	}
	w.WriteBody("%s%d", string(letter), i)
}

func (w XMLWriter) Take(params Take) {
	w.WriteTake(
		params.Opponent,
		params.Instance,
		params.Suggest,
		true)
}

func (w XMLWriter) Reach(params Reach) {
	w.Begin("REACH").
		WriteWho(params.Opponent)
	if len(params.Score) > 0 {
		w.WriteArg("ten", util.ScoreString(params.Score))
	}
	w.WriteIntArg("step", params.Step).
		End()
}

func (w XMLWriter) Declare(params Declare) {
	wr := w.Begin("N").
		WriteWho(params.Opponent).
		WriteIntArg("m", int(params.Meld))

	if params.Suggest > 0 {
		wr.WriteIntArg("t", int(params.Suggest))
	}
	wr.End()
}

func (w XMLWriter) writeInit(in Init) {
	w.WriteArg("seed", parser.SeedString(&in.Seed)).
		WriteArg("ten", util.ScoreString(in.Scores)).
		WriteListInt("chip", in.Chip)

	w.WriteDealer(in.Dealer).
		WriteArg("hai", util.InstanceString(in.Hand))
}

func (w XMLWriter) Init(params Init) {
	w.Begin("INIT")
	w.writeInit(params)
	w.End()
}

func (w XMLWriter) Reinit(params Reinit) {
	w.Begin("REINIT")
	w.writeInit(params.Init)
	for k, v := range params.Melds {
		if len(v) == 0 {
			continue
		}
		w.WriteArg("m"+strconv.Itoa(k), util.MeldString(v))
	}
	for k, v := range params.Discard {
		tiles := v
		if len(tiles) == 0 {
			continue
		}
		r := params.Riichi[k]
		if r > 0 {
			tiles = make(tile.Instances, 0, len(v)+1)
			tiles = append(tiles, v[:r]...)
			tiles = append(tiles, 255)
			tiles = append(tiles, v[r:]...)
		}
		w.WriteArg("kawa"+strconv.Itoa(k), util.InstanceString(tiles))
	}
	w.End()
}

func (w XMLWriter) LogInfo(params LogInfo) {
	w.Begin("TAIKYOKU")
	w.WriteDealer(params.Dealer)
	// No hash for game with bots
	if params.Hash != "" {
		w.WriteArg("log", params.Hash)
	}
	w.End()
}

func (w XMLWriter) Go(params Go) {
	w.Begin("GO")
	w.WriteIntArg("type", params.LobbyType)
	if params.LobbyNumber != -1 {
		w.WriteIntArg("lobby", params.LobbyNumber)
	}
	if params.GpID != "" {
		w.WriteArg("gpid", params.GpID)
	}
	w.End()
}

func (w XMLWriter) Hello(params Hello) {
	w.Begin("HELO")
	w.WriteArg("uname", util.Escape(params.Name))
	w.WriteArg("auth", params.Auth)
	if params.Message != "" {
		em := util.EscapeWithExceptions(params.Message, " :0123456789")
		w.WriteArg("nintei", em)
	}
	if params.PF4 != "" {
		w.WriteArg("PF4", params.PF4)
	}
	if params.Expire != 0 {
		w.WriteIntArg("expire", params.Expire)
	}
	if params.ExpireDays != 0 {
		w.WriteIntArg("expiredays", params.ExpireDays)
	}
	if params.RatingScale != "" {
		w.WriteArg("ratingscale", params.RatingScale)
	}
	if params.RR != "" {
		w.WriteArg("rr", params.RR)
	}
	w.End()
}

func (w XMLWriter) UserList(params UserList) {
	w.WriteUserList(params.Users)
}

func (w XMLWriter) WriteUserList(ul tbase.UserList) {
	w.Begin("UN")
	var sex []string
	for _, v := range ul.Sex {
		sex = append(sex, v.Letter())
	}

	for k, v := range ul.Names {
		w.WriteArg("n"+strconv.Itoa(k), util.Escape(v))
	}
	w.WriteListInt("dan", ul.Dan)
	w.WriteListInt("rc", ul.RC)
	w.WriteListInt("gold", ul.Gold)
	w.WriteListFloat("rate", ul.Rate)
	w.WriteList("sx", sex)
	w.End()
}

func (w XMLWriter) LobbyStats(params LobbyStats) {
	w.Begin("LN").
		WriteArg("n", params.N).
		WriteArg("j", params.J).
		WriteArg("g", params.G).
		End()
}

func (w XMLWriter) Agari(a tbase.Agari) {
	w.WriteAgari(&a)
}

func (w XMLWriter) WriteAgari(a *tbase.Agari) {
	w.Begin("AGARI")
	w.WriteTableStatus(a.Status)
	w.WriteArg("hai", util.InstanceString(a.Hand))
	if len(a.Melds) > 0 {
		w.WriteArg("m", util.MeldString(a.Melds))
	}
	w.WriteInstance("machi", a.WinTile)
	w.WriteFmtArg("ten", "%d,%d,%d", a.Score.Fu, a.Score.Total, a.Score.Riichi)
	if len(a.Yakus) > 0 {
		w.WriteArg("yaku", util.YakuString(a.Yakus))
	}
	if len(a.Yakumans) > 0 {
		w.WriteArg("yakuman", util.YakumanString(a.Yakumans))
	}
	w.WriteListInt("chip", a.Chips)
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
		w.WriteArg("owari", util.FinalsString(a.FinalScores))
	}
	if a.Ratio != "" {
		w.WriteArg("ratio", a.Ratio)
	}
	w.AddTrailingSpace().End()
}

func (w XMLWriter) Indicator(params WithInstance) {
	w.Begin("DORA").
		WriteInstance("hai", params.Instance).
		AddTrailingSpace().End()
}

func (w XMLWriter) EndButton(params EndButton) {
	w.Begin("PROF").
		WriteIntArg("lobby", params.LobbyNumber).
		WriteIntArg("type", params.LobbyType).
		WriteArg("add", params.Add).
		End()
}

func (w XMLWriter) Furiten(params Furiten) {
	val := 0
	if params.Furiten {
		val = 1
	}
	w.Begin("FURITEN").
		WriteIntArg("show", val).
		AddTrailingSpace().
		End()
}

func (w XMLWriter) Ryuukyoku(a tbase.Ryuukyoku) {
	w.WriteRyuukyoku(&a)
}

func (w XMLWriter) WriteRyuukyoku(a *tbase.Ryuukyoku) {
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
		w.WriteArg("owari", util.FinalsString(a.Finals))
	}
	w.WriteListInt("ratio", a.Ratio)
	w.AddTrailingSpace().End()
}

func (w XMLWriter) Rejoin(params Rejoin) {
	r := ""
	if params.Rejoin {
		r = ",r"
	}
	w.WriteBody(`REJOIN t="%d,%d%s"`, params.LobbyNumber, params.LobbyType, r)
}

func (w XMLWriter) Disconnect(params WithOpponent) {
	w.Begin("BYE").WriteWho(params.Opponent).End()
}

func (w XMLWriter) Reconnect(params Reconnect) {
	w.Begin("UN")
	w.WriteArg("n"+strconv.Itoa(int(params.Opponent)), util.Escape(params.Name))
	w.AddTrailingSpace().End()
}

func (w XMLWriter) Chat(params Chat) {
	w.Begin("CHAT")
	w.WriteArg("uname", util.Escape(params.Name))
	w.WriteArg("text", util.Escape(params.Message))
	w.AddTrailingSpace().End()
}

func (w XMLWriter) Ranking(params Ranking) {
	w.WriteBody(`RANKING v2="%s"`, params.V2)
}

func (w XMLWriter) Recover(params Recover) {
	w.Begin("SAIKAI")
	w.WriteTableStatus(params.Status)
	w.WriteDealer(params.Dealer)
	w.WriteScoreChanges(params.Changes)
	w.End()
}
