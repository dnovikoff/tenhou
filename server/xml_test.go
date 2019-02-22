package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinWrite(t *testing.T) {
	c := NewXMLWriter()
	c.CancelJoin()
	c.Join(0, 9, false)
	c.Join(1, 2, true)

	assert.Equal(t, `<JOIN /> <JOIN t="0,9" /> <JOIN t="1,2,r" />`, c.String())
	c2 := NewXMLWriter()
	require.NoError(t, ProcessXMLMessage(c.String(), c2))
	assert.Equal(t, `<JOIN /> <JOIN t="0,9" /> <JOIN t="1,2,r" />`, c2.String())
}
