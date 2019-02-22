package game

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/hand/effective"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/network"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/server"
	"github.com/dnovikoff/tenhou/tbase"
)

var rules = &yaku.RulesStruct{
	IsOpenTanyao:         true,
	RenhouLimit:          yaku.LimitYakuman,
	IsHaiteiFromLiveOnly: true,
	AkaDoras:             tile.Instances{},
}

var scoring = &score.RulesStruct{
	IsManganRound:  false,
	IsKazoeYakuman: true,
	IsYakumanSum:   true,
	HonbaValue:     100,
}

type Game struct {
	Context    context.Context
	Dealer     bool
	Human      *Player
	Robot      *Player
	Wall       tile.Instances
	Client     client.Controller
	Connection network.XMLConnection
	rnd        *rand.Rand
	Turn       bool
	Rinshan    bool
	reader     *parser.NodeReader
	err        error
	logger     *log.Logger
}

var indicator = tile.Red.Instance(0)
var indicators = tile.Instances{indicator}

func NewGame(c network.XMLConnection) *Game {
	src := rand.NewSource(time.Now().UTC().UnixNano())
	rnd := rand.New(src)
	this := &Game{
		Connection: c,
		Dealer:     true,
		Turn:       true,
		Human:      NewPlayer(120000),
		Robot:      NewPlayer(120000),
		rnd:        rnd,
		Context:    context.Background(),
	}
	writer := client.NewXMLWriter()
	writer.Commit = this.commit
	this.Client = writer
	this.logger = log.New(os.Stdout, "", log.LstdFlags)
	this.reader = parser.NewNodeReader()
	return this
}

func ctxTimeout(ctx context.Context) (context.Context, func()) {
	return context.WithTimeout(ctx, time.Second*30)
}

func (g *Game) commit(data string) {
	if data == "" {
		return
	}
	ctx, cancel := ctxTimeout(g.Context)
	defer cancel()
	g.check(g.Connection.Write(ctx, data))
}

func (g *Game) GetHand() compact.Instances {
	x := compact.NewInstances()
	hand, wall := g.Wall[:13], g.Wall[13:]
	g.Wall = wall
	return x.Add(hand)
}

func (g *Game) canTake() bool {
	return len(g.Wall) > 0
}

func (g *Game) take() tile.Instance {
	x := g.Wall[0]
	g.Wall = g.Wall[1:]
	return x
}

func (g *Game) Player(x bool) *Player {
	if x {
		return g.Human
	}
	return g.Robot
}

func (g *Game) opp(x bool) tbase.Opponent {
	if x {
		return tbase.Self
	}
	return tbase.Front
}

func (g *Game) scores() tbase.Scores {
	return tbase.Scores{g.Human.Score, 0, g.Robot.Score, 0}
}

func (g *Game) diff(who bool, m score.Money) tbase.ScoreChanges {
	if !who {
		m *= -1
	}
	g.Human.Score += m
	g.Robot.Score -= m
	return tbase.ScoreChanges{
		tbase.ScoreChange{tbase.MoneyToInt(g.Human.Score + -m), tbase.MoneyToInt(m)},
		tbase.ScoreChange{},
		tbase.ScoreChange{tbase.MoneyToInt(g.Robot.Score + m), tbase.MoneyToInt(-m)},
		tbase.ScoreChange{},
	}
}

func (g *Game) ProcessOne(cb *server.Callbacks) bool {
	repeat := false
	badRequest := false
	cb.CbPing = func() {
		g.logger.Print("Got ping")
		repeat = true
	}
	cb.CbRequestLobbyStatus = func(int, int) {
		g.logger.Print("Got lobby status")
		repeat = true
	}
	cb.Default = func() {
		badRequest = true
	}
	repeat = true
	var node *parser.Node
	var err error
	for repeat {
		repeat = false
		node, err = g.reader.Next()
		if !g.check(err) {
			return false
		}
		err = server.ProcessXMLNode(node, cb)
		if !g.check(err) {
			return false
		}
	}
	if badRequest {
		g.check(fmt.Errorf("Unexpected message %s", node.Name))
		return false
	}
	return true
}

