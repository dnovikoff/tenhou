package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/dnovikoff/tenhou/tools/tentool/logs"
	"github.com/dnovikoff/tenhou/tools/tentool/stats"
)

func main() {
	rootCmd := &cobra.Command{Use: "tentool"}
	rootCmd.AddCommand(stats.CMD(), logs.CMD())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
