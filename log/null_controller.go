package log

import (
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/tbase"
)

type NullController struct{}

var _ Controller = NullController{}

func (this NullController) Open(Info) bool                  { return true }
func (this NullController) Close()                          {}
func (this NullController) SetFloatFormat()                 {}
func (this NullController) Shuffle(Shuffle)                 {}
func (this NullController) Go(client.WithLobby)             {}
func (this NullController) Start(client.WithDealer)         {}
func (this NullController) Init(Init)                       {}
func (this NullController) Draw(WithOpponentAndInstance)    {}
func (this NullController) Discard(WithOpponentAndInstance) {}
func (this NullController) Declare(Declare)                 {}
func (this NullController) Ryuukyoku(tbase.Ryuukyoku)       {}
func (this NullController) Reach(client.Reach)              {}
func (this NullController) Agari(tbase.Agari)               {}
func (this NullController) Indicator(client.WithInstance)   {}
func (this NullController) Disconnect(client.WithOpponent)  {}
func (this NullController) UserList(client.UserList)        {}
func (this NullController) Reconnect(client.Reconnect)      {}
