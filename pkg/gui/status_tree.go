package gui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
)

func GetTreeFromStatusFiles(files []*models.File) *models.StatusLineNode {
	root := &models.StatusLineNode{}

	sort.SliceStable(files, func(i, j int) bool {
		return files[i].Name < files[j].Name
	})

	var curr *models.StatusLineNode
	for _, file := range files {
		split := strings.Split(file.Name, string(os.PathSeparator))
		curr = root
	outer:
		for i, dir := range split {
			var setFile *models.File
			if i == len(split)-1 {
				setFile = file
			}
			for _, existingChild := range curr.Children {
				if existingChild.Name == dir {
					curr = existingChild
					continue outer
				}
			}
			newChild := &models.StatusLineNode{
				Name: dir,
				File: setFile,
			}
			curr.Children = append(curr.Children, newChild)

			curr = newChild
		}
	}

	return root
}

func RenderStatusTree(root *models.StatusLineNode) string {
	return RenderStatusTreeAux(root, 0)
}

func RenderStatusTreeAux(root *models.StatusLineNode, depth int) string {
	// start at the root and go depth first
	indentation := strings.Repeat(" ", depth*2)

	if root.File != nil {
		return indentation + root.GetShortStatus() + " " + filepath.Base(root.File.Name) + "\n"
	}

	output := indentation + root.GetShortStatus() + " " + root.Name + "/\n"
	for _, child := range root.Children {
		output += RenderStatusTreeAux(child, depth+1)
	}

	return output
}
