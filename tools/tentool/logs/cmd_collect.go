package logs

import (
	"fmt"
	"io/ioutil"

	"github.com/dnovikoff/tenhou/tools/utils"
)

func Collect(args []string) {
	index, err := LoadIndex()
	utils.Check(err)
	originalLen := index.Len()
	for _, v := range args {
		bytes, err := ioutil.ReadFile(v)
		utils.Check(err)
		links := ParseIDs(string(bytes))
		index.Add(links)
	}
	newLen := index.Len()
	fmt.Printf("Found %v new links provided files. New database size: %v\n", newLen-originalLen, newLen)
	utils.Check(index.Save())
}
