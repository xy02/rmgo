package rmgo

import (
	"reflect"
	"strings"
)

type (
	Condition map[Field]Exp

	//Exp means Expression
	Exp map[Operator]Message

	Field string

	Operator int

	Message interface{}

	Unpipe func()
)

func (c Condition) Match(target Message) bool {
	for field, exp := range c {
		if !exp.Match(target, field) {
			return false
		}
	}
	return true
}

func (exp Exp) Match(target Message, f Field) bool {
	value := f.PathValue(target)
	for op, v := range exp {
		if !op.Exec(value, v) {
			return false
		}
	}
	return true
}

func (f Field) PathValue(target Message) Message {
	names := strings.Split(string(f), ".")
	v := reflect.ValueOf(target)
	v = pathValue(v, names, 0)
	if v == (reflect.Value{}) {
		return v
	}
	r := v.Interface()
	return r
}

func pathValue(v reflect.Value, names []string, index int) reflect.Value {
	v = reflect.Indirect(v)
	name := names[index]
	switch v.Kind() {
	case reflect.Struct:
		v = v.FieldByName(name)
	case reflect.Map:
		key := reflect.ValueOf(name)
		v = v.MapIndex(key)
	}
	next := index + 1
	if len(names) == next {
		return v
	}
	return pathValue(v, names, next)
}
