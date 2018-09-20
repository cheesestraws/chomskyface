package main

import (
	"fmt"
	"strings"

	"github.com/hironobu-s/go-corenlp/document"
)

type TreeNode struct {
	Label string
	Depth int
	
	estimatedSizes bool
	eWidth int
	eHeight int
	x int
	y int
	
	Children []*TreeNode
}

func tfp(parse *document.Parse, level int) *TreeNode {
	// prune unnecessary root nodes
	if parse.Pos == "ROOT" && len(parse.Children) == 1 {
		return tfp(parse.Children[0], level)
	}
	
	node := &TreeNode{
		Depth: level,
	}
	
	if parse.Text != "" && parse.Pos != "" {
		node.Label = fmt.Sprintf("(%s) %s", parse.Pos, parse.Text)
	} else if parse.Text != "" {
		node.Label = parse.Text
	} else {
		node.Label = parse.Pos
	}
	
	// Now generate children
	for _, c := range parse.Children {
		node.Children = append(node.Children, tfp(c, level + 1))
	}
	
	return node
}

func TreeFromParse(parse *document.Parse) *TreeNode {
	return tfp(parse, 0)
}

func (t *TreeNode) PrettyPrint() {
	indent := strings.Repeat(" ", t.Depth)
	
	var dimensions string
	if t.estimatedSizes {
		dimensions = fmt.Sprintf("[x:%d, y: %d, w: %d, h: %d]", t.x, t.y, t.eWidth, t.eHeight)
	}
	
	fmt.Printf("%s%s %s\n", indent, t.Label, dimensions)
	for _, c := range t.Children {
		c.PrettyPrint()
	}
}

