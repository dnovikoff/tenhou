package log

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dnovikoff/tenhou/tbase"

	"github.com/facebookgo/stackerr"
)

const (
	LogUrlPrefix = "http://tenhou.net/0/?log="
	XmlUrlPrefix = "http://e.mjv.jp/0/log/plainfiles.cgi?"
)

var magic []uint64 = []uint64{
	22136, 52719, 55146, 42104,
	59591, 46934, 9248, 28891,
	49597, 52974, 62844, 4015,
	18311, 50730, 43056, 17939,
	64838, 38145, 27008, 39128,
	35652, 63407, 65535, 23473,
	35164, 55230, 27536, 4386,
	64920, 29075, 42617, 17294,
	18868, 2081}

type Info struct {
	Time      time.Time
	Rules     tbase.LobbyRules
	Lobby     int
	Id        string
	FixedId   string
	FixedName string
	FullName  string
	LogUrl    string
	XmlUrl    string
}

func filenameFix(filename string) (fixedId string, fixedName string) {
	idx := strings.LastIndex(filename, "-")
	if idx == -1 {
		return
	}
	id := filename[idx+1:]
	if id[0] != 'x' {
		return id, filename
	}
	id = id[1:]
	var a, b, c uint64
	fmt.Sscanf(id, "%04x%04x%04x", &a, &b, &c)
	index := 0
	if filename > "2010041111gm" {
		x, _ := strconv.ParseUint("3"+filename[4:10], 10, 64)
		y, _ := strconv.ParseUint(filename[9:10], 10, 64)
		index = int(x % (33 - y))
	}
	prefix := filename[:idx+1]

	first := (a ^ b ^ magic[index]) & 0xFFFF
	second := (b ^ c ^ magic[index] ^ magic[index+1]) & 0xFFFF
	fixedId = fmt.Sprintf("%04x%04x", first, second)
	fixedName = prefix + fixedId
	return
}

func (i *Info) DebugString() string {
	return fmt.Sprintf("[%v][%v][%v][%v]",
		i.Time,
		i.Rules.DebugString(),
		i.Lobby,
		i.Id)
}

func extractId(in string) string {
	first := strings.LastIndexByte(in, '/')
	last := strings.LastIndexByte(in, '.')
	if last == -1 {
		last = len(in)
	}
	if first == -1 {
		first = 0
	} else {
		first++
	}
	if last < first {
		u, err := url.Parse(in)
		if err != nil {
			return "not parsed"
		}
		values, _ := url.ParseQuery(u.RawQuery)
		x := values["log"]
		if len(x) != 1 {
			return "len is not 1"
		}
		return x[0]
	}
	return in[first:last]
}

//2009061806gm-00a1-0000-6d13c207
///2009/03/17/2009031702gm
func ParseLogInfo(fileName string) (ret *Info, err error) {
	fileName = extractId(fileName)
	groups := strings.Split(fileName, "-")
	if len(groups) != 4 {
		err = stackerr.Newf("'%v' should contain 4 groups", fileName)
		return
	}
	info := &Info{}
	info.Id = groups[3]
	info.Lobby, err = strconv.Atoi(groups[2])
	if err != nil {
		err = stackerr.Wrap(err)
		return
	}
	info.Rules, err = DecodeLobby(groups[1])
	if err != nil {
		return
	}
	info.Time, err = time.Parse("2006010215gm", groups[0])
	if err != nil {
		err = stackerr.Wrap(err)
		return
	}
	info.FullName = fileName
	info.LogUrl = LogUrlPrefix + fileName
	fixedId, fixedName := filenameFix(fileName)
	info.XmlUrl = XmlUrlPrefix + fixedName
	info.FixedId = fixedId
	info.FixedName = fixedName
	ret = info
	return
}

func DecodeLobby(in string) (res tbase.LobbyRules, err error) {
	i64, err := strconv.ParseInt(in, 16, 0)
	if err != nil {
		err = stackerr.Wrap(err)
		return
	}
	res = tbase.LobbyRules(i64)
	return
}
