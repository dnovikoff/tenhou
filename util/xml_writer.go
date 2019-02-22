package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type XMLWriter struct {
	buf       *bytes.Buffer
	AddSpaces bool
	Commit    func(string)
}

func NewXMLWriter() XMLWriter {
	return XMLWriter{&bytes.Buffer{}, true, nil}
}

func (w XMLWriter) Reset() {
	w.buf.Reset()
}

func (w XMLWriter) String() string {
	return w.buf.String()
}

func (w XMLWriter) Begin(tag string) XMLWriter {
	if w.AddSpaces && w.buf.Len() > 0 {
		w.buf.WriteByte(' ')
	}
	w.buf.WriteByte('<')
	w.buf.WriteString(tag)
	return w
}

func (w XMLWriter) End() {
	w.buf.WriteString("/>")
	if w.Commit != nil {
		w.Commit(w.String())
		w.Reset()
	}
}

func (w XMLWriter) AddTrailingSpace() XMLWriter {
	w.buf.WriteByte(' ')
	return w
}

func (w XMLWriter) Buffer() *bytes.Buffer {
	return w.buf
}

func (w XMLWriter) WriteListInt(key string, values []int) XMLWriter {
	x := make([]string, len(values))
	for k, v := range values {
		x[k] = strconv.Itoa(v)
	}
	return w.WriteList(key, x)
}

func (w XMLWriter) WriteListFloat(key string, values []tbase.Float) XMLWriter {
	x := make([]string, len(values))
	for k, v := range values {
		if v.IsInt {
			x[k] = strconv.Itoa(int(v.Value))
		} else {
			x[k] = fmt.Sprintf("%.2f", v.Value)
		}
	}
	return w.WriteList(key, x)
}

func (w XMLWriter) WriteList(key string, values []string) XMLWriter {
	if len(values) == 0 {
		return w
	}
	return w.WriteArg(key, strings.Join(values, ","))
}

func (w XMLWriter) WriteArg(key string, value string) XMLWriter {
	w.buf.WriteString(" " + key + `="` + value + `"`)
	return w
}

func (w XMLWriter) WriteFmtArg(key string, format string, args ...interface{}) XMLWriter {
	return w.WriteArg(key, fmt.Sprintf(format, args...))
}

func (w XMLWriter) WriteIntArg(key string, value int) XMLWriter {
	w.WriteArg(key, strconv.Itoa(value))
	return w
}

func (w XMLWriter) WriteInstance(key string, value tile.Instance) XMLWriter {
	return w.WriteIntArg(key, InstanceToTenhou(value))
}

func (w XMLWriter) WriteOpponent(key string, d tbase.Opponent) XMLWriter {
	return w.WriteIntArg(key, int(d))
}

func (w XMLWriter) WriteDealer(d tbase.Opponent) XMLWriter {
	return w.WriteOpponent("oya", d)
}

func (w XMLWriter) WriteWho(d tbase.Opponent) XMLWriter {
	return w.WriteOpponent("who", d)
}

func (w XMLWriter) WriteTableStatus(status tbase.TableStatus) XMLWriter {
	return w.WriteArg("ba", fmt.Sprintf("%d,%d", status.Honba, status.Sticks))
}

func (w XMLWriter) WriteScoreChanges(sc tbase.ScoreChanges) XMLWriter {
	return w.WriteArg("sc", changesToString(sc))
}

func (w XMLWriter) Write(format string, args ...interface{}) {
	fmt.Fprintf(w.buf, format, args...)
}

func (w XMLWriter) WriteBody(format string, args ...interface{}) {
	w.Begin("")
	w.Write(format, args...)
	w.End()
}

func changesToString(ch tbase.ScoreChanges) string {
	tmp := make([]string, len(ch))
	for k, v := range ch {
		tmp[k] = fmt.Sprintf("%d,%d", v.Score, v.Diff)
	}
	return strings.Join(tmp, ",")
}
