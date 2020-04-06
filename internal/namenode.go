package internal

import (
	//	"encoding/json"

	"fmt"
	"strings"
)

type Node struct {
	Name     string
	Children []Node
	size     int
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
				fmt.Println(names)
				root = append(root, Node{Name: names[0], size: size})
			} else {
				root = append(root, Node{Name: names[0]})
			}
		}

		root[i].Children = AddToTree(root[i].Children, names[1:], size)

	}
	return root
}

func main() {
	s := []string{
		"test/1/1.jpg",
		"test/2/2.jpg",
		"test/3/3.jpg",
	}
	//	var tree []Node
	mytree := &NodeTree{}
	for i := range s {
		mytree.Nodes = AddToTree(mytree.Nodes, strings.Split(s[i], "/"), 500)

	}
	fmt.Println(mytree.Nodes)
	//b, err := json.Marshal(tree)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(b))
}
