package logs

import (
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	XmlUrlPrefix = "http://e.mjv.jp/0/log/?"
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

type ParsedID struct {
	// Year,month,day,hour
	Time       string
	Type       string
	Number     string
	OriginalID string
	DownloadID string
}

func GetFilePath(id *ParsedID) (string, error) {
	time, err := time.Parse("2006010215gm", id.Time)
	if err != nil {
		return "", err
	}
	formatted := time.Format("2006/01/02/15")
	return path.Join(id.Type, id.Number, formatted, id.OriginalID), nil
}

func GetDownloadLink(id *ParsedID) string {
	return strings.Join(
		[]string{XmlUrlPrefix +
			id.Time,
			id.Type,
			id.Number,
			id.DownloadID}, "-")
}

func ParseID(input string) *ParsedID {
	groups := strings.Split(input, "-")
	if len(groups) != 4 {
		return nil
	}
	result := &ParsedID{
		Time:       groups[0],
		Type:       groups[1],
		Number:     groups[2],
		OriginalID: groups[3],
	}
	if result.OriginalID[0] != 'x' {
		result.DownloadID = result.OriginalID
		return result
	}
	id := result.OriginalID[1:]
	var a, b, c uint64
	fmt.Sscanf(id, "%04x%04x%04x", &a, &b, &c)
	index := 0
	if input > "2010041111gm" {
		x, _ := strconv.ParseUint("3"+input[4:10], 10, 64)
		y, _ := strconv.ParseUint(input[9:10], 10, 64)
		index = int(x % (33 - y))
	}
	first := (a ^ b ^ magic[index]) & 0xFFFF
	second := (b ^ c ^ magic[index] ^ magic[index+1]) & 0xFFFF
	result.DownloadID = fmt.Sprintf("%04x%04x", first, second)
	return result
}
