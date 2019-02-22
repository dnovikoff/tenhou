package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnovikoff/tenhou/parser"
)

func TestParseXML(t *testing.T) {
	ret, err := parser.ParseXML(`<N type="1" hai0="132" hai1="134" />`)
	require.NoError(t, err)
	require.Equal(t, 1, len(ret))
	assert.Equal(t, parser.Node{
		Name: "N",
		Attributes: map[string]string{
			"type": "1",
			"hai0": "132",
			"hai1": "134",
		},
	}, ret[0])
}

func TestParseXMLList(t *testing.T) {
	ret, err := parser.ParseXML(`<A/><B/>`)
	require.NoError(t, err)
	assert.Equal(t, parser.Nodes{
		parser.Node{Name: "A"},
		parser.Node{Name: "B"},
	}, ret)
}

func TestParseXMLStrictCheck(t *testing.T) {
	ret, err := parser.ParseXML(`<HELO ratingscale="PF3=1.000000&PF4=1.000000&PF01C=0.582222&PF02C=0.501632&PF03C=0.414869&PF11C=0.823386&PF12C=0.709416&PF13C=0.586714&PF23C=0.378722&PF33C=0.535594&PF1C00=8.000000"/>`)
	require.NoError(t, err)
	require.Equal(t, 1, len(ret))
	assert.Equal(t, parser.Node{
		Name: "HELO",
		Attributes: map[string]string{
			"ratingscale": "PF3=1.000000&PF4=1.000000&PF01C=0.582222&PF02C=0.501632&PF03C=0.414869&PF11C=0.823386&PF12C=0.709416&PF13C=0.586714&PF23C=0.378722&PF33C=0.535594&PF1C00=8.000000",
		},
	}, ret[0])
}
