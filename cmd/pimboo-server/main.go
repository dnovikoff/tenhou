package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/dnovikoff/tenhou/cmd/pimboo-server/game"
	"github.com/dnovikoff/tenhou/network"
)

var addrFlag = flag.String("addr", ":10080", "listen addr")

func checkSuccess(err error) bool {
	if err == nil {
		return true
	}
	log.Printf("Error: %v", err)
	return false
}

func checkError(err error) {
	if err == nil {
		return
	}
	log.Fatal(err)
}

func handle(sConn net.Conn) {
	defer sConn.Close()
	handleImpl(sConn)
}

func handleImpl(sConn net.Conn) {
	impl := network.NewXMLConnection(sConn)
	logger := func(format string, args ...interface{}) {
		fmt.Fprintf(os.Stdout, format+"\n", args...)
	}
	con := network.NewXMLConnectionDebugger(impl, logger)

	parentCtx, _ := context.WithCancel(context.Background())
	ctx, _ := context.WithTimeout(parentCtx, time.Second*10)
	str, err := con.Read(ctx)
	if !checkSuccess(err) {
		return
	}
	if str == "<policy-file-request/>" {
		err = con.Write(ctx, `<cross-domain-policy><allow-access-from domain="*.mjv.jp" to-ports="80,843,10080" /><allow-access-from domain="*.tenhou.net" to-ports="80,843,10080" /></cross-domain-policy>`)
		checkSuccess(err)
		return
	}
	if str == "<GETSWF />" {
		err = con.Write(ctx, `<SWF src="0/app/1430/welcome.swf" />`)
		checkSuccess(err)
		return
	}
	game := game.NewGame(con)
	game.Run()
}

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", *addrFlag)
	checkError(err)
	log.Printf("Started server on addr '%v'", *addrFlag)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handle(conn)
	}
}
