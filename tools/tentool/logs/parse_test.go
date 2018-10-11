package logs

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIDs(t *testing.T) {
	assert.Equal(t, []string{
		"2018092823gm-00a9-0000-2735720b",
		"2018092823gm-00a9-0000-3e233540",
		"2018092823gm-00a9-0000-445a08e4",
		"2018092823gm-00a9-0000-5fc2633e",
		"2018092823gm-00a9-0000-64b7a22b",
		"2018092823gm-00a9-0000-9b1b0e33",
		"2018092823gm-00a9-0000-c10a110c",
		"2018092823gm-00b9-0000-a2899a33",
		"2018092823gm-00b9-0000-c5512392",
		"2018092823gm-00e1-0000-af14e5c0"}, ParseIDs(loadTestData(t, "scc.html")))
}

func TestParseIDs2(t *testing.T) {
	assert.Equal(t, []string{
		"2018052202gm-0029-0000-4d37104e",
		"2018052202gm-0029-0000-5337edec",
		"2018052203gm-0029-0000-e57c0825",
	}, ParseIDs(`
		http://tenhou.net/0/?log=2018052202gm-0029-0000-5337edec&tw=0
http://tenhou.net/0/?log=2018052202gm-0029-0000-4d37104e&tw=2
http://tenhou.net/0/?log=2018052203gm-0029-0000-e57c0825&tw=0		
		`))
}

func loadTestData(t require.TestingT, filename string) string {
	data, err := ioutil.ReadFile("test_data/" + filename)
	require.NoError(t, err)
	return string(data)
}
