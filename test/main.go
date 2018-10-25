package main

import (
	"fmt"
	"time"

	"github.com/xy02/rmgo"
)

func main() {
	fac := rmgo.MessagerFactory{}
	a := fac.Create()
	a.Input(Inner{Name: "ok"})
	b := fac.Create()
	fmt.Printf("a: %p  b: %p\n", a, b)
	a.Pipe2(&rmgo.Condition{"Name": rmgo.Is("ok")}, b)
	time.Sleep(time.Second)
	a.Input(Inner{Name: "ok"})
}

type Inner struct {
	Name string
}
