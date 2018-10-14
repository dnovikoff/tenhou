package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"

	"github.com/spf13/cobra"

	"github.com/dnovikoff/tenhou/tools/tentool/logs"
	"github.com/dnovikoff/tenhou/tools/tentool/stats"
)

func main() {
	cpuProfile := os.Getenv("CPUPROFILE")
	httpProfile := os.Getenv("HTTPPROFILE")
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer func() {
			pprof.StopCPUProfile()
			f.Close()
		}()
	}
	if httpProfile != "" {
		go func() {
			log.Println(http.ListenAndServe(httpProfile, nil))
		}()
	}
	rootCmd := &cobra.Command{Use: "tentool"}
	rootCmd.AddCommand(stats.CMD(), logs.CMD())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
