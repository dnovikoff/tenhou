package stats

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

var reg = regexp.MustCompile(`<a href="(scraw[0-9]+\.zip)">`)

// ParseMain parses tenhou main page http://tenhou.net/sc/raw/
// and returns list of year index archives. Like scraw2017.zip
func ParseMain(data string) []string {
	sub := reg.FindAllStringSubmatch(data, -1)
	res := make([]string, len(sub))
	for k, v := range sub {
		res[k] = v[1]
	}
	sort.Strings(res)
	return res
}

type ListItem struct {
	File string `yaml:"file"`
	Size int    `yaml:"size"`
}

func MustParseList(data string) []ListItem {
	out, err := ParseList(data)
	if err != nil {
		panic(err)
	}
	return out
}

func ParseList(data string) ([]ListItem, error) {
	var out []ListItem
	first := strings.Index(data, "[")
	last := strings.LastIndex(data, "]")
	if first == -1 || last == -1 {
		return nil, fmt.Errorf("Incorrect string")
	}
	data = data[first : last+1]
	data = strings.Replace(data, ":", ": ", -1)
	err := yaml.Unmarshal([]byte(data), &out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
