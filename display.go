package main

import (
	"strconv"
	"strings"
	"errors"
	"fmt"
)

func DisplayBList(node *BList, sb *strings.Builder, readably bool) error {
	sb.WriteRune('(')
	end := len(node.Value) - 1
	for i, expr := range node.Value {
		err := Display(expr, sb, readably)
		if err != nil {
			return err
		}

		if i != end {
			sb.WriteRune(' ')
		}
	}
	sb.WriteRune(')')

	return nil
}

func DisplayBName(node *BName, sb *strings.Builder, readably bool) error {
	sb.WriteString(node.Value)
	return nil
}

func DisplayBString(node *BString, sb *strings.Builder, readably bool) error {
	sb.WriteRune('"')
	if readably == true {
		PrStr(sb, node.Value, readably)
	} else {
		sb.WriteString(node.Value)
	}
	sb.WriteRune('"')
	return nil
}

func DisplayBNil(node *BNil, sb *strings.Builder, readably bool) error {
	sb.WriteString("nil")
	return nil
}

func DisplayBInt(node *BInt, sb *strings.Builder, readably bool) error {
	s := strconv.FormatInt(node.Value, 10)
	sb.WriteString(s)
	return nil
}

func DisplayBBool(node *BBool, sb *strings.Builder, readably bool) error {
	if node.Value == true {
		sb.WriteString("true")
	} else {
		sb.WriteString("false")
	}
	return nil
}

func DisplayKeyword(node *BKeyword, sb *strings.Builder, readably bool) error {
	sb.WriteRune(':')
	sb.WriteString(node.Value)
	return nil
}

func Display(node *BType, sb *strings.Builder, readably bool) error {
	switch v := (*node).(type) {
	case BInt:
		return DisplayBInt(&v, sb, readably)
	case BBool:
		return DisplayBBool(&v, sb, readably)
	case BNil:
		return DisplayBNil(&v, sb, readably)
	case BList:
		return DisplayBList(&v, sb, readably)
	case BName:
		return DisplayBName(&v, sb, readably)
	case BKeyword:
		return DisplayKeyword(&v, sb, readably)
	case BString:
		return DisplayBString(&v, sb, readably)
	default:
		return errors.New(fmt.Sprintf("Unable to find type for btype %+v\n", *node ))
	}
}



/*
func Display(x *BType, readably bool) (string, error) {
	if x.BBoolean != nil {
		if x.BBoolean.Value {
			return "true", nil
		} else {
			return "false", nil
		}
	}

	case BFunc:
		return "#<function>"

	if x.BList != nil {
		var res strings.Builder

		for _, s := range x.BList.Values {
			o, err := Display(s, readably)
			if err != nil {
				return "", err
			}
			res.WriteString(o)
			res.WriteString(" ")
		}

		r, size := utf8.DecodeLastRuneInString(res.String())
		if r == utf8.RuneError && (size == 0 || size == 1) {
			size = 0
		}

		return fmt.Sprintf("(%s)", res.String()[:len(res.String())-size]), nil
	}

	if x.BNumber != nil {
		return betterFormat(x.BNumber.Value), nil
	}

	if x.BString != nil {
		if readably {
			return fmt.Sprintf("\"%s\"", PrStr(x.BString.Value, readably)), nil
		} else {
			return x.BString.Value, nil
		}
	}

	if x.BVector != nil {
		var res strings.Builder

		for _, s := range x.BVector.Values {
			o, err := Display(s, readably)
			if err != nil {
				return "", err
			}
			res.WriteString(o)
			res.WriteString(" ")
		}

		r, size := utf8.DecodeLastRuneInString(res.String())
		if r == utf8.RuneError && (size == 0 || size == 1) {
			size = 0
		}

		return fmt.Sprintf("[%s]", res.String()[:len(res.String())-size]), nil
	}

	if x.BHashMap != nil {
		var res strings.Builder

		for _, s := range x.BHashMap.Map {
			o, err := Display(s.Key, readably)
			res.WriteString(o)
			if err != nil {
				return "", err
			}
			res.WriteString(" ")
			o1, err1 := Display(s.Value, readably)
			res.WriteString(o1)
			if err1 != nil {
				return "", err1
			}
			res.WriteString(" ")
		}

		r, size := utf8.DecodeLastRuneInString(res.String())
		if r == utf8.RuneError && (size == 0 || size == 1) {
			size = 0
		}

		return fmt.Sprintf("{%s}", res.String()[:len(res.String())-size]), nil
	}
	
	if x.BKeyword != nil {
		return fmt.Sprintf(":%s", x.BKeyword.Value), nil
	}

	if x.BUnquote != nil {
		o3, err3 := Display(x.BUnquote.Value, readably)
		if err3 != nil {
			return "", err3
		}

		return fmt.Sprintf("'%s", o3), nil
	}

	if x.BQuasiquote != nil {
		o2, err2 := Display(x.BQuasiquote.Value, readably)
		if err2 != nil {
			return "", err2
		}

		return fmt.Sprintf("`%s", o2), nil
	}

	case BException:
		return Display(val.Message, readably)

	if x.BName != nil {
		return x.BName.Value, nil
	}

	case BAtom:
		return "(atom {})".format(Display(val.data, readably))

	if x.BNil != nil {
		return "nil", nil
	}

	panic(fmt.Sprintf("Unable to display: %+v", x))
}
*/

func Unescape(sb *strings.Builder, s string) {
	// strconv.Unquote
	esc := false

	for _, i := range s {
		if i == '\\' && esc == false {
			esc = true
		} else if i == '\\' && esc == true {
			sb.WriteRune('\\')
			esc = false
		} else if i == 'n' && esc == true {
			sb.WriteRune('\n')
			esc = false
		} else if i == '"' && esc == true {
			sb.WriteRune('"')
			esc = false
		} else {
			sb.WriteRune(i)
			esc = false
		}
	}
}

func Escape(sb *strings.Builder, s string) {
	for _, i := range s {
		if i == '\\' {
			sb.WriteString("\\\\")
		} else if i == '"' {
			sb.WriteString("\\\"")
		} else if i == '\n' {
			sb.WriteString("\\n")
		} else {
			sb.WriteRune(i)
		}
	}
}

func PrStr(sb *strings.Builder, x string, readably bool) {
     if readably {
	     Escape(sb, x)
     } else {
	     sb.WriteString(x)
     }
}

func betterFormat(num float64) string {
    s := fmt.Sprintf("%.6f", num)
    return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

/*
func ParseFile(data string) (*Program) {
	tree := &Program{}
	err := parser.ParseString(data, tree)
	if err != nil {
		return err, nil
	}

	return nil, tree
}
*/

func Prnt(e *BType) (error) {
	var sb strings.Builder

	err := Display(e, &sb, true)
	if err != nil {
		return err
	}
	fmt.Print(sb.String() + "\n")
	return nil
}
