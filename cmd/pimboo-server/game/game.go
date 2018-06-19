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
	IsManganRound:   false,
	IsKazoeYakuman:  true,
	IsYakumanDouble: false,
	IsYakumanSum:    true,
	HonbaValue:      100,
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

func (this *Game) commit(data string) {
	if data == "" {
		return
	}
	ctx, cancel := ctxTimeout(this.Context)
	defer cancel()
	this.check(this.Connection.Write(ctx, data))
}

func (this *Game) GetHand() compact.Instances {
	x := compact.NewInstances()
	hand, wall := this.Wall[:13], this.Wall[13:]
	this.Wall = wall
	return x.Add(hand)
}

func (this *Game) canTake() bool {
	return len(this.Wall) > 0
}

func (this *Game) take() tile.Instance {
	x := this.Wall[0]
	this.Wall = this.Wall[1:]
	return x
}

func (this *Game) Player(x bool) *Player {
	if x {
		return this.Human
	}
	return this.Robot
}

func (this *Game) opp(x bool) tbase.Opponent {
	if x {
		return tbase.Self
	}
	return tbase.Front
}

func (this *Game) scores() tbase.Scores {
	return tbase.Scores{this.Human.Score, 0, this.Robot.Score, 0}
}

func (this *Game) diff(who bool, m score.Money) tbase.ScoreChanges {
	if !who {
		m *= -1
	}
	this.Human.Score += m
	this.Robot.Score -= m
	return tbase.ScoreChanges{
		tbase.ScoreChange{this.Human.Score + -m, m},
		tbase.ScoreChange{},
		tbase.ScoreChange{this.Robot.Score + m, -m},
		tbase.ScoreChange{},
	}
}

func (this *Game) ProcessOne(cb *server.Callbacks) bool {
	repeat := false
	badRequest := false
	cb.CbPing = func() {
		this.logger.Print("Got ping")
		repeat = true
	}
	cb.CbRequestLobbyStatus = func(int, int) {
		this.logger.Print("Got lobby status")
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
		node, err = this.reader.Next()
		if !this.check(err) {
			return false
		}
		err = server.ProcessXMLNode(node, cb)
		if !this.check(err) {
			return false
		}
	}
	if badRequest {
		this.check(fmt.Errorf("Unexpected message %s", node.Name))
		return false
	}
	return true
}

func (this *Game) isDead() bool {
	return this.Human.Score < 0 || this.Robot.Score < 0 || this.err != nil
}

func inds(head tile.Tile, data compact.Tiles) (ret tile.Instances) {
	ret = tile.Instances{head.Instance(0)}
	for _, v := range data.Tiles() {
		ret = append(ret, v.Instance(0))
	}
	return
}

func (this *Game) doAgari(
	who bool,
	isTsumo bool,
	hand compact.Instances,
	melds tbase.CalledList,
	winTile tile.Instance,
	indicators tile.Instances,
	money score.Money,
	fu yaku.FuPoints,
	yaku yaku.YakuSet,
	yakuman yaku.YakumanSet) {
	diff := this.diff(who, money)
	var fin tbase.ScoreChanges
	if this.isDead() {
		fin = this.diff(who, 0)
		for k, v := range fin {
			fin[k].Diff = v.Score
		}
	}
	op := who
	if !isTsumo {
		op = !op
	}
	a := tbase.Agari{
		Who:            this.opp(who),
		From:           this.opp(op),
		Score:          tbase.Score{fu, money, 0},
		FinalScores:    fin,
		Changes:        diff,
		Hand:           hand.Instances(),
		DoraIndicators: indicators,
		WinTile:        winTile,
		Melds:          tbase.EncodeCalledList(melds),
	}
	var err error
	a.Yakus, err = tbase.YakusFromCore(yaku)
	if err != nil {
		this.logger.Printf("Error converting yakus: %v", err)
	}
	a.Yakumans, err = tbase.YakumansFromCore(yakuman)
	if err != nil {
		this.logger.Printf("Error converting yakumans: %v", err)
	}
	if len(indicators) > 5 {
		a.DoraIndicators, a.UraIndicators = indicators[:5], indicators[5:]
	}
	this.Client.Agari(a)
	this.wait()
}

