package log

import (
	"encoding/xml"
	"io/ioutil"
	"testing"

	"github.com/dnovikoff/tenhou/parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXML(t *testing.T) {
	data, err := ioutil.ReadFile("test_data/example.xml")
	require.NoError(t, err)
	x := &parser.Root{}
	err = xml.Unmarshal(data, &x)
	require.NoError(t, err)

	assert.NoError(t, ProcessXMLNodes(x.Nodes, NullController{}))
}
