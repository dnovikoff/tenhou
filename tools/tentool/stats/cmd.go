package stats

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/spf13/cobra"

	"github.com/dnovikoff/tenhou/tools/tentool/utils"
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
	var cleanup bool
	downloadCMD := &cobra.Command{
		Use:   "download",
		Short: "Download all stat files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			(&downloadCommand{interactive: interactiveFlag, cleanup: cleanup}).Run()
		},
	}
	downloadCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive downloader progress.")
	downloadCMD.Flags().BoolVar(&cleanup, "cleanup", false, "Clean old files after download.")

	yadiskCMD := &cobra.Command{
		Use:   "yadisk",
		Short: "Download year stats archives from yadisk. That might be faster.",
		Run: func(cmd *cobra.Command, args []string) {
			(&yadisk{
				interactive: interactiveFlag,
			}).Run(args)
		},
	}
	yadiskCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive downloader progress.")

	rootCMD.AddCommand(initCMD)
	rootCMD.AddCommand(validateCMD)
	rootCMD.AddCommand(downloadCMD)
	rootCMD.AddCommand(yadiskCMD)
	return rootCMD
}

type downloadCommand struct {
	interactive bool
	cleanup     bool
	index       *FileIndex
	cleanIndex  *FileIndex
}

func (d *downloadCommand) Init() {
	index, err := LoadIndex()
	utils.Check(err)
	fmt.Println("Validation index file")
	utils.Check(index.Validate())
	fmt.Println("Validation OK")
	d.index = index
	d.cleanIndex = NewIndex()
}

func (d *downloadCommand) download(v string) {
	u := MakeFullURL(v)
	if v, ok := d.index.data[u]; ok {
		d.cleanIndex.data[u] = v
	}
	path := filepath.Join(Location, v)
	if d.index.Check(u) {
		fmt.Printf("Url '%v' already downloaded to '%v'\n", u, path)
		return
	}
	exits, err := utils.FileExists(path)
	utils.Check(err)
	if exits {
		fmt.Printf("File '%v' already downloaded. Adding to index. \n", path)
	} else {
		utils.Check(
			utils.NewDownloader(utils.AddTracker(
				utils.NewInteractiveTracker(u, path, d.interactive),
			)).WriteFile(context.TODO(), u, path),
		)
	}
	d.cleanIndex.JustAdd(u, path)
	utils.Check(d.index.Add(u, path))
}

func (d *downloadCommand) Run() {
	d.Init()
	d.downloadMain()
	d.downloadList(ListOldURL)
	d.downloadList(ListURL)

	unused := d.findUnusedFiles()
	if len(unused) > 0 {
		fmt.Printf("There are %v (of total %v) unused files.\n", len(unused), len(d.index.data))
		if !d.cleanup {
			fmt.Printf("Run with --cleanup to remove them\n")
		} else {
			fmt.Printf("Running cleanup.\n")
			for _, v := range unused {
				fmt.Printf("Removing file %v\n", v)
				utils.Check(os.Remove(v))
			}
			d.index.data = d.cleanIndex.data
			utils.Check(d.index.Save())
			fmt.Printf("Cleanup done.\n")
		}
	}
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

func (d *downloadCommand) findUnusedFiles() []string {
	usedFiles := make(map[string]bool, len(d.index.data))
	for k, v := range d.index.data {
		if _, found := d.cleanIndex.data[k]; found {
			usedFiles[v] = true
		} else {
			if !usedFiles[v] {
				usedFiles[v] = false
			}
		}
	}
	var unusedFiles []string
	for k, v := range usedFiles {
		if !v {
			unusedFiles = append(unusedFiles, k)
		}
	}
	sort.Strings(unusedFiles)
	return unusedFiles
}

func (d *downloadCommand) downloadString(location string) string {
	return utils.MustDownload(context.TODO(), location, utils.AddTracker(
		utils.NewInteractiveTracker(location, "", d.interactive)))
}
