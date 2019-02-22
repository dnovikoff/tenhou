package server

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

type exampleWriter struct {
	XMLWriter
	t *testing.T
}

const goodID = "IDDEADBEAF-xxxxxxxx"

func (w exampleWriter) Hello(name string, tid string, sex tbase.Sex) {
	if name != goodID {
		require.False(w.t, strings.HasPrefix(name, "ID"), "remove sensetive data "+name)
	}
	w.XMLWriter.Hello(name, tid, sex)
}

func TestExamples(t *testing.T) {
	w := exampleWriter{NewXMLWriter(), t}
	util.ProcessExampleLogs(t, "Send:", func(line string) {
		w.Reset()
		require.NoError(t, ProcessXMLMessage(line, w), line)
		util.CompareLines(t, line, w.String())
	})
}
