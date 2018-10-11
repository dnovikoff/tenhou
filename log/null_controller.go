package log

import (
	"github.com/dnovikoff/tenhou/client"
	"github.com/dnovikoff/tenhou/tbase"
)

type NullController struct{}

var _ Controller = NullController{}

func (NullController) Open(Info) bool                  { return true }
func (NullController) Close()                          {}
func (NullController) SetFloatFormat()                 {}
func (NullController) Shuffle(Shuffle)                 {}
func (NullController) Go(client.WithLobby)             {}
func (NullController) Start(client.WithDealer)         {}
func (NullController) Init(Init)                       {}
func (NullController) Draw(WithOpponentAndInstance)    {}
func (NullController) Discard(WithOpponentAndInstance) {}
func (NullController) Declare(Declare)                 {}
func (NullController) Ryuukyoku(tbase.Ryuukyoku)       {}
func (NullController) Reach(client.Reach)              {}
func (NullController) Agari(tbase.Agari)               {}
func (NullController) Indicator(client.WithInstance)   {}
func (NullController) Disconnect(client.WithOpponent)  {}
func (NullController) UserList(client.UserList)        {}
func (NullController) Reconnect(client.Reconnect)      {}