func (g *Game) isDead() bool {
	return g.Human.Score < 0 || g.Robot.Score < 0 || g.err != nil
}

func inds(head tile.Tile, data compact.Tiles) (ret tile.Instances) {
	ret = tile.Instances{head.Instance(0)}
	for _, v := range data.Tiles() {
		ret = append(ret, v.Instance(0))
	}
	return
}

func (g *Game) doAgari(
	who bool,
	isTsumo bool,
	hand compact.Instances,
	melds tbase.CalledList,
	winTile tile.Instance,
	indicators tile.Instances,
	money score.Money,
	fu yaku.FuPoints,
	yaku yaku.YakuSet,
	yakuman yaku.Yakumans) {
	diff := g.diff(who, money)
	var fin tbase.ScoreChanges
	if g.isDead() {
		fin = g.diff(who, 0)
		for k, v := range fin {
			fin[k].Diff = v.Score
		}
	}
	op := who
	if !isTsumo {
		op = !op
	}
	a := tbase.Agari{
		Who:            g.opp(who),
		From:           g.opp(op),
		Score:          tbase.Score{fu, money, 0},
		FinalScores:    fin.ToFinal(true),
		Changes:        diff,
		Hand:           hand.Instances(),
		DoraIndicators: indicators,
		WinTile:        winTile,
		Melds:          tbase.EncodeCalledList(melds),
	}
	var err error
	a.Yakus, err = tbase.YakusFromCore(yaku)
	if err != nil {
		g.logger.Printf("Error converting yakus: %v", err)
	}
	a.Yakumans, err = tbase.YakumansFromCore(yakuman)
	if err != nil {
		g.logger.Printf("Error converting yakumans: %v", err)
	}
	if len(indicators) > 5 {
		a.DoraIndicators, a.UraIndicators = indicators[:5], indicators[5:]
	}
	g.Client.Agari(a)
	g.wait()
}

func (g *Game) wait() {
	g.logger.Print("Waiting...")
	cb := &server.Callbacks{}
	cb.CbNextReady = func() {}
	cb.CbGoOK = func() {}
	cb.CbBye = func() {
		g.Connection.Close()
	}
	cb.CbRequestLobbyStatus = func(int, int) {}
	g.ProcessOne(cb)
}

func wind(x bool) base.Wind {
	if x {
		return base.WindEast
	}
	return base.WindWest
}

func (g *Game) tryWin(t tile.Instance, who, isTsumo bool) (done bool) {
	p := g.Player(who)
	s := score.GetYakumanScore(scoring, 1, 0)
	penalty := -s.PayRon
	if g.Dealer {
		penalty = -s.PayRonDealer
	}
	// Noten ron
	if p.tempai == nil {
		if who {
			g.doAgari(
				who,
				isTsumo,
				p.Hand,
				p.Melds,
				t,
				inds(tile.North, 0),
				penalty,
				13,
				nil,
				yaku.Yakumans{yaku.YakumanKokushi},
			)
			return true
		}
		return false
	}
	// Wrong tile ron
	waits := tempai.GetWaits(p.tempai)
	if !waits.Check(t.Tile()) {
		if who {
			g.doAgari(
				who,
				isTsumo,
				p.Hand,
				p.Melds,
				t,
				inds(tile.White, waits),
				penalty,
				13,
				nil,
				yaku.Yakumans{yaku.YakumanKokushi},
			)
			return true
		}
		return false
	}
	if !isTsumo && p.furiten {
		tls := (p.Discard.UniqueTiles() & waits)
		if who {
			g.doAgari(
				who,
				isTsumo,
				p.Hand,
				p.Melds,
				t,
				inds(tile.Red, tls),
				penalty,
				13,
				nil,
				yaku.Yakumans{yaku.YakumanKokushi},
			)
			return true
		}
		return false
	}
	ctx := &yaku.Context{
		Tile:        t,
		Rules:       rules,
		IsLastTile:  !g.canTake(),
		IsFirstTake: p.first,
		IsTsumo:     isTsumo,
		SelfWind:    wind(g.Dealer == g.Turn),
		IsRinshan:   g.Rinshan,
	}
	win := yaku.Win(p.tempai, ctx, nil)
	s = score.GetScoreByResult(scoring, win, 0)
	pay := s.PayRon
	if who == g.Dealer {
		pay = s.PayRonDealer
	}

	g.doAgari(
		who,
		isTsumo,
		p.Hand,
		p.Melds,
		t,
		inds(tile.East, waits),
		pay,
		yaku.FuPoints(win.Fus.Sum().Round()),
		win.Yaku,
		win.Yakumans,
	)
	return true
}

