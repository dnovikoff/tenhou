package logs

import (
	"fmt"
	"io/ioutil"

	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

func Collect(args []string) {
	index, err := LoadIndex()
	utils.Check(err)
	newLinks := 0
	for _, v := range args {
		bytes, err := ioutil.ReadFile(v)
		utils.Check(err)
		links := ParseIDs(string(bytes))
		newLinks += index.AddIDs(links)
	}
	fmt.Printf("Found %v new links provided files. New database size: %v\n", newLinks, index.Len())
	utils.Check(index.Save())
}
