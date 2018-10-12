package logs

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dnovikoff/tenhou/tools/tentool/stats"
	"github.com/dnovikoff/tenhou/tools/utils"
	"github.com/spf13/cobra"
)

const (
	Location = "./tenhou/logs/"
)

func NewIndex() *FileIndex {
	return NewFileIndex(filepath.Join(Location, "index.json"))
}

func NewParsedIndex() *stats.FileIndex {
	return stats.NewFileIndex(filepath.Join(Location, "parsed.json"))
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
		Use:   "logs",
		Short: "Work with tenhou logs",
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
				fmt.Println("Index initialized for ", index.Path)
				return
			}
			if err == nil {
				fmt.Println("Index already initialized. Delete " + index.Path + " to reinit.")
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
	updateCMD := &cobra.Command{
		Use:   "update",
		Short: "Update links index from downloaded stat files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			(&updater{interactive: interactiveFlag}).Run()
		},
	}
	updateCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive progress.")

	var parallel int
	downloadCMD := &cobra.Command{
		Use:   "download",
		Short: "Download log files",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			(&downloader{
				interactive: interactiveFlag,
				parallel:    parallel,
			}).Run()
		},
	}
	downloadCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive downloader progress.")
	downloadCMD.Flags().IntVar(&parallel, "parallel", 1, "Number of parallel downloads.")

	importCMD := &cobra.Command{
		Use:   "import",
		Short: "Import zip files with logs",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			(&importer{
				interactive: interactiveFlag,
			}).Run(args)
		},
	}
	importCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive progress.")

	var dryFlag bool
	var forceFlag bool
	exportCMD := &cobra.Command{
		Use:   "export",
		Short: "Export zip files by config",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			(&exporter{
				interactive: interactiveFlag,
				dry:         dryFlag,
				force:       forceFlag,
			}).Run(args)
		},
	}
	exportCMD.Flags().BoolVar(&interactiveFlag, "interactive", true, "Use interactive export progress.")
	exportCMD.Flags().BoolVar(&dryFlag, "dry", false, "Do not create files. Just show expected results.")
	exportCMD.Flags().BoolVar(&forceFlag, "force", false, "Rewrite already existing files.")

	collectCMD := &cobra.Command{
		Use:   "collect",
		Short: "Collect links from text files",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			Collect(args)
		},
	}

	statusCMD := &cobra.Command{
		Use:   "status",
		Short: "Show downloaded status",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			Status()
		},
	}

	rootCMD.AddCommand(initCMD)
	rootCMD.AddCommand(validateCMD)
	rootCMD.AddCommand(updateCMD)
	rootCMD.AddCommand(downloadCMD)
	rootCMD.AddCommand(importCMD)
	rootCMD.AddCommand(exportCMD)
	rootCMD.AddCommand(collectCMD)
	rootCMD.AddCommand(statusCMD)
	return rootCMD
}
