package game

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/dnovikoff/tenhou/network"
	"github.com/dnovikoff/tenhou/parser"
	"github.com/dnovikoff/tenhou/server"
	"github.com/dnovikoff/tenhou/tbase"
	"github.com/dnovikoff/tenhou/util"

	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/effective"
	"github.com/dnovikoff/tempai-core/hand/tempai"
	"github.com/dnovikoff/tempai-core/meld"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tempai-core/yaku"
	"github.com/dnovikoff/tenhou/client"
)

type Player struct {
	Hand    compact.Instances
	Melds   meld.Melds
	Discard compact.Instances
	Score   score.Money
	tempai  tempai.IndexedResult
	furiten bool
	first   bool
}

var rules = yaku.Rules{
	OpenTanyao:           true,
	Renhou:               yaku.LimitYakuman,
	HaiteiIsFromLiveOnly: true,
	AkaDoras:             tile.Instances{},
}

var scoring = score.Rules{
	ManganRound:   false,
	KazoeYakuman:  true,
	DoubleYakuman: false,
	YakumanSum:    true,
	HonbaValue:    100,
}

type Game struct {
	Context    context.Context
	Dealer     bool
	Human      *Player
	Robot      *Player
	Wall       tile.Instances
	Client     client.XMLWriter
	Connection network.XMLConnection
	rnd        *rand.Rand
	Turn       bool
	Rinshan    bool
	reader     *util.NodeReader
	err        error
	logger     *log.Logger
}

var indicator = tile.Red.Instance(0)
var indicators = tile.Instances{indicator}

func NewGame(c network.XMLConnection) *Game {
	src := rand.NewSource(time.Now().UTC().UnixNano())
	rnd := rand.New(src)
	w := client.NewXMLWriter()
	this := &Game{
		Connection: c,
		Client:     w,
		Dealer:     true,
		Turn:       true,
		Human:      NewPlayer(120000),
		Robot:      NewPlayer(120000),
		rnd:        rnd,
		Context:    context.Background(),
	}
	this.logger = log.New(os.Stdout, "", log.LstdFlags)
	this.reader = &util.NodeReader{}
	return this
}

func (this *Game) ctx() context.Context {
	c, _ := context.WithTimeout(this.Context, time.Second*30)
	return c
}

func (this *Game) Send() {
	data := this.Client.String()
	if data == "" {
		return
	}
	this.check(this.Connection.Write(this.ctx(), data))
	this.Client.Reset()
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

func (this *Game) opp(x bool) base.Opponent {
	if x {
		return base.Self
	}
	return base.Front
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
	this.reader.Read = func() (string, error) {
		return this.Connection.Read(this.ctx())
	}
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
	melds meld.Melds,
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
		Melds:          tbase.NewTenhouMelds(melds),
		Yakumans:       tbase.YakumansFromCore(yakuman),
		Yakus:          tbase.YakusFromCore(yaku),
	}
	if len(indicators) > 5 {
		a.DoraIndicators, a.UraIndicators = indicators[:5], indicators[5:]
	}
	this.Client.Agari(&a)
	this.Send()
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
	s := scoring.GetYakumanScore(1, 0)
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
	if !p.tempai.Waits().Check(t.Tile()) {
		if who {
			this.doAgari(
				who,
				isTsumo,
				p.Hand,
				p.Melds,
				t,
				inds(tile.White, p.tempai.Waits()),
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
		tls := (p.Discard.UniqueTiles() & p.tempai.Waits())
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
		Rules:       &rules,
		IsLastTile:  !this.canTake(),
		IsFirstTake: p.first,
		IsTsumo:     isTsumo,
		SelfWind:    wind(this.Dealer == this.Turn),
		IsRinshan:   this.Rinshan,
	}
	win := yaku.Win(p.tempai, ctx)
	s = scoring.GetScoreByResult(win, 0)
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
		inds(tile.East, p.tempai.Waits()),
		pay,
		yaku.FuPoints(win.Fus.Sum().Round()),
		win.Yaku,
		win.Yakuman,
	)
	return true
}

func (this *Game) RobotTurn() (result bool) {
	this.Client.Take(base.Front, 0, 0)
	this.Send()
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
	p.Melds.AddTo(visible)
	this.Human.Melds.AddTo(visible)

	res := effective.Calculate(p.Hand, len(p.Melds), visible)
	x := res.Sorted(visible)
	bestTile := x.Best().Tile
	toDrop := p.Hand.GetMask(bestTile).First()
	p.drop(toDrop)
	sg := client.SuggestRon
	this.Client.Drop(base.Front, toDrop, toDrop == t, sg)
	this.Send()
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
	sg := client.SuggestTsumo
	this.Client.Take(base.Self, t, sg)
	this.Send()
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
		this.Client.Drop(base.Self, i, i == t, client.SuggestNone)
		this.Send()
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
			if !p.Hand.CheckFull(first.Tile()) {
				cb.Default()
				return
			}
			m := meld.NewKan(first.Tile(), first.CopyId())
			p.first = false
			this.Robot.first = false
			this.Rinshan = true
			p.Melds = append(p.Melds, m.Meld())
			p.Hand.SetCount(first.Tile(), 0)
			this.Client.Declare(base.Self, tbase.NewTenhouMeld(m.Meld()), client.SuggestNone)
			this.Send()
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
	this.Client.Ryuukyoku(&rk)
	this.Send()
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

func (this *Game) Run() {
	this.wait()
	this.Client.Go("", 11, 0)
	this.Send()
	this.Client.UserList(tbase.UserList{
		tbase.User{Num: 0, Name: "Player", Sex: tbase.SexMale, Rate: 1500},
		tbase.User{Num: 1, Name: "_", Sex: tbase.SexFemale, Rate: 1500},
		tbase.User{Num: 2, Name: "Robot", Sex: tbase.SexComputer, Rate: 1500},
		tbase.User{Num: 3, Name: "_", Sex: tbase.SexFemale, Rate: 1500},
	})
	this.Send()
	this.Client.LogInfo(base.Self, "")
	this.Send()
	this.wait() // Ok
	//	this.wait() // Ready
	rnd := 0
	startTile := tile.Tiles{tile.Pin1, tile.Man1, tile.Sou1}
	for this.RunOne(rnd, startTile[rnd%len(startTile)]) {
		rnd++
		this.Dealer = !this.Dealer
	}
	this.wait()
	this.Connection.Close()
}

func (this *Game) RunOne(rnd int, startTile tile.Tile) bool {
	this.logger.Printf("Round %v START", rnd)
	tiles := compact.NewAllInstancesFromTo(startTile, startTile+9).Instances()
	tile.Shuffle(tiles, this.rnd)
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
	this.Send()

	for this.MakeTurn() {
		if this.err != nil {
			return false
		}
	}
	this.logger.Printf("Round %v END", rnd)
	return !this.isDead()
}