func (g *Game) RobotTurn() (result bool) {
	{
		params := client.Take{}
		params.Opponent = tbase.Front
		g.Client.Take(params)
	}
	t := g.take()
	p := g.Robot
	p.take(t)
	if g.tryWin(t, false, true) {
		return false
	}
	visible := compact.NewInstances()
	visible.
		Merge(p.Hand).
		Merge(p.Discard).
		Merge(g.Human.Discard)
	p.Melds.Add(visible)
	g.Human.Melds.Add(visible)

	res := effective.Calculate(
		p.Hand,
		calc.Declared(p.Melds.Core()),
		calc.Used(visible),
	)
	x := res.Sorted(visible)
	bestTile := x.Best().Tile
	toDrop := p.Hand.GetMask(bestTile).First()
	p.drop(toDrop)
	params := client.Drop{}
	params.Instance = toDrop
	params.Suggest = client.SuggestRon
	params.IsTsumogiri = (toDrop == t)
	params.Opponent = tbase.Front
	g.Client.Drop(params)
	cb := &server.Callbacks{}
	cb.CbCall = func(x server.Answer, t tile.Instances) {
		switch x {
		case server.AnswerSkip:
			result = true
		case server.AnswerRon:
			g.tryWin(toDrop, true, false)
			result = false
		default:
			// Unexpected answer
			cb.Default()
		}
	}
	g.ProcessOne(cb)
	return
}

func (g *Game) HumanTurn() (result bool) {
	t := g.take()
	params := client.Take{}
	params.Opponent = tbase.Self
	params.Instance = t
	params.Suggest = client.SuggestTsumo
	g.Client.Take(params)
	p := g.Human
	p.take(t)

	cb := &server.Callbacks{}
	extraTurn := false
	cb.CbDrop = func(i tile.Instance) {
		if !p.Hand.Check(i) {
			cb.Default()
			return
		}
		p.drop(i)
		params := client.Drop{}
		params.Opponent = tbase.Self
		params.Instance = t
		params.IsTsumogiri = (i == t)
		g.Client.Drop(params)
		result = !g.tryWin(i, false, false)
	}
	cb.CbCall = func(x server.Answer, i tile.Instances) {
		switch x {
		case server.AnswerClosedKan:
			if !g.canTake() || len(i) != 1 {
				cb.Default()
				return
			}
			first := i[0]
			if !p.Hand.GetMask(first.Tile()).IsFull() {
				cb.Default()
				return
			}
			m := &tbase.Called{
				Type:   tbase.ClosedKan,
				Called: first,
				Core:   calc.Kan(t.Tile()),
				Tiles:  compact.NewMask(0, t.Tile()).SetCount(4).Instances(),
			}
			p.first = false
			g.Robot.first = false
			g.Rinshan = true
			p.Melds = append(p.Melds, m)
			p.Hand.SetCount(first.Tile(), 0)
			params := client.Declare{}
			params.Opponent = tbase.Self
			params.Meld = tbase.EncodeCalled(m)
			g.Client.Declare(params)
			extraTurn = true
		case server.AnswerSkip:
			result = true
		case server.AnswerTsumo:
			g.tryWin(t, true, true)
		default:
			// Unexpected answer
			cb.Default()
		}
	}
	g.ProcessOne(cb)

	if extraTurn {
		return g.HumanTurn()
	}
	return
}

