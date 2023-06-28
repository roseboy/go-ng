package main

import (
	"fmt"

	"github.com/roseboy/go-ng/ng"
)

func main() {
	fmt.Println("open 'http://localhost:8000' in browser")
	fmt.Println()

	err := ng.NewServer(8000).RegisterPlugins(&DemoPlugin{}).Start()
	if err != nil {
		panic(err)
	}
}
