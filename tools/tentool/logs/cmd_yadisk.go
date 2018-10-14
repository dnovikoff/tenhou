package logs

import (
	"github.com/dnovikoff/tenhou/tools/tentool/utils"
)

type yadisk struct {
	interactive bool
	index       *FileIndex
}

func (y *yadisk) Run(args []string) {
	index, err := LoadIndex()
	utils.Check(err)
	y.index = index
	if len(args) == 0 {
		args = []string{"https://yadi.sk/d/FIIkaucSNjR3Kw"}
	}
	for _, v := range args {
		utils.Check(y.download(v))
	}
}

func (y *yadisk) download(u string) error {
	return utils.YaDiskDownload(u, Location, y.interactive, func(_, p string) error {
		return addZipToIndex(y.index, p)
	})
}
