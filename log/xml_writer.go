package log

import (
	"strconv"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

type XMLWriter struct {
	client.XMLWriter
	FloatFormat bool
}

var _ Controller = &XMLWriter{}

func (w XMLWriter) Open(*Info) bool {
	w.Write(`<mjloggm ver="2.3">`)
	return true
}

func (w XMLWriter) Close() {
	w.Write(`</mjloggm>`)
}

func (this *XMLWriter) SetFloatFormat() {
	this.FloatFormat = true
}

func (this XMLWriter) client() client.XMLWriter {
	return this.XMLWriter
}

func (this XMLWriter) Shuffle(seed, ref string) {
	this.Begin("SHUFFLE").
		WriteArg("seed", seed).
		WriteArg("ref", ref).
		End()
}

func (this XMLWriter) Go(t int, lobby int) {
	this.client().Go("", t, lobby)
}

func (this XMLWriter) UserList(users tbase.UserList) {
	this.WriteUserList(users, this.FloatFormat)
}

func (this XMLWriter) Start(d base.Opponent) {
	this.LogInfo(d, "")
}

func (this XMLWriter) Init(in *Init) {
	this.Begin("INIT").
		WriteArg("seed", in.Seed.String()).
		WriteArg("ten", util.ScoreString(in.Scores)).
		WriteDealer(in.Dealer)
	for k, v := range in.Hands {
		this.WriteArg("hai"+strconv.Itoa(k), util.InstanceString(v))
	}
	if in.Shuffle != "" {
		this.WriteArg("shuffle", in.Shuffle)
	}
	this.End()
}

func (this XMLWriter) Draw(o base.Opponent, t tile.Instance) {
	this.WriteTake(o, t, 0, false)
}

func (this XMLWriter) Discard(o base.Opponent, t tile.Instance) {
	this.client().Drop(o, t, false, 0)
}

func (this XMLWriter) Declare(o base.Opponent, m tbase.Meld) {
	this.client().Declare(o, m, 0)
}

func (this XMLWriter) Ryuukyoku(a *tbase.Ryuukyoku) {
	this.WriteRyuukyoku(a, this.FloatFormat)
}

func (this XMLWriter) Reach(o base.Opponent, step int, score []score.Money) {
	this.client().Reach(o, step, score)
}

func (this XMLWriter) Agari(a *tbase.Agari) {
	this.WriteAgari(a, this.FloatFormat)
}

func (this XMLWriter) Indicator(x tile.Instance) {
	this.client().Indicator(x)
}

func (this XMLWriter) Disconnect(x base.Opponent) {
	this.Begin("BYE").
		WriteOpponent("who", x).
		End()
}

func (this XMLWriter) Reconnect(x base.Opponent, name string) {
	this.client().Reconnect(x, name)
}
