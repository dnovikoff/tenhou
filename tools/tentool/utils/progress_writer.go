package utils

import (
	"fmt"
	"io"
	"time"
)

type ProgressWriter struct {
	w        *InteractiveWriter
	total    int
	skipped  int
	progress int
	prefix   string
	eta      bool
	start    time.Time

	delay     time.Duration
	lastDelay time.Time
	disabled  bool
}

func NewProgressWriter(w io.Writer, prefix string, total int) *ProgressWriter {
	return &ProgressWriter{
		w:      NewInteractiveWriter(w),
		prefix: prefix,
		total:  total,
		start:  time.Now(),
	}
}

func (w *ProgressWriter) Disable() *ProgressWriter {
	w.disabled = true
	return w
}

func (w *ProgressWriter) Progress() int {
	return w.progress
}

func (w *ProgressWriter) SetETA() *ProgressWriter {
	w.eta = true
	return w
}

func (w *ProgressWriter) SetDelay(d time.Duration) *ProgressWriter {
	w.delay = d
	return w
}

func (w *ProgressWriter) Skip() {
	w.Inc()
	w.skipped++
}

func (w *ProgressWriter) Inc() {
	w.Advance(1)
}

func (w *ProgressWriter) Advance(x int) {
	w.progress += x
}

func (w *ProgressWriter) timeLeft() time.Duration {
	currentTime := time.Now()
	elapsed := currentTime.Sub(w.start)
	itemsLeft := w.total - w.progress
	var speed float64
	nanos := elapsed.Nanoseconds()
	if nanos != 0 {
		speed = float64(elapsed.Nanoseconds()) / float64(w.progress-w.skipped)
	}
	left := time.Nanosecond * time.Duration(speed*float64(itemsLeft))
	left = left.Truncate(time.Second)
	return left
}

func (w *ProgressWriter) timeString() string {
	if !w.eta || w.progress == 0 {
		return ""
	}
	return fmt.Sprintf(" Time left: %v", w.timeLeft())
}

func (w *ProgressWriter) Start() {
	if w.disabled {
		return
	}
	w.display()
}

func (w *ProgressWriter) Display() {
	if w.disabled {
		return
	}
	if w.delay != 0 {
		now := time.Now()
		if !w.lastDelay.Add(w.delay).Before(now) {
			return
		}
		w.lastDelay = now
	}
	w.display()
}

func (w *ProgressWriter) display() {
	w.w.Printf(w.prefix+" %v/%v (%v%%)"+w.timeString(), w.progress, w.total, w.progress*100/w.total)
}

func (w *ProgressWriter) Done() {
	if w.disabled {
		return
	}
	w.display()
	w.w.Println()
}
