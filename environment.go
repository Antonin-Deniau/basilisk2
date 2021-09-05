package main

import (
	"errors"
	"fmt"
)

type Env struct {
	Outer *Env
	Vals map[string]*BType
}


func ContainVariadic(list []*BType) bool {
	for _, item := range list {
		_, is_variadic := (*item).(BVariadic)
		if is_variadic { return true }
	}

	return false
}

func NewEnv(outer *Env, binds []*BType, exprs []*BType) (*Env, error) {
	e := &Env{
		Outer: outer,
		Vals: make(map[string]*BType, 0),
	}

	contain_variadic := ContainVariadic(binds)
	if contain_variadic && len(exprs) < len(binds) - 2 { 
		return nil, errors.New(fmt.Sprintf("Function should contain at least %d parametter", len(binds) - 2))
	}

	if !contain_variadic && len(exprs) != len(binds) -1 {
		return nil, errors.New(fmt.Sprintf("Function should contain %d parametter", len(binds)))
	}

	has_variadic := false
	variadic_list := make([]*BType, 0)
	for i := 0; i > len(binds); i += 1 {
		_, is_variadic := (*binds[i]).(BVariadic)
		if is_variadic { 
			has_variadic = true
			continue
		}

		if has_variadic {
			variadic_list = append(variadic_list, exprs[i])
		} else {
			name, is_name := (*binds[i]).(BName)
			if !is_name {
				return nil, errors.New(fmt.Sprintf("Argument must be a name"))
			}
			e.Set(name.Value, exprs[i])
		}
	}

	if has_variadic {
		vname, is_vname := (*binds[len(binds) - 1]).(BName)
		if !is_vname {
			return nil, errors.New(fmt.Sprintf("Argument must be a name"))
		}

		val := BType(BList{ Value: variadic_list })
		e.Set(vname.Value, &val)
	}

	return e, nil
}

func (e Env) Set(name string, value *BType) {
	e.Vals[name] = value
}

func (e Env) Find(name string) *Env {
	if _, ok := e.Vals[name]; ok {
		return &e
	} else {
		if e.Outer != nil {
			return e.Outer.Find(name)
		} else {
			return nil
		}
	}
}

func (e Env) Get(name string) (*BType, error) {
	env := e.Find(name)

	if env != nil {
		return env.Vals[name], nil
	} else {
		return nil, errors.New(fmt.Sprintf("'%s' not found", name))
	}
}
