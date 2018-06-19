package tbase

import (
	"github.com/dnovikoff/tempai-core/compact"
	"github.com/dnovikoff/tempai-core/hand/calc"
	"github.com/dnovikoff/tempai-core/tile"
)

type Called struct {
	Type     CallType
	Opponent Opponent
	Tiles    tile.Instances
	Called   tile.Instance
	Upgraded tile.Instance
	// Core representation of meld for tempai calculator
	Core calc.Meld
}

type CalledList []*Called

func (cl CalledList) Core() calc.Melds {
	melds := make(calc.Melds, len(cl))
	for k, v := range cl {
		melds[k] = v.Core
	}
	return melds
}

func (cl CalledList) Add(x compact.Instances) {
	for _, v := range cl {
		v.Add(x)
	}
}

func (c *Called) IsKan() bool {
	switch c.Type {
	case Kan, UpgrdedKan, ClosedKan:
		return true
	}
	return false
}

func (c *Called) Add(x compact.Instances) {
	x.Add(c.Tiles)
	if c.Called != tile.InstanceNull {
		x.Set(c.Called)
	}
	if c.Type == UpgrdedKan && c.Upgraded != tile.InstanceNull {
		x.Set(c.Upgraded)
	}
}

type CallType int

const (
	Chi CallType = iota + 1
	Pon
	Kan
	ClosedKan
	UpgrdedKan
)

type Opponent int

const (
	Self Opponent = iota
	Right
	Front
	Left
)