func (this Game) wait() {
	this.logger.Print("Waiting...")
	cb := &server.Callbacks{}
	cb.CbNextReady = func() {}
	cb.CbGoOK = func() {}
	cb.CbBye = func() {
		this.Connection.Close()
	}
	cb.CbRequestLobbyStatus = func(int, int) {}
	this.ProcessOne(cb)
}

func wind(x bool) base.Wind {
	if x {
		return base.WindEast
	}
	return base.WindWest
}

func (this *Game) tryWin(t tile.Instance, who, isTsumo bool) (done bool) {
	p := this.Player(who)
	s := score.GetYakumanScore(scoring, 1, 0)
	penalty := -s.PayRon
	if this.Dealer {
		penalty = -s.PayRonDealer
	}
	// Noten ron
	if p.tempai == nil {
		if who {
			this.doAgari(
				who,
				isTsumo,
				p.Hand,
				p.Melds,
				t,
				inds(tile.North, 0),
				penalty,
				13,
				nil,
				yaku.YakumanSet{yaku.YakumanKokushi: 1},
			)
			return true
		}
		return false
	}
	// Wrong tile ron
	waits := tempai.GetWaits(p.tempai)
	if !waits.Check(t.Tile()) {
		if who {
			this.doAgari(
				who,
				isTsumo,
				p.Hand,
				p.Melds,
				t,
				inds(tile.White, waits),
				penalty,
				13,
				nil,
				yaku.YakumanSet{yaku.YakumanKokushi: 1},
			)
			return true
		}
		return false
	}
	if !isTsumo && p.furiten {
		tls := (p.Discard.UniqueTiles() & waits)
		if who {
			this.doAgari(
				who,
				isTsumo,
				p.Hand,
				p.Melds,
				t,
				inds(tile.Red, tls),
				penalty,
				13,
				nil,
				yaku.YakumanSet{yaku.YakumanKokushi: 1},
			)
			return true
		}
		return false
	}
	ctx := &yaku.Context{
		Tile:        t,
		Rules:       rules,
		IsLastTile:  !this.canTake(),
		IsFirstTake: p.first,
		IsTsumo:     isTsumo,
		SelfWind:    wind(this.Dealer == this.Turn),
		IsRinshan:   this.Rinshan,
	}
	win := yaku.Win(p.tempai, ctx, nil)
	s = score.GetScoreByResult(scoring, win, 0)
	pay := s.PayRon
	if who == this.Dealer {
		pay = s.PayRonDealer
	}

	this.doAgari(
		who,
		isTsumo,
		p.Hand,
		p.Melds,
		t,
		inds(tile.East, waits),
		pay,
		yaku.FuPoints(win.Fus.Sum().Round()),
		win.Yaku,
		win.Yakuman,
	)
	return true
}

func (this *Game) RobotTurn() (result bool) {
	{
		params := client.Take{}
		params.Opponent = tbase.Front
		this.Client.Take(params)
	}
	t := this.take()
	p := this.Robot
	p.take(t)
	if this.tryWin(t, false, true) {
		return false
	}
	visible := compact.NewInstances()
	visible.
		Merge(p.Hand).
		Merge(p.Discard).
		Merge(this.Human.Discard)
	p.Melds.Add(visible)
	this.Human.Melds.Add(visible)

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
	this.Client.Drop(params)
	cb := &server.Callbacks{}
	cb.CbCall = func(x server.Answer, t tile.Instances) {
		switch x {
		case server.AnswerSkip:
			result = true
		case server.AnswerRon:
			this.tryWin(toDrop, true, false)
			result = false
		default:
			// Unexpected answer
			cb.Default()
		}
	}
	this.ProcessOne(cb)
	return
}

