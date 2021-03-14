package models

import (
	"fmt"
	"sort"
	"strings"
)

const EXPANDED_ARROW = "▼"
const COLLAPSED_ARROW = "►"

type StatusLineNode struct {
	Children  []*StatusLineNode
	File      *File
	Name      string
	Collapsed bool
}

func (s *StatusLineNode) GetShortStatus() string {
	// need to see if any child has unstaged changes.
	if s.IsLeaf() {
		return s.File.ShortStatus
	}

	firstChar := " "
	secondChar := " "
	if s.HasUnstagedChanges() {
		firstChar = "M"
	}
	if s.HasStagedChanges() {
		secondChar = "M"
	}

	return firstChar + secondChar
}

func (s *StatusLineNode) HasUnstagedChanges() bool {
	if s.IsLeaf() {
		return s.File.HasUnstagedChanges
	}

	for _, child := range s.Children {
		if child.HasUnstagedChanges() {
			return true
		}
	}

	return false
}

func (s *StatusLineNode) HasStagedChanges() bool {
	if s.IsLeaf() {
		return s.File.HasStagedChanges
	}

	for _, child := range s.Children {
		if child.HasStagedChanges() {
			return true
		}
	}

	return false
}

func (s *StatusLineNode) GetNodeAtIndex(index int) *StatusLineNode {
	node, _ := s.getNodeAtIndexAux(index)

	return node
}

func (s *StatusLineNode) getNodeAtIndexAux(index int) (*StatusLineNode, int) {
	offset := 1

	if index == 0 {
		return s, offset
	}

	for _, child := range s.Children {
		node, offsetChange := child.getNodeAtIndexAux(index - offset)
		offset += offsetChange
		if node != nil {
			return node, offset
		}
	}

	return nil, offset
}

func (s *StatusLineNode) IsLeaf() bool {
	return len(s.Children) == 0
}

func (s *StatusLineNode) Size() int {
	output := 1

	for _, child := range s.Children {
		output += child.Size()
	}

	return output
}

func (s *StatusLineNode) Flatten() []*StatusLineNode {
	arr := []*StatusLineNode{s}

	for _, child := range s.Children {
		arr = append(arr, child.Flatten()...)
	}

	return arr
}

// need a function which will return a formatted array of lines.

func (s *StatusLineNode) Render() []string {
	return s.renderAux(-1)
}

func (s *StatusLineNode) renderAux(depth int) []string {
	if s == nil {
		return []string{}
	}

	name := fmt.Sprintf("%s %s", s.GetShortStatus(), s.Name)

	if s.IsLeaf() {
		if depth == -1 {
			return []string{}
		}
		return []string{fmt.Sprintf("%s%s", strings.Repeat("  ", depth), name)}
	}

	if s.Collapsed {
		return []string{fmt.Sprintf("%s%s %s", strings.Repeat("  ", depth), name, COLLAPSED_ARROW)}
	}

	arr := []string{}
	if depth > -1 {
		arr = append(arr, fmt.Sprintf("%s%s %s", strings.Repeat("  ", depth), name, EXPANDED_ARROW))
	}

	for _, child := range s.Children {
		arr = append(arr, child.renderAux(depth+1)...)
	}

	return arr
}

func (s *StatusLineNode) SortTree() {
	s.sortChildren()

	for _, child := range s.Children {
		child.SortTree()
	}
}

func (s *StatusLineNode) sortChildren() {
	if s.IsLeaf() {
		return
	}

	sortedChildren := make([]*StatusLineNode, len(s.Children))
	copy(sortedChildren, s.Children)

	sort.Slice(sortedChildren, func(i, j int) bool {
		if !sortedChildren[i].IsLeaf() && sortedChildren[j].IsLeaf() {
			return true
		}
		if sortedChildren[i].IsLeaf() && !sortedChildren[j].IsLeaf() {
			return false
		}

		return sortedChildren[i].Name < sortedChildren[j].Name
	})

	// TODO: think about making this in-place
	s.Children = sortedChildren
}
