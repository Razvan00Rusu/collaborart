package main

import (
	"collaborart/backend"
	"collaborart/frontend"
	"fmt"
)

func main() {
	fmt.Println("Hello World!")
	frontend.Frontend()
	backend.StartServer()
}
