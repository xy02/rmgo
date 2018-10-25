package rmgo

import (
	"errors"
	"fmt"
	"reflect"
)

type Filter interface {
	Input(Message) error
	//input message from filter
	InputFrom(Filter, Message) error
	//pipe another Filter to send filtered message
	Pipe(condtion *Condition, filter Filter, savedMessages interface{}) error
	//pipe another Filter to send filtered message
	Pipe2(condtion *Condition, filter Filter) error
	//unpipe
	Unpipe(*Condition)
	//destroy the Filter
	Destroy()
}

type MessagerFactory struct {
	MessageKeeperFactory
}

func (f MessagerFactory) Create() Filter {
	return &Messager{
		keeper: f.MessageKeeperFactory.Create(),
	}
}

type MessageType interface{}

type Messager struct {
	lastMsg   Message
	collector Collector
	keeper    Keeper
}

func (f *Messager) InputFrom(filter Filter, msg Message) error {
	return f.Input(msg)
}

func (f *Messager) Input(msg Message) error {
	fmt.Printf("%p onInput: %+v \n", f, msg)
	if msg == nil {
		return errors.New("message can not be nil")
	}
	f.lastMsg = msg
	//save msg
	if f.collector != nil {
		err := f.collector.Insert(msg)
		if err != nil {
			return err
		}
	}
	f.Broadcast(msg)
	return nil
}

//Broadcast msg to all piped filter which matchs condition
func (f *Messager) Broadcast(msg Message) {
	filters := f.keeper.Select(msg)
	for _, filter := range filters {
		filter.InputFrom(f, msg)
	}
}

func (f *Messager) Pipe2(conditionPtr *Condition, filter Filter) error {
	return f.Pipe(conditionPtr, filter, nil)
}

func (f *Messager) Pipe(conditionPtr *Condition, filter Filter, savedMessages interface{}) error {
	if conditionPtr == nil {
		return errors.New("condition is nil")
	}
	if savedMessages == nil {
		//ouput last message from memory
		if f.lastMsg != nil && conditionPtr.Match(f.lastMsg) {
			err := filter.InputFrom(f, f.lastMsg)
			if err != nil {
				return err
			}
		}
	} else if f.collector != nil {
		//output saved message from database
		//verify
		v := reflect.ValueOf(savedMessages)
		if v.Kind() != reflect.Ptr {
			return errors.New("savedMessages must be a slice address")
		}
		v = v.Elem()
		if v.Kind() != reflect.Slice {
			return errors.New("savedMessages must be a slice address")
		}
		err := f.collector.Select(*conditionPtr, savedMessages)
		if err != nil {
			return err
		}
		for i, length := 0, v.Len(); i < length; i++ {
			msg := v.Index(i).Interface()
			err = filter.InputFrom(f, msg)
			if err != nil {
				return err
			}
		}
	}
	f.keeper.Insert(conditionPtr, filter)
	return nil
}

func (f *Messager) Unpipe(conditionPtr *Condition) {
	f.keeper.Remove(conditionPtr)
}

func (f *Messager) Destroy() {
	f.keeper.Destroy()
}
