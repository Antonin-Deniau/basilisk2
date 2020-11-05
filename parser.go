package main

import (
	"github.com/alecthomas/participle/lexer/stateful"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle"
	"fmt"
	"os/user"
	"strings"
	"errors"
	"log"
	"github.com/chzyer/readline"
	"unicode/utf8"
)

type Program struct {
	Pos lexer.Position

	Lines []*BType `@@*`
}

func Unescape(s string) string {
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

func Display(x interface{}, readably bool) (string, error) {
	switch val := x.(type) {
	case BBoolean:
		if val.Value {
			return "true", nil
		} else {
			return "false", nil
		}

	/*
	case BFunc:
		return "#<function>"
	*/

	case BList:
		var res strings.Builder

		for _, s := range val.Values {
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

	case BNumber:
		return fmt.Sprintf("%f", val.Value), nil

	case BString:
		if readably {
			return fmt.Sprintf("\"%s\"", PrStr(val.Value, readably)), nil
		} else {
			return val.Value, nil
		}

	case BVector:
		var res strings.Builder

		for _, s := range val.Values {
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

	case BHashMap:
		var res strings.Builder

		for _, s := range val.Map {
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
	
	case BKeyword:
		return fmt.Sprintf(":%s", val.Value), nil

	/*
	case BException:
		return Display(val.Message, readably)
	*/

	case BName:
		return val.Value, nil

	/*
	case BAtom:
		return "(atom {})".format(Display(val.data, readably))
	*/
	case BNil:
		return "nil", nil

	default:
		return "", errors.New(fmt.Sprintf("Unable to display: %+v", val))
	}
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

func Parse(data string, parser *participle.Parser) (*BType, error) {
	tree := &BType{}
	err := parser.ParseString(data, tree)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func Prnt(e *BType) (error) {
	val, err := Display(e, true)
	if err != nil {
		return err
	}
	fmt.Print(val + "\n")
	return nil
}

func getParser() (*participle.Parser) {
	graphQLLexer := lexer.Must(stateful.New(stateful.Rules{
		"Root": {
			{"Comment", `;.*(?=(\n|$))`, nil },
			{"Ident", "[^\"^.@~`\\[\\]:{}&'0-9\\s,();][^\"^@~`\\[\\]:{}\\s();]*", nil },
			{"Number", `-?\d+(\.\d+)?`, nil },
			{"Whitespace", `(\s|\t|\n|,|\r)+`, nil },
			{"String", `"(\.|[^"\\])*"`, nil },
			{"Nil", `nil`, nil },
			{"Bool", `true|false`, nil },
		},
	}))

	return participle.MustBuild(&Program{},
	    participle.Lexer(graphQLLexer),
	    participle.Elide("Comment", "Whitespace"),
	)
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

	parser := getParser()


	for {
		line, err := rl.Readline()
		if err != nil { // io.EOF
			log.Print(fmt.Sprintf("Exception: %s", err))
			break
		}

		err1 := Prnt(Parse(line, parser))
		if err1 != nil {
			log.Print(fmt.Sprintf("Exception: %s", err1))
		}

		rl.SaveHistory(line)
	}
}

