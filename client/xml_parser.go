package client

import (
	"net/url"
	"strconv"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"
)

func ProcessXMLMessage(message string, c Controller) (err error) {
	nodes, err := parser.ParseXML(message)
	if err != nil {
		return
	}
	for _, v := range nodes {
		if err = ProcessXMLNode(&v, c); err != nil {
			return
		}
		if err = v.ValidateUnused(); err != nil {
			return
		}
	}
	return
}

func GetWithOpponent(node *parser.Node) WithOpponent {
	return WithOpponent{node.GetWho()}
}

func GetWithDealer(node *parser.Node) WithDealer {
	return WithDealer{node.GetDealer()}
}

func getWithSuggest(node *parser.Node) WithSuggest {
	return WithSuggest{Suggest(node.Int("t"))}
}

func GetWithLobby(node *parser.Node) WithLobby {
	l := -1
	if node.Check("lobby") {
		l = node.Int("lobby")
	}
	return WithLobby{
		LobbyType:   node.Int("type"),
		LobbyNumber: l,
	}
}

func ProcessXMLNode(node *parser.Node, c Controller) (err error) {
	switch node.Name {
	case "REACH":
		c.Reach(Reach{
			GetWithOpponent(node),
			node.GetScores(),
			node.Int("step")})
	case "N":
		c.Declare(Declare{
			GetWithOpponent(node),
			getWithSuggest(node),
			node.GetMeld()})
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
		c.LogInfo(LogInfo{
			GetWithDealer(node),
			node.String("log"),
		})
	case "GO":
		c.Go(Go{
			WithLobby: GetWithLobby(node),
			GpID:      node.String("gpid"),
		})
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
		c.Hello(Hello{
			Name:        name,
			Auth:        node.String("auth"),
			RatingScale: node.String("ratingscale"),
			PF4:         node.String("PF4"),
			RR:          node.String("rr"),
			Expire:      node.Int("expire"),
			ExpireDays:  node.Int("expiredays"),
			Message:     message,
		})
	case "UN":
		err = ProcessUserList(node, c)
	case "LN":
		c.LobbyStats(LobbyStats{
			N: node.String("n"),
			J: node.String("j"),
			G: node.String("g"),
		})
	case "AGARI":
		result, err := parser.ParseAgari(node)
		if err != nil {
			return err
		}
		c.Agari(*result)
	case "DORA":
		c.Indicator(WithInstance{
			node.GetInstance("hai"),
		})
	case "PROF":
		c.EndButton(EndButton{
			WithLobby: GetWithLobby(node),
			Add:       node.String("add"),
		})
	case "FURITEN":
		c.Furiten(Furiten{
			node.String("show") == "1"})
	case "RYUUKYOKU":
		r, err := parser.ParseRyuukyoku(node)
		if err != nil {
			return err
		}
		c.Ryuukyoku(*r)
	case "REJOIN":
		n, t, rejoin, err := util.ParseJoinString(node.String("t"))
		if err != nil {
			return err
		}
		c.Rejoin(Rejoin{
			WithLobby: WithLobby{
				LobbyNumber: n,
				LobbyType:   t,
			},
			Rejoin: rejoin,
		})
	case "BYE":
		c.Disconnect(GetWithOpponent(node))
	case "RANKING":
		c.Ranking(Ranking{node.String("v2")})
	case "CHAT":
		name, err := url.PathUnescape(node.String("uname"))
		if err != nil {
			return stackerr.Wrap(err)
		}
		text, err := url.PathUnescape(node.String("text"))
		if err != nil {
			return stackerr.Wrap(err)
		}
		c.Chat(Chat{name, text})
	case "SAIKAI":
		ch := node.GetScoreChanges()
		status, err := node.GetTableStatus()
		if err != nil {
			return err
		}
		c.Recover(Recover{
			GetWithDealer(node),
			status,
			ch})
	default:
		_, o, t := parseLetterNode(node.Name, 'D')
		if t != tile.InstanceNull {
			c.Drop(Drop{
				WithOpponent{o},
				WithInstance{t},
				getWithSuggest(node),
				false,
			})
			return
		}
		_, o, t = parseLetterNode(node.Name, 'd')
		if t != tile.InstanceNull {
			c.Drop(Drop{
				WithOpponent{o},
				WithInstance{t},
				getWithSuggest(node),
				true,
			})
			return
		}
		ok, o, t := parseLetterNode(node.Name, 'T')
		if ok {
			c.Take(Take{
				WithOpponent{o},
				WithInstance{t},
				getWithSuggest(node),
			})
			return
		}
		return stackerr.Newf("Unexpectd node '%v'", node.Name)
	}
	return
}

func ProcessUserList(node *parser.Node, c UNController) error {
	if len(node.Attributes) == 1 {
		o := tbase.Self
		for k := range node.Attributes {
			switch k {
			case "n0":
				o = tbase.Self
			case "n1":
				o = tbase.Right
			case "n2":
				o = tbase.Front
			case "n3":
				o = tbase.Left
			default:
				return stackerr.Newf("Unexpected key '%v' for UN (RECONNECT)", k)
			}
			name, err := url.QueryUnescape(node.String(k))
			if err != nil {
				return stackerr.Wrap(err)
			}
			c.Reconnect(Reconnect{
				WithOpponent{o}, name})
		}
		return nil
	}
	ul := tbase.UserList{
		Dan:  node.IntList("dan"),
		RC:   node.IntList("rc"),
		Rate: node.FloatList("rate"),
		Gold: node.IntList("gold"),
	}
	sx := node.StringList("sx")
	for k, v := range sx {
		sex := tbase.ParseSexLetter(v)
		ul.Sex = append(ul.Sex, sex)
		name, err := url.QueryUnescape(node.String("n" + strconv.Itoa(k)))
		if err != nil {
			return stackerr.Wrap(err)
		}
		ul.Names = append(ul.Names, name)
	}
	c.UserList(UserList{ul})
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

func parseLetterNode(in string, first byte) (ok bool, o tbase.Opponent, t tile.Instance) {
	t = tile.InstanceNull
	if len(in) < 1 {
		return
	}
	firstLetter := in[0]
	if firstLetter < first {
		return
	}
	o = tbase.Opponent(firstLetter - first)
	if o > tbase.Left {
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
	t = util.InstanceFromTenhou(num)
	return
}
