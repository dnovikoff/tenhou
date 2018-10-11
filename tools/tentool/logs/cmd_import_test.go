package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	f := func(x string) string {
		r, err := makeZipPath(x)
		require.NoError(t, err)
		return r
	}
	assert.Equal(t, "hello.zip", f("tenhou/logs/hello.zip"))
	assert.Equal(t, "hello.zip", f("./tenhou/logs/hello.zip"))
	assert.Equal(t, "subfolder/hello.zip", f("./tenhou/logs/subfolder/hello.zip"))
	assert.Equal(t, "../hello.zip", f("./tenhou/hello.zip"))
	assert.Equal(t, "../../hello.zip", f("hello.zip"))
	assert.Equal(t, "/root/hello.zip", f("/root/hello.zip"))
}
