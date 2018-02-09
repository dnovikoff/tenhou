package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/cmd/pimboo-server/game"
	"github.com/dnovikoff/tenhou/network"
	"github.com/dnovikoff/tenhou/server"
	"github.com/dnovikoff/tenhou/tbase"
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
	con := network.NewXMLConnectionDebugger(impl, os.Stdout)

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
	cb := &server.Callbacks{}
	game := game.NewGame(con)
	cb.CbHello = func(name string, tid string, sex tbase.Sex) {
		game.Client.Hello(name, "20180117-e7b5e83e", client.HelloStats{})
		game.Send()
	}
	if !game.ProcessOne(cb) {
		return
	}
	cb.CbHello = nil
	cb.CbAuth = func(string) {}
	if !game.ProcessOne(cb) {
		return
	}
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
