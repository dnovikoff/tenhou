package server

import (
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type Answer int

const (
	AnswerSkip      Answer = 0
	AnswerPon       Answer = 1
	AnswerOpenedKan Answer = 2
	AnswerChi       Answer = 3
	AnswerClosedKan Answer = 4
	AnswerChankan   Answer = 5
	AnswerRon       Answer = 6
	AnswerTsumo     Answer = 7
	// What is 8???
	AnswerDraw Answer = 9
)

// This is how server looks from client point of view
type Controller interface {
	Hello(name string, tid string, sex tbase.Sex)
	Auth(value string)
	// TODO: research values
	// Could be both big and small
	RequestLobbyStatus(v, V int) // PXR
	Join(lobbyNumber int, lobbyType int, rejoin bool)
	CancelJoin()
	Drop(t tile.Instance)
	Call(a Answer, tiles tile.Instances)
	Reach(t tile.Instance)
	Ping() // Z
	GoOK()
	NextReady()
	Bye()
	Chat(message string)
}
