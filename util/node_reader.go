package util

import "github.com/dnovikoff/tenhou/parser"

type NodeReader struct {
	Read  func() (string, error)
	nodes parser.Nodes
}

func (this *NodeReader) Next() (node *parser.Node, err error) {
	for len(this.nodes) == 0 {
		var message string
		message, err = this.Read()
		if err != nil {
			return
		}
		this.nodes, err = ParseXML(message)
		if err != nil {
			return
		}
	}
	node = &this.nodes[0]
	this.nodes = this.nodes[1:]
	return
}
