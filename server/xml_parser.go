package server

import (
	"fmt"
	"net/url"

	"github.com/dnovikoff/tenhou/tbase"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/util"
)

func ProcessXMLMessage(message string, c Controller) (err error) {
	nodes, err := util.ParseXML(message)
	if err != nil {
		return
	}
	for _, v := range nodes {
		err = ProcessXMLNode(&v, c)
		if err != nil {
			return
		}
		unused := v.Keys()
		if len(unused) > 0 {
			err = fmt.Errorf("Unused keys for node '%v': %v", v.Name, unused)
			return
		}
	}
	return
}

func ProcessXMLNode(node *parser.Node, c Controller) (err error) {
	switch node.Name {
	//<HELO name="NoName" tid="f0" sx="M" />
	case "HELO":
		sex := tbase.ParseSexLetter(node.String("sx"))
		c.Hello(node.String("name"), node.String("tid"), sex)
	//<AUTH val="20180117-c1ebb26f"/>
	case "AUTH":
		c.Auth(node.String("val"))
	//<PXR V="1" />
	case "PXR":
		c.RequestLobbyStatus(node.Int("v"), node.Int("V"))
	//<JOIN t="0,9" />
	case "JOIN":
		if len(node.Attributes) == 0 {
			c.CancelJoin()
			break
		}
		n, t, rejoin, err := util.ParseJoinString(node.String("t"))
		if err != nil {
			return err
		}
		c.Join(n, t, rejoin)
	case "D":
		c.Drop(tile.Instance(node.Int("p")))
	//<N type="3" hai0="96" hai1="102" />
	case "N":
		tiles := node.GetHai("hai")
		if len(tiles) == 0 {
			for i := 0; i < 4; i++ {
				tiles = append(tiles, node.GetHaiNum(i)...)
			}
		}
		c.Call(Answer(node.Int("type")), tiles)
	case "REACH":
		c.Reach(tile.Instance(node.Int("hai")))
	case "Z":
		c.Ping()
	case "GOK":
		c.GoOK()
	case "NEXTREADY":
		c.NextReady()
	case "BYE":
		c.Bye()
	case "CHAT":
		text, err := url.QueryUnescape(node.String("text"))
		if err != nil {
			return stackerr.Wrap(err)
		}
		c.Chat(text)
	default:
		return stackerr.Newf("Unexpected node '%v'", node.Name)
	}
	return
}
