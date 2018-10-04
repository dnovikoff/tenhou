package stats

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dnovikoff/tenhou/tools/utils"
	"github.com/spf13/cobra"
)

const (
	Location = "./tenhou/stats/"
)

func NewIndex() *FileIndex {
	return NewFileIndex(filepath.Join(Location, "index.json"))
}

func LoadIndex() (*FileIndex, error) {
	i := NewIndex()
	err := i.Load()
	if err != nil {
		return nil, err
	}
	return i, err
}

func CMD() *cobra.Command {
	rootCMD := &cobra.Command{
		Use:   "stats",
		Short: "Work with tenhou stat files",
	}
	initCMD := &cobra.Command{
		Use:   "init",
		Short: "Initialize database in current folder",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			index := NewIndex()
			err := index.Load()
			if os.IsNotExist(err) {
				utils.Check(index.Save())
				fmt.Println("Index initialized for ", Location)
				return
			}
			if err == nil {
				fmt.Println("Index already initialized. Delete " + Location + " to reinit.")
				os.Exit(1)
			}
		},
	}
	var repairFlag bool
	validateCMD := &cobra.Command{
		Use:   "validate",
		Short: "Validate index file, and existance of mentioned files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			index, err := LoadIndex()
			utils.Check(err)
			err = index.Validate()
			if err == nil {
				fmt.Println("Index is valid")
				return
			}
			if err != nil {
				if repairFlag {
					utils.Check(index.Save())
					fmt.Println("Index repaired")
				} else {
					utils.Check(err)
				}
			}
		},
	}
	validateCMD.Flags().BoolVar(&repairFlag, "repair", false, "Repair index if broken. Delete all records of not found files.")

	var interactiveFlag bool
	downloadCMD := &cobra.Command{
		Use:   "download",
		Short: "Download all stat files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			(&downloadCommand{interactive: interactiveFlag}).Run()
		},
	}
	downloadCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive downloader progress.")

	rootCMD.AddCommand(initCMD)
	rootCMD.AddCommand(validateCMD)
	rootCMD.AddCommand(downloadCMD)
	return rootCMD
}

type downloadCommand struct {
	interactive bool
	index       *FileIndex
}

func (d *downloadCommand) Init() {
	index, err := LoadIndex()
	utils.Check(err)
	fmt.Println("Validation index file")
	utils.Check(index.Validate())
	fmt.Println("Validation OK")
	d.index = index
}

func (d *downloadCommand) download(v string) {
	u := MakeFullURL(v)
	path := filepath.Join(Location, v)
	if d.index.Check(u) {
		fmt.Printf("Url '%v' already downloaded to '%v'\n", u, path)
		return
	}
	utils.Check(
		utils.NewDownloader(utils.AddTracker(
			utils.NewInteractiveTracker(u, path, d.interactive),
		)).WriteFile(u, path),
	)
	d.index.Add(u, path)
}

func (d *downloadCommand) Run() {
	d.Init()
	d.downloadMain()
	d.downloadList(ListOldURL)
	d.downloadList(ListURL)
}

func (d *downloadCommand) downloadList(url string) {
	lst := d.downloadString(url)
	items, err := ParseList(lst)
	utils.Check(err)
	fmt.Printf("Downloaded %v items from %v\n", len(items), url)
	for _, v := range items {
		d.download("dat/" + v.File)
	}
}

func (d *downloadCommand) downloadMain() {
	mainPage := d.downloadString(MainPageURL)
	files := ParseMain(mainPage)
	fmt.Printf("%v files parsed from main page\n", len(files))
	for _, v := range files {
		d.download(v)
	}
}

func (d *downloadCommand) downloadString(location string) string {
	return utils.MustDownload(location, utils.AddTracker(
		utils.NewInteractiveTracker(location, "", d.interactive)))
}
