package log

import (
	"strconv"
	"strings"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/util"
)

func ProcessXMLNodes(nodes parser.Nodes, c Controller) (err error) {
	for _, v := range nodes {
		err = ProcessXMLNode(&v, c)
		if err != nil {
			return
		}
	}
	return
}

func ProcessXMLNode(node *parser.Node, c Controller) (err error) {
	name := node.Name
	switch name {
	case "SHUFFLE":
		c.Shuffle(Shuffle{
			Seed: node.String("seed"),
			Ref:  node.String("ref"),
		})
	case "GO":
		c.Go(client.GetWithLobby(node))
	case "UN":
		floatFormat := strings.Contains(node.Attributes["rate"], ".")
		if floatFormat {
			c.SetFloatFormat()
		}
		err = client.ProcessUserList(node, c)
	case "REACH":
		params := client.Reach{Step: node.Int("step"), Score: node.GetScores()}
		params.Opponent = node.GetWho()
		c.Reach(params)
	case "TAIKYOKU":
		c.Start(client.WithDealer{node.GetDealer()})
	case "INIT":
		x, err := node.GetInit()
		if err != nil {
			return err
		}
		c.Init(Init{
			x,
			node.GetHands(),
			node.String("shuffle"),
		})
	case "N":
		params := Declare{}
		params.Opponent = node.GetWho()
		params.Meld = node.GetMeld()
		c.Declare(params)
	case "AGARI":
		floatFormat := strings.Contains(node.Attributes["owari"], ".")
		agari, err := parser.ParseAgari(node)
		if err != nil {
			return err
		}
		if floatFormat {
			c.SetFloatFormat()
		}
		c.Agari(*agari)
	case "RYUUKYOKU":
		floatFormat := strings.Contains(node.Attributes["owari"], ".")
		r, err := parser.ParseRyuukyoku(node)
		if err != nil {
			return err
		}
		if floatFormat {
			c.SetFloatFormat()
		}
		c.Ryuukyoku(*r)
	case "BYE":
		params := client.WithOpponent{node.GetWho()}
		c.Disconnect(params)
	case "DORA":
		i := node.GetInstance("hai")
		if i == tile.InstanceNull {
			return stackerr.Newf("Unexpected hai value")
		}
		c.Indicator(client.WithInstance{i})
	default:
		if !tryDefault(node, c) {
			return stackerr.Newf("Unexpected node %v", name)
		}
	}
	if err != nil {
		return
	}
	keys := node.Keys()
	if len(keys) > 0 {
		return stackerr.Newf("Unprocessed attributes %v for %v [%v]", keys, name, node)
	}
	return
}

func tryDefault(in *parser.Node, c Controller) bool {
	first := in.Name[0]
	tileStr := in.Name[1:]
	x, err := strconv.Atoi(tileStr)
	if err != nil {
		return false
	}
	params := WithOpponentAndInstance{}
	params.Instance = util.InstanceFromTenhou(x)
	if first >= 'T' && first <= 'W' {
		params.Opponent = base.Opponent(first - 'T')
		c.Draw(params)
		return true
	} else if first >= 'D' && first <= 'G' {
		params.Opponent = base.Opponent(first - 'D')
		c.Discard(params)
		return true
	}
	return false
}
