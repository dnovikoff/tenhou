package tbase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestYakuParse(t *testing.T) {
	ret, err := ParseYakuList("35,6,52,2,54,1")
	require.NoError(t, err)
	assert.Equal(t, "YakuAkaDora: 1, YakuChinitsu: 6, YakuDora: 2", ret.String())
}

func TestYakuParseWrong(t *testing.T) {
	_, err := ParseYakuList("35,6,52,2,a54,1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")
	_, err = ParseYakuList("35,6,52,2,54,")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid syntax")
	_, err = ParseYakuList("35,6,52,2,154,1")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Yaku with value 154 not found in map")
}
