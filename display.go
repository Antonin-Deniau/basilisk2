package main

import (
	"strings"
	"errors"
	"fmt"
)

func DisplayBList(list *BList, sb *strings.Builder, readably bool) error {
	sb.WriteRune('(')
	for _, expr := range list.Value {
		err := Display(expr, sb, readably)
		if err != nil {
			return err
		}
	}
	sb.WriteRune(')')

	return nil
}

func DisplayBString(string *BString, sb *strings.Builder, readably bool) error {
	sb.WriteRune('"')
	sb.WriteString(string.Value)
	sb.WriteRune('"')
	return nil
}

func Display(expr *BType, sb *strings.Builder, readably bool) error {
	switch v := (*expr).(type) {
	case BList:
		return DisplayBList(&v, sb, readably)
	case BString:
		return DisplayBString(&v, sb, readably)
	default:
		return errors.New(fmt.Sprintf("Unable to find type for btype %+v\n", expr))
	}
}
