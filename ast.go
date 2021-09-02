package main

import (
	"errors"
	"fmt"
	"strings"
)

type Node struct {
	Value string
	ParserRule string
	Type string
	Repeat bool
	Validated bool
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

	for _, expr := range node.Childs {
		res_node, err := ProcessNode(expr)
		if err != nil {
			return nil, err
		}

		if res_node == nil {
			continue
		}
		list.Value = append(list.Value, &res_node)

		return list, nil
	}


	return nil, errors.New("No expression after quote")
}

func ProcessList(node *Node) (BType, error) {
	list := BList{Value: make([]*BType, 0)}
	for _, child := range node.Childs {
		res_node, err := ProcessNode(child)
		if err != nil {
			return nil, err
		}

		if res_node == nil {
			continue
		}

		list.Value = append(list.Value, &res_node)
	}

	return list, nil
}

func ProcessString(node *Node) (BType, error) {
	var sb strings.Builder
	Unescape(&sb, node.Value)

	string := BString{Value: sb.String()}

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

func ProcessKeyWord(node *Node) (BType, error) {
	return BKeyword{Value: node.Value}, nil
}

func ProcessNode(node *Node) (BType, error) {
	//fmt.Printf("=> %s = %s\n", node.Type, node.Value)
	switch node.Type {
	case "Comment":
		return nil, nil
	case "Whitespace":
		return nil, nil
	case "Quote":
		return ProcessQuote(node)
	case "Expr":
		for _, expr := range node.Childs {
			res_node, err := ProcessNode(expr)
			if err != nil {
				return nil, err
			}

			return res_node, nil
		}
	case "Keyword":
		return ProcessKeyWord(node)
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
