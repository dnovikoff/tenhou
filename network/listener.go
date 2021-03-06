package network

import (
	"context"
	"net"
	"sync"
	"time"
)

const DefaultPolicy = `<cross-domain-policy><allow-access-from domain="*.mjv.jp" to-ports="80,843,10080" /><allow-access-from domain="*.tenhou.net" to-ports="80,843,10080" /></cross-domain-policy>`
const DefaultSWF = `<SWF src="0/app/1430/welcome.swf" />`
const DefaultNetwork = "tcp"
const DefaultAddress = ":10080"

type Listener struct {
	Policy  string
	SWF     string
	OnError func(err error)
	Handler func(XMLConnection)
}

func NewListener() *Listener {
	return &Listener{
		Policy: DefaultNetwork,
		SWF:    DefaultSWF,
	}
}

func (lis *Listener) checkError(err error) bool {
	if err == nil || lis.OnError == nil {
		return true
	}
	lis.OnError(err)
	return false
}

func (lis *Listener) handle(parentCtx context.Context, sConn net.Conn) {
	con := NewXMLConnection(sConn)
	ctx, cancel := context.WithTimeout(parentCtx, time.Second*10)
	defer cancel()
	str, err := con.Read(ctx)
	if !lis.checkError(err) {
		return
	}
	if str == "<policy-file-request/>" {
		err = con.Write(ctx, lis.Policy)
		if !lis.checkError(err) {
			sConn.Close()
		}
		return
	} else if str == "<GETSWF />" {
		err = con.Write(ctx, lis.SWF)
		if !lis.checkError(err) {
			sConn.Close()
		}
		return
	}
	lis.Handler(con)
}

func (lis *Listener) Start(ctx context.Context, network, address string) (waitForExit func()) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err := lis.ListenAndServe(ctx, network, address)
		lis.checkError(err)
		wg.Done()
	}()
	return wg.Wait
}

func (lis *Listener) ListenAndServe(ctx context.Context, network, address string) (err error) {
	ln, err := net.Listen(network, address)
	if err != nil {
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go lis.handle(ctx, conn)
	}
}
