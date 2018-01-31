package log

import (
	"github.com/dnovikoff/tempai-core/base"
	"github.com/dnovikoff/tempai-core/score"
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type NullController struct{}

var _ Controller = NullController{}

func (this NullController) Open(*Info) bool                         { return true }
func (this NullController) Close()                                  {}
func (this NullController) SetFloatFormat()                         {}
func (this NullController) Shuffle(seed, ref string)                {}
func (this NullController) Go(t int, lobby int)                     {}
func (this NullController) Start(base.Opponent)                     {}
func (this NullController) Init(*Init)                              {}
func (this NullController) Draw(base.Opponent, tile.Instance)       {}
func (this NullController) Discard(base.Opponent, tile.Instance)    {}
func (this NullController) Declare(base.Opponent, tbase.Meld)       {}
func (this NullController) Ryuukyoku(*tbase.Ryuukyoku)              {}
func (this NullController) Reach(base.Opponent, int, []score.Money) {}
func (this NullController) Agari(*tbase.Agari)                      {}
func (this NullController) Indicator(tile.Instance)                 {}
func (this NullController) Disconnect(base.Opponent)                {}
func (this NullController) UserList(users tbase.UserList)           {}
func (this NullController) Reconnect(o base.Opponent, name string)  {}
