package internal

import (
	"fmt"
	"strings"
)

type Node struct {
	Name     string
	Children []Node
	Size     int
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

func main() {
	s := []string{
		"test/1/3/1.jpg",
		"test/1/3/2.jpg",
		"test/1/3/3.jpg",
		"test/2/5/2.jpg",
		"test3/3/3.jpg",
	}
	//	var tree []Node
	mytree := &NodeTree{}
	for i := range s {
		mytree.Nodes = AddToTree(mytree.Nodes, strings.Split(s[i], "/"), 500)

	}
	//tmpNode := []Node{}
	fmt.Println(mytree.Nodes[0].Children)
	a := "test/1"
	c := GetNodeChildren(mytree.Nodes, strings.Split(a, "/"))
	c1 := GetNodeChildren(mytree.Nodes, strings.Split(a, "/"))

	fmt.Println("c is:", c)
	fmt.Println("c is:", c1)
	//fmt.Println(mytree.Nodes)
	//b, err := json.Marshal(tree)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(b))
}
