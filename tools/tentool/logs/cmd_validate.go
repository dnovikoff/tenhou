package logs

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"time"

	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/log"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

type validate struct {
	readFlag   bool
	indexFlag  bool
	repairFlag bool
	parseFlag  bool
	index      *FileIndex
}

func (v *validate) loadIndex() *FileIndex {
	if v.index == nil {
		index, err := LoadIndex()
		utils.Check(err)
		v.index = index
	}
	return v.index
}

func (v *validate) Run() {
	v.readFlag = v.readFlag || v.parseFlag
	if v.indexFlag {
		v.stepIndex()
	}
	if v.readFlag {
		v.stepParse()
	}
}

func (v *validate) stepIndex() {
	index := v.loadIndex()
	removed, err := index.Validate()
	indexFixed := false
	utils.Check(err)
	if len(removed) == 0 {
		fmt.Println("Index is valid")
	} else {
		indexFixed = true
	}
	infos := index.ValidateNames()
	if len(infos) > 0 {
		fmt.Printf("%v incorrectly donwloaded logs (detected by names)\n", len(infos))
		bytes, err := json.MarshalIndent(infos, "", " ")
		utils.Check(err)
		utils.Check(ioutil.WriteFile(fileName("names.json"), bytes, 0644))
		indexFixed = true
	} else {
		fmt.Printf("All names are correct\n")
	}
	if len(removed) > 0 {
		fmt.Printf("%v index records to remove:\n", len(removed))
		for _, v := range removed {
			fmt.Println(v)
		}
	}
	if v.repairFlag && indexFixed {
		utils.Check(index.Save())
		fmt.Println("Index repaired")
	}
}

type errorForFile struct {
	name string
	err  error
}

func (e *errorForFile) Error() string {
	return e.name + ": " + e.err.Error()
}

func wrapVisit(f func(*FileInfo, Opener) error) func(*FileInfo, Opener) error {
	return func(i *FileInfo, o Opener) error {
		err := f(i, o)
		if err == nil {
			return nil
		}
		return &errorForFile{i.File, err}
	}
}

type nameExtractor struct {
	log.NullController
	names []string
}

func fixNames(x []string) []string {
	if len(x) == 0 {
		return nil
	}
	if x[len(x)-1] != "" {
		return x
	}
	return x[:len(x)-1]
}

func (e *nameExtractor) UserList(ul client.UserList) {
	e.names = fixNames(ul.Users.GetNames())
}

func parseNames(data []byte) ([]string, error) {
	c := &nameExtractor{}
	x := &parser.Root{}
	err := xml.Unmarshal(data, &x)
	if err != nil {
		return nil, err
	}
	err = log.ProcessXMLNodes(x.Nodes, c)
	if err != nil {
		return nil, err
	}
	return c.names, nil
}

type taskContext struct {
	pool    chan *bytes.Buffer
	updated int
	index   *FileIndex
}

type task struct {
	info  *FileInfo
	buf   *bytes.Buffer
	names []string
	ctx   *taskContext
}

func (t *task) Run() error {
	names, err := parseNames(t.buf.Bytes())
	t.done()
	t.names = names
	return err
}

func (t *task) done() {
	t.buf.Reset()
	t.ctx.pool <- t.buf
}

func (t *task) Save(blocked bool) error {
	t.info.LogNames = t.names
	t.ctx.updated++
	if !blocked || t.ctx.updated < 10000 {
		return nil
	}
	t.ctx.updated = 0
	return t.ctx.index.Save()
}

func (v *validate) stepParse() {
	index := v.loadIndex()
	total := v.index.Len()
	w := utils.NewProgressWriter(os.Stdout, "Processing", total).
		SetETA().
		SetDelay(time.Millisecond * 300)
	w.Start()
	cnt := runtime.NumCPU()
	ctx := &taskContext{
		pool:  make(chan *bytes.Buffer, cnt),
		index: index,
	}
	for i := 0; i < cnt; i++ {
		ctx.pool <- &bytes.Buffer{}
	}
	s := utils.Scheduler{}
	s.Start(cnt, 100000)
	utils.Check(index.Visit(wrapVisit(func(i *FileInfo, o Opener) error {
		w.Display()
		if s.Error() != nil {
			return s.Stop()
		}
		if i.LogNames != nil {
			w.Skip()
			return nil
		}
		w.Inc()
		_, err := log.ParseLogInfo(i.File)
		if err != nil {
			return err
		}
		t := &task{
			info: i,
			buf:  <-ctx.pool,
			ctx:  ctx,
		}
		r, err := o.Open()
		if err != nil {
			if v.repairFlag && !i.IsInsideZip() {
				fmt.Println("Removing file " + i.File + " because of error: " + err.Error())
				os.Remove(i.File)
				return nil
			}
			return err
		}
		_, err = io.Copy(t.buf, r)
		if err != nil {
			r.Close()
			if v.repairFlag && !i.IsInsideZip() {
				fmt.Println("Removing file " + i.File + " because of error: " + err.Error())
				os.Remove(i.File)
				return nil
			}
			return err
		}
		err = r.Close()
		if err != nil {
			return err
		}
		if v.parseFlag {
			s.Push(t)
		} else {
			t.done()
		}
		return nil
	})))
	w.Done()
	utils.Check(s.Stop())
	if ctx.updated > 0 {
		utils.Check(index.Save())
	}
}
