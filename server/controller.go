package server

import (
	"github.com/dnovikoff/tempai-core/tile"
	"github.com/dnovikoff/tenhou/tbase"
)

type Answer int

const (
	AnswerSkip      Answer = iota
	AnswerPon              // 1
	AnswerOpenedKan        // 2
	AnswerChi              // 3
	AnswerClosedKan        // 4
	AnswerChankan          // 5
	AnswerRon              // 6
	AnswerTsumo            // 7
	Answer8                // What is 8???
	AnswerDraw             // 9
	AnswerSanmaDora        // 10
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
