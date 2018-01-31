package tbase

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLobbyEncode(t *testing.T) {
	assert.Equal(t, "0841", RulesDzjanso.String())
}
