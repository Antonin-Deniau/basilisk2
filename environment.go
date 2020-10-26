package main

import (
	"fmt"
)

type Env struct {
	Outer *Env
	Vals map[string]*BaslType
}

func NewEnv(outer *Env, binds []Name, exprs []BaslType) *Env {
	e := &Env{
		Outer: outer,
		Vals: make([]BaslType),
	}

	if len(binds) != len(exprs) {
		if val, ok := binds[Name("&")]; !ok {
			raise BaslException(fmt.Sprintf("Function should contain %d parametter", len(binds)))
		}

		if len(exprs) < len(binds) - 2 {
			raise BaslException(fmt.Sprintf("Function should contain at least %d parametter", len(binds) - 2))
		}
	}

	for i in zip(binds, exprs) {
		if Name("&") == i[0] { break }
		e.Set(i[0], i[1])
	}

	if Name("&") in binds {
		e.Set(binds[-1], tuple(exprs[len(binds) - 2::]))
	}
}

func (e Env) Set(name Name, value BaslType) {
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

func (e Env) Get(name) {
	env = self.find(name)

	if env is not None {
		return env.vals[name]
	} else {
		raise BaslException("{} not found".format("'{}'".format(name)))
	}
}
