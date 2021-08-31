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
