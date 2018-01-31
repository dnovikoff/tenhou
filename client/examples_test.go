package client

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

type ExampleWriter struct {
	XMLWriter
	num    int
	hash   string
	camera base.Opponent
	w      io.Writer
}

func (this *ExampleWriter) Init(i Init) {
	this.num++
	this.XMLWriter.Init(i)
}

func (this *ExampleWriter) write(s Suggest, format string, args ...interface{}) {
	fmt.Fprintf(this.w, "http://tenhou.net/0/?log=%s&ts=%d&tw=%d SUGGEST=%d ",
		this.hash,
		this.num,
		this.camera,
		s,
	)
	fmt.Fprintf(this.w, format+"\n", args...)
}

func (this *ExampleWriter) Drop(o base.Opponent, t tile.Instance, isTsumogiri bool, s Suggest) {
	this.XMLWriter.Drop(o, t, isTsumogiri, s)
	if s == 0 {
		return
	}
	this.write(s, "DROP=%s", t.Tile())
}

func (this *ExampleWriter) Take(o base.Opponent, t tile.Instance, s Suggest) {
	this.XMLWriter.Take(o, t, s)
	if s == 0 {
		return
	}
	this.write(s, "TAKE=%s", t.Tile())
}

func (this *ExampleWriter) Declare(o base.Opponent, m tbase.Meld, s Suggest) {
	this.XMLWriter.Declare(o, m, s)
	if s == 0 {
		return
	}
	i := compact.NewInstances()
	m.Convert().AddTo(i)

	this.write(s, "Call=%s", i.Instances())
}

func (this *ExampleWriter) LogInfo(dealer base.Opponent, hash string) {
	this.hash = hash
	this.num = -1
	this.camera = (4 - dealer) % 4
	this.XMLWriter.LogInfo(dealer, hash)
}

var _ Controller = &ExampleWriter{}

func TestExamples(t *testing.T) {
	// f, err := os.Create("suggest.txt")
	// require.NoError(t, err)
	// defer f.Close()
	f := &bytes.Buffer{}
	w := &ExampleWriter{XMLWriter: NewXMLWriter(), w: f}
	util.ProcessExampleLogs(t, "Get:", func(line string) {
		w.Reset()
		require.NoError(t, ProcessXMLMessage(line, w), line)
		util.CompareLines(t, line, w.String())
	})
}
