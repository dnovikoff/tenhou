package log

import (
	"strconv"

	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

type XMLWriter struct {
	client.XMLWriter
	FloatFormat bool
}

var _ Controller = &XMLWriter{}

func (w *XMLWriter) Open(Info) bool {
	w.Write(`<mjloggm ver="2.3">`)
	return true
}

func (w *XMLWriter) Close() {
	w.Write(`</mjloggm>`)
}

func (w *XMLWriter) SetFloatFormat() {
	w.FloatFormat = true
}

func (w *XMLWriter) client() client.XMLWriter {
	return w.XMLWriter
}

func (w *XMLWriter) Shuffle(params Shuffle) {
	w.Begin("SHUFFLE").
		WriteArg("seed", params.Seed).
		WriteArg("ref", params.Ref).
		End()
}

func (w *XMLWriter) Go(params client.WithLobby) {
	w.client().Go(client.Go{WithLobby: params})
}

func (w XMLWriter) UserList(params client.UserList) {
	w.WriteUserList(params.Users, w.FloatFormat)
}

func (w *XMLWriter) Start(params client.WithDealer) {
	w.LogInfo(client.LogInfo{WithDealer: params})
}

func (w *XMLWriter) Init(params Init) {
	w.Begin("INIT").
		WriteArg("seed", parser.SeedString(&params.Seed)).
		WriteArg("ten", util.ScoreString(params.Scores)).
		WriteDealer(params.Dealer)
	for k, v := range params.Hands {
		w.WriteArg("hai"+strconv.Itoa(k), util.InstanceString(v))
	}
	if params.Shuffle != "" {
		w.WriteArg("shuffle", params.Shuffle)
	}
	w.End()
}

func (w *XMLWriter) Draw(params WithOpponentAndInstance) {
	w.WriteTake(params.Opponent, params.Instance, 0, false)
}

func (w *XMLWriter) Discard(params WithOpponentAndInstance) {
	x := client.Drop{}
	x.WithInstance = params.WithInstance
	x.WithOpponent = params.WithOpponent
	w.client().Drop(x)
}

func (w *XMLWriter) Declare(params Declare) {
	x := client.Declare{}
	x.WithOpponent = params.WithOpponent
	x.Meld = params.Meld
	w.client().Declare(x)
}

func (w *XMLWriter) Ryuukyoku(params tbase.Ryuukyoku) {
	w.WriteRyuukyoku(&params, w.FloatFormat)
}

func (w *XMLWriter) Agari(params tbase.Agari) {
	w.WriteAgari(&params, w.FloatFormat)
}

func (w *XMLWriter) Disconnect(params client.WithOpponent) {
	w.Begin("BYE").
		WriteOpponent("who", params.Opponent).
		End()
}
