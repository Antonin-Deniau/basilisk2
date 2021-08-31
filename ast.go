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

func ProcessQuote(node *Node) (BType, error) {
	list := BList{Value: make([]*BType, 0)}

	name := BType(BName{ Value: "quote" })

	list.Value = append(list.Value, &name)

	for _, child := range node.Childs {
		res_node, err := ProcessNode(child)
		if res_node == nil {
			return nil, errors.New("No expression after quote")
		}

		if err != nil {
			return nil, err
		}

		list.Value = append(list.Value, &res_node)
		return list, nil
	}

	return nil, nil
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
	string := BString{Value: Unescape(node.Value)}

	return string, nil
}

func ProcessName(node *Node) (BType, error) {
	return BName{Value: node.Value}, nil
}

func ProcessBool(node *Node) (BType, error) {
	return BBool{Value: node.Value == "true"}, nil
}

func ProcessNil(node *Node) (BType, error) {
	return BNil{}, nil
}

func ProcessNode(node *Node) (BType, error) {
	switch node.Type {
	case "Comment":
		return nil, nil
	case "Quote":
		return ProcessQuote(node)
	case "Expr":
		for _, expr := range node.Childs {
			return ProcessNode(expr)
		}
	case "Nil":
		return ProcessNil(node)
	case "Name":
		return ProcessName(node)
	case "Bool":
		return ProcessBool(node)
	case "List":
		return ProcessList(node)
	case "String":
		return ProcessString(node)
	default:
		return nil, errors.New(fmt.Sprintf("Unable to find type %s", node.Type))
	}

	return nil, errors.New("WTF")
}
