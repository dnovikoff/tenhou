package parser

import (
	"encoding/xml"

	"github.com/facebookgo/stackerr"

	"github.com/dnovikoff/tenhou/tbase"
)

type Nodes []Node

type Root struct {
	XMLName xml.Name `xml:"mjloggm"`
	Nodes   Nodes    `xml:",any"`
}

func (node *Node) GetInit() (x tbase.Init, err error) {
	x.Seed, err = tbase.ParseSeed(node.String("seed"))
	if err != nil {
		return
	}
	x.Dealer = node.GetDealer()
	x.Scores = node.GetScores()
	return
}

func ParseRyuukyoku(node *Node) (r *tbase.Ryuukyoku, err error) {
	t := node.String("type")
	dt := tbase.DrawMap[t]
	if dt == tbase.DrawUnknown {
		err = stackerr.Newf("Unknown draw type '%s'", t)
		return
	}
	status, err := node.GetTableStatus()
	if err != nil {
		return
	}
	r = &tbase.Ryuukyoku{}
	r.ScoreChanges = node.GetScoreChanges()
	r.Hands = node.GetHands()
	r.Finals = node.GetFinalScores()
	r.TableStatus = status
	r.DrawType = dt
	return
}

func (this *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	switch t.(type) {
	case xml.EndElement:
	default:
		return stackerr.Newf("Unexpected element %v", t)
	}

	this.Name = start.Name.Local
	attrs := make(map[string]string, len(start.Attr))
	for _, x := range start.Attr {
		attrs[x.Name.Local] = x.Value
	}
	this.Attributes = attrs
	return nil
}