func (this *Game) HumanTurn() (result bool) {
	t := this.take()
	params := client.Take{}
	params.Opponent = tbase.Self
	params.Instance = t
	params.Suggest = client.SuggestTsumo
	this.Client.Take(params)
	p := this.Human
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
		this.Client.Drop(params)
		result = !this.tryWin(i, false, false)
	}
	cb.CbCall = func(x server.Answer, i tile.Instances) {
		switch x {
		case server.AnswerClosedKan:
			if !this.canTake() || len(i) != 1 {
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
			this.Robot.first = false
			this.Rinshan = true
			p.Melds = append(p.Melds, m)
			p.Hand.SetCount(first.Tile(), 0)
			params := client.Declare{}
			params.Opponent = tbase.Self
			params.Meld = tbase.EncodeCalled(m)
			this.Client.Declare(params)
			extraTurn = true
		case server.AnswerSkip:
			result = true
		case server.AnswerTsumo:
			this.tryWin(t, true, true)
		default:
			// Unexpected answer
			cb.Default()
		}
	}
	this.ProcessOne(cb)

	if extraTurn {
		return this.HumanTurn()
	}
	return
}

func (this *Game) check(err error) bool {
	if err != nil {
		this.logger.Printf("New error: %v", err)
	}
	if this.err != nil {
		return false
	}
	if err == nil {
		return true
	}
	this.err = err
	return false
}

func (this *Game) MakeDraw() {
	rk := tbase.Ryuukyoku{
		DrawType: tbase.DrawEnd,
	}
	this.Client.Ryuukyoku(rk)
	this.wait()
}

func (this *Game) MakeTurn() (x bool) {
	if !this.canTake() {
		this.MakeDraw()
		return false
	}
	if this.Turn {
		x = this.HumanTurn()
	} else {
		x = this.RobotTurn()
	}
	this.Turn = !this.Turn
	return
}

func (this *Game) auth() bool {
	cb := &server.Callbacks{}
	cb.CbHello = func(name string, tid string, sex tbase.Sex) {
		this.Client.Hello(client.Hello{Name: name, Auth: "20180117-e7b5e83e"})
	}
	if !this.ProcessOne(cb) {
		return false
	}
	cb.CbHello = nil
	cb.CbAuth = func(string) {}
	if !this.ProcessOne(cb) {
		return false
	}
	return true
}

func (this *Game) Run() {
	ctx, stop := context.WithCancel(this.Context)
	this.reader.ReadCallback = func(ctx context.Context) (string, error) {
		rCtx, cancel := ctxTimeout(ctx)
		defer cancel()
		return this.Connection.Read(rCtx)
	}
	waitForExit := this.reader.Start(ctx)
	defer func() {
		stop()
		waitForExit()
		this.Connection.Close()
	}()
	if !this.auth() {
		return
	}
	this.wait()
	params := client.Go{}
	params.LobbyType = 11
	this.Client.Go(params)
	this.Client.UserList(client.UserList{tbase.UserList{
		tbase.User{Num: 0, Name: "Player", Sex: tbase.SexMale, Rate: 1500},
		tbase.User{Num: 1, Name: "_", Sex: tbase.SexFemale, Rate: 1500},
		tbase.User{Num: 2, Name: "Robot", Sex: tbase.SexComputer, Rate: 1500},
		tbase.User{Num: 3, Name: "_", Sex: tbase.SexFemale, Rate: 1500},
	}})
	this.Client.LogInfo(client.LogInfo{})
	this.wait() // Ok
	//	this.wait() // Ready
	rnd := 0
	startTile := tile.Tiles{tile.Pin1, tile.Man1, tile.Sou1}
	for this.RunOne(rnd, startTile[rnd%len(startTile)]) {
		rnd++
		this.Dealer = !this.Dealer
	}
	this.wait()
}

func (this *Game) RunOne(rnd int, startTile tile.Tile) bool {
	this.logger.Printf("Round %v START", rnd)
	tiles := compact.AllInstancesFromTo(startTile, startTile+9).Instances()
	this.rnd.Shuffle(len(tiles), func(i, j int) {
		tiles[i], tiles[j] = tiles[j], tiles[i]
	})
	this.Wall = tiles
	this.Human.Init(this.GetHand())
	this.Robot.Init(this.GetHand())
	this.Client.Init(client.Init{
		Init: tbase.Init{
			Seed: tbase.Seed{
				RoundNumber: rnd,
				Indicator:   indicator,
				Dice:        [2]int{1, 2},
			},
			Scores: this.scores(),
			Dealer: this.opp(this.Dealer),
		},
		Hand: this.Human.Hand.Instances(),
	})

	for this.MakeTurn() {
		if this.err != nil {
			return false
		}
	}
	this.logger.Printf("Round %v END", rnd)
	return !this.isDead()
}
