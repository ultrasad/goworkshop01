package controllers

import (
	"fmt"

	"github.com/labstack/echo"
)

//CallHelloRoutine is call hello with go routine
func CallHelloRoutine(c echo.Context) (err error) {
	fmt.Println("First Call")

	name := "Hanajung"
	//go hello(name)

	CallHello := make(chan string)
	go hello(name, CallHello)

	fmt.Println("Finish Call")
	//time.Sleep(time.Second)

	fmt.Println("Call from Hello: ", <-CallHello)

	return nil
}

func hello(name string, result chan<- string) {
	output := "You name: " + name
	fmt.Printf("In fnc Hello: %s\n", output)
	result <- output
}
