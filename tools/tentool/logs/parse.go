package logs

import (
	"regexp"
	"sort"
)

var reg = regexp.MustCompile(`tenhou\.net/0/\?log=([0-9a-z-]+)`)

// ParseIDs parses stat file to find ids of logs
func ParseIDs(data string) []string {
	sub := reg.FindAllStringSubmatch(data, -1)
	res := make([]string, len(sub))
	for k, v := range sub {
		res[k] = v[1]
	}
	sort.Strings(res)
	return res
}
