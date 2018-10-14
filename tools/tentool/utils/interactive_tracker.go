package utils

import (
	"os"

	"github.com/c2h5oh/datasize"
)

type InteractiveTracker struct {
	u, p        string
	total       int64
	progress    int64
	interactive bool
	w           *InteractiveWriter
}

func NewInteractiveTracker(u, p string, i bool) *InteractiveTracker {
	return &InteractiveTracker{u: u, p: p, interactive: i, w: NewInteractiveWriter(os.Stdout)}
}

func (t *InteractiveTracker) SetPath(p string) {
	t.p = p
}

func (t *InteractiveTracker) Start(total int64) {
	t.total = total
	t.writeHeader()
}

func (t *InteractiveTracker) writeHeader() {
	if t.p == "" {
		t.w.Printf("Downloading %v", t.u)
	} else {
		t.w.Printf("Downloading %v to %v", t.u, t.p)
	}
	t.w.Println()
}

func (t *InteractiveTracker) Write(bytes int) {
	if !t.interactive {
		return
	}
	t.progress += int64(bytes)
	if t.total == 0 {
		t.w.Printf("%v of unknown total size", datasize.ByteSize(t.progress))
	} else {
		t.w.Printf("%s/%s (%v%%)",
			datasize.ByteSize(t.progress).HumanReadable(),
			datasize.ByteSize(t.total).HumanReadable(),
			t.progress*100/t.total,
		)
	}
}

func (t *InteractiveTracker) Done(total int64, err error) {
	if !t.interactive {
		return
	}
	t.w.Println()
}