func (g *Game) check(err error) bool {
	if err != nil {
		g.logger.Printf("New error: %v", err)
	}
	if g.err != nil {
		return false
	}
	if err == nil {
		return true
	}
	g.err = err
	return false
}

func (g *Game) MakeDraw() {
	rk := tbase.Ryuukyoku{
		DrawType: tbase.DrawEnd,
	}
	g.Client.Ryuukyoku(rk)
	g.wait()
}

func (g *Game) MakeTurn() (x bool) {
	if !g.canTake() {
		g.MakeDraw()
		return false
	}
	if g.Turn {
		x = g.HumanTurn()
	} else {
		x = g.RobotTurn()
	}
	g.Turn = !g.Turn
	return
}

func (g *Game) auth() bool {
	cb := &server.Callbacks{}
	cb.CbHello = func(name string, tid string, sex tbase.Sex) {
		g.Client.Hello(client.Hello{Name: name, Auth: "20180117-e7b5e83e"})
	}
	if !g.ProcessOne(cb) {
		return false
	}
	cb.CbHello = nil
	cb.CbAuth = func(string) {}
	if !g.ProcessOne(cb) {
		return false
	}
	return true
}

func (g *Game) Run() {
	ctx, stop := context.WithCancel(g.Context)
	g.reader.ReadCallback = func(ctx context.Context) (string, error) {
		rCtx, cancel := ctxTimeout(ctx)
		defer cancel()
		return g.Connection.Read(rCtx)
	}
	waitForExit := g.reader.Start(ctx)
	defer func() {
		stop()
		waitForExit()
		g.Connection.Close()
	}()
	if !g.auth() {
		return
	}
	g.wait()
	params := client.Go{}
	params.LobbyType = 11
	g.Client.Go(params)
	g.Client.UserList(client.UserList{tbase.UserList{
		Names: []string{"Player", "_", "Robot", "_"},
		Sex:   []tbase.Sex{tbase.SexMale, tbase.SexFemale, tbase.SexComputer, tbase.SexFemale},
		Rate:  []tbase.Float{{1500, true}, {1500, true}, {1500, true}, {1500, true}},
	}})
	g.Client.LogInfo(client.LogInfo{})
	g.wait() // Ok
	//	this.wait() // Ready
	rnd := 0
	startTile := tile.Tiles{tile.Pin1, tile.Man1, tile.Sou1}
	for g.RunOne(rnd, startTile[rnd%len(startTile)]) {
		rnd++
		g.Dealer = !g.Dealer
	}
	g.wait()
}

func (g *Game) RunOne(rnd int, startTile tile.Tile) bool {
	g.logger.Printf("Round %v START", rnd)
	tiles := compact.AllInstancesFromTo(startTile, startTile+9).Instances()
	g.rnd.Shuffle(len(tiles), func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})
	g.Wall = tiles
	g.Human.Init(g.GetHand())
	g.Robot.Init(g.GetHand())
	g.Client.Init(client.Init{
		Init: tbase.Init{
			Seed: tbase.Seed{
				RoundNumber: rnd,
				Indicator:   indicator,
				Dice:        [2]int{1, 2},
			},
			Scores: g.scores(),
			Dealer: g.opp(g.Dealer),
		},
		Hand: g.Human.Hand.Instances(),
	})

	for g.MakeTurn() {
		if g.err != nil {
			return false
		}
	}
	g.logger.Printf("Round %v END", rnd)
	return !g.isDead()
}
