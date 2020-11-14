package models

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
	node, _ := s.GetNodeAtIndexAux(index)

	return node
}

func (s *StatusLineNode) GetNodeAtIndexAux(index int) (*StatusLineNode, int) {
	offset := 1

	if index == 0 {
		return s, offset
	}

	for _, child := range s.Children {
		node, offsetChange := child.GetNodeAtIndexAux(index - offset)
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
