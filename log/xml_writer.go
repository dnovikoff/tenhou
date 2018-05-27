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

func (w XMLWriter) Open(Info) bool {
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

func (this XMLWriter) Shuffle(params Shuffle) {
	this.Begin("SHUFFLE").
		WriteArg("seed", params.Seed).
		WriteArg("ref", params.Ref).
		End()
}

func (this XMLWriter) Go(params client.WithLobby) {
	this.client().Go(client.Go{WithLobby: params})
}

func (this XMLWriter) UserList(params client.UserList) {
	this.WriteUserList(params.Users, this.FloatFormat)
}

func (this XMLWriter) Start(params client.WithDealer) {
	this.LogInfo(client.LogInfo{WithDealer: params})
}

func (this XMLWriter) Init(params Init) {
	this.Begin("INIT").
		WriteArg("seed", parser.SeedString(&params.Seed)).
		WriteArg("ten", util.ScoreString(params.Scores)).
		WriteDealer(params.Dealer)
	for k, v := range params.Hands {
		this.WriteArg("hai"+strconv.Itoa(k), util.InstanceString(v))
	}
	if params.Shuffle != "" {
		this.WriteArg("shuffle", params.Shuffle)
	}
	this.End()
}

func (this XMLWriter) Draw(params WithOpponentAndInstance) {
	this.WriteTake(params.Opponent, params.Instance, 0, false)
}

func (this XMLWriter) Discard(params WithOpponentAndInstance) {
	x := client.Drop{}
	x.WithInstance = params.WithInstance
	x.WithOpponent = params.WithOpponent
	this.client().Drop(x)
}

func (this XMLWriter) Declare(params Declare) {
	x := client.Declare{}
	x.WithOpponent = params.WithOpponent
	x.Meld = params.Meld
	this.client().Declare(x)
}

func (this XMLWriter) Ryuukyoku(params tbase.Ryuukyoku) {
	this.WriteRyuukyoku(&params, this.FloatFormat)
}

func (this XMLWriter) Agari(params tbase.Agari) {
	this.WriteAgari(&params, this.FloatFormat)
}

func (this XMLWriter) Disconnect(params client.WithOpponent) {
	this.Begin("BYE").
		WriteOpponent("who", params.Opponent).
		End()
}
