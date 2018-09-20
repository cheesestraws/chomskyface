package main

import (
	"github.com/fogleman/gg"
)

const NODE_PADDING = 5
const NODE_MARGIN = 5

func (t *TreeNode) calculateSizes(gg *gg.Context) {
	// estimate our size
	tw, th := gg.MeasureString(t.Label)
	
	w := NODE_MARGIN * 2 + NODE_PADDING * 2 + int(tw)
	h := NODE_MARGIN * 2 + NODE_PADDING * 2 + int(th)
	
	t.eWidth = w
	t.eHeight = h
	
	t.estimatedSizes = true

	// and do the same for each child
	for _, c := range t.Children {
		c.calculateSizes(gg)
	}
}

func (t *TreeNode) calculateXes(gg *gg.Context) {
	// first fill in the x positions of all the leaf nodes
	t.childlessCX(gg, nil)
	
	// then work up the tree
	t.otherCX(gg)
}

// calculate the X coordinates for childless nodes
func (t *TreeNode) childlessCX(gg *gg.Context, lastLeaf *TreeNode) *TreeNode {
	// if I have children, run myself recursively over all of them,
	// threading the state through as needed
	leaf := lastLeaf
	if len(t.Children) > 0 {
		for _, c := range t.Children {
			leaf = c.childlessCX(gg, leaf)
		}
		
		return leaf
	} else {
		// I do not have children
		if leaf == nil {
			// we are the first childless node.  X=0
			t.x = 0
		} else {
			// we are the previous childless node's x + its width
			t.x = leaf.x + leaf.eWidth
		}
		
		return t
	}
}

func (t *TreeNode) otherCX(gg *gg.Context) {
	if len(t.Children) == 0 {
		// leaf nodes have already been done in childlessCX
		return
	}
	
	// make sure all our children have X co-ordinates
	for _, c := range t.Children {
		c.otherCX(gg)
	}
	
	// the furthest left of our children is the first child's x
	furthestLeft := t.Children[0].x
	
	// the furthest right is the last child's x + its width
	furthestRight := t.Children[len(t.Children)-1].x + t.Children[len(t.Children)-1].eWidth
	
	// our midpoint is the (furthest right - furthest left / 2)	relative to our first child
	// then + furthest left to offset it back into our real coordinate space thing
	midpoint := ((furthestRight - furthestLeft) / 2) + furthestLeft
	
	t.x = midpoint - (t.eWidth / 2)
}

func (t *TreeNode) draw(gg *gg.Context) {
	// Draw my children
	for _, c := range t.Children {
		c.draw(gg)
	}

	// temporary bodge
	t.y = t.Depth * 50

	gg.SetRGB(0,0,0)

	// Draw me
	gg.DrawRectangle(float64(t.x + NODE_MARGIN), 
		float64(t.y + NODE_MARGIN),
		float64(t.eWidth - (NODE_MARGIN * 2)), 
		float64(t.eHeight - (NODE_MARGIN * 2)))
	gg.Stroke()
	
	gg.DrawStringAnchored(t.Label, float64(t.x + NODE_MARGIN + NODE_PADDING),
		float64(t.y + NODE_MARGIN + NODE_PADDING), 0, 1)
	
	// draw the lines to my children
	// start point always the same: my bottom and halfway along my width
	x1 := t.x + (t.eWidth / 2)
	y1 := t.y + t.eHeight - NODE_MARGIN
	
	for _, c := range t.Children {
		x2 := c.x + (c.eWidth / 2)
		y2 := c.y + NODE_MARGIN
		
		gg.DrawLine(float64(x1), float64(y1), float64(x2), float64(y2))
		gg.Stroke()
	}
}

func (t *TreeNode) PrepareToDraw(gg *gg.Context) {
	t.calculateSizes(gg)
	t.calculateXes(gg)
}

func (t *TreeNode) Draw(gg *gg.Context) {
	t.draw(gg)
}

func (t *TreeNode) GWidth() int {
	// are we childless?  If so, our width is just our own
	if len(t.Children) == 0 {
		return t.eWidth
	}
	
	// if not, it's the sum of all our children's widths
	var accum int
	for _, c := range t.Children {
		accum += c.GWidth()
	}
	return accum
}

func (t *TreeNode) MaxDepth() int {
	if len(t.Children) == 0 {
		return t.Depth
	} else {
		var maxDepth int
		for _, c := range t.Children {
			if c.MaxDepth() > maxDepth {
				maxDepth = c.MaxDepth()
			}
		}
		return maxDepth
	}
}

func (t *TreeNode) GHeight() int {
	return (t.MaxDepth() + 1) * 50
}
