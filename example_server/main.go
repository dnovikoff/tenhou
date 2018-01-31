package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/tile"

	"github.com/dnovikoff/tenhou/client"
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

type Server struct {
	server.XMLWriter
	client client.XMLWriter
	name   string
	conn   *network.XMLConnection
	ctx    context.Context
	cancel func()
}

func (w *Server) send() {
	m := w.client.Buffer().String()
	fmt.Println("Send: " + m)
	err := w.conn.Write(w.ctx, m)
	if !checkSuccess(err) {
		w.cancel()
	}
	w.client.Reset()
}

func (w *Server) Hello(name string, tid string, sex tbase.Sex) {
	w.name = name
	w.client.Hello(name, "20180117-e7b5e83e", client.HelloStats{})
	w.send()
}

func (w *Server) RequestLobbyStatus(v, V int) {
	w.client.LobbyStats(
		"BWQ1BOi1Xk1Ep",
		"4D4C4B8D4B4B8C4B12B4C9C1B3B2C1C2C1C2B",
		"Dc3E1Io1I2U1BE4Q8JA3M1Ck12GA4BM4E8BU4E8V1Ci2BC2BC2G1e1V1S1t2b",
	)
	w.send()
}

func (w *Server) Join(lobbyNumber int, lobbyType int, rejoin bool) {
	w.client.Go("", lobbyType, lobbyNumber)
	users := tbase.UserList{
		tbase.User{Name: w.name, Dan: 10, Rate: 3400, Sex: tbase.SexFemale},
		tbase.User{Name: "Second", Dan: 12, Rate: 2200, Sex: tbase.SexFemale},
		tbase.User{Name: "Third", Dan: 13, Rate: 900, Sex: tbase.SexComputer},
		tbase.User{Name: "Fourth", Dan: 15, Rate: 1801, Sex: tbase.SexMale},
	}
	w.send()
	w.client.UserList(users)
	w.send()
	w.client.LogInfo(base.Self, "")
	w.send()
}

func (w *Server) Bye() {
	w.cancel()
}

func (w *Server) NextReady() {
	g := compact.NewTileGenerator()
	hand, _ := g.InstancesFromString("1367m1566p1699s1z")
	w.client.Init(client.Init{
		Init: tbase.Init{
			Seed: tbase.Seed{
				RoundNumber: 0,
				Honba:       11,
				Sticks:      6,
				Dice:        [2]int{1, 2},
				Indicator:   g.Instance(tile.Red),
			},
			Dealer: base.Self,
			Scores: tbase.Scores{12000, 13000, 25000, 99900},
		},
		Hand: hand,
	})
	w.send()
	w.client.Take(base.Self, g.Instance(tile.Green), client.SuggestDraw|client.SuggestTsumo)
	w.send()
}

func handle(sConn net.Conn) {
	defer sConn.Close()
	handleImpl(sConn)
}

func handleImpl(sConn net.Conn) {
	con := network.NewXMLConnection(sConn)
	ctx, cancel := context.WithCancel(context.Background())
	s := &Server{
		server.NewXMLWriter(),
		client.NewXMLWriter(),
		"",
		con, ctx, cancel,
	}
	for {
		ctx, _ := context.WithTimeout(s.ctx, time.Second*10)
		str, err := con.Read(ctx)
		fmt.Println("Get: " + str)
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
		err = server.ProcessXMLMessage(str, s)
		if !checkSuccess(err) {
			return
		}
	}
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
