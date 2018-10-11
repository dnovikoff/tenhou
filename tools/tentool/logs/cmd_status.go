package logs

import (
	"fmt"

	"github.com/dnovikoff/tenhou/tools/utils"
)

func Status() {
	index, err := LoadIndex()
	utils.Check(err)
	total := index.Len()
	if total == 0 {
		fmt.Println("No logs downloaded")
		return
	}
	downloaded := 0
	for _, v := range index.data {
		if len(v) > 0 {
			downloaded++
		}
	}
	fmt.Printf("Downloaded %v ot of %v files (%v%%)\n", downloaded, total, downloaded*100/total)
}
