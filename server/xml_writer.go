package server

import (
	"net/url"
	"strings"

	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"

	"github.com/dnovikoff/tempai-core/tile"
)

type XMLWriter struct {
	util.XMLWriter
}

var _ Controller = XMLWriter{}

func NewXMLWriter() XMLWriter {
	return XMLWriter{util.NewXMLWriter()}
}

func (w XMLWriter) Hello(name string, tid string, sex tbase.Sex) {
	w.Begin("HELO").
		WriteArg("name", url.PathEscape(name)).
		WriteArg("tid", tid).
		WriteArg("sx", sex.Letter()).
		AddTrailingSpace().End()
}

func (w XMLWriter) Auth(value string) {
	w.WriteBody(`AUTH val="%s"`, value)
}

func (w XMLWriter) RequestLobbyStatus(v, V int) {
	w.Begin("PXR")
	if v > 0 {
		w.WriteIntArg("v", v)
	}
	if V > 0 {
		w.WriteIntArg("V", V)
	}
	w.AddTrailingSpace().End()
}

func (w XMLWriter) Join(lobbyNumber int, lobbyType int, rejoin bool) {
	r := ""
	if rejoin {
		r = ",r"
	}
	w.WriteBody(`JOIN t="%d,%d%s" `, lobbyNumber, lobbyType, r)
}

func (w XMLWriter) CancelJoin() {
	w.WriteBody(`JOIN `)
}

func (w XMLWriter) Drop(t tile.Instance) {
	w.WriteBody(`D p="%d"`, t)
}

func (w XMLWriter) Call(t Answer, tiles tile.Instances) {
	if t == AnswerSkip {
		// Decline
		w.WriteBody("N ")
		return
	}
	w.Begin("N")
	w.Write(` type="%d"`, t)
	if len(tiles) == 1 {

		w.Write(` hai="%d" `, util.InstanceToTenhou(tiles[0]))
	} else {
		for k, v := range tiles {
			w.Write(` hai%d="%d"`, k, util.InstanceToTenhou(v))
		}
		w.AddTrailingSpace()
	}
	w.End()
}

func (w XMLWriter) Reach(t tile.Instance) {
	w.Begin("REACH").
		WriteInstance("hai", t).
		AddTrailingSpace().
		End()
}

func (w XMLWriter) Ping() {
	w.WriteBody("Z ")
}

func (w XMLWriter) GoOK() {
	w.WriteBody("GOK ")
}

func (w XMLWriter) NextReady() {
	w.WriteBody("NEXTREADY ")
}

func (w XMLWriter) Bye() {
	w.WriteBody("BYE ")
}

func (w XMLWriter) Chat(message string) {
	w.Begin("CHAT")
	prefix := ""
	if strings.HasPrefix(message, "/") {
		prefix, message = message[:1], message[1:]
	}
	w.WriteArg("text", prefix+url.PathEscape(message))
	w.AddTrailingSpace().End()
}
