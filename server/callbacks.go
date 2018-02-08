// Generated with cb-generator. DO NOT EDIT
package server

import (
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type Callbacks struct {
	CbAuth               func(string)
	CbBye                func()
	CbCall               func(Answer, tile.Instances)
	CbCancelJoin         func()
	CbChat               func(string)
	CbDrop               func(tile.Instance)
	CbGoOK               func()
	CbHello              func(string, string, tbase.Sex)
	CbJoin               func(int, int, bool)
	CbNextReady          func()
	CbPing               func()
	CbReach              func(tile.Instance)
	CbRequestLobbyStatus func(int, int)

	Default func()
}

var _ Controller = &Callbacks{}

func (c *Callbacks) Auth(a0 string) {
	if c.CbAuth != nil {
		c.CbAuth(a0)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Bye() {
	if c.CbBye != nil {
		c.CbBye()
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Call(a0 Answer, a1 tile.Instances) {
	if c.CbCall != nil {
		c.CbCall(a0, a1)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) CancelJoin() {
	if c.CbCancelJoin != nil {
		c.CbCancelJoin()
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Chat(a0 string) {
	if c.CbChat != nil {
		c.CbChat(a0)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Drop(a0 tile.Instance) {
	if c.CbDrop != nil {
		c.CbDrop(a0)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) GoOK() {
	if c.CbGoOK != nil {
		c.CbGoOK()
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Hello(a0 string, a1 string, a2 tbase.Sex) {
	if c.CbHello != nil {
		c.CbHello(a0, a1, a2)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Join(a0 int, a1 int, a2 bool) {
	if c.CbJoin != nil {
		c.CbJoin(a0, a1, a2)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) NextReady() {
	if c.CbNextReady != nil {
		c.CbNextReady()
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Ping() {
	if c.CbPing != nil {
		c.CbPing()
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) Reach(a0 tile.Instance) {
	if c.CbReach != nil {
		c.CbReach(a0)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}

func (c *Callbacks) RequestLobbyStatus(a0 int, a1 int) {
	if c.CbRequestLobbyStatus != nil {
		c.CbRequestLobbyStatus(a0, a1)
		return
	}
	if c.Default != nil {
		c.Default()
	}
}
