package client

import (
	"net/url"
	"strconv"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/tbase"
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
			err = stackerr.Newf("Unused keys for node '%v': %v", v.Name, unused)
			return
		}
	}
	return
}

func parseLetterNode(in string, first byte) (ok bool, o base.Opponent, t tile.Instance) {
	t = tile.InstanceNull
	if len(in) < 1 {
		return
	}
	firstLetter := in[0]
	if firstLetter < first {
		return
	}
	o = base.Opponent(firstLetter - first)
	if o > base.Left {
		return
	}
	if len(in) == 1 {
		ok = true
		return
	}
	num, err := strconv.Atoi(in[1:])
	if err != nil {
		return
	}
	ok = true
	t = tile.Instance(num)
	return
}

func ProcessXMLNode(node *parser.Node, c Controller) (err error) {
	switch node.Name {
	case "REACH":
		c.Reach(node.GetWho(), node.Int("step"), node.GetScores())
	case "N":
		s := Suggest(node.Int("t"))
		c.Declare(node.GetWho(), node.GetMeld(), s)
	case "INIT":
		i, err := getInit(node)
		if err != nil {
			return err
		}
		c.Init(i)
	case "REINIT":
		r := Reinit{}
		r.Init, err = getInit(node)
		if err != nil {
			return
		}
		r.Melds = make([]tbase.Melds, len(r.Scores))
		r.Riichi = make([]int, len(r.Scores))
		r.Discard = make([]tile.Instances, len(r.Scores))
		for k := range r.Scores {
			m := node.IntList("m" + strconv.Itoa(k))
			var melds tbase.Melds
			for _, x := range m {
				melds = append(melds, tbase.Meld(x))
			}
			r.Melds[k] = melds
			d := node.IntList("kawa" + strconv.Itoa(k))
			drop := make(tile.Instances, 0, len(d))
			for n, x := range d {
				if x == 255 {
					r.Riichi[k] = n
					continue
				}
				drop = append(drop, tile.Instance(x))
			}
			r.Discard[k] = drop
		}
		c.Reinit(r)
	case "TAIKYOKU":
		c.LogInfo(node.GetDealer(), node.String("log"))
	case "GO":
		c.Go(node.String("gpid"), node.Int("type"), node.Int("lobby"))
	case "HELO":
		name := node.String("uname")
		name, err := url.QueryUnescape(name)
		if err != nil {
			return stackerr.Wrap(err)
		}
		message, err := url.QueryUnescape(node.String("nintei"))
		if err != nil {
			return stackerr.Wrap(err)
		}
		stats := HelloStats{
			RatingScale: node.String("ratingscale"),
			PF4:         node.String("PF4"),
			RR:          node.String("rr"),
			Expire:      node.Int("expire"),
			ExpireDays:  node.Int("expiredays"),
			Message:     message,
		}
		c.Hello(name, node.String("auth"), stats)
	case "UN":
		err = ProcessUserList(node, c)
	case "LN":
		c.LobbyStats(node.String("n"), node.String("j"), node.String("g"))
	case "AGARI":
		result, err := parser.ParseAgari(node)
		if err != nil {
			return err
		}
		c.Agari(result)
	case "DORA":
		c.Indicator(tile.Instance(node.Int("hai")))
	case "PROF":
		c.EndButton(node.Int("lobby"), node.Int("type"), node.String("add"))
	case "FURITEN":
		c.Furiten(node.String("show") == "1")
	case "RYUUKYOKU":
		r, err := parser.ParseRyuukyoku(node)
		if err != nil {
			return err
		}
		c.Ryuukyoku(r)
	case "REJOIN":
		n, t, rejoin, err := util.ParseJoinString(node.String("t"))
		if err != nil {
			return err
		}
		c.Rejoin(n, t, rejoin)
	case "BYE":
		c.Disconnect(node.GetWho())
	case "RANKING":
		c.Ranking(node.String("v2"))
	case "CHAT":
		name, err := url.PathUnescape(node.String("uname"))
		if err != nil {
			return stackerr.Wrap(err)
		}
		text, err := url.PathUnescape(node.String("text"))
		if err != nil {
			return stackerr.Wrap(err)
		}
		c.Chat(name, text)
	case "SAIKAI":
		ch := node.GetScoreChanges()
		status, err := node.GetTableStatus()
		if err != nil {
			return err
		}
		c.Recover(status, node.GetDealer(), ch)
	default:
		_, o, t := parseLetterNode(node.Name, 'D')
		if !t.IsNull() {
			c.Drop(o, t, false, Suggest(node.Int("t")))
			return
		}
		_, o, t = parseLetterNode(node.Name, 'd')
		if !t.IsNull() {
			c.Drop(o, t, true, Suggest(node.Int("t")))
			return
		}
		ok, o, t := parseLetterNode(node.Name, 'T')
		if ok {
			c.Take(o, t, Suggest(node.Int("t")))
			return
		}
		return stackerr.Newf("Unexpectd node '%v'", node.Name)
	}
	return
}

func ProcessUserList(node *parser.Node, c UNController) error {
	if len(node.Attributes) == 1 {
		o := base.Self
		for k := range node.Attributes {
			switch k {
			case "n0":
				o = base.Self
			case "n1":
				o = base.Right
			case "n2":
				o = base.Front
			case "n3":
				o = base.Left
			default:
				return stackerr.Newf("Unexpected key '%v' for UN (RECONNECT)", k)
			}
			name, err := url.QueryUnescape(node.String(k))
			if err != nil {
				return stackerr.Wrap(err)
			}
			c.Reconnect(o, name)
		}
		return nil
	}
	dan := node.IntList("dan")
	rate := node.FloatList("rate")
	sx := node.StringList("sx")
	if len(dan) != 4 || len(rate) != 4 || len(sx) != 4 {
		return stackerr.Newf("Bad lens for arrays in UN")
	}
	ul := make(tbase.UserList, 4)
	for k := range ul {
		sex := tbase.ParseSexLetter(sx[k])
		name, err := url.QueryUnescape(node.String("n" + strconv.Itoa(k)))
		if err != nil {
			return stackerr.Wrap(err)
		}

		ul[k] = tbase.User{
			Num:  k,
			Name: name,
			Dan:  dan[k],
			Rate: rate[k],
			Sex:  sex,
		}
	}
	c.UserList(ul)
	return nil
}

func getInit(node *parser.Node) (ret Init, err error) {
	ret.Init, err = node.GetInit()
	if err != nil {
		return
	}
	ret.Hand = node.GetHai("hai")
	return
}
