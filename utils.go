package main

import (
	"fmt"
	"strings"
)

// ------------------------------------------DEV UTILS------------------------------------


func DisplayNode(node *Node, deep int) {
	if len(node.Childs) == 0 {
		fmt.Printf("%s(%s = [%s])\n", strings.Repeat(" ", deep), node.Type, node.Value)
	} else {
		fmt.Printf("%s(%s = [%s]\n", strings.Repeat(" ", deep), node.Type, node.Value)
		for _, child := range node.Childs {
			DisplayNode(child, deep + 2)
		}
		fmt.Printf("%s)\n", strings.Repeat(" ", deep))
	}
}

// ------------------------------------------GENERAL UTILS------------------------------------