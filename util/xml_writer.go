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

func (this XMLWriter) Reset() {
	this.buf.Reset()
}

func (this XMLWriter) String() string {
	return this.buf.String()
}

func (this XMLWriter) Begin(tag string) XMLWriter {
	if this.AddSpaces && this.buf.Len() > 0 {
		this.buf.WriteByte(' ')
	}
	this.buf.WriteByte('<')
	this.buf.WriteString(tag)
	return this
}

func (this XMLWriter) End() {
	this.buf.WriteString("/>")
	if this.Commit != nil {
		this.Commit(this.String())
		this.Reset()
	}
}

func (this XMLWriter) AddTrailingSpace() XMLWriter {
	this.buf.WriteByte(' ')
	return this
}

func (this XMLWriter) Buffer() *bytes.Buffer {
	return this.buf
}

func (this XMLWriter) WriteListInt(key string, values []int) XMLWriter {
	x := make([]string, len(values))
	for k, v := range values {
		x[k] = strconv.Itoa(v)
	}
	return this.WriteList(key, x)
}

func (this XMLWriter) WriteListFloat(key string, values []tbase.Float) XMLWriter {
	x := make([]string, len(values))
	for k, v := range values {
		if v.IsInt {
			x[k] = strconv.Itoa(int(v.Value))
		} else {
			x[k] = fmt.Sprintf("%.2f", v.Value)
		}
	}
	return this.WriteList(key, x)
}

func (this XMLWriter) WriteList(key string, values []string) XMLWriter {
	if len(values) == 0 {
		return this
	}
	return this.WriteArg(key, strings.Join(values, ","))
}

func (this XMLWriter) WriteArg(key string, value string) XMLWriter {
	this.buf.WriteString(" " + key + `="` + value + `"`)
	return this
}

func (this XMLWriter) WriteFmtArg(key string, format string, args ...interface{}) XMLWriter {
	return this.WriteArg(key, fmt.Sprintf(format, args...))
}

func (this XMLWriter) WriteIntArg(key string, value int) XMLWriter {
	this.WriteArg(key, strconv.Itoa(value))
	return this
}

func (this XMLWriter) WriteInstance(key string, value tile.Instance) XMLWriter {
	return this.WriteIntArg(key, InstanceToTenhou(value))
}

func (this XMLWriter) WriteOpponent(key string, d tbase.Opponent) XMLWriter {
	return this.WriteIntArg(key, int(d))
}

func (this XMLWriter) WriteDealer(d tbase.Opponent) XMLWriter {
	return this.WriteOpponent("oya", d)
}

func (this XMLWriter) WriteWho(d tbase.Opponent) XMLWriter {
	return this.WriteOpponent("who", d)
}

func (this XMLWriter) WriteTableStatus(status tbase.TableStatus) XMLWriter {
	return this.WriteArg("ba", fmt.Sprintf("%d,%d", status.Honba, status.Sticks))
}

func (this XMLWriter) WriteScoreChanges(sc tbase.ScoreChanges) XMLWriter {
	return this.WriteArg("sc", changesToString(sc))
}

func (this XMLWriter) Write(format string, args ...interface{}) {
	fmt.Fprintf(this.buf, format, args...)
}

func (this XMLWriter) WriteBody(format string, args ...interface{}) {
	this.Begin("")
	this.Write(format, args...)
	this.End()
}

func changesToString(ch tbase.ScoreChanges) string {
	tmp := make([]string, len(ch))
	for k, v := range ch {
		tmp[k] = fmt.Sprintf("%d,%d", v.Score, v.Diff)
	}
	return strings.Join(tmp, ",")
}
