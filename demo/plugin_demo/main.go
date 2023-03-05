package main

import (
	"fmt"
	"github.com/roseboy/go-ng/ng"
)

func main() {
	fmt.Println("open 'http://localhost:8000' in browser")
	fmt.Println()

	err := ng.NewServer().RegisterPlugins(&DemoPlugin{}).Start(8000)
	if err != nil {
		panic(err)
	}
}
