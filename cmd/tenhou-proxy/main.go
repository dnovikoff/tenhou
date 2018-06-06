package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/dnovikoff/tenhou/network"
	"github.com/dnovikoff/tenhou/util"
)

var addrFlag = flag.String("addr", ":10080", "listen addr")
var tenhouAddr = flag.String("remote-host", "133.242.10.78:10080", "tenhou flash client port")
var sexComputer = flag.Bool("sex-c", false, "change gender to C, send by client")

func proxy(ctx context.Context, from, to network.XMLConnection) {
	for proxyOne(ctx, from, to) {
	}
}

func proxyOne(ctx context.Context, from, to network.XMLConnection) bool {
	nctx, cancel := context.WithTimeout(ctx, time.Second*60)
	defer cancel()
	message, err := from.Read(nctx)
	if strings.HasPrefix(message, `<cross-domain-policy>`) {
		ports := "80,443,843,10080"
		w := util.NewXMLWriter()
		w.Write("<cross-domain-policy>")
		for _, domain := range []string{
			"*.mjv.jp",
			"*.tenhou.net",
			"localhost",
			// Add your domain here if needed
		} {
			w.Begin("allow-access-from").
				WriteArg("domain", domain).
				WriteArg("to-ports", ports).
				End()
		}
		w.Write("</cross-domain-policy>")
		message = w.String()
	} else if *sexComputer && strings.HasPrefix(message, `<HELO name="`) {
		sxStr := ` sx="`
		n := strings.Index(message, sxStr)
		if n > 0 {
			idx := n + len(sxStr)
			message = message[:idx] + "C" + message[idx+1:]
		}
	}

	if !checkSuccess(err) {
		return false
	}
	err = to.Write(nctx, message)
	if !checkSuccess(err) {
		return false
	}
	return true
}

func handle(sConn net.Conn, filename string) {
	defer sConn.Close()
	file, err := os.Create(filename)
	if !checkSuccess(err) {
		return
	}
	defer file.Close()
	cConn, err := net.Dial("tcp", *tenhouAddr)
	if !checkSuccess(err) {
		return
	}
	defer cConn.Close()

	server := network.NewXMLConnection(sConn)
	clientImpl := network.NewXMLConnection(cConn)
	client := network.NewXMLConnectionDebugger(clientImpl, network.NewMutexLogger(file))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		proxy(ctx, server, client)
		cancel()
	}()
	proxy(ctx, client, server)
	cancel()

}

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

func main() {
	flag.Parse()
	ln, err := net.Listen("tcp", *addrFlag)
	checkError(err)
	u1 := uuid.NewV4()
	preifx := hex.EncodeToString(u1.Bytes()[:4])
	log.Printf("Started server on addr '%v'. Sequence id is '%v'", *addrFlag, preifx)

	num := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		num++
		filename := fmt.Sprintf(preifx+"_%04d.log", num)
		log.Printf("File for new connection is '%v'", filename)
		go func(num int, filename string) {
			handle(conn, filename)
			log.Printf("Done with %d", num)
		}(num, filename)
	}
}
