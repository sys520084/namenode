package internal

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
)

type Node struct {
	Name     string
	Children []Node
	Size     int
}

func (node Node) String() string {
	bs, _ := json.MarshalIndent(node, "", "  ")
	return string(bs)
}

func SplitPath(path string) []string {
	return strings.Split(path, string(filepath.Separator))
}

func DeleteNode(root []Node, names []string) []Node {
	if len(names) == 0 {
		return root
	}
	var (
		i int
	)

	for i = 0; i < len(root); i++ {
		if root[i].Name == names[0] { //already in tree
			if len(names) == 1 {
				return del(root, i)
			} else {
				root[i].Children = DeleteNode(root[i].Children, names[1:])
			}
		}
	}

	return root
}

func Copy(root []Node, source, target []string) ([]Node, error) {
	sourceNode, isExist := Get(root, source)
	if !isExist {
		return root, fmt.Errorf("not found source")
	}

	if sourceNode.Size == 0 || len(sourceNode.Children) != 0 {
		return root, fmt.Errorf("source is directory, not support")
	}

	targetNode, isExist := Get(root, target)
	if isExist {
		if targetNode.Size == 0 || len(targetNode.Children) != 0 {
			return root, fmt.Errorf("target is directory, not support")
		}
		root = DeleteNode(root, target)
	}

	root = AddToTree(root, target, sourceNode.Size)

	return root, nil
}

func Move(root []Node, source []string, target []string) ([]Node, error) {
	tmp, err := Copy(root, source, target)
	if err != nil {
		return root, fmt.Errorf("Copy failed:%s", err)
	}

	return DeleteNode(tmp, source), nil
}

//
func Get(root []Node, names []string) (Node, bool) {
	if len(names) == 0 {
		return Node{}, false
	}

	var (
		i int
	)

	for i = 0; i < len(root); i++ {
		currentNode := root[i]
		if currentNode.Name == names[0] { //already in tree
			if len(names) == 1 {
				return currentNode, true
			} else {
				return Get(currentNode.Children, names[1:])
			}
		}
	}

	return Node{}, false
}

// del index
func del(s []Node, i int) []Node {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

type NodeTree struct {
	Nodes []Node
}

func AddToTree(root []Node, names []string, size int) []Node {
	if len(names) > 0 {
		var i int
		for i = 0; i < len(root); i++ {
			if root[i].Name == names[0] { //already in tree
				break
			}
		}
		if i == len(root) {
			if len(names) == 1 {
				root = append(root, Node{Name: names[0], Size: size})
			} else {
				root = append(root, Node{Name: names[0]})
			}
		}

		root[i].Children = AddToTree(root[i].Children, names[1:], size)

	}
	return root
}

//var result []Node

func GetNodeChildren(root []Node, names []string) []Node {
	result := []Node{}
	if len(names) > 0 {
		var (
			i       int
			isExist bool
		)
		for i = 0; i < len(root); i++ {
			if root[i].Name == names[0] {
				isExist = true
				break
			}
		}

		if !isExist {
			return root
		}

		if len(names) == 1 {
			if isExist {
				//get reesult
				result = root[i].Children
			}
			return result
		}
		return GetNodeChildren(root[i].Children, names[1:])
	}

	return result

}
