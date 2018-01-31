package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dnovikoff/tenhou/tbase"
)

func TestNodeScoreChanged(t *testing.T) {
	node := Node{Name: "AGARI", Attributes: map[string]string{"sc": "302,-323,236,0,149,323,313,0"}}
	assert.Equal(t, tbase.ScoreChanges{
		tbase.ScoreChange{Score: 30200, Diff: -32300},
		tbase.ScoreChange{Score: 23600, Diff: 0},
		tbase.ScoreChange{Score: 14900, Diff: 32300},
		tbase.ScoreChange{Score: 31300, Diff: 0}},
		node.GetScoreChanges())
}

func TestNodeScore(t *testing.T) {
	node := Node{Name: "AGARI", Attributes: map[string]string{"ten": "266,250,221,263"}}
	assert.Equal(t, tbase.Scores{26600, 25000, 22100, 26300}, node.GetScores())
}

// func TestNodeUserList3(t *testing.T) {
// 	node := Node{Name: "UL", Attributes: map[string]string{
// 		"n0": "n0", "n1": "n1", "n2": "n2", "n3": "",
// 		"dan":  "16,20,16,0",
// 		"rate": "2128,2331,2111,1500",
// 		"sx":   "M,M,F,C",
// 	}}
// 	users, err := node.UserList()
// 	require.NoError(t, err)
// 	assert.Equal(t,
// 		tbase.UserList{
// 			tbase.User{Num: 0, Name: "n0", Dan: 16, Rate: 2128, Sex: 1},
// 			tbase.User{Num: 1, Name: "n1", Dan: 20, Rate: 2331, Sex: 1},
// 			tbase.User{Num: 2, Name: "n2", Dan: 16, Rate: 2111, Sex: 2},
// 		}, users)
// }

// func TestNodeUserList4(t *testing.T) {
// 	node := Node{Name: "UL", Attributes: map[string]string{
// 		"n0":   "%E4%B8%80%E6%96%B9%E9%80%9A%E8%A1%8C",
// 		"n1":   "%E5%88%BB%E5%AD%90%E5%B8%82%E6%B0%91",
// 		"n2":   "%E5%AF%82%E5%A4%9C%E9%9C%BD%E6%9C%88",
// 		"n3":   "%E8%80%81%E9%A6%AC",
// 		"dan":  "16,16,16,16",
// 		"rate": "2054,2051,2062,2019",
// 		"sx":   "M,M,M,M",
// 	}}
// 	users, err := node.UserList()
// 	require.NoError(t, err)
// 	assert.Equal(t, 4, len(users))
// }
