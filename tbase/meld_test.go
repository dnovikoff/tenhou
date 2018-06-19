package tbase

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tempai-core/tile"
)

func TestMeldsBase(t *testing.T) {
	for _, v := range []struct {
		expected tile.Tile
		m        Meld
	}{
		{tile.North, 46633},
		{tile.White, 47690},
		{tile.Green, 49705},
	} {
		t.Run(v.expected.String()+strconv.Itoa(int(v.m)), func(t *testing.T) {
			assert.Equal(t, v.expected, v.m.Decode().Core.Tile())
		})
	}
}

func TestMelds(t *testing.T) {
	for _, v := range []struct {
		expected string
		m        Meld
	}{
		{"444z Right", 46633},
		{"555z Front", 47690},
		{"666z Right", 49705},
		// Upgraded
		{"5555z Front", 47666},
		{"2222m Left", 1539},
		{"1111z Front", 27906},
		{"1111p Right", 9217},
		{"444s Front", 33354},
	} {
		t.Run(v.expected, func(t *testing.T) {
			d := v.m.Decode()
			actual := d.Core.Tiles().String() + " " + d.Opponent.String()
			assert.Equal(t, v.expected, actual)
		})
	}
}

func TestMeldsReconvert(t *testing.T) {
	for _, v := range []Meld{
		// chi [6]78p
		37199,
		33354,
		// Closed green kan with copyId = 1
		33024,
		// Closed green kan with copyId = 0
		32768,
	} {
		t.Run(strconv.Itoa(int(v)), func(t *testing.T) {
			assert.Equal(t, v, v.Decode().Encode())
		})
	}
}
