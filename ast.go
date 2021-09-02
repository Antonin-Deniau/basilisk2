package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	Value string
	ParserRule string
	Type string
	Counter int
	Limit int
	Parent *Node
	Childs []*Node 
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

func ProcessInt(node *Node) (BType, error) {
	i, err := strconv.ParseInt(node.Value, 10, 64)
	if err != nil {
		return nil, err
	}

	return BInt{Value: int64(i)}, nil
}

func ProcessMeta(node *Node) (BType, error) {
	list := BList{Value: make([]*BType, 0)}

	name := BType(BName{ Value: "with-meta" })
	list.Value = append(list.Value, &name)

	for i := range node.Childs {
        i = len(node.Childs) - 1 - i
        expr := node.Childs[i]

		res_node, err := ProcessNode(expr)
		if err != nil {
			return nil, err
		}

		list.Value = append(list.Value, &res_node)
	}

	return list, nil
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
	case "Int":
		return ProcessInt(node)
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
	case "Meta":
		return ProcessMeta(node)
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
