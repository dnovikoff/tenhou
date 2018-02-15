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
)

var addrFlag = flag.String("addr", ":10080", "listen addr")
var tenhouAddr = flag.String("remote-host", "133.242.10.78:10080", "tenhou flash client port")

func proxy(ctx context.Context, from, to network.XMLConnection, prefix string, logs chan string) {
	for {
		nctx, _ := context.WithTimeout(ctx, time.Second*60)
		message, err := from.Read(nctx)
		message = strings.Replace(message, `<HELO name="NoName" tid="f0" sx="M" />`, `<HELO name="NoName" tid="f0" sx="C" />`, 1)
		if !checkSuccess(err) {
			return
		}
		logs <- prefix + message
		err = to.Write(nctx, message)
		if !checkSuccess(err) {
			return
		}
	}
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
	client := network.NewXMLConnection(cConn)
	logs := make(chan string, 1024)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		proxy(ctx, server, client, "Send: ", logs)
		cancel()
	}()
	go func() {
		proxy(ctx, client, server, "Get: ", logs)
		cancel()
	}()
	for {
		select {
		case log := <-logs:
			fmt.Fprintln(file, log)
		case <-ctx.Done():
			return
		}
	}
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
