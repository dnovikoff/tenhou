package utils

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectFilename(t *testing.T) {
	var handler func(w http.ResponseWriter, req *http.Request)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		handler(w, req)
	}))
	defer server.Close()
	c := server.Client()
	d := NewDownloader(Client(c))
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	t.Run("Disposition", func(t *testing.T) {
		handler = func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Disposition", `attachment;filename="foo.png"`)
			w.Write([]byte(`Some content`))
		}
		filename, err := d.Filename(ctx, server.URL+"/folder/other-folder/bar.jpg")
		require.NoError(t, err)
		assert.Equal(t, "foo.png", filename)
	})
	t.Run("NoDisposition", func(t *testing.T) {
		var userAgent string
		handler = func(w http.ResponseWriter, req *http.Request) {
			userAgent = req.UserAgent()
			w.Write([]byte(`Some content`))
		}
		filename, err := d.Filename(ctx, server.URL+"/folder/other-folder/bar.jpg")
		require.NoError(t, err)
		assert.Equal(t, "bar.jpg", filename)
		assert.Equal(t, "TenToolBot (+https://github.com/dnovikoff/tenhou/tools/tentool)", userAgent)
	})
}
