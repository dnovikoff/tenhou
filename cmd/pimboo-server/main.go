package main

import (
	"context"
	"flag"
	"log"

	"github.com/dnovikoff/tenhou/cmd/pimboo-server/game"
	"github.com/dnovikoff/tenhou/network"
)

var addrFlag = flag.String("addr", ":10080", "listen addr")

func checkError(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}

func main() {
	flag.Parse()
	ln := network.NewListener()
	ln.OnError = func(err error) {
		checkError(err)
	}
	ln.Handler = func(c network.XMLConnection) {
		log.Printf("New game connection")
		game := game.NewGame(c)
		game.Run()
		c.Close()
		log.Printf("Game connection closed")
	}
	log.Printf("Starting server on addr '%v'", *addrFlag)
	err := ln.ListenAndServe(context.Background(), "tcp", *addrFlag)
	checkError(err)
}
