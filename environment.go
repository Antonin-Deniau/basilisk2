package main

/*
import (
	"fmt"
)

type Env struct {
	Outer *Env
	Vals map[string]*BType
}

func NewEnv(outer *Env, binds []BName, exprs []BType) *Env {
	e := &Env{
		Outer: outer,
		Vals: make([]BType),
	}

	if len(binds) != len(exprs) {
		if val, ok := binds[BName("&")]; !ok {
			raise BaslException(fmt.Sprintf("Function should contain %d parametter", len(binds)))
		}

		if len(exprs) < len(binds) - 2 {
			raise BaslException(fmt.Sprintf("Function should contain at least %d parametter", len(binds) - 2))
		}
	}

	for i in zip(binds, exprs) {
		if BName("&") == i[0] { break }
		e.Set(i[0], i[1])
	}

	if BName("&") in binds {
		e.Set(binds[-1], tuple(exprs[len(binds) - 2::]))
	}
}

func (e Env) Set(name BName, value BType) {
	e.vals[name] = value
	return value
}

func (e Env) Find(name) *Env {
	if val, ok := e.vals[name]; ok {
		return e
	} else {
		if e.Outer != nil {
			return e.Outer.Find(name)
		} else {
			return nil
		}
	}
}

func (e Env) Get(name) (BType, error) {
	env = self.find(name)

	if env != nil {
		return env.vals[name]
	} else {
		return nil, BException(fmt.Sprintf("'%s' not found", name))
	}
}
*/