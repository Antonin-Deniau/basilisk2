package main

import (
	"fmt"
	"os/user"
	"log"
	"github.com/chzyer/readline"
	"strings"
	"unicode/utf8"
)

func Unescape(s string) string {
	// strconv.Unquote
	var res strings.Builder
	esc := false

	for _, i := range s {
		if i == '\\' && esc == false {
			esc = true
		} else if i == '\\' && esc == true {
			res.WriteRune('\\')
			esc = false
		} else if i == 'n' && esc == true {
			res.WriteRune('\n')
			esc = false
		} else if i == '"' && esc == true {
			res.WriteRune('"')
			esc = false
		} else {
			res.WriteRune(i)
			esc = false
		}
	}

	return res.String()
}

func Escape(s string) string {
	var res strings.Builder
	for _, i := range s {
		if i == '\\' {
			res.WriteString("\\\\")
		} else if i == '"' {
			res.WriteString("\\\"")
		} else if i == '\n' {
			res.WriteString("\\n")
		} else {
			res.WriteRune(i)
		}
	}

	return res.String()
}


func PrStr(x string, readably bool) string {
     if readably {
	     return Escape(x)
     } else {
	     return x
     }
}

func InitParser() *BType {
	blist := NewExpr("(", Repeat(Token("expr")), ")")
	bmetadata := NewExpr("^", Token("expr"), Token("expr"))
	bderef := NewExpr("@", Token("expr"))
	bhashmap := NewExpr("{", Repeat(Token("expr")), "}")
	bvector := NewExpr("[", Repeat(Token("expr")), "]")
	bkeyword := NewRegex(":[^\"^.@~`\\[\\]:{}&'0-9\\s,();][^\"^@~`\\[\\]:{}\\s();]*")
	bquote := NewExpr("'", Token("expr"))
	bquasiquote := NewExpr("`", Token("expr"))

	root := NewRoot(
		blist,
		bmetadata,
		bderef,
		bhashmap,
		bvector,
		bkeyword,
		bquote,
		bquasiquote,
		bspliceunquote,
		bunquote,
		bunquote,
		bcomma,
		btoken,
		bcomment,
		bnumber,
		bboolean,
		bstring,
		bvariadic,
		bnil,
	)
	ignore := NewIgnore()

	return NewParser(root, ignore)
}

func Parse(data *string, parser *BType) (*BType, error) {
	index := 0
	parserState := InitParser()
	totalLen := len(*data)

	for index < totalLen {
		for _, x := range parser.BList.Values {

		}
	}

	return nil, nil// TODO
}

func betterFormat(num float64) string {
    s := fmt.Sprintf("%.6f", num)
    return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func Display(x *BType, readably bool) (string, error) {
	if x.BBoolean != nil {
		if x.BBoolean.Value {
			return "true", nil
		} else {
			return "false", nil
		}
	}

	/*
	case BFunc:
		return "#<function>"
	*/

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

	/*
	case BException:
		return Display(val.Message, readably)
	*/

	if x.BName != nil {
		return x.BName.Value, nil
	}

	/*
	case BAtom:
		return "(atom {})".format(Display(val.data, readably))
	*/
	if x.BNil != nil {
		return "nil", nil
	}

	panic(fmt.Sprintf("Unable to display: %+v", x))
}

func Prnt(e *BType) (error) {
	val, err := Display(e, true)
	if err != nil {
		return err
	}
	fmt.Print(val + "\n")
	return nil
}

func main() {
	// INIT READLINE CONFIG
	usr, err := user.Current()
	if err != nil {
		log.Fatal( err )
	}

	config := &readline.Config{
		Prompt: "basilisk> ",
		HistoryFile: usr.HomeDir + "/.basilisk_history",
	}
	rl, err := readline.NewEx(config)

	if err != nil {
		panic(err)
	}
	defer rl.Close()

	parser := InitParser()

	for {
		line, err := rl.Readline()
		if err != nil {
			log.Print(fmt.Sprintf("Exception: %s", err))
			break
		}

		res, err1 := Parse(&line, parser)
		if err1 != nil {
			log.Print(fmt.Sprintf("Exception: %s", err1))
		} else {

			err2 := Prnt(res)
			if err2 != nil {
				log.Print(fmt.Sprintf("Exception: %s", err2))
			}
		}

		rl.SaveHistory(line)
	}
}

