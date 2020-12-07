package main

import "fmt"

func main() {
	d := &Dispatcher{}
	d.Init("arni")

	handler1 := func(e Event) error {
		fmt.Println("In Handler1")
		fmt.Println(e.Value().(string))
		return nil
	}
	handler2 := func(e Event) error {
		fmt.Println("In Handler2")
		fmt.Println(e.Value().(string))
		return nil
	}

	d.Register(MessageEvent, handler1)
	d.Register(MessageEvent, handler2)

	d.Dispatch(MessageEvent, "First Hello World!")
	d.Remove(MessageEvent, handler1)

	d.Dispatch(MessageEvent, "Second Hello World!")

}
