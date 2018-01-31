package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLobbyDecode(t *testing.T) {
	decode := func(in string) string {
		res, err := DecodeLobby(in)
		require.NoError(t, err)
		return res.DebugString()
	}
	assert.Equal(t, "AKT4fX", decode("00a1"))
	assert.Equal(t, "AKH4fX", decode("00a9"))
	assert.Equal(t, "AKT4F0", decode("0841"))
}

func TestInfoParse(t *testing.T) {
	parse := func(in string) string {
		ret, err := ParseLogInfo(in)
		require.NoError(t, err)
		return ret.DebugString()
	}

	assert.Equal(t, "[2009-06-18 06:00:00 +0000 UTC][AKT4fX][0][6d13c207]", parse("2009061806gm-00a1-0000-6d13c207"))
	assert.Equal(t, "[2009-03-17 02:00:00 +0000 UTC][AKT4fX][0][336ee82a]", parse("00a1/2009/03/17/2009031702gm-00a1-0000-336ee82a.mjlog"))
	assert.Equal(t, "[2009-09-15 15:00:00 +0000 UTC][AKH4fX][0][1753da1d]", parse("http://tenhou.net/0/?log=2009091515gm-00a9-0000-1753da1d"))
}

func TestInfoParsefilename(t *testing.T) {
	ret, err := ParseLogInfo("http://tenhou.net/0/?log=2017022002gm-000b-2873-xc0294a92421e&tw=3")
	require.NoError(t, err)
	assert.Equal(t, "http://tenhou.net/0/?log=2017022002gm-000b-2873-xc0294a92421e", ret.LogUrl)
	assert.Equal(t, "http://e.mjv.jp/0/log/plainfiles.cgi?2017022002gm-000b-2873-fb28df66", ret.XmlUrl)
}
