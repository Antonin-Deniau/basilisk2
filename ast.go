package main

import (
	"errors"
	"fmt"
)

type Node struct {
	Value string
	Type string
	Parent *Node
	Childs []*Node 
}


func FindRootNode(node *Node) *Node {
	curr_node := node
	for curr_node.Parent != nil {
		curr_node = curr_node.Parent
	}

	return curr_node
}


func ProcessList(node *Node) (BType, error) {
	list := BList{Value: make([]*BType, 0)}
	for _, child := range node.Childs {
		res_node, err := ProcessNode(child)
		if res_node == nil {
			continue
		}

		if err != nil {
			return nil, err
		}

		list.Value = append(list.Value, &res_node)
	}

	return list, nil
}

func ProcessString(node *Node) (BType, error) {
	string := BString{Value: node.Value}

	return string, nil
}

func ProcessNode(node *Node) (BType, error) {
	switch node.Type {
	case "comment":
		return nil, nil
	case "expr":
		for _, expr := range node.Childs {
			return ProcessNode(expr)
		}
	case "list":
		return ProcessList(node)
	case "string":
		return ProcessString(node)
	default:
		return nil, errors.New(fmt.Sprintf("Unable to find type %s", node.Type))
	}

	return nil, errors.New("WTF")
}
