package main

import (
	"fmt"
)

type t struct {
	f string `this is not a valid struct tag`
}

func (*t) MarshalJSON() string { // nonstandard definition
	return "This is not valid JSON (probably)."
}

func main() {
	fmt.Printf("Hello, %s!\n") // missing string argument
	return
	fmt.Println("Are you still there?") // unreachable
}
