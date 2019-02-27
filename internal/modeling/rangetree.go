package modeling

import (
	"github.com/lavalamp-/ipv666/internal/logging"
	"github.com/spf13/viper"
)

type RangeTree struct {
	ChildrenCount		uint64
	Children			map[uint16]*RangeTreeNode
}

type RangeTreeNode struct {
	ChildrenCount		uint64
	Children			map[uint16]*RangeTreeNode
	Depth				int
}

func NewRangeTree() *RangeTree {
	return &RangeTree{
		Children:		make(map[uint16]*RangeTreeNode),
		ChildrenCount:	0,
	}
}

func NewRangeTreeFromRanges(toAdd []*GenRange) *RangeTree {
	toReturn := NewRangeTree()
	added, skipped := toReturn.AddRanges(toAdd)
	logging.Infof("Added %d ranges to range tree (skipped %d).", added, skipped)
	return toReturn
}

func newRangeTreeNode(depth int) *RangeTreeNode {
	return &RangeTreeNode{
		Children:		make(map[uint16]*RangeTreeNode),
		ChildrenCount:	0,
		Depth:			depth,
	}
}

func (rangeTree *RangeTree) AddRange(toAdd *GenRange) bool {
	addNybbles := toAdd.GetTreeNybbles()
	if rangeTree.containsRangeByNybbles(addNybbles) {
		return false
	}
	if _, ok := rangeTree.Children[addNybbles[0]]; !ok {
		rangeTree.Children[addNybbles[0]] = newRangeTreeNode(1)
	}
	rangeTree.Children[addNybbles[0]].addNybbles(addNybbles[1:])
	rangeTree.ChildrenCount++
	return true
}

func (rangeTree *RangeTree) AddRanges(toAdd []*GenRange) (int, int) {
	added, skipped := 0, 0
	for i, curAdd := range toAdd {
		if i % viper.GetInt("LogLoopEmitFreq") == 0 {
			logging.Infof("Adding range %d out of %d to RangeTree.", i, len(toAdd))
		}
		if rangeTree.AddRange(curAdd) {
			added++
		} else {
			skipped++
		}
	}
	return added, skipped
}

func (rangeTree *RangeTree) ContainsRange(toCheck *GenRange) bool {
	return rangeTree.containsRangeByNybbles(toCheck.GetTreeNybbles())
}

func (rangeTree *RangeTree) containsRangeByNybbles(toCheck []uint16) bool {
	if toCheck[0] == 16 {
		if val, ok := rangeTree.Children[16]; ok {
			return val.containsNybbles(toCheck[1:])
		} else {
			return false
		}
	} else {
		if val, ok := rangeTree.Children[16]; ok {
			if val.containsNybbles(toCheck[1:]) {
				return true
			}
		}
		if val, ok := rangeTree.Children[toCheck[0]]; ok {
			if val.containsNybbles(toCheck[1:]) {
				return true
			}
		}
		return false
	}
}

func (rangeTreeNode *RangeTreeNode) containsNybbles(toCheck []uint16) bool {
	if len(toCheck) == 1 {
		if _, ok := rangeTreeNode.Children[16]; ok {
			return true
		} else if _, ok := rangeTreeNode.Children[toCheck[0]]; ok {
			return true
		} else {
			return false
		}
	} else {
		if val, ok := rangeTreeNode.Children[16]; ok {
			if val.containsNybbles(toCheck[1:]) {
				return true
			}
		}
		if val, ok := rangeTreeNode.Children[toCheck[0]]; ok {
			if val.containsNybbles(toCheck[1:]) {
				return true
			}
		}
		return false
	}
}

func (rangeTreeNode *RangeTreeNode) addNybbles(nybbles []uint16) {
	if len(nybbles) == 0 {
		return
	} else if _, ok := rangeTreeNode.Children[nybbles[0]]; !ok {
		rangeTreeNode.Children[nybbles[0]] = newRangeTreeNode(rangeTreeNode.Depth + 1)
	}
	rangeTreeNode.Children[nybbles[0]].addNybbles(nybbles[1:])
	rangeTreeNode.ChildrenCount++
}