# router

Router is a simple and lightweight request router library that is not especific for HTTP requests.

It supports routing variables that match the pattern. `Method` is whatever fits your project.

Idea comes after the usage of [httprouter](https://github.com/julienschmidt/httprouter) for many projects and need a simpler pattern matching router for a non-HTTP server. _Thanks for the inspiration!_

## Usage

Let's make a typical `hello world!`:

```
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
```


**NOTE**: There is no explicit matches prevention, so `example/:a` and `example/:b` with different handlers will make unexpected results (maps are *not ordered*, expect funny errors). 