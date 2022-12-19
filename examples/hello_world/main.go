package main

import (
	"fmt"
	"time"

	"github.com/fernandezvara/router"
)

func main() {

	r := router.New()

	r.Method("TEST").Insert("hello", func(_ *router.Params) error {
		fmt.Print("Hello!\n")
		return nil
	})
	r.Method("TEST").Insert("hello/:name", helloNameFunc)
	r.Method("TEST").Insert("hello/:name/:surname", helloNameSurnameFunc)

	t := r.Method("TEST")

	time.Sleep(2 * time.Second)
	t.Execute("hello")
	time.Sleep(2 * time.Second)
	t.Execute("hello/Antonio")
	time.Sleep(2 * time.Second)
	t.Execute("hello/Antonio/Fernandez")

}

func helloNameFunc(p *router.Params) error {

	fmt.Printf("Hello, %s!\n", p.Param("name"))
	return nil

}

func helloNameSurnameFunc(p *router.Params) error {

	fmt.Printf("Hello, %s %s!\n", p.Param("name"), p.Param("surname"))
	return nil

}
