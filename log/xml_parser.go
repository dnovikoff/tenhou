package log

import (
	"strconv"
	"strings"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/facebookgo/stackerr"
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
		c.Shuffle(node.String("seed"), node.String("ref"))
	case "GO":
		l := -1
		if node.Check("lobby") {
			l = node.Int("lobby")
		}
		c.Go(node.Int("type"), l)
	case "UN":
		floatFormat := strings.Contains(node.Attributes["rate"], ".")
		if floatFormat {
			c.SetFloatFormat()
		}
		err = client.ProcessUserList(node, c)
	case "REACH":
		c.Reach(node.GetWho(), node.Int("step"), node.GetScores())
	case "TAIKYOKU":
		c.Start(node.GetDealer())
	case "INIT":
		x, err := node.GetInit()
		if err != nil {
			return err
		}
		c.Init(&Init{
			x,
			node.GetHands(),
			node.String("shuffle"),
		})
	case "N":
		c.Declare(node.GetWho(), node.GetMeld())
	case "AGARI":
		floatFormat := strings.Contains(node.Attributes["owari"], ".")
		agari, err := parser.ParseAgari(node)
		if err != nil {
			return err
		}
		if floatFormat {
			c.SetFloatFormat()
		}
		c.Agari(agari)
	case "RYUUKYOKU":
		floatFormat := strings.Contains(node.Attributes["owari"], ".")
		r, err := parser.ParseRyuukyoku(node)
		if err != nil {
			return err
		}
		if floatFormat {
			c.SetFloatFormat()
		}
		c.Ryuukyoku(r)
	case "BYE":
		c.Disconnect(node.GetWho())
	case "DORA":
		i := node.GetInstance("hai")
		if i.IsNull() {
			return stackerr.Newf("Unexpected hai value")
		}
		c.Indicator(i)
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
	instance := tile.Instance(x)
	if first >= 'T' && first <= 'W' {
		c.Draw(base.Opponent(first-'T'), instance)
		return true
	} else if first >= 'D' && first <= 'G' {
		c.Discard(base.Opponent(first-'D'), instance)
		return true
	}
	return false
}
