package logs

import (
	"fmt"
	"strings"

	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

func RemoveIndex(prefix string) {
	index, err := LoadIndex()
	utils.Check(err)
	removed := 0

	index.data.Files = removeIndex(index.data.Files, prefix, &removed)
	for _, v := range index.data.Zips {
		v.Files = removeIndex(v.Files, prefix, &removed)
	}
	fmt.Printf("Removed %v links. New database size: %v\n", removed, index.Len())
	utils.Check(index.Save())
}

func removeIndex(infos []*FileInfo, prefix string, removed *int) []*FileInfo {
	next := make([]*FileInfo, 0, len(infos))
	for _, v := range infos {
		if !strings.HasPrefix(v.ID, prefix) {
			next = append(next, v)
		} else {
			*removed++
		}
	}
	return next
}
