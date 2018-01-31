package log

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/util"
)

func loadTestData(t require.TestingT, filename string) []byte {
	data, err := ioutil.ReadFile("test_data/" + filename)
	require.NoError(t, err)
	return data
}

func loadTestXml(t require.TestingT, data []byte) parser.Nodes {
	x := &parser.Root{}
	require.NoError(t, xml.Unmarshal(data, &x))
	return x.Nodes
}

func TestAkasValidate(t *testing.T) {
	assert.Equal(t, []tile.Instance{tile.Instance(16), tile.Instance(52), tile.Instance(88)}, yaku.RulesTenhouRed.AkaDoras)
}

func processXMLFiles(t util.TestingT, f func(data string, nodes parser.Nodes)) {
	dir := "./test_data/"
	infos, err := ioutil.ReadDir(dir)
	require.NoError(t, err)
	tested := false
	for _, info := range infos {
		name := info.Name()
		if strings.HasSuffix(name, ".xml") {
			t.Log("Processing " + name)
			data := loadTestData(t, name)
			nodes := loadTestXml(t, data)
			f(string(data), nodes)
			tested = true
		}
	}
	require.True(t, tested)
}

func TestLogValidate(t *testing.T) {
	processXMLFiles(t, func(_ string, nodes parser.Nodes) {
		var validateError error
		c := NewValidator(&validateError)
		c.Info = &Info{}
		require.NoError(t, ProcessXMLNodes(nodes, c))
		require.NoError(t, validateError)
	})
}

func TestLogReadAndWrite(t *testing.T) {
	processXMLFiles(t, func(data string, nodes parser.Nodes) {
		w := &XMLWriter{client.NewXMLWriter(), false}
		w.AddSpaces = false
		w.Open(nil)
		require.NoError(t, ProcessXMLNodes(nodes, w))
		w.Close()

		expected := util.FixLine(data)
		actual := util.FixLine(w.String())
		if !assert.Equal(t, expected, actual) {
			// cmp := expected + "\n" + actual
			// ioutil.WriteFile("cmp.txt", []byte(cmp), 0644)
			t.Fail()
		}
	})
}

func TestLogCheckOpenCondition(t *testing.T) {
	info, err := ParseLogInfo("2009061806gm-00a9-0000-6d13c207")
	require.NoError(t, err)
	x := NewValidator(nil)
	require.True(t, x.Open(info))
}
