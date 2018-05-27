package parser

import (
	"encoding/xml"
	"strings"

	"github.com/facebookgo/stackerr"
)

func ParseXML(input string) (ret Nodes, err error) {
	// Dirty hack
	input = "<mjloggm>" + input + "</mjloggm>"
	d := xml.NewDecoder(strings.NewReader(input))
	d.Strict = false
	var root Root
	err = stackerr.Wrap(d.Decode(&root))
	for k, v := range root.Nodes {
		if len(v.Attributes) == 0 {
			root.Nodes[k].Attributes = nil
		}
	}
	ret = root.Nodes
	return
}
