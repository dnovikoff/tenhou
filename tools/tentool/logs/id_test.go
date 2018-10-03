package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseLogID(t *testing.T) {
	parsed := ParseID("2018092823gm-00a9-0000-2735720b")
	assert.Equal(t, &ParsedID{
		Time:       "2018092823gm",
		Type:       "00a9",
		Number:     "0000",
		OriginalID: "2735720b",
		DownloadID: "2735720b"},
		parsed)
	assert.Equal(t, "http://e.mjv.jp/0/log/?2018092823gm-00a9-0000-2735720b", GetDownloadLink(parsed))
	path, err := GetFilePath(parsed)
	require.NoError(t, err)
	assert.Equal(t, "00a9/0000/2018/09/28/23/2735720b", path)

}
