package main

import (
	"os/user"
	"os"
	"fmt"
	"strings"
	"log"
	"github.com/chzyer/readline"
)

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

/*
	// REPL
	replEnv := Env(nil, make([]BName), make([]BType))
	for k, v := range ns {
		replEnv.Set(k, v)
	}

	replEnv.set("eval", NewBFunc(func(e string) {
		return Evl(e string, replEnv)
	}))
	replEnv.set("*ARGV*", NewBList(os.Args[1:]))
	replEnv.set("*host-language*", "basilisk")

	LoadStr("(def! not (fn* (a) (if a false true)))", replEnv)
	LoadStr('(def! load-file (fn* (f) (eval (read-string (str "(do " (slurp f) "\\nnil)")))))', replEnv)
	LoadStr("(defmacro! cond (fn* (& xs) (if (> (count xs) 0) (list 'if (first xs) (if (> (count xs) 1) (nth xs 1) (throw \"odd number of forms to cond\")) (cons 'cond (rest (rest xs)))))))", replEnv)
*/

	if len(os.Args) >= 2 {
		//LoadStr(fmt.Sprintf("(load-file \"%s\" )", Escape(os.Args[1])), replEnv)
		
		// TEST PARSER
		dat, err := os.ReadFile(os.Args[2])
	    if err != nil {
	        panic(err)
	    }
	    
    	bexpr, parse_err := Parse(string(dat))
	    if parse_err != nil {
	        panic(parse_err)
	    }

		var sb strings.Builder
		disp_err := Display(bexpr, &sb, true)
		if disp_err != nil {
			fmt.Printf("ERROR %s\n", disp_err)
			return
		}

		fmt.Printf("BType expression: %+v\n", sb.String())

	} else {
		// LoadStr('(println (str "Mal [" *host-language* "]"))', replEnv)

		for {
			line, err := rl.Readline()
			if err != nil { // io.EOF
				log.Print(fmt.Sprintf("Exception: %s", err))
				break
			}

			//err2 := Rep(line, replEnv)
			err2 := Rep(line)

			if err2 != nil {
				log.Print(fmt.Sprintf("Exception: %s", err2))
			}

			rl.SaveHistory(line)
		}
	}
}

func Read(e string) (*BType, error) {
	bexpr, parse_err := Parse(e)
    if parse_err != nil {
        return nil, parse_err
    }

	return bexpr, nil
}

func Print(e *BType) error {
	var sb strings.Builder

	 err := Display(e, &sb, true)
	 if err != nil {
	 	return err
	 }

	os.Stdout.Write([]byte(sb.String()))
	os.Stdout.Write([]byte("\n"))

	return nil
}

func Rep(e string) (error) {
//func Rep(e string, env *Env) (error) {
	b, err := Read(fmt.Sprintf("%s \n", e))
	if err != nil {
		return err
	}

	/*
	c, err := Evl(b, env)
	if err != nil {
		return err
	}*/

	print_err := Print(b)
	if print_err != nil {
		return print_err
	}

	return nil
}

/*
func Rep(e string, env *Env) (error) {
	b, err := Read(e)
	if err != nil {
		return err
	}

	c, err := Evl(b, env)
	if err != nil {
		return err
	}

	Print(c)

	return nil
}

func LoadStr(e string, env *Env) {
	b, err := Read(e)

	if err != nil {
		log.Print(fmt.Sprintf("Error in LoadString,Read: %s", err))
		return
	}


	_, err := Evl(b, env)
	if err != nil {
		log.Print(fmt.Sprintf("Error in LoadString,Evl: %s", err))
		return
	}
}


*/
