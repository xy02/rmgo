package rmgo

import (
	"errors"
	"reflect"
)

type Collector interface {
	Insert(Message) error
	Select(Condition, result interface{}) error
	Disconnect()
	Delete()
}

type MessageCollector struct {
	lastMsg Message
}

func (mc MessageCollector) Insert(msg Message) error {
	if msg == nil {
		return errors.New("message can not be nil")
	}
	mc.lastMsg = msg
	return nil
}

var ErrBadArg = errors.New("savedMessages argument must be a slice address")

func (mc MessageCollector) Select(condition Condition, savedMessages interface{}) error {
	//verify
	v := reflect.ValueOf(savedMessages)
	if v.Kind() != reflect.Ptr {
		return ErrBadArg
	}
	v = v.Elem()
	if v.Kind() != reflect.Slice {
		return ErrBadArg
	}

	return errors.New("not implement")
}
