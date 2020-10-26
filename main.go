package main

import (
	"os/user"
	"os/args"
	"os"
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

	// REPL
	replEnv := Env(nil, make([]BName), make([]BaslType))
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

	if len(os.Args) >= 2 {
		LoadStr(fmt.Sprintf("(load-file \"%s\" )", Escape(os.Args[1])), replEnv)
	} else {
		LoadStr('(println (str "Mal [" *host-language* "]"))', replEnv)

		for {
			line, err := rl.Readline()
			if err != nil { // io.EOF
				log.Print(fmt.Sprintf("Exception: %s", err))
				break
			}

			err2 := Rep(line, replEnv)
			if err2 != nil {
				log.Print(fmt.Sprintf("Exception: %s", err2))
			}

			rl.SaveHistory(line)
		}
	}

}

func Read(e string) (error) {
	return Parse(e)
}

func Print(e string) {
	os.Stdout.Write(Display(e, true))
	os.Stdout.Write("\n")
}

func Rep(e string, env *Env) (error) {
	b, err := Read(e)
	if err != nil {
		return err
	}

	c, err2 := Evl(b, env)
	if err2 != nil {
		return err2
	}

	Print(c)

	return nil
}

func LoadStr(e string) {
	b := Read(e)
	Evl(b, env)
}

