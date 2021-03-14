package gui

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/sirupsen/logrus"
)

type StatusLineManager struct {
	Files    []*models.File
	Tree     *models.StatusLineNode
	TreeMode bool
	Log      *logrus.Entry
}

func (m *StatusLineManager) GetItemAtIndex(index int) *models.StatusLineNode {
	if m.TreeMode {
		// need to traverse the three depth first until we get to the index.
		return m.Tree.GetNodeAtIndex(index + 1) // ignoring root
	}

	m.Log.Warn(index)
	if index > len(m.Files)-1 {
		return nil
	}

	return &models.StatusLineNode{File: m.Files[index]}
}

func (m *StatusLineManager) GetAllItems() []*models.StatusLineNode {
	return m.Tree.Flatten()[1:] // ignoring root
}

func (m *StatusLineManager) GetItemsLength() int {
	return m.Tree.Size() - 1 // ignoring root
}

func (m *StatusLineManager) GetAllFiles() []*models.File {
	return m.Files
}

func (m *StatusLineManager) SetFiles(files []*models.File) {
	m.Files = files
	m.Tree = GetTreeFromStatusFiles(files)
}

func (m *StatusLineManager) Render() []string {
	return m.Tree.Render() // in this case the root is ignored in Tree.Render() itself
}
